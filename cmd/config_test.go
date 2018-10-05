package cmd

import "testing"

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
		expect := "sx_export_range 1.00 2.00 3.00"
		Equal(actual, expect, t)
	})

}
