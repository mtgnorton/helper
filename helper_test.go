package helper

import (
	"fmt"
	"testing"
)

func TestRollingAverage_Update(t *testing.T) {
	initialValue := 0.0
	RtBeta := 0.9

	ra := NewRollingAverage(initialValue, RtBeta)

	// Simulate updating the rolling average with new values
	newValues := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}

	for _, value := range newValues {
		ra.Add(value)
		count, average, ma, cma := ra.Get()
		fmt.Printf("count:%v, average:%v, ma:%v, cma:%v \n", count, average, ma, cma)
	}
}
