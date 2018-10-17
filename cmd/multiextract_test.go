package cmd

import (
	"github.com/mitchellh/go-homedir"
	"testing"
)

func TestAllMultiExtract(t *testing.T) {
	home, _ := homedir.Dir()
	dst := PathJoin(home, "Dst")
	src := PathJoin(home, "Src")
	t.Run("001_NewMultiTask", func(t *testing.T) {
		actual := NewMultiTask([]string{"A", "B"}, 2.5, 7.5, 17.5, "Dst", 10, false)
		expect := MultiTask{
			Signals:  []string{"A", "B"},
			Start:    2.5,
			Step:     7.5,
			Stop:     17.5,
			BaseDir:  "Dst",
			Parallel: 10,
			GC:       false,
		}

		Assert(actual.CompareTo(expect), t)
	})

	t.Run("002_GetTargetDirectories", func(t *testing.T) {
		TryMkdirAll(src)
		TryMkdirAll(PathJoin(src, "SEED001"))
		TryMkdirAll(PathJoin(src, "SEED002"))
		CreateFile(PathJoin(src, "Dummy.csv"))
		actual := GetTargetDirectories(src)
		expect := []string{
			"SEED001",
			"SEED002",
		}

		for key, value := range expect {
			Equal(actual[key], value, t)
		}

	})

	mt := NewMultiTask([]string{"A", "B"}, 2.5, 7.5, 17.5, src, 10, false)

	t.Run("003_GenerateTaskSlice", func(t *testing.T) {
		actual := mt.GenerateTaskSlice()
		expect := []Task{
			{
				DstDir:  PathJoin(home, "Result/SEED001"),
				SrcDir:  PathJoin(src, "SEED001"),
				Range:   NewRange(2.5, 7.5, 17.5),
				Signals: mt.Signals,
			},
			{
				DstDir:  PathJoin(home, "Result/SEED002"),
				SrcDir:  PathJoin(src, "SEED002"),
				Range:   NewRange(2.5, 7.5, 17.5),
				Signals: mt.Signals,
			},
		}
		for key, value := range actual {
			Assert(value.CompareTo(expect[key]), t)
		}
	})

	t.Run("XXX_Remove", func(t *testing.T) {
		RemoveFile(dst)
		RemoveFile(src)
	})

}
