package textarea

// A simple program demonstrating the textarea component from the Bubbles
// component library.

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// func main() {
// 	p := tea.NewProgram(initialModel())

// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func New() tea.Model {
// 	return initialModel()
// }

type errMsg error

type Model struct {
	textarea textarea.Model
	Value    string
	err      error
}

func New() Model {
	ti := textarea.New()
	ti.Placeholder = `name:`
	ti.Focus()

	return Model{
		textarea: ti,
		err:      nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:

			return m, tea.Quit
		case tea.KeyCtrlS:
			value := m.textarea.Value()
			m.Value = value
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"Paste your Podinate profile file.\n\n%s\n\n%s",
		m.textarea.View(),
		"(ctrl+c to quit, ctrl+s to save profile)",
	) + "\n\n"
}
