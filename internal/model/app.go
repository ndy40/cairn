package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndy40/cairn/internal/clipboard"
	"github.com/ndy40/cairn/internal/fetcher"
	"github.com/ndy40/cairn/internal/search"
	"github.com/ndy40/cairn/internal/store"
)

// AppState represents the active interaction mode.
type AppState int

const (
	StateBrowse AppState = iota
	StateSearch
	StateAdd
	StateConfirmDelete
	StateTagFilter
	StateArchive
	StateEdit
)

// Internal tea.Msg types.

type bookmarksLoadedMsg struct{ bookmarks []*store.Bookmark }
type archivedBookmarksLoadedMsg struct{ bookmarks []*store.Bookmark }
type bookmarkDeletedMsg struct{ id int64 }
type bookmarkUpdatedMsg struct{}
type deleteErrorMsg struct{ err error }
type fetchSaveResultMsg struct {
	fetchErr bool
	isDup    bool
	err      error
}

// App is the root bubbletea model.
type App struct {
	store  *store.Store
	state  AppState
	width  int
	height int

	// Sub-models.
	browse  BrowseModel
	srch    SearchModel
	add     AddModel
	archive ArchiveModel

	startupArchiveCount int

	// All bookmarks in memory (source of truth for searching).
	allBookmarks    []*store.Bookmark
	activeTagFilter []string
	tagFilterModel  TagFilterModel

	// Confirm-delete state.
	deleteID    int64
	deleteTitle string

	// Edit-tags state.
	editModel EditModel
	editID    int64

	// Footer status message (transient).
	footerMsg string

	// Help overlay.
	showHelp bool
}

// New creates a new App model with the given store and startup archive count.
func New(s *store.Store, archiveCount int) App {
	app := App{
		store:               s,
		state:               StateBrowse,
		browse:              newBrowseModel(),
		add:                 newAddModel(),
		archive:             newArchiveModel(),
		startupArchiveCount: archiveCount,
	}
	if archiveCount > 0 {
		app.footerMsg = fmt.Sprintf("%d bookmark(s) archived", archiveCount)
	}
	return app
}

// Init loads bookmarks from the store on startup.
func (a App) Init() tea.Cmd {
	return loadBookmarks(a.store)
}

// Update handles all incoming messages.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		listH := msg.Height - 2 // subtract footer row
		a.browse.setSize(msg.Width, listH)
		a.srch.setSize(msg.Width, listH-2)
		return a, nil

	case bookmarksLoadedMsg:
		a.allBookmarks = msg.bookmarks
		a.browse.load(applyTagFilter(msg.bookmarks, a.activeTagFilter))
		return a, nil

	case fetchSaveResultMsg:
		switch {
		case msg.isDup:
			a.add.setStatus("Already bookmarked")
		case msg.err != nil:
			a.add.setStatus("Error: " + msg.err.Error())
		case msg.fetchErr:
			a.state = StateBrowse
			return a, loadBookmarks(a.store)
		default:
			a.state = StateBrowse
			return a, loadBookmarks(a.store)
		}
		return a, nil

	case archivedBookmarksLoadedMsg:
		a.archive.load(msg.bookmarks)
		a.archive.setSize(a.width, a.height-2)
		return a, nil

	case bookmarkUpdatedMsg:
		if a.state == StateArchive {
			return a, loadArchivedBookmarks(a.store)
		}
		return a, loadBookmarks(a.store)

	case bookmarkDeletedMsg:
		a.state = StateBrowse
		return a, loadBookmarks(a.store)

	case deleteErrorMsg:
		a.footerMsg = "Delete failed: " + msg.err.Error()
		a.state = StateBrowse
		return a, nil

	case openURLErrMsg:
		a.footerMsg = "Could not open: " + msg.err.Error()
		return a, nil

	case tea.KeyMsg:
		a.footerMsg = ""

		// Global quit.
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}

		// Help overlay toggle.
		if a.showHelp {
			a.showHelp = false
			return a, nil
		}
		if msg.String() == "?" && a.state != StateAdd {
			a.showHelp = true
			return a, nil
		}

		// Global Ctrl+P — paste from clipboard.
		if msg.String() == "ctrl+p" && (a.state == StateBrowse || a.state == StateSearch) {
			text, err := clipboard.Read()
			if err != nil {
				a.footerMsg = err.Error()
				return a, nil
			}
			a.add = newAddModel()
			a.add.setURL(text)
			a.state = StateAdd
			return a, a.add.Init()
		}

		switch a.state {
		case StateBrowse:
			return a.updateBrowse(msg)
		case StateSearch:
			return a.updateSearch(msg)
		case StateAdd:
			return a.updateAdd(msg)
		case StateConfirmDelete:
			return a.updateConfirmDelete(msg)
		case StateTagFilter:
			return a.updateTagFilter(msg)
		case StateArchive:
			return a.updateArchive(msg)
		case StateEdit:
			return a.updateEdit(msg)
		}
	}

	// Delegate non-key messages to active sub-models.
	switch a.state {
	case StateBrowse:
		bm, cmd := a.browse.Update(msg)
		a.browse = bm
		return a, cmd
	case StateSearch:
		sm, cmd := a.srch.Update(msg)
		a.srch = sm
		return a, cmd
	case StateTagFilter:
		tf, cmd := a.tagFilterModel.Update(msg)
		a.tagFilterModel = tf
		return a, cmd
	case StateArchive:
		ar, cmd := a.archive.Update(msg)
		a.archive = ar
		return a, cmd
	case StateEdit:
		em, cmd := a.editModel.Update(msg)
		a.editModel = em
		return a, cmd
	}
	return a, nil
}

