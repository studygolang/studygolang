package process

import (
	"os/exec"
	"syscall"
)

func BackgroundRun(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	return cmd.Start()
}
