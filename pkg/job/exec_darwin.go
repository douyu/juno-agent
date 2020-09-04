package job

import "syscall"

func makeCmdAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func killProcess(pid int) error {
	return syscall.Kill(-pid, syscall.SIGKILL)
}
