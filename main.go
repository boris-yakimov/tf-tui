package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: split the init, update and veiw parts into a separate package -- only do this once we have things working already !

type model struct {
	output string // Stores the command output
	err    error  // Stores any error that occurs
}

type commandOutputMsg struct {
	output string
	err    error
}

func tfInit(backendPath string) tea.Msg {
	return func() tea.Msg {
		// TODO: figure out in which dir to run this
		tf := "terraform"
		arg1 := "init"
		arg2 := fmt.Sprintf("--backend-config=%s", backendPath) // backends/dev.tfbackend

		// Execute the command
		out, err := exec.Command(tf, arg1, arg2).CombinedOutput()
		return commandOutputMsg{output: string(out), err: err}
	}
}

func (m model) Init() tea.Cmd {
	// TODO: make sure we are in the right directory of the lz project
	// TODO: render ascii logo
	// TODO: provide options for plan
	return func() tea.Msg {
		return tfInit("backends/dev.tfbackend")
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case commandOutputMsg:
		if msg.err != nil {
			m.err = msg.err
			m.output = "Error: " + msg.err.Error()
		} else {
			m.output = msg.output
		}
		return m, tea.Quit // exit after receiving the output

	case tea.KeyMsg:
		// TODO: does q work here ?
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	// return updated model to the Bubble Tea runtime for processing
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("An error occurred:\n%s\n", m.err)
	}
	return fmt.Sprintf("Tf cmd output:\n%s\n", m.output)
}

func main() {
	// TODO: validate that debug log works
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("lz-tui-debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
