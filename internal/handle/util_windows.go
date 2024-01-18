package handle

import (
	"os/exec"
	"syscall"
)

func decorateCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
