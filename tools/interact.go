package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"

	"github.com/fatih/color"
)

// ReadFromStdin returns the user input. You can forbid the empty value by
// setting the parameter to false.
func ReadFromStdin(allowEmptyValue bool) string {
	s := bufio.NewScanner(os.Stdin)
	var t string
	for {
		if s.Scan() {
			t = s.Text()
			if (t != "") || (t == "" && allowEmptyValue) {
				return t
			}
			Warning("Empty values are not allowed. Please re-type the value.")
			QuestionWithPrompt("")
		}
	}
}

// ReadConstraintedValueFromStdin checks that stdin receives one of the allowed
// value passed as parameter.
// The function returns the value supplied in the stdin.
func ReadConstraintedValueFromStdin(allowed ...string) string {
	values := make(map[string]int)
	for _, s := range allowed {
		values[s] = 1
	}

	var t string
	s := bufio.NewScanner(os.Stdin)
	for {
		if s.Scan() {
			t = s.Text()
			if values[t] == 1 {
				return t
			}
			Warning(fmt.Sprintf("The value you typed %q is not an allowed value. Allowed values are:", t))
			fmt.Println()
			for _, v := range allowed {
				fmt.Printf("  * %s\n", v)
			}
			QuestionWithPrompt("")
		}
	}
}

// ReadYesOrNoFromStdin checks that the user enter any of the Y,y,N,n value.
// Any other value is rejected with a message asking for entering new values.
// If the value filled is Y or y, the function returns true.
// If the value filled is N or n, the function returns false.
func ReadYesOrNoFromStdin() bool {
	s := bufio.NewScanner(os.Stdin)
	var t string
	for {
		if s.Scan() {
			t = s.Text()
			switch t {
			case "Y", "y":
				return true
			case "N", "n":
				return false
			default:
				Info("We only support Y,y,N,n as input. Please re-type your choice")
			}
			QuestionWithPrompt("")
		}
	}
}

// ExecCmdCaptureStreamsPublishOutput runs a command and returns the stdout, stderr
// and the error to the user.
// Note that the stdout and stderr are dumped so the user sees the progression
// of the command. Returning the stdout and stderr makes possible to log what
// happened with the exact output. So the user can have a complete log of what
// happened.
func ExecCmdCaptureStreamsPublishOutput(cmd exec.Cmd) ([]byte, []byte, int, error) {
	copyStderr := bytes.Buffer{}
	copyStdout := bytes.Buffer{}
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	errorCode := 0
	err := cmd.Start()
	if err != nil {
		return nil, nil, errorCode, errors.Wrap(err, "failed to start process")
	}

	multiWriterOut := io.MultiWriter(os.Stdout, &copyStdout)
	multiWriterErr := io.MultiWriter(os.Stderr, &copyStderr)

	go func() {
		io.Copy(multiWriterOut, stdoutIn)
	}()

	go func() {
		io.Copy(multiWriterErr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// Code retrieved here: http://www.nljb.net/default/Golang-Get-exit-code/
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				errorCode = status.ExitStatus()
			}
		}
		return copyStdout.Bytes(), copyStderr.Bytes(), errorCode, err
	}

	if IsDebugActivated() {
		green := color.New(color.FgGreen)
		green.Printf("Done in %s\n", cmd.ProcessState.UserTime())
	}

	return copyStdout.Bytes(), copyStderr.Bytes(), errorCode, err
}

// ExecCmdCaptureStreamsNoOutput executes a command and captures the output.
// But, on the contrary of ExecuteCmdAndCaptureOutputs, the command shows
// nothing on the stdout. stdout and stderr are captures in the same flow
// and returned as the first returned parameter.
func ExecCmdCaptureStreamsNoOutput(cmd exec.Cmd) ([]byte, []byte, int, error) {
	copyStderr := bytes.Buffer{}
	copyStdout := bytes.Buffer{}
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	err := cmd.Start()
	errorCode := 0
	if err != nil {
		return nil, nil, errorCode, errors.Wrap(err, "failed to start process")
	}

	multiWriterOut := io.MultiWriter(&copyStdout)
	multiWriterErr := io.MultiWriter(&copyStderr)

	go func() {
		io.Copy(multiWriterOut, stdoutIn)
	}()

	go func() {
		io.Copy(multiWriterErr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// Code retrieved here: http://www.nljb.net/default/Golang-Get-exit-code/
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				errorCode = status.ExitStatus()
			}
		}
		return copyStdout.Bytes(), copyStderr.Bytes(), errorCode, err
	}
	if IsDebugActivated() {
		green := color.New(color.FgGreen)
		green.Printf("Done in %s\n", cmd.ProcessState.UserTime())
	}

	return copyStdout.Bytes(), copyStderr.Bytes(), errorCode, err
}
