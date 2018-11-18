// Copyright © 2018 xztaityoz
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:        "count",
	Aliases:    []string{"cnt"},
	ArgAliases: []string{"rule", "target"},
	Short:      "",
	Long:       ``,
	Run: func(cmd *cobra.Command, args []string) {
		ct := NewCountTask(args)
		ct.Run()
	},
}

type CountTask struct {
	IsRangeSEEDCountUp bool
	Rule               string
	Target             string
	Parallel           int
}

func (ct CountTask) Run() {
	if ct.IsRangeSEEDCountUp {
		res := ct.RangeSEEDCountUp()
		for _, v := range res {
			fmt.Println(v)
		}
	} else {
		res := ct.CountUp()
		fmt.Println(res)
	}
}

// new count task
func NewCountTask(args []string) CountTask {
	var rt = CountTask{
		IsRangeSEEDCountUp: IsRangeSEEDCountUp,
	}
	// specified rule string from command line
	if RuleFile == "" && len(args) == 2 {
		rt.Rule = args[0]
		rt.Target = args[1]
	} else if RuleFile != "" && len(args) > 0 {
		rt.Rule = Cat(RuleFile)
		rt.Target = args[0]
	} else {
		Fatal("awx count: unexpected command line arguments")
	}
	return rt
}

func (ct CountTask) GetRuleScript() string {
	return fmt.Sprintf("BEGIN{s=0}%s{s++}END{print s}", ct.Rule)
}

func (ct CountTask) CountUp() int64 {
	command := exec.Command("awk", ct.GetRuleScript(), ct.Target)
	out, err := command.Output()
	if err != nil {
		Fatal("awx error: awk command failed\n\tScript: ", ct.GetRuleScript(), "\n", err)
	}

	str := strings.Trim(string(out), "\n")
	rt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		Fatal(err)
	}

	return rt
}

func (ct CountTask) RangeSEEDCountUp() []int64 {
	wd, err := os.Getwd()
	if err != nil {
		Fatal(err)
	}

	if filepath.Base(wd) != "Result" {
		Fatal("awx count: Invalid Directory.\n\t", filepath.Base(wd))
	}

	var dirs []string
	for _, v := range Ls(wd) {
		if v.IsDir() && len(v.Name()) == 7 && v.Name()[0:4] == "SEED" {
			dirs = append(dirs, PathJoin(wd, v.Name()))
		}
	}

	result := make([]int64, len(dirs))
	receiver := ct.CountUpWorker(dirs)
	for i := 0; i < len(dirs); i++ {
		pair := <-receiver
		result[pair.Key] = pair.Value
	}

	return result
}

type Pair struct {
	Key   int
	Value int64
}

func (ct CountTask) CountUpWorker(dirs []string) <-chan Pair {
	receiver := make(chan Pair, ct.Parallel)
	for i, v := range dirs {
		src := PathJoin(v, ct.Target)
		go func(i int, p string) {
			if !IsExistsFile(p) {
				receiver <- Pair{
					Key:   i,
					Value: -1,
				}
				return
			}
			command := exec.Command("awk", ct.GetRuleScript(), p)
			out, err := command.Output()
			if err != nil {
				Fatal("awx error: awk command failed.\n\t", ct.GetRuleScript(), "\n", err)
			}
			str := strings.Trim(string(out), "\n")
			res, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				Fatal(err)
			}

			receiver <- Pair{
				Key:   i,
				Value: res,
			}
		}(i, src)
	}

	return receiver
}

var IsRangeSEEDCountUp bool

var RuleFile string

func init() {
	rootCmd.AddCommand(countCmd)
	countCmd.Flags().BoolVarP(&IsRangeSEEDCountUp, "RangeSEED", "R", false, "RangeSEEDシミュレーションの結果を数え上げます")
	countCmd.Flags().StringVarP(&RuleFile, "ruleFile", "t", "", "ルールを記述したファイルを指定できます")
	countCmd.Flags().Int32P("Parallel", "P", 1, "並列実行する個数です")
}
