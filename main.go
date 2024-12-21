package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
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
type doneMsg struct{}

func (m model) Init() tea.Cmd {
	// TODO: make sure we are in the right directory of the lz project before doing anything
	// TODO: make sure we are logged in to aws account
	// TODO: render ascii logo - maybe this should be in view() ?
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

		// fmt.Printf("%s %s %s\n", tf, tf_cmd, tf_backend)

		cmd := exec.Command(tf, tf_cmd, tf_backend)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout pipe: %w", err)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("failed to get stderr pipe: %w", err)
		}

		if err := cmd.Start(); err != nil {
			return errMsg(err)
		}

		// stream and process stdout
		go func() {
			var outputLines []string
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				raw := scanner.Text()
				cleaned := cleanOutput(raw)

				// Skip empty lines
				if cleaned != "" {
					outputLines = append(outputLines, cleaned)
				}
			}

			finalOutput := strings.Join(outputLines, "\n")
			fmt.Println(finalOutput)
		}()

		// stream and process stderr
		go func() {
			var outputLines []string
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				raw := scanner.Text()
				cleaned := cleanOutput(raw)

				// Skip empty lines
				if cleaned != "" {
					outputLines = append(outputLines, cleaned)
				}
			}

			finalOutput := strings.Join(outputLines, "\n")
			fmt.Println(finalOutput)
		}()

		// cmd.Stdout = &stdout
		// cmd.Stderr = &stderr

		// TODO: add disclaimer for which environemnt - red if mgmt or prod
		fmt.Printf("%s %s %s ...", tf, tf_cmd, tf_backend)
		if err := cmd.Wait(); err != nil {
			return errMsg(err)
		}

		return doneMsg{}
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

		// TODO: check if it may not be better to refactor this to use tea.Batch for better chaining of commands
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

		// TODO: add some tests to make sure that the init has completed successfully before moving forward

		// plan or apply
		cmd := exec.Command(tf, tf_cmd, tf_vars)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout pipe: %w", err)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("failed to get stderr pipe: %w", err)
		}

		// cmd.Stdout = &stdout
		// cmd.Stderr = &stderr

		// TODO: add disclaimer for which environemnt - red if mgmt or prod
		// fmt.Printf("%s %s %s ...", tf, tf_cmd, tf_vars)
		if err := cmd.Start(); err != nil {
			return errMsg(err)
		}

		// stream and process stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				raw := scanner.Text()
				cleaned := cleanOutput(raw)

				// Skip completely empty lines
				if cleaned != "" {
					fmt.Println(cleaned)
				}
			}
		}()

		// stream and process stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				raw := scanner.Text()
				cleaned := cleanOutput(raw)

				// Skip completely empty lines
				if cleaned != "" {
					fmt.Println(cleaned)
				}
			}
		}()

		if err := cmd.Wait(); err != nil {
			return errMsg(err)
		}

		return doneMsg{}
	}
}

// cleanOutput processes raw terminal output to remove unwanted artifacts.
func cleanOutput(input string) string {
	// Remove ANSI escape codes
	ansiEscapePattern := `\x1b\[[0-9;]*[a-zA-Z]`
	re := regexp.MustCompile(ansiEscapePattern)
	cleaned := re.ReplaceAllString(input, "")

	// Normalize multiple spaces/tabs into a single space
	cleaned = regexp.MustCompile(`[ \t]+`).ReplaceAllString(cleaned, " ")

	// Remove leading and trailing spaces
	cleaned = strings.TrimSpace(cleaned)

	// Remove empty lines
	if cleaned == "" {
		return ""
	}

	return cleaned
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
