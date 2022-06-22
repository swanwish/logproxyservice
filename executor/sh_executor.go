package executor

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/swanwish/go-common/logs"
)

type shExecutor struct {
}

func NewShExecutor() shExecutor {
	executor := shExecutor{}
	return executor
}

func (e *shExecutor) RunScript(script string) (string, error) {
	output := make([]string, 0)
	cmd := exec.Command("sh", "-c", script)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		logs.Errorf("Failed to get stdout pipe, the error is %#v", err)
		output = append(output, err.Error())
		return "", err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		logs.Errorf("Failed to get stderr pipe, the error is %#v", err)
		output = append(output, err.Error())
		return "", err
	}
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		stdoutReader := bufio.NewReader(stdoutPipe)
		for {
			line, _, err := stdoutReader.ReadLine()
			if err != nil {
				if err != io.EOF {
					logs.Errorf("Failed to read line from stdoutPipe, the error is %#v, message %s", err, err.Error())
				}
				break
			}
			output = append(output, string(line))
			// logs.Debugf("stdout: %s, %t", string(line), isPrefix)
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		stderrReader := bufio.NewReader(stderrPipe)
		for {
			line, isPrefix, err := stderrReader.ReadLine()
			if err != nil {
				if err != io.EOF {
					logs.Errorf("Failed to read line from stderrPipe, the error is %#v, error message: %s", err, err.Error())
				}
				break
			}
			output = append(output, string(line))
			logs.Debugf("stderr: %s, %t", string(line), isPrefix)
		}
		wg.Done()
	}()
	err = cmd.Start()
	if err != nil {
		logs.Errorf("Failed to start command, the error is %#v", err)
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		logs.Errorf("Failed to wait, the error is %#v", err)
		return "", err
	}
	wg.Wait()

	result := strings.Join(output, "\n")
	// logs.Debugf("The output is %s", result)
	return result, nil
}
