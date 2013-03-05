package process

import (
	"os/exec"
)

func BackgroundRun(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	err := cmd.Start()
	go cmd.Wait()
	return err
}
