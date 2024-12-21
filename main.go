package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: split the init, update and veiw parts into a separate package -- only do this once we have things working already !
// TODO: how do we test that we are already logged in to the right account - maybe a simple s3 ls

// TUI model
type model struct {
	choices []string // stores choices that will appear in the menu
	cursor  int      // position where we are in the menu

	output string // stores command outputs
	err    error  // stores any error that occurs
}

type errMsg error

// type doneMsg struct{}
type outputMsg string

func (m model) Init() tea.Cmd {
	// TODO: make sure we are in the right directory of the lz project before doing anything
	// TODO: make sure we are logged in to aws account
	// TODO: render ascii logo - maybe this should be in view() ?
	// tfBackendPath := "backends/dev.tfbackend"
	// return tfInit(tfBackendPath)
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

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

func (m model) View() string {
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

func tfInit(backendPath string) tea.Cmd {
	// TODO: error if no file exists on path
	return func() tea.Msg {
		// var stdout bytes.Buffer
		// var stderr bytes.Buffer

		tf := "terraform"
		tf_cmd := "init"
		tf_backend := fmt.Sprintf("--backend-config=%s", backendPath)

		fmt.Printf("%s %s %s\n", tf, tf_cmd, tf_backend)

		cmd := exec.Command(tf, tf_cmd, tf_backend)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		// TODO: add disclaimer for which environemnt - red if mgmt or prod
		fmt.Printf("%s %s %s ...", tf, tf_cmd, tf_backend)
		if err := cmd.Run(); err != nil {
			return errMsg(err)
		}

		// return doneMsg{}
		return outputMsg(out.String())
	}
}

func tfAction(tfAction string, varFilePath string, tfBackendPath string) tea.Cmd {
	if tfAction != "plan" && tfAction != "apply" {
		// TODO: handle this as an error
		fmt.Println("Invalid Terraform action, should be plan or apply")
	}

	return func() tea.Msg {
		// var stdout bytes.Buffer
		// var stderr bytes.Buffer

		tf := "terraform"
		tf_cmd := tfAction
		tf_vars := fmt.Sprintf("-var-file=%s", varFilePath)

		// TODO: why is this not getting executed
		// init
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			initCmd := tfInit(tfBackendPath)
			if initCmd != nil {
				initCmd()
			}
		}()

		wg.Wait() // wait for tf init to finish before running plan or apply

		// plan or apply
		cmd := exec.Command(tf, tf_cmd, tf_vars)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		// TODO: add disclaimer for which environemnt - red if mgmt or prod
		fmt.Printf("%s %s %s ...", tf, tf_cmd, tf_vars)
		err := cmd.Run()
		if err != nil {
			return errMsg(err)
		}

		// return doneMsg{}
		return outputMsg(out.String())
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

	environments := []string{
		"dev",
		"stage",
		"prod",
	}

	p := tea.NewProgram(model{
		cursor:  0,
		choices: environments,
	})

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
