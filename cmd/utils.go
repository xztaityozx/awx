package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func Print(v ...interface{}) {
	log.Print(v...)
}

func IsExistsFile(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func Cat(p string) string {
	if b, err := ioutil.ReadFile(p); err != nil {
		Fatal(err)
	} else {
		return string(b)
	}
	return ""
}

func Write(p string, data string) {
	if err := ioutil.WriteFile(p, []byte(data), 0644); err != nil {
		Fatal(err)
	}
}

func RemoveFile(p string) {
	if err := os.RemoveAll(p); err != nil {
		Fatal(err)
	}
}

func TryMkdirAll(path string) {
	if _, ste := os.Stat(path); ste != nil {
		if err := os.MkdirAll(path, 0755); err != nil {
			Fatal(err)
		}
	}
}

func PathJoin(p ...string) string {
	return filepath.Join(p...)
}
