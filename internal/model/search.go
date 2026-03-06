package model

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ndy40/cairn/internal/store"
)

// SearchModel manages the search text input.
// Filtering logic lives in app.go to enable the two-stage FTS5+fuzzy pipeline.
type SearchModel struct {
	input  textinput.Model
	width  int
	height int
}

func newSearchModel(_ []*store.Bookmark) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search bookmarks…"
	ti.CharLimit = 256
	ti.Focus()
	return SearchModel{input: ti}
}

// Init starts the blinking cursor.
func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

// setSize stores the available terminal dimensions.
func (m *SearchModel) setSize(w, h int) {
	m.width = w
	m.height = h
	m.input.Width = w - 4
}

// Update handles keypresses in search mode.
func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+a" {
			m.input.SetValue("")
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// term returns the current search term.
func (m SearchModel) term() string {
	return m.input.Value()
}

// inputView renders the text input for embedding in the parent view.
func (m SearchModel) inputView() string {
	return m.input.View()
}

// noResultsView renders a centred "no results" message.
func (m SearchModel) noResultsView() string {
	msg := "No results for «" + m.input.Value() + "»"
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(msg),
	)
}
