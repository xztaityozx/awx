package cmd

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	var sum = Summary{
		Files:  []string{"A.csv", "B.csv", "C.csv"},
		Status: true,
	}
	t.Run("001_ToString", func(t *testing.T) {
		actual := sum.ToString()
		expect := fmt.Sprint(`AWX Summary:
	Status: Success
	Files:
		A.csv
		B.csv
		C.csv
`)
		Equal(actual, expect, t)
	})
}
