package job

import (
	"os/exec"
	"strconv"
	"syscall"
)

func makeCmdAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func killProcess(pid int) error {
	return exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid)).Run()
}

