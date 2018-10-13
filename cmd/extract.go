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
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//signals, _ := cmd.Flags().GetStringSlice("signalName")
		//dst, _ := cmd.Flags().GetString("dst")
		//src, _ := cmd.Flags().GetString("src")
		//start, _ := cmd.Flags().GetFloat64("start")
		//step, _ := cmd.Flags().GetFloat64("step")
		//stop, _ := cmd.Flags().GetFloat64("stop")
	},
}

func Run(task Task) error {
	if len(task.Signals) == 0 {
		return errors.New("awx task error: signals length is 0")
	}

	if !task.IsValidDirectory() {
		return errors.New("awx task error: Invalid directory")
	}

	response := task.Extract()
	if len(response) != len(task.Signals) {
		return errors.New("awx task error: Unexpected result task.Extract()")
	}

	return nil
}

func (this Task) Extract() []string {
	command := exec.Command("bash", "-c", fmt.Sprintf("cd %s && wv -k -ace_no_gui %s &> ./log.wv", this.SrcDir, this.AcePath))
	if err := command.Run(); err != nil {
		Fatal(err)
	}

	var rt []string
	for _, v := range this.Signals {
		src := PathJoin(this.SrcDir, v+".csv")
		dst := PathJoin(this.DstDir, v+".csv")

		// 整形
		ShapingCSV(src, dst, this.Range.Count())

		rt = append(rt, dst)
	}
	return rt
}

// CSVから必要なところだけ取り出す
func ShapingCSV(p string, dst string, step int) {
	box := strings.Split(Cat(p), "\n")
	var store []string
	for _, v := range box {
		if len(v) == 0 || v[0] == '#' || v[0] == 'T' {
			continue
		}

		t := strings.Split(strings.Replace(v, " ", "", -1), ",")
		if len(t) < 2 {
			Fatal("awx error: Invalid csv file", p)
		}

		store = append(store, t[1])
	}
	if len(store)%step != 0 {
		Fatal("awx error: Invalid data length", p)
	}

	CreateFile(dst)
	for i := 0; i < len(store); i += step {
		for j := 0; j < step-1; j++ {
			WriteAppend(dst, fmt.Sprintf("%s, ", store[i+j]))
		}
		WriteAppend(dst, fmt.Sprintf("%s\n", store[i+step-1]))
	}

}

// Check Valid Directory
func (this Task) IsValidDirectory() bool {
	return IsExistsFile(PathJoin(this.SrcDir, "resultsMap.xml")) &&
		IsExistsFile(PathJoin(this.SrcDir, "results.xml")) &&
		IsExistsFile(PathJoin(this.AcePath))
}

// write ACE script to SrcDir
func (task *Task) writeAce() {
	ace := mkExtractAce(*task)
	path := ""
	for _, v := range task.Signals {
		path = path + v + "_"
	}
	path = fmt.Sprintf("%s%.2f_%.2f_%.2f.ace", path, task.Range.Start, task.Range.Step, task.Range.Stop)
	path = PathJoin(task.SrcDir, path)
	Write(path, ace)
	task.AcePath = path
}

// make ACE script for wave viewer
func mkExtractAce(task Task) string {
	var rt string = fmt.Sprintf(`set xml [ sx_open_wdf "resultsMap.xml" ]
sx_export_csv on
%s`, task.Range.ToCommandString())

	for _, v := range task.Signals {
		rt = fmt.Sprintf("%s\nset www [ sx_find_wave_in_file $xml %s ]\nsx_export_data \"%s.csv\" $www",
			rt, v, v)
	}
	return rt
}

func init() {
	rootCmd.AddCommand(extractCmd)
	wd, _ := os.Getwd()
	extractCmd.Flags().StringP("dst", "d", config.BaseDir, "書き出しディレクトリです")
	extractCmd.Flags().StringP("src", "s", wd, "ターゲットのディレクトリです")
	extractCmd.Flags().StringSliceP("signalName", "S", []string{"N1", "N2", "BLB", "BL"}, "取り出したい波形のリストです")
	extractCmd.Flags().Float64("start", 0, "プロットの最初の時間[ns]です")
	extractCmd.Flags().Float64("step", 0, "プロットの刻み幅[ns]です")
	extractCmd.Flags().Float64("stop", 0, "プロットの最後の時間[ns]です")
}
