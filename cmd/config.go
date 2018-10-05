package cmd

import "fmt"

type AwxConfig struct {
	BaseDir        string
	Range          Range
	ResultNameRule string
}

type Range struct {
	Start float64
	Step  float64
	Stop  float64
}

// new Range struct
func NewRange(start float64, step float64, stop float64) Range {
	return Range{
		Start: start,
		Step:  step,
		Stop:  stop,
	}
}

// to command string
func (r Range) ToCommandString() string {
	return fmt.Sprintf("sx_export_range %.2f %.2f %.2f", r.Start, r.Step, r.Stop)
}
