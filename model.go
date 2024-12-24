package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// TUI model
type Model struct {
	choices []string // stores choices that will appear in the menu
	cursor  int      // position where we are in the menu

	// terminal dimentions
	width  int
	height int

	// output string // stores command outputs
	// err    error  // stores any error that occurs
}

// a pointer here because we want to keep seting the right dimentions of the terminal in the model itself every time that m.Update() is called
func (m *Model) updateTerminalDimentions(width, height int) {
	m.width = width
	m.height = height
}

func (m Model) Init() tea.Cmd {
	// TODO: make sure we are in the right directory of the lz project before doing anything
	// TODO: make sure we are logged in to aws account
	// TODO: render ascii logo - maybe this should be in view() ?
	// tfBackendPath := "backends/dev.tfbackend"
	// return tfInit(tfBackendPath)
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// keep track of terminal dimensions
	case tea.WindowSizeMsg:
		m.updateTerminalDimentions(msg.Width, msg.Height)

	case tea.KeyMsg:
		// keep track of which key was pressed
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			tfBackendPath := "backends/dev.tfbackend"
			tfVarsPath := "vars/dev.tfvars"
			// TODO: show full init output
			// TODO: wait for init to finish
			// TODO: only than proceed to plan
			// TODO: show full output in trail mode as it appears on the screen

			return m, (tfAction("plan", tfVarsPath, tfBackendPath))
		}
	}

	// return updated model to the Bubble Tea runtime for processing
	return m, nil
}

func (m Model) View() string {
	s := "Choose an option:\n\n"
	for i, choice := range m.choices {
		cursor := " " // cursor indicator
		if i == m.cursor {
			cursor = ">" // highlight current choice
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\nPress q to quit.\n"
	return s
}
