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
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "count",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("count called")
	},
}

type CountTask struct {
	IsRangeSEEDCountUp bool
	IgnoreSigma        bool
	Rule               string
	Target             string
}

// new count task
func NewCountTask(args []string) CountTask {
	var rt = CountTask{
		IsRangeSEEDCountUp: IsRangeSEEDCountUp,
		IgnoreSigma:        IgnoreSigma,
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

var IsRangeSEEDCountUp bool
var IgnoreSigma bool

var RuleFile string

func init() {
	rootCmd.AddCommand(countCmd)
	countCmd.Flags().BoolVarP(&IsRangeSEEDCountUp, "RangeSEED", "R", false, "RangeSEEDシミュレーションの結果を数え上げます")
	countCmd.Flags().BoolVarP(&IgnoreSigma, "ignoreSigma", "G", false, "Sigmaの値も一緒に出力します")
	countCmd.Flags().StringVarP(&RuleFile, "ruleFile", "t", "", "ルールを記述したファイルを指定できます")
	countCmd.Flags().Bool("make", false, "ルールスクリプトを対話形式で作成します")
}
