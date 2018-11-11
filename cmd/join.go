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
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			Fatal("joinしたいファイルをスペース区切りで入力してね")
		}

		if !IsExistsFiles(args) {
			Fatal("見つからないファイルがあります")
		}

		jt := JoinTask{
			Files: args,
		}

		dst := jt.Join()
		if !IsExistsFile(dst) {
			Fatal("Joinに失敗しました")
		}

		Print("awx join: join complete:", dst)
	},
}

type JoinTask struct {
	Files []string
}

func (this JoinTask) Join() string {
	var box [][]string
	dst := ""
	var max = 0
	var length = len(this.Files)
	for key, value := range this.Files {
		base := filepath.Base(value)
		dst += strings.Replace(base, ".csv", "_", -1)
		cat := Cat(value)
		box = append(box, strings.Split(cat, "\n"))
		if max < len(box[key]) {
			max = len(box[key])
		}
	}

	dst = PathJoin(filepath.Dir(this.Files[0]), fmt.Sprintf("%s.csv", dst))

	for i := 0; i < max; i++ {
		line := fmt.Sprint(i)
		for j := 0; j < length; j++ {
			if len(box[j]) < i {
				continue
			}
			line += "," + strings.Trim(box[j][i], "\n")
		}
		line += "\n"

		WriteAppend(dst, line)
	}

	return dst
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
