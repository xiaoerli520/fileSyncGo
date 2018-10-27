package utils

import (
	"os"
	"io/ioutil"
	"sinago/Error"
	"os/exec"
	"path/filepath"
	"strings"
	"errors"
	"runtime"
	"bytes"
	"strconv"
)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

func CreateFile(path string) string {
	f, err := os.Create(path)
	Error.CheckErr(err)
	defer f.Close()
	return path
}

func LoadFile(path string) string {
	ft, err := ioutil.ReadFile(path)
	Error.CheckErr(err)
	fileString := string(ft)
	return fileString
}

func WriteFile(path string, content string, isTruncate bool) {
	var ft *os.File
	var err error
	if isTruncate {
		ft, err = os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0755)
	} else {
		ft, err = os.OpenFile(path, os.O_RDWR | os.O_CREATE , 0755)
	}

	Error.CheckErr(err)
	ft.WriteString(content)
	defer ft.Close()
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}