// Copyright © 2018 xztaityozx
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
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

// multiextractCmd represents the multiextract command
var multiextractCmd = &cobra.Command{
	Use:     "multiextract",
	Aliases: []string{"mex"},
	Short:   "",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		gc, _ := cmd.Flags().GetBool("GC")
		p, _ := cmd.Flags().GetInt("parallel")
		sigName, _ := cmd.Flags().GetStringSlice("signalName")
		start, _ := cmd.Flags().GetFloat64("start")
		step, _ := cmd.Flags().GetFloat64("step")
		stop, _ := cmd.Flags().GetFloat64("stop")

		wd, _ := os.Getwd()
		mt := NewMultiTask(sigName, start, step, stop, wd, p, gc)
		res := mt.Run()
		Print(res.ToString())
	},
}

type MultiTask struct {
	Signals  []string
	Start    float64
	Step     float64
	Stop     float64
	BaseDir  string
	Parallel int
	GC       bool
}

// Compare MultiTask
func (s MultiTask) CompareTo(t MultiTask) bool {
	if len(s.Signals) != len(t.Signals) {
		return false
	}

	for key, value := range s.Signals {
		if value != t.Signals[key] {
			return false
		}
	}

	return s.Parallel == t.Parallel &&
		s.BaseDir == t.BaseDir &&
		s.Start == t.Start &&
		s.Step == t.Step &&
		s.Stop == t.Stop &&
		s.GC == t.GC
}

func (this MultiTask) GenerateTaskSlice() []Task {
	var rt []Task
	for _, value := range GetTargetDirectories(this.BaseDir) {
		dst := PathJoin(this.BaseDir, "../Result", value)
		src := PathJoin(this.BaseDir, value)
		r := NewRange(this.Start, this.Step, this.Stop)
		rt = append(rt, NewTask(dst, src, r, this.Signals))
	}
	return rt
}


func (this MultiTask) Run() Summary {
	sum := Summary{
		Files:  []string{},
		Status: false,
	}

	
	worker := func(tasks []Task) <- chan Summary {
		rt := make(chan Summary, len(tasks))
		for _, value := range tasks {
			go func(t Task){
				res,err := t.Run()
				if err != nil {
					Fatal(err)
				}
				rt<-res
			}(value)
		}
		return rt
	}
	tasks := this.GenerateTaskSlice()
	receiver := worker(tasks)

	for i := 0; i< len(tasks); i++ {
		res := <-receiver
		sum.Status = sum.Status && res.Status
		sum.Files = append(sum.Files, res.Files...)
	}
		
	
	return sum
}

// 処理対象のディレクトリをリストアップする
func GetTargetDirectories(p string) []string {
	list, err := ioutil.ReadDir(p)
	if err != nil {
		Fatal(err)
	}

	var rt []string
	for _, value := range list {
		if value.IsDir() {
			rt = append(rt, value.Name())
		}
	}

	return rt
}

// Constructor for MultiTask
func NewMultiTask(sig []string, start float64, step float64, stop float64, base string, para int, gc bool) MultiTask {
	return MultiTask{
		Signals:  sig,
		Start:    start,
		Step:     step,
		Stop:     stop,
		BaseDir:  base,
		Parallel: para,
		GC:       gc,
	}
}

func init() {
	rootCmd.AddCommand(multiextractCmd)

	multiextractCmd.Flags().Bool("GC", false, "作業後波形データを削除します")
	multiextractCmd.Flags().IntP("parallel", "P", 1, "並列して実行される数を指定します")
	multiextractCmd.Flags().StringSliceP("signalName", "S", []string{"N1", "N2", "BLB", "BL"}, "取り出したい波形のリストです")
	multiextractCmd.Flags().Float64("start", 0, "プロットの最初の時間[ns]です")
	multiextractCmd.Flags().Float64("step", 0, "プロットの刻み幅[ns]です")
	multiextractCmd.Flags().Float64("stop", 0, "プロットの最後の時間[ns]です")
}
