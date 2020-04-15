package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GenRandomHash() string {
	rdata := make([]byte, 64)
	rand.Read(rdata)
	hash := sha256.Sum256(rdata)
	token := fmt.Sprintf("%x", hash)
	return token
}

func GenRandomString(n int) string {
	const lb = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		t := make([]byte, 1)
		rand.Read(t)
		b[i] = lb[int(t[0])%len(lb)]
	}
	return string(b)
}

func GenRandomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

func ReadFile(path string) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GetExecDir() string {
	exe_path, _ := os.Executable()
	return filepath.Dir(exe_path)
}

func ExecPath(name string) string {
	return filepath.Join(GetExecDir(), name)
}
