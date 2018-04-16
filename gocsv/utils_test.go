package main

import "testing"

func TestGetIndicesForColumns(t *testing.T) {
	testCases := []struct {
		headers []string
		columns []string
		indices []int
	}{
		{[]string{"what", "me", "worry"}, []string{"me"}, []int{1}},
		{[]string{"what", "me", "worry"}, []string{"me", "me"}, []int{1, 1}},
		{[]string{"what", "me", "worry"}, []string{"1"}, []int{0}},
		{[]string{"what", "me", "worry"}, []string{"1-2"}, []int{0, 1}},
		{[]string{"what", "me", "worry"}, []string{"2-1"}, []int{1, 0}},
		{[]string{"what", "me", "worry"}, []string{"1-2", "1-3"}, []int{0, 1, 0, 1, 2}},
		{[]string{"what", "me", "worry"}, []string{"1-3"}, []int{0, 1, 2}},
		{[]string{"what", "me", "worry"}, []string{"-2"}, []int{0, 1}},
		{[]string{"what", "me", "worry"}, []string{"2-"}, []int{1, 2}},
		{[]string{"what", "4", "worry"}, []string{"4"}, []int{1}},
		{[]string{"what", "4-", "worry"}, []string{"4-"}, []int{1}},
		{[]string{"what", "me", "me"}, []string{"me"}, []int{1, 2}},
	}
	for i, testCase := range testCases {
		indices, err := GetIndicesForColumns(testCase.headers, testCase.columns)
		if err != nil {
			t.Errorf("Test case %d: unexpected error: %v", i, err)
		}
		if len(indices) != len(testCase.indices) {
			t.Errorf("Test case %d: expected length of %d but got length of %d", i, len(testCase.indices), len(indices))
		}
		for j, index := range indices {
			if index != testCase.indices[j] {
				t.Errorf("Test case %d: expected value %d at index %d but got %d", i, testCase.indices[j], j, index)
			}
		}
	}
}
