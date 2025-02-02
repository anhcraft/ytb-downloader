package request

import (
	"os/exec"
	"syscall"
)

func DecorateCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