func (a App) updateBrowse(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "/":
		a.state = StateSearch
		a.srch = newSearchModel(a.allBookmarks)
		a.srch.setSize(a.width, a.height-4)
		return a, a.srch.Init()
	case "p":
		if sel := a.browse.selected(); sel != nil {
			id := sel.ID
			permanent := !sel.IsPermanent
			return a, func() tea.Msg {
				if err := a.store.SetPermanent(id, permanent); err != nil {
					return deleteErrorMsg{err: err}
				}
				return bookmarkUpdatedMsg{}
			}
		}
		return a, nil
	case "e":
		if sel := a.browse.selected(); sel != nil {
			a.editModel = newEditModel(sel)
			a.editID = sel.ID
			a.state = StateEdit
			return a, textinput.Blink
		}
		return a, nil
	case "t":
		a.tagFilterModel = newTagFilterModel(a.allBookmarks, a.activeTagFilter)
		a.tagFilterModel = a.tagFilterModel.setSize(a.width, a.height-4)
		a.state = StateTagFilter
		return a, nil
	case "a":
		a.state = StateArchive
		return a, loadArchivedBookmarks(a.store)
	case "enter":
		if sel := a.browse.selected(); sel != nil {
			return a, openBookmarkCmd(a.store, sel)
		}
		return a, nil
	default:
		bm, cmd := a.browse.Update(msg)
		a.browse = bm
		if a.browse.deleteRequested {
			if sel := a.browse.selected(); sel != nil {
				a.deleteID = sel.ID
				a.deleteTitle = sel.Title
				a.state = StateConfirmDelete
				a.browse.deleteRequested = false
			}
		}
		return a, cmd
	}
}

func (a App) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.state = StateBrowse
		a.browse.load(a.allBookmarks)
		return a, nil
	case "enter":
		if sel := a.browse.selected(); sel != nil {
			return a, openBookmarkCmd(a.store, sel)
		}
		return a, nil
	default:
		sm, cmd := a.srch.Update(msg)
		a.srch = sm
		// Two-stage search: FTS5 pre-filter then fuzzy ranking.
		filtered := a.twoStageSearch(a.srch.term())
		a.browse.load(filtered)
		return a, cmd
	}
}

// twoStageSearch applies tag filter, then FTS5 pre-filtering (for terms >= 3 chars),
// then fuzzy ranks the candidates.
func (a *App) twoStageSearch(term string) []*store.Bookmark {
	base := applyTagFilter(a.allBookmarks, a.activeTagFilter)
	if term == "" {
		return base
	}
	var candidates []*store.Bookmark
	if len([]rune(term)) >= 3 {
		ids, err := a.store.FTSSearch(term)
		if err == nil && len(ids) > 0 {
			idSet := make(map[int64]bool, len(ids))
			for _, id := range ids {
				idSet[id] = true
			}
			for _, b := range base {
				if idSet[b.ID] {
					candidates = append(candidates, b)
				}
			}
		}
	}
	if len(candidates) == 0 {
		candidates = base
	}
	return search.Search(term, candidates)
}

// applyTagFilter returns bookmarks matching any of the selected tags (OR logic).
// Returns the full slice unchanged when no tags are selected.
func applyTagFilter(bookmarks []*store.Bookmark, tags []string) []*store.Bookmark {
	if len(tags) == 0 {
		return bookmarks
	}
	tagSet := make(map[string]bool, len(tags))
	for _, t := range tags {
		tagSet[t] = true
	}
	var out []*store.Bookmark
	for _, b := range bookmarks {
		for _, bt := range b.Tags {
			if tagSet[bt] {
				out = append(out, b)
				break
			}
		}
	}
	return out
}

