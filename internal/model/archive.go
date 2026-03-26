package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndy40/cairn/internal/store"
)

// ArchiveBookmarkItem wraps a store.Bookmark for the archive list.
type ArchiveBookmarkItem struct {
	b *store.Bookmark
}

func (i ArchiveBookmarkItem) Title() string {
	if i.b.Title != "" {
		return i.b.Title
	}
	return i.b.URL
}

func (i ArchiveBookmarkItem) Description() string {
	archivedOn := ""
	if i.b.ArchivedAt != nil {
		archivedOn = i.b.ArchivedAt.Format("2006-01-02")
	}
	lastVisited := "Never visited"
	if i.b.LastVisitedAt != nil {
		lastVisited = "Last: " + i.b.LastVisitedAt.Format("2006-01-02")
	}
	return fmt.Sprintf("%s · %s · Archived: %s", i.b.Domain, lastVisited, archivedOn)
}

func (i ArchiveBookmarkItem) FilterValue() string { return i.b.Title + " " + i.b.Domain }

// ArchiveModel is the archive list view.
type ArchiveModel struct {
	list             list.Model
	restoreRequested bool
}

func newArchiveModel() ArchiveModel {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true

	l := list.New(nil, delegate, 0, 0)
	l.Title = "Archive"
	l.SetShowHelp(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return ArchiveModel{list: l}
}

func (m *ArchiveModel) load(bookmarks []*store.Bookmark) {
	items := make([]list.Item, len(bookmarks))
	for i, b := range bookmarks {
		items[i] = ArchiveBookmarkItem{b: b}
	}
	m.list.SetItems(items)
}

func (m *ArchiveModel) setSize(w, h int) {
	m.list.SetSize(w, h)
}

func (m *ArchiveModel) selected() *store.Bookmark {
	item, ok := m.list.SelectedItem().(ArchiveBookmarkItem)
	if !ok {
		return nil
	}
	return item.b
}

// Update handles key events for the archive list.
func (m ArchiveModel) Update(msg tea.Msg) (ArchiveModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.restoreRequested = true
			return m, nil
		case "esc":
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the archive list or empty state.
func (m ArchiveModel) View() string {
	if len(m.list.Items()) == 0 {
		return lipgloss.Place(
			m.list.Width(), m.list.Height(),
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(
				"No archived bookmarks.",
			),
		)
	}
	return m.list.View()
}
