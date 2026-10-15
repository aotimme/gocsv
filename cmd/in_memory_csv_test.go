package cmd

import (
	"slices"
	"testing"
)

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

func TestSampleRowIndices(t *testing.T) {
	// a 5-row CSV
	var imc = InMemoryCsv{rows: [][]string{{}, {}, {}, {}, {}}}

	t.Run("seeding", func(t *testing.T) {
		var (
			// random

			got_0 = imc.SampleRowIndices(5, false, 0)

			// seeded

			got_1  = imc.SampleRowIndices(5, false, 1)
			want_1 = []int{0, 4, 2, 3, 1}

			got_2  = imc.SampleRowIndices(5, false, 2)
			want_2 = []int{3, 0, 1, 2, 4}
		)

		if !slices.Equal(got_1, want_1) {
			t.Errorf("for seed=1, got %v; want %v", got_1, want_1)
		}
		if !slices.Equal(got_2, want_2) {
			t.Errorf("for seed=2, got %v; want %v", got_2, want_2)
		}
		if slices.Equal(got_0, got_1) || slices.Equal(got_0, got_2) {
			t.Errorf("got_0=%v, should not match got_1=%v or got_2=%v", got_0, got_1, got_2)
		}
	})

	t.Run("without_replacement", func(t *testing.T) {
		got := imc.SampleRowIndices(5, false, 1)
		want := []int{0, 4, 2, 3, 1}
		if !slices.Equal(got, want) {
			t.Errorf("got %v; want %v", got, want)
		}
	})

	t.Run("with_replacement", func(t *testing.T) {
		got := imc.SampleRowIndices(10, true, 1)
		want := []int{1, 2, 2, 4, 1, 3, 0, 0, 1, 0}
		if !slices.Equal(got, want) {
			t.Errorf("got %v; want %v", got, want)
		}
	})
}
