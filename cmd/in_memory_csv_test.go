package cmd

import "testing"

func TestMedian(t *testing.T) {
	testCases := []struct {
		values []int64
		want   float64
	}{
		{[]int64{0}, 0},
		{[]int64{0, 1}, 0.5},
		{[]int64{0, 1, 2}, 1},
		{[]int64{0, 1, 2, 3}, 1.5},

		{[]int64{-1, 0, 1}, 0}, // negative

		{[]int64{1, 0, -1}, 0}, // unsorted
	}

	t.Run("int", func(t *testing.T) {
		for _, tc := range testCases {
			stats := NewIntColumnsStats(tc.values)
			stats.CalculateMedian()
			if stats.median != tc.want {
				t.Errorf("for %v, got %v; want %v", tc.values, stats.median, tc.want)
			}
		}
	})

	t.Run("float", func(t *testing.T) {
		for _, tc := range testCases {
			values := make([]float64, len(tc.values))
			for i := range tc.values {
				values[i] = float64(tc.values[i])
			}

			stats := NewFloatColumnsStats(values)
			stats.CalculateMedian()
			if stats.median != tc.want {
				t.Errorf("for %v, got %v; want %v", values, stats.median, tc.want)
			}
		}
	})
}
