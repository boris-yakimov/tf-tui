package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: split the init, update and veiw parts into a separate package -- only do this once we have things working already !

// TUI model
type model struct {
	output string // Stores the command output
	err    error  // Stores any error that occurs
}

func (m model) Init() tea.Cmd {
	// TODO: make sure we are in the right directory of the lz project
	// TODO: render ascii logo
	// TODO: provide options for plan
	return tfInit()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case commandOutputMsg:
		m.output = msg.output
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		// keep track of which key was pressed
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	// return updated model to the Bubble Tea runtime for processing
	return m, nil
}

func (m model) View() string {
	// TODO: have to make the view autorefresh to see the full tf output scrolling
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit.", m.err)
	}

	if m.output == "" {
		return "No output received from the command.\n\nPress 'q' to quit."
	}

	return fmt.Sprintf("Output: \n\n%s\n\nPress 'q' to quit.", m.output)
}

// tf and AWS commands
type commandOutputMsg struct {
	output string
	err    error
}

func tfInit() tea.Cmd {
	return func() tea.Msg {
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		tf := "terraform"
		arg1 := "init"

		// TODO: add additional backend path once tested
		// arg2 := fmt.Sprintf("--backend-config=%s", backendPath) // backends/dev.tfbackend

		cmd := exec.Command(tf, arg1) // run in separate shell
		// cmd := exec.Command("ls", "-lah") // run in separate shell
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// TODO: add disclaimer for which environemnt - red if mgmt or prod
		fmt.Printf("Running Terraform Init... \n\n")
		err := cmd.Run()

		return commandOutputMsg{
			output: stdout.String() + stderr.String(),
			err:    err,
		}
	}
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
		fmt.Printf("Error: %v\n", err)
		// os.Exit(1)
	}
}
