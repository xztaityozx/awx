package cmd

import (
	"testing"
)

func Assert(status bool, t *testing.T) {
	if !status {
		t.Fatal(t.Name(), "Faild")
	}
}

func Equal(actual interface{}, expect interface{}, t *testing.T) {
	if actual != expect {
		t.Fatal(actual, "is not", expect)
	}
}

func TestAllConfig(t *testing.T) {
	t.Run("001_NewRange", func(t *testing.T) {
		actual := NewRange(0, 0, 0)
		expect := Range{
			Start: 0,
			Step:  0,
			Stop:  0,
		}

		Equal(actual, expect, t)
	})

	t.Run("002_NewRange", func(t *testing.T) {
		actual := NewRange(1, 2, 3)
		expect := Range{
			Start: 1,
			Step:  2,
			Stop:  3,
		}

		Equal(actual, expect, t)
	})

	t.Run("003_ToCommandString", func(t *testing.T) {
		actual := NewRange(1, 2, 3).ToCommandString()
		expect := "sx_export_range 1.00ns 2.00ns 3.00ns"
		Equal(actual, expect, t)
	})

	t.Run("004_NewTask", func(t *testing.T) {
		actual := NewTask("dst", "src", Range{
			Start: 10,
			Step:  20,
			Stop:  30,
		}, []string{""})
		expect := Task{
			DstDir: "dst",
			SrcDir: "src",
			Range: Range{
				Start: 10,
				Step:  20,
				Stop:  30,
			},
			Signals: []string{""},
			AcePath: "",
		}
		Assert(expect.CompareTo(actual), t)
	})

	t.Run("005_Count", func(t *testing.T) {
		actual := Range{
			Start: 2.5,
			Step:  7.5,
			Stop:  17.5,
		}.Count()
		expect := 3

		Equal(actual, expect, t)
	})
}