func (a App) updateAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.state = StateBrowse
		return a, nil
	case "enter":
		rawURL := a.add.URL()
		if rawURL == "" {
			a.add.setStatus("Please enter a URL")
			return a, nil
		}
		rawTags := store.NormaliseTagsFromString(a.add.Tags())
		return a, fetchAndSave(a.store, rawURL, rawTags)
	default:
		am, cmd := a.add.Update(msg)
		a.add = am
		return a, cmd
	}
}

func (a App) updateArchive(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.state = StateBrowse
		return a, nil
	case "r":
		if sel := a.archive.selected(); sel != nil {
			id := sel.ID
			return a, func() tea.Msg {
				if err := a.store.RestoreByID(id); err != nil {
					return deleteErrorMsg{err: err}
				}
				return bookmarkUpdatedMsg{}
			}
		}
		return a, nil
	default:
		ar, cmd := a.archive.Update(msg)
		a.archive = ar
		if a.archive.restoreRequested {
			a.archive.restoreRequested = false
			if sel := a.archive.selected(); sel != nil {
				id := sel.ID
				return a, func() tea.Msg {
					if err := a.store.RestoreByID(id); err != nil {
						return deleteErrorMsg{err: err}
					}
					return bookmarkUpdatedMsg{}
				}
			}
		}
		return a, cmd
	}
}

func (a App) updateTagFilter(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "t":
		a.activeTagFilter = a.tagFilterModel.SelectedTags()
		a.state = StateBrowse
		a.browse.load(applyTagFilter(a.allBookmarks, a.activeTagFilter))
		return a, nil
	case "c":
		tf, cmd := a.tagFilterModel.Update(msg)
		a.tagFilterModel = tf
		a.activeTagFilter = nil
		a.browse.load(a.allBookmarks)
		a.state = StateBrowse
		return a, cmd
	default:
		tf, cmd := a.tagFilterModel.Update(msg)
		a.tagFilterModel = tf
		return a, cmd
	}
}

func (a App) updateConfirmDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		id := a.deleteID
		return a, func() tea.Msg {
			if err := a.store.DeleteByID(id); err != nil {
				return deleteErrorMsg{err: err}
			}
			return bookmarkDeletedMsg{id: id}
		}
	case "n", "esc":
		a.state = StateBrowse
		return a, nil
	}
	return a, nil
}

func (a App) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		id := a.editID
		tags := a.editModel.Tags()
		a.state = StateBrowse
		return a, func() tea.Msg {
			if err := a.store.UpdateTags(id, tags); err != nil {
				return deleteErrorMsg{err: err}
			}
			return bookmarkUpdatedMsg{}
		}
	case "esc":
		a.state = StateBrowse
		return a, nil
	default:
		em, cmd := a.editModel.Update(msg)
		a.editModel = em
		return a, cmd
	}
}

// ── View ──────────────────────────────────────────────────────────────────────

var (
	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingLeft(1)

	errorFooterStyle = footerStyle.Copy().Foreground(lipgloss.Color("196"))

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Width(66)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62")).
			MarginBottom(1)

	dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	searchBarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)
)

func (a App) View() string {
	if a.showHelp {
		return a.helpView()
	}
	switch a.state {
	case StateAdd:
		return a.addView()
	case StateConfirmDelete:
		return a.confirmDeleteView()
	case StateSearch:
		return a.searchView()
	case StateTagFilter:
		return a.tagFilterView()
	case StateArchive:
		return a.archiveView()
	case StateEdit:
		return a.editView()
	default:
		return a.browseView()
	}
}

func (a App) footer(text string) string {
	if a.footerMsg != "" {
		return errorFooterStyle.Render(a.footerMsg)
	}
	return footerStyle.Render(text)
}

func (a App) browseView() string {
	filterLine := ""
	if fs := filterStatusLine(a.activeTagFilter); fs != "" {
		filterLine = "\n" + footerStyle.Render(fs)
	}
	return a.browse.View() + filterLine + "\n" +
		a.footer("[/] Search  [Ctrl+P] Add  [Enter] Open  [e] Edit  [d] Delete  [p] Pin  [t] Tags  [a] Archive  [?] Help  [Ctrl+C] Quit")
}

func (a App) archiveView() string {
	return a.archive.View() + "\n" +
		footerStyle.Render("[r] Restore  [Esc] Back  [Ctrl+C] Quit")
}

func (a App) tagFilterView() string {
	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center,
		modalStyle.Render(a.tagFilterModel.View()))
}

func (a App) searchView() string {
	bar := searchBarStyle.Width(a.width - 2).Render(a.srch.inputView())
	var listView string
	if a.srch.term() != "" && len(a.browse.list.Items()) == 0 {
		listView = a.srch.noResultsView()
	} else {
		listView = a.browse.View()
	}
	footer := footerStyle.Render("[Esc] Clear  [Enter] Open  [Ctrl+P] Add  [Ctrl+C] Quit")
	return bar + "\n" + listView + "\n" + footer
}

