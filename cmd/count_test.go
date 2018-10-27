package cmd

import (
	"github.com/mitchellh/go-homedir"
	"testing"
)

func TestAllCount(t *testing.T) {
	home, _ := homedir.Dir()
	src := PathJoin(home, "Src")
	TryMkdirAll(src)
	rule := PathJoin(src, "rule")
	Write(rule, "This is dummy rule")
	files := []string{
		PathJoin(src, "A.csv"),
		PathJoin(src, "B.csv"),
		PathJoin(src, "C.csv"),
	}
	var target string

	t.Run("000_Prepare", func(t *testing.T) {
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

		target = JoinTask{
			Files: files,
		}.Join()
	})

	t.Run("001_NewCountTask", func(t *testing.T) {
		IsRangeSEEDCountUp = true
		IgnoreSigma = false
		RuleFile = ""
		args := []string{"This is dummy rule", target}

		actual := NewCountTask(args)
		expect := CountTask{
			IgnoreSigma:        false,
			IsRangeSEEDCountUp: true,
			Rule:               "This is dummy rule",
			Target:             target,
		}

		Equal(actual, expect, t)
	})
	t.Run("002_NewCountTask", func(t *testing.T) {
		IsRangeSEEDCountUp = false
		IgnoreSigma = true
		RuleFile = rule
		args := []string{target}

		actual := NewCountTask(args)
		expect := CountTask{
			IgnoreSigma:        true,
			IsRangeSEEDCountUp: false,
			Rule:               "This is dummy rule",
			Target:             target,
		}

		Equal(actual, expect, t)
	})

	t.Run("003_GetRuleScript", func(t *testing.T) {
		actual := CountTask{
			Rule: "$1>=1 && $2>=2",
		}.GetRuleScript()
		expect := "BEGIN{s=0}$1>=1 && $2>=2{s++}END{print s}"

		Equal(actual, expect, t)
	})
}
