package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndy40/cairn/internal/store"
)

const (
	editFieldTitle = 0
	editFieldURL   = 1
	editFieldTags  = 2
)

// EditModel handles the edit panel for an existing bookmark.
type EditModel struct {
	titleInput  textinput.Model
	titleErr    string
	urlInput    textinput.Model
	tagsInput   textinput.Model
	activeField int
}

func newEditModel(b *store.Bookmark) EditModel {
	titleTI := textinput.New()
	titleTI.Placeholder = "Bookmark title"
	titleTI.CharLimit = 500
	titleTI.Width = 58
	titleTI.SetValue(b.Title)
	titleTI.Focus()

	urlTI := textinput.New()
	urlTI.Placeholder = "https://example.com"
	urlTI.CharLimit = 2048
	urlTI.Width = 58
	urlTI.SetValue(b.URL)
	urlTI.Blur()

	tagsTI := textinput.New()
	tagsTI.Placeholder = "work, go, tools  (comma-separated, max 3)"
	tagsTI.CharLimit = 200
	tagsTI.Width = 58
	tagsTI.SetValue(strings.Join(b.Tags, ", "))
	tagsTI.Blur()

	return EditModel{
		titleInput:  titleTI,
		urlInput:    urlTI,
		tagsInput:   tagsTI,
		activeField: editFieldTitle,
	}
}

// Update handles keypresses delegated from the parent.
func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "tab":
			return m.advanceField(1), nil
		case "shift+tab":
			return m.advanceField(-1), nil
		}
	}

	var cmd tea.Cmd
	switch m.activeField {
	case editFieldTitle:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case editFieldURL:
		m.urlInput, cmd = m.urlInput.Update(msg)
	default:
		m.tagsInput, cmd = m.tagsInput.Update(msg)
	}
	return m, cmd
}

func (m EditModel) advanceField(delta int) EditModel {
	const numFields = 3
	m.activeField = ((m.activeField + delta) + numFields) % numFields
	m.titleInput.Blur()
	m.urlInput.Blur()
	m.tagsInput.Blur()
	switch m.activeField {
	case editFieldTitle:
		m.titleInput.Focus()
	case editFieldURL:
		m.urlInput.Focus()
	default:
		m.tagsInput.Focus()
	}
	return m
}

// View renders the title, URL, and tags input fields.
func (m EditModel) View() string {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	activeLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	titleLabel := labelStyle.Render("Title")
	urlLabel := labelStyle.Render("URL")
	tagsLabel := labelStyle.Render("Tags")

	switch m.activeField {
	case editFieldTitle:
		titleLabel = activeLabel.Render("Title")
	case editFieldURL:
		urlLabel = activeLabel.Render("URL")
	default:
		tagsLabel = activeLabel.Render("Tags")
	}

	titleSection := titleLabel + "\n" + m.titleInput.View()
	if m.titleErr != "" {
		titleSection += "\n" + errStyle.Render(m.titleErr)
	}

	return titleSection + "\n\n" +
		urlLabel + "\n" +
		m.urlInput.View() + "\n\n" +
		tagsLabel + "\n" +
		m.tagsInput.View() + "\n\n" +
		labelStyle.Render("Tab: switch fields  Enter: save  Esc: cancel")
}

// Title returns the trimmed title input value.
func (m EditModel) Title() string {
	return strings.TrimSpace(m.titleInput.Value())
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
