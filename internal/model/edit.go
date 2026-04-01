package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndy40/cairn/internal/store"
)

const (
	editFieldURL  = 0
	editFieldTags = 1
)

// EditModel handles the edit panel for an existing bookmark.
type EditModel struct {
	title       string
	urlInput    textinput.Model
	tagsInput   textinput.Model
	activeField int
}

func newEditModel(b *store.Bookmark) EditModel {
	urlTI := textinput.New()
	urlTI.Placeholder = "https://example.com"
	urlTI.CharLimit = 2048
	urlTI.Width = 58
	urlTI.SetValue(b.URL)
	urlTI.Focus()

	tagsTI := textinput.New()
	tagsTI.Placeholder = "work, go, tools  (comma-separated, max 3)"
	tagsTI.CharLimit = 200
	tagsTI.Width = 58
	tagsTI.SetValue(strings.Join(b.Tags, ", "))
	tagsTI.Blur()

	return EditModel{
		title:       b.Title,
		urlInput:    urlTI,
		tagsInput:   tagsTI,
		activeField: editFieldURL,
	}
}

// Update handles keypresses delegated from the parent.
func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "tab" || keyMsg.String() == "shift+tab" {
			if m.activeField == editFieldURL {
				m.activeField = editFieldTags
				m.urlInput.Blur()
				m.tagsInput.Focus()
			} else {
				m.activeField = editFieldURL
				m.tagsInput.Blur()
				m.urlInput.Focus()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.activeField == editFieldURL {
		m.urlInput, cmd = m.urlInput.Update(msg)
	} else {
		m.tagsInput, cmd = m.tagsInput.Update(msg)
	}
	return m, cmd
}

// View renders the title, URL input, and tags input fields.
func (m EditModel) View() string {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	activeLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	titleStyle := lipgloss.NewStyle().Bold(true)

	urlLabel := labelStyle.Render("URL")
	tagsLabel := labelStyle.Render("Tags")
	if m.activeField == editFieldURL {
		urlLabel = activeLabel.Render("URL")
	} else {
		tagsLabel = activeLabel.Render("Tags")
	}

	return titleStyle.Render(m.title) + "\n\n" +
		urlLabel + "\n" +
		m.urlInput.View() + "\n\n" +
		tagsLabel + "\n" +
		m.tagsInput.View() + "\n\n" +
		labelStyle.Render("Tab: switch fields  Enter: save  Esc: cancel")
}

// URL returns the trimmed URL input value.
func (m EditModel) URL() string {
	return strings.TrimSpace(m.urlInput.Value())
}

// Tags splits the current input value on commas and returns the raw slice.
// Normalisation is applied by the store's UpdateTags method.
func (m EditModel) Tags() []string {
	return strings.Split(m.tagsInput.Value(), ",")
}
