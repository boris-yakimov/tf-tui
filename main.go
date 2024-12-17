package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: split the init, update and veiw parts into a separate package -- only do this once we have things working already !
// TODO: how do we test that we are already logged in to the right account - maybe a simple s3 ls

// TUI model
type model struct {
	choices []string // stores choices that will appear in the menu of the TUI
	cursor  int      // position where we are in the mnu

	output string // stores the command output
	err    error  // stores any error that occurs
}

func (m model) Init() tea.Cmd {
	// TODO: make sure we are in the right directory of the lz project before doing anything
	// TODO: render ascii logo - maybe this should be in view() ?
	// tfBackendPath := "backends/dev.tfbackend"
	// return tfInit(tfBackendPath)
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// TODO: add an option to trigger a plan or apply

	case commandOutputMsg:
		m.output = msg.output
		m.err = msg.err
		return m, nil

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
			tfInit(tfBackendPath)
			// TODO: show full output in trail mode as it appears on the screen
			tfAction("plan", tfVarsPath)

			// TODO:: how do we do this ?
			// command := m.commands[m.cursor]
			// return m, runShellCommand(command)
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

	// if m.output == "" {
	// 	return "No output received from the command.\n\nPress 'q' to quit."
	// }

	return fmt.Sprintf("waiting for command output \n\n%s\n\nPress 'q' to quit.", m.output)
}

// tf and AWS commands
type commandOutputMsg struct {
	output string
	err    error
}

func tfInit(backendPath string) tea.Cmd {
	// TODO: error if no file exists on path
	return func() tea.Msg {
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		tf := "terraform"
		tf_cmd := "init"
		tf_backend := fmt.Sprintf("--backend-config=%s", backendPath)

		fmt.Printf("%s %s %s\n", tf, tf_cmd, tf_backend)

		cmd := exec.Command(tf, tf_cmd, tf_backend)
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

func tfAction(tfAction string, varFilePath string) tea.Cmd {
	if tfAction != "plan" && tfAction != "apply" {
		// TODO: handle this as an error
		fmt.Println("Invalid Terraform action, should be plan or apply")
	}

	return func() tea.Msg {
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		tf := "terraform"
		tf_cmd := tfAction
		tf_vars := fmt.Sprintf("-var-file=%s", varFilePath)

		cmd := exec.Command(tf, tf_cmd, tf_vars)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// TODO: add disclaimer for which environemnt - red if mgmt or prod
		fmt.Printf("Running Terraform tfAction ... \n\n")
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
