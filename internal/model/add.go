package model

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	focusURL = iota
	focusTags
)

// AddModel handles the add-bookmark modal.
type AddModel struct {
	urlInput  textinput.Model
	tagsInput textinput.Model
	focus     int
	status    string
}

func newAddModel() AddModel {
	url := textinput.New()
	url.Placeholder = "https://example.com"
	url.CharLimit = 2048
	url.Width = 58

	tags := textinput.New()
	tags.Placeholder = "work, go, tools  (comma-separated, max 3)"
	tags.CharLimit = 200
	tags.Width = 58

	return AddModel{urlInput: url, tagsInput: tags, focus: focusURL}
}

// Init focuses the URL input.
func (m AddModel) Init() tea.Cmd {
	m.urlInput.Focus()
	return textinput.Blink
}

// Update handles keypresses delegated from the parent.
func (m AddModel) Update(msg tea.Msg) (AddModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.focus == focusURL {
				m.focus = focusTags
				m.urlInput.Blur()
				m.tagsInput.Focus()
				return m, textinput.Blink
			}
			m.focus = focusURL
			m.tagsInput.Blur()
			m.urlInput.Focus()
			return m, textinput.Blink
		case "shift+tab":
			if m.focus == focusTags {
				m.focus = focusURL
				m.tagsInput.Blur()
				m.urlInput.Focus()
				return m, textinput.Blink
			}
			m.focus = focusTags
			m.urlInput.Blur()
			m.tagsInput.Focus()
			return m, textinput.Blink
		}
	}

	var cmd tea.Cmd
	if m.focus == focusURL {
		m.urlInput, cmd = m.urlInput.Update(msg)
	} else {
		m.tagsInput, cmd = m.tagsInput.Update(msg)
	}
	return m, cmd
}

// View renders both input fields and the status line.
func (m AddModel) View() string {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))

	out := labelStyle.Render("URL") + "\n" +
		m.urlInput.View() + "\n\n" +
		labelStyle.Render("Tags") + "\n" +
		m.tagsInput.View()

	if m.status != "" {
		if m.status == "Already bookmarked" || m.status == "Please enter a URL" {
			out += "\n" + statusStyle.Render(m.status)
		} else {
			out += "\n" + okStyle.Render(m.status)
		}
	}
	return out
}

// URL returns the current URL input value.
func (m AddModel) URL() string {
	return m.urlInput.Value()
}

// Tags returns the current tags input value.
func (m AddModel) Tags() string {
	return m.tagsInput.Value()
}

// setURL pre-fills the URL input with a URL and focuses it.
func (m *AddModel) setURL(url string) {
	m.urlInput.SetValue(url)
	m.urlInput.Focus()
	m.urlInput.CursorEnd()
	m.focus = focusURL
}

// setStatus sets the status message shown below the inputs.
func (m *AddModel) setStatus(msg string) {
	m.status = msg
}
