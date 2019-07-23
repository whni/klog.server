package main

import (
	"bytes"
	"fmt"
	"github.com/snabb/isoweek"
	"os"
	"time"
)

/* list equal check */
func isStringListEqual(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for index := range x {
		if x[index] != y[index] {
			return false
		}
	}
	return true
}

func isIntListEqual(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}
	for index := range x {
		if x[index] != y[index] {
			return false
		}
	}
	return true
}

func isInt64ListEqual(x, y []int64) bool {
	if len(x) != len(y) {
		return false
	}
	for index := range x {
		if x[index] != y[index] {
			return false
		}
	}
	return true
}

func isUintListEqual(x, y []uint) bool {
	if len(x) != len(y) {
		return false
	}
	for index := range x {
		if x[index] != y[index] {
			return false
		}
	}
	return true
}

func isUint64ListEqual(x, y []uint64) bool {
	if len(x) != len(y) {
		return false
	}
	for index := range x {
		if x[index] != y[index] {
			return false
		}
	}
	return true
}

func isFloat32ListEqual(x, y []float32) bool {
	if len(x) != len(y) {
		return false
	}
	for index := range x {
		if x[index] != y[index] {
			return false
		}
	}
	return true
}

func isFloat64ListEqual(x, y []float64) bool {
	if len(x) != len(y) {
		return false
	}
	for index := range x {
		if x[index] != y[index] {
			return false
		}
	}
	return true
}

func isFloat64List2DEqual(x, y [][]float64) bool {
	if len(x) != len(y) {
		return false
	}
	for index := range x {
		if !isFloat64ListEqual(x[index], y[index]) {
			return false
		}
	}
	return true
}

/* time operations */
func getCurrentTimeStampUnixMillis() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

func getSecondsOfWeek(ts time.Time) uint64 {
	year, week := ts.ISOWeek()
	weekStart := isoweek.StartTime(year, week, time.Local)
	weekDuration := ts.Sub(weekStart)
	return uint64(weekDuration.Seconds())
}

/* string operations */
func cToGoString(c []byte) string {
	n := bytes.IndexByte(c, 0)
	if n < 0 {
		n = len(c)
	}
	return string(c[:n])
}

/* file operations */
func saveFile(name string, data []byte) error {
	if data == nil {
		return fmt.Errorf("empty data given")
	}

	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func moveFile(source string, destination string) error {
	err := os.Rename(source, destination)
	return err
}

func getFileSize(file string) int64 {
	fi, e := os.Stat(file)
	if e != nil {
		return 0
	}
	return fi.Size()
}
