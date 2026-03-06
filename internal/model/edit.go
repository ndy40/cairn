package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndy40/cairn/internal/store"
)

// EditModel handles the edit-tags panel for an existing bookmark.
type EditModel struct {
	title     string
	tagsInput textinput.Model
}

func newEditModel(b *store.Bookmark) EditModel {
	ti := textinput.New()
	ti.Placeholder = "work, go, tools  (comma-separated, max 3)"
	ti.CharLimit = 200
	ti.Width = 58
	ti.SetValue(strings.Join(b.Tags, ", "))
	ti.Focus()
	return EditModel{title: b.Title, tagsInput: ti}
}

// Update handles keypresses delegated from the parent.
func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd) {
	var cmd tea.Cmd
	m.tagsInput, cmd = m.tagsInput.Update(msg)
	return m, cmd
}

// View renders the read-only title and the editable tags field.
func (m EditModel) View() string {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	titleStyle := lipgloss.NewStyle().Bold(true)
	return titleStyle.Render(m.title) + "\n\n" +
		labelStyle.Render("Tags") + "\n" +
		m.tagsInput.View()
}

// Tags splits the current input value on commas and returns the raw slice.
// Normalisation is applied by the store's UpdateTags method.
func (m EditModel) Tags() []string {
	return strings.Split(m.tagsInput.Value(), ",")
}