func (a App) addView() string {
	content := headerStyle.Render("Add Bookmark") + "\n" +
		a.add.View() + "\n\n" +
		dimStyle.Render("[Enter] Save  [Esc] Cancel")
	modal := modalStyle.Render(content)
	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, modal)
}

func (a App) editView() string {
	content := headerStyle.Render("Edit Tags") + "\n" +
		a.editModel.View() + "\n\n" +
		dimStyle.Render("[Enter] Save  [Esc] Cancel")
	modal := modalStyle.Render(content)
	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, modal)
}

func (a App) confirmDeleteView() string {
	title := a.deleteTitle
	if len([]rune(title)) > 50 {
		runes := []rune(title)
		title = string(runes[:47]) + "..."
	}
	content := headerStyle.Render("Delete Bookmark") + "\n" +
		fmt.Sprintf("Delete «%s»?\n\n", title) +
		dimStyle.Render("[y/Enter] Confirm Delete  [n/Esc] Cancel")
	modal := modalStyle.Render(content)
	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, modal)
}

func (a App) helpView() string {
	content := headerStyle.Render("Keyboard Shortcuts") + "\n\n" +
		"Browse Mode\n" +
		"  /          Enter search\n" +
		"  Ctrl+P     Add bookmark from clipboard\n" +
		"  Enter      Open selected in browser\n" +
		"  e          Edit tags on selected bookmark\n" +
		"  d, Delete  Delete selected bookmark\n" +
		"  p          Toggle permanent (pin) flag\n" +
		"  t          Open tag filter overlay\n" +
		"  a          Open archive view\n" +
		"  j / ↓      Move down\n" +
		"  k / ↑      Move up\n" +
		"  g          Jump to top\n" +
		"  G          Jump to bottom\n\n" +
		"Tag Filter Overlay\n" +
		"  j / ↓      Move down\n" +
		"  k / ↑      Move up\n" +
		"  Space/Enter Toggle selected tag\n" +
		"  c          Clear all tag filters\n" +
		"  Esc / t    Close overlay\n\n" +
		"Archive View\n" +
		"  r          Restore selected bookmark\n" +
		"  Esc        Return to browse\n\n" +
		"Search Mode\n" +
		"  Type       Filter in real time\n" +
		"  Esc        Clear search, return to browse\n" +
		"  Ctrl+A     Clear search term\n" +
		"  Enter      Open selected result\n\n" +
		"Add Mode\n" +
		"  Tab        Next field (URL → Tags)\n" +
		"  Enter      Save bookmark\n" +
		"  Esc        Cancel\n\n" +
		"Delete Confirmation\n" +
		"  y, Enter   Confirm deletion\n" +
		"  n, Esc     Cancel\n\n" +
		"Global\n" +
		"  ?          Toggle this help\n" +
		"  Ctrl+C     Quit\n\n" +
		dimStyle.Render("Press any key to close")
	box := modalStyle.Width(72).Render(content)
	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, box)
}

// ── Commands ──────────────────────────────────────────────────────────────────

// openBookmarkCmd opens the bookmark URL in the default browser and reloads
// the bookmark list on success.
//
// Flow:
//  1. openURLRaw(url) — starts the browser process
//  2. On success: returns loadBookmarks(s)() — reloads the list
//
// On browser start failure, openURLErrMsg is returned.
func openBookmarkCmd(s *store.Store, b *store.Bookmark) tea.Cmd {
	return func() tea.Msg {
		if err := openURLRaw(b.URL); err != nil {
			return openURLErrMsg{err: err}
		}
		return loadBookmarks(s)()
	}
}

func loadBookmarks(s *store.Store) tea.Cmd {
	return func() tea.Msg {
		bookmarks, _ := s.List()
		return bookmarksLoadedMsg{bookmarks: bookmarks}
	}
}

func loadArchivedBookmarks(s *store.Store) tea.Cmd {
	return func() tea.Msg {
		bookmarks, _ := s.ListArchived()
		return archivedBookmarksLoadedMsg{bookmarks: bookmarks}
	}
}

func fetchAndSave(s *store.Store, rawURL string, tags []string) tea.Cmd {
	return func() tea.Msg {
		title, description, fetchErr := fetcher.Fetch(rawURL)
		_, insertErr := s.Insert(rawURL, title, description, tags)
		if insertErr != nil {
			if insertErr == store.ErrDuplicate {
				return fetchSaveResultMsg{isDup: true}
			}
			return fetchSaveResultMsg{err: insertErr}
		}
		return fetchSaveResultMsg{fetchErr: fetchErr != nil}
	}
}
