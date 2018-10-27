package cmd

import (
	"os/exec"
	"bytes"
)

// exec shell command and get return
func ExecShell(command string, arg ...string) (string, string) {
	cmd := exec.Command(command, arg...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Run()

	return out.String(), stderr.String()
}

// exec shell command at Designated directory and get return
func ExecShellAt(dir string, command string, arg []string) (string, string) {
	cmd := exec.Command(command, arg...)
	cmd.Dir = dir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Run()

	return out.String(), stderr.String()
}
