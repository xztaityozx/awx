package cmd

import (
	"fmt"
)

type AwxConfig struct {
	BaseDir string
	Range   Range
	WVPath  string
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
	return fmt.Sprintf("sx_export_range %.2fns %.2fns %.2fns", r.Start, r.Step, r.Stop)
}

type Task struct {
	DstDir  string
	SrcDir  string
	Range   Range
	Signals []string
	AcePath string
}

// make Task struct
func NewTask(dst string, src string, r Range, signals []string) Task {
	return Task{
		DstDir:  dst,
		SrcDir:  src,
		Range:   r,
		Signals: signals,
		AcePath: "",
	}
}

func (t Task) CompareTo(s Task) bool {
	if t.Range != s.Range {
		return false
	}
	if t.DstDir != s.DstDir {
		return false
	}
	if t.SrcDir != s.SrcDir {
		return false
	}
	if t.AcePath != s.AcePath {
		return false
	}
	for i, v := range t.Signals {
		if v != s.Signals[i] {
			return false
		}
	}
	return true
}
