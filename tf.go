package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error
type doneMsg struct{}

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
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// cmd.Stdout = &stdout
		// cmd.Stderr = &stderr

		// TODO: add disclaimer for which environemnt - red if mgmt or prod
		fmt.Printf("%s %s %s ...", tf, tf_cmd, tf_backend)
		err := cmd.Run()
		if err != nil {
			return errMsg(err)
		}

		return doneMsg{}
	}
}

func tfAction(tfAction string, varFilePath string) tea.Cmd {
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
		// var wg sync.WaitGroup
		// wg.Add(1)
		//
		// go func() {
		// 	defer wg.Done()
		// 	initCmd := tfInit(tfBackendPath)
		// 	if initCmd != nil {
		// 		initCmd()
		// 	}
		// }()
		//
		// wg.Wait() // wait for tf init to finish before running plan or apply

		// plan or apply
		cmd := exec.Command(tf, tf_cmd, tf_vars)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// cmd.Stdout = &stdout
		// cmd.Stderr = &stderr

		// TODO: add disclaimer for which environemnt - red if mgmt or prod
		fmt.Printf("%s %s %s ...", tf, tf_cmd, tf_vars)
		err := cmd.Run()
		if err != nil {
			return errMsg(err)
		}

		return doneMsg{}
	}
}
