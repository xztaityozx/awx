package cmd

import (
	"github.com/mitchellh/go-homedir"
	"strings"
	"testing"
)

func TestAllJoin(t *testing.T) {
	home, _ := homedir.Dir()
	src := PathJoin(home, "Src")
	files := []string{
		PathJoin(src, "A.csv"),
		PathJoin(src, "B.csv"),
		PathJoin(src, "C.csv"),
	}
	t.Run("000_Prepare", func(t *testing.T) {
		TryMkdirAll(src)

		WriteAppend(files[0], "1.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[0], "2.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[0], "3.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[0], "4.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[0], "5.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[1], "1.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[1], "2.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[1], "3.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[1], "4.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[1], "5.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[2], "1.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[2], "2.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[2], "3.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[2], "4.0000, 2.0000, 3.0000, 4.0000\n")
		WriteAppend(files[2], "5.0000, 2.0000, 3.0000, 4.0000\n")
	})

	t.Run("001_Join", func(t *testing.T) {
		actual := JoinTask{
			Files: files,
		}.Join()
		expect := PathJoin(src, "A_B_C_.csv")

		Equal(actual, expect, t)
		Equal(strings.Split(Cat(actual), "\n")[0], "0,1.0000, 2.0000, 3.0000, 4.0000,1.0000, 2.0000, 3.0000, 4.0000,1.0000, 2.0000, 3.0000, 4.0000", t)
	})

}
