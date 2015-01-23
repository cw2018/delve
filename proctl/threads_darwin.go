package proctl

import "syscall"

func (t *ThreadContext) Kill(sig syscall.Signal) error {
	return syscall.Kill(t.Id, sig)
}
