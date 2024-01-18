package handle

func decorateCmd(cmd *Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
