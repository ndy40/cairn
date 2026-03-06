package model

import (
	"fmt"
	"os/exec"
	"runtime"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndy40/cairn/internal/store"
)

// BookmarkItem wraps a store.Bookmark to satisfy list.Item.
type BookmarkItem struct {
	b *store.Bookmark
}

func (i BookmarkItem) Title() string {
	title := i.b.Title
	if title == "" {
		title = i.b.URL
	}
	if i.b.IsPermanent {
		return "[pin] " + title
	}
	return title
}

func (i BookmarkItem) Description() string {
	date := i.b.CreatedAt.Format("2006-01-02")
	desc := fmt.Sprintf("%s · %s", i.b.Domain, date)

	if len(i.b.Tags) > 0 {
		tagStr := ""
		for _, t := range i.b.Tags {
			tagStr += " #" + t
		}
		desc += " ·" + tagStr
	}
	return desc
}

func (i BookmarkItem) FilterValue() string { return i.b.Title + " " + i.b.Domain }

// BrowseModel wraps a bubbles/list for the bookmark browse view.
type BrowseModel struct {
	list            list.Model
	deleteRequested bool
	openErrMsg      string
	extraKeys       []key.Binding
}

var (
	jumpTopKey = key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "jump to top"),
	)
	jumpBottomKey = key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "jump to bottom"),
	)
	deleteKey = key.NewBinding(
		key.WithKeys("d", "delete"),
		key.WithHelp("d", "delete"),
	)
	openKey = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open in browser"),
	)
	permanentKey = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "toggle permanent"),
	)
	editKey = key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit tags"),
	)
)

func newBrowseModel() BrowseModel {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true

	l := list.New(nil, delegate, 0, 0)
	l.Title = "Bookmarks"
	l.SetShowHelp(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false) // We use custom fuzzy search.
	l.DisableQuitKeybindings()

	// Add vim-style and jump keys.
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{jumpTopKey, jumpBottomKey, deleteKey, permanentKey, editKey}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{jumpTopKey, jumpBottomKey, deleteKey, openKey, permanentKey, editKey}
	}

	return BrowseModel{list: l}
}

// load replaces the list items with the given bookmarks.
func (m *BrowseModel) load(bookmarks []*store.Bookmark) {
	items := make([]list.Item, len(bookmarks))
	for i, b := range bookmarks {
		items[i] = BookmarkItem{b: b}
	}
	m.list.SetItems(items)
}

// setSize resizes the list.
func (m *BrowseModel) setSize(w, h int) {
	m.list.SetSize(w, h)
}

// selected returns the currently highlighted bookmark, or nil.
func (m *BrowseModel) selected() *store.Bookmark {
	item, ok := m.list.SelectedItem().(BookmarkItem)
	if !ok {
		return nil
	}
	return item.b
}

// Update handles key events for the browse list.
func (m BrowseModel) Update(msg tea.Msg) (BrowseModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "g":
			m.list.Select(0)
			return m, nil
		case "G":
			m.list.Select(len(m.list.Items()) - 1)
			return m, nil
		case "d", "delete":
			m.deleteRequested = true
			return m, nil
		case "enter":
			// Handled at App level via openBookmarkCmd so the store is accessible.
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the list or empty state.
func (m BrowseModel) View() string {
	if len(m.list.Items()) == 0 {
		return lipgloss.Place(
			m.list.Width(), m.list.Height(),
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(
				"No bookmarks yet.\nPress Ctrl+P to add your first bookmark.",
			),
		)
	}
	return m.list.View()
}

// openURLErrMsg is sent when the browser open command fails.
type openURLErrMsg struct{ err error }

// openURLRaw runs the OS-appropriate command to open a URL and returns any error.
func openURLRaw(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

