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
sx_export_range 1.00ns 2.00ns 3.00ns
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
sx_export_range 1.00ns 2.00ns 3.00ns
set www [ sx_find_wave_in_file $xml A ]
sx_export_data "A.csv" $www
set www [ sx_find_wave_in_file $xml B ]
sx_export_data "B.csv" $www
set www [ sx_find_wave_in_file $xml C ]
sx_export_data "C.csv" $www`
		Equal(Cat(task.AcePath), expect, t)
		RemoveFile(task.AcePath)
	})
}
