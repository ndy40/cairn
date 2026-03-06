package model

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndy40/cairn/internal/store"
)

// tagItem wraps a tag string as a list.Item.
type tagItem struct {
	tag      string
	selected bool
}

func (i tagItem) Title() string {
	if i.selected {
		return "[x] " + i.tag
	}
	return "[ ] " + i.tag
}
func (i tagItem) Description() string { return "" }
func (i tagItem) FilterValue() string { return i.tag }

// TagFilterModel is the tag selection overlay.
type TagFilterModel struct {
	list     list.Model
	selected map[string]bool
	allTags  []string
}

func newTagFilterModel(bookmarks []*store.Bookmark, current []string) TagFilterModel {
	// Collect unique tags from active bookmarks.
	seen := make(map[string]bool)
	for _, b := range bookmarks {
		for _, t := range b.Tags {
			seen[t] = true
		}
	}
	var allTags []string
	for t := range seen {
		allTags = append(allTags, t)
	}
	sort.Strings(allTags)

	// Build selected set from current active filter.
	selected := make(map[string]bool, len(current))
	for _, t := range current {
		selected[t] = true
	}

	// Build list items.
	items := make([]list.Item, len(allTags))
	for i, t := range allTags {
		items[i] = tagItem{tag: t, selected: selected[t]}
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	l := list.New(items, delegate, 40, 15)
	l.Title = "Tag Filter  [Space/Enter] toggle  [c] clear  [Esc] close"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return TagFilterModel{list: l, selected: selected, allTags: allTags}
}

func (m TagFilterModel) setSize(w, h int) TagFilterModel {
	m.list.SetSize(w, h)
	return m
}

// SelectedTags returns the currently selected tag names.
func (m TagFilterModel) SelectedTags() []string {
	var out []string
	for _, t := range m.allTags {
		if m.selected[t] {
			out = append(out, t)
		}
	}
	return out
}

// Update handles key events for the tag filter overlay.
func (m TagFilterModel) Update(msg tea.Msg) (TagFilterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "enter":
			// Toggle the selected tag.
			item, ok := m.list.SelectedItem().(tagItem)
			if !ok {
				break
			}
			item.selected = !item.selected
			m.selected[item.tag] = item.selected
			// Refresh the list item in place.
			idx := m.list.Index()
			items := m.list.Items()
			items[idx] = item
			cmd := m.list.SetItems(items)
			return m, cmd
		case "c":
			// Clear all selections.
			m.selected = make(map[string]bool)
			items := m.list.Items()
			for i, it := range items {
				ti := it.(tagItem)
				ti.selected = false
				items[i] = ti
			}
			cmd := m.list.SetItems(items)
			return m, cmd
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the tag filter overlay.
func (m TagFilterModel) View() string {
	if len(m.allTags) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("No tags found. Add tags to your bookmarks first.")
	}
	return m.list.View()
}

// filterStatusLine returns a footer note showing active tag filters, or empty string.
func filterStatusLine(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	line := "Filter:"
	for _, t := range tags {
		line += fmt.Sprintf(" #%s", t)
	}
	return line
}
