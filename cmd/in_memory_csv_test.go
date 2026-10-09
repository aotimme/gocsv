package cmd

import "testing"

func TestIntColumnStatsCalculateMedian(t *testing.T) {
	testCases := []struct {
		name   string
		values []int64
		median float64
	}{
		{"two values", []int64{1, 3}, 2},
		{"fractional median", []int64{1, 2}, 1.5},
		{"unsorted negative values", []int64{-1, -9, 7, -5}, -3},
		{"repeated middle values", []int64{9, 1, 5, 5}, 5},
		{"odd count", []int64{9, 1, 5}, 5},
		{"single value", []int64{-2}, -2},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			stats := NewIntColumnsStats(tt.values)
			stats.CalculateMedian()
			if stats.median != tt.median {
				t.Errorf("expected median %v but got %v", tt.median, stats.median)
			}
		})
	}
}

func TestFloatColumnStatsCalculateMedian(t *testing.T) {
	testCases := []struct {
		name   string
		values []float64
		median float64
	}{
		{"two values", []float64{1.5, 3.5}, 2.5},
		{"fractional median", []float64{1.5, 2.5}, 2},
		{"unsorted negative values", []float64{-1.5, -9.5, 7.5, -5.5}, -3.5},
		{"repeated middle values", []float64{9.5, 1.5, 5.5, 5.5}, 5.5},
		{"odd count", []float64{9.5, 1.5, 5.5}, 5.5},
		{"single value", []float64{-2.5}, -2.5},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			stats := NewFloatColumnsStats(tt.values)
			stats.CalculateMedian()
			if stats.median != tt.median {
				t.Errorf("expected median %v but got %v", tt.median, stats.median)
			}
		})
	}
}
