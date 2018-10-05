package cmd

import (
	"log"
	"os"
	"path/filepath"
)

func Fatal(v ...interface{}) {
	log.Fatal(v)
}

func Print(v ...interface{}) {
	log.Print(v)
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
