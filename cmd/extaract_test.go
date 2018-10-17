package cmd

import (
	"os"
	"testing"
)

func TestAllExtract(t *testing.T) {
	home := os.Getenv("HOME")
	dst := PathJoin(home, "Dst")
	TryMkdirAll(dst)
	src := PathJoin(home, "Src")
	TryMkdirAll(src)
	r := Range{
		Start: 1,
		Step:  2,
		Stop:  3,
	}
	t.Run("001_mkExtractAce", func(t *testing.T) {
		actual := mkExtractAce(NewTask("", "", r, []string{"A", "B", "C"}))
		expect := `set xml [ sx_open_wdf "resultsMap.xml" ]
sx_export_csv on
sx_export_range 1.00ns 3.00ns 2.00ns
set www [ sx_find_wave_in_file $xml A ]
sx_export_data "A.csv" $www
set www [ sx_find_wave_in_file $xml B ]
sx_export_data "B.csv" $www
set www [ sx_find_wave_in_file $xml C ]
sx_export_data "C.csv" $www`
		Equal(actual, expect, t)
	})

	t.Run("002_writeAce", func(t *testing.T) {
		signals := []string{"A", "B", "C"}
		var task Task = NewTask(dst, src, r, signals)
		task.writeAce()

		Equal(PathJoin(src, "A_B_C_1.00_2.00_3.00.ace"), task.AcePath, t)
		Assert(IsExistsFile(task.AcePath), t)
		expect := `set xml [ sx_open_wdf "resultsMap.xml" ]
sx_export_csv on
sx_export_range 1.00ns 3.00ns 2.00ns
set www [ sx_find_wave_in_file $xml A ]
sx_export_data "A.csv" $www
set www [ sx_find_wave_in_file $xml B ]
sx_export_data "B.csv" $www
set www [ sx_find_wave_in_file $xml C ]
sx_export_data "C.csv" $www`
		Equal(Cat(task.AcePath), expect, t)
		RemoveFile(task.AcePath)
	})

	t.Run("003_Extract", func(t *testing.T) {
		task := NewTask(dst, src, Range{
			Start: 2.5,
			Step:  7.5,
			Stop:  17.5,
		}, []string{"A"})

		actual := task.Extract()
		expect := []string{
			PathJoin(task.DstDir, "A.csv"),
		}

		for i, v := range actual {
			Equal(v, expect[i], t)
		}
	})

	t.Run("004_Run", func(t *testing.T) {
		CreateFile(PathJoin(src, "resultsMap.xml"))
		CreateFile(PathJoin(src, "results.xml"))
		task := NewTask(dst, src, Range{
			Start: 2.5,
			Step:  7.5,
			Stop:  17.5,
		}, []string{"A"})
		sum, err := task.Run()
		if err != nil {
			t.Fatal(err)
		}
		Assert(sum.Status, t)
	})

	t.Run("XXX_Remove", func(t *testing.T) {
		RemoveFile(dst)
		RemoveFile(src)
	})
}
