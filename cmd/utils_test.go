package cmd

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

func TestGetBaseFilenameWithoutExtension(t *testing.T) {
	testCases := []struct {
		filename     string
		baseFilename string
	}{
		{"../test-files/simple.csv", "simple"},
		{"../test-files/simple", "simple"},
		{"simple.csv", "simple"},
		{"/simple.csv", "simple"},
		{"./simple.csv", "simple"},
	}
	for _, tt := range testCases {
		t.Run(tt.filename, func(t *testing.T) {
			baseFilename := GetBaseFilenameWithoutExtension(tt.filename)
			if baseFilename != tt.baseFilename {
				t.Errorf("Expected %q but got %q", tt.baseFilename, baseFilename)
			}
		})
	}
}

func TestValidGetDelimiterFromString(t *testing.T) {
	testCases := []struct {
		delimiter string
		comma     rune
	}{
		{",", ','},
		{";", ';'},
		{"\\t", '\t'},
		{"|", '|'},
		{"\\x01", '\x01'},
		{"\\u0001", '\x01'},
	}
	for _, tt := range testCases {
		t.Run(tt.delimiter, func(t *testing.T) {
			delimiterRune, err := GetDelimiterFromString(tt.delimiter)
			if err != nil {
				t.Errorf("Expected \"%#U\" but instead got an error: %v", tt.comma, err)
			}
			if delimiterRune != tt.comma {
				t.Errorf("Expected \"%#U\" but got \"%#U\"", tt.comma, delimiterRune)
			}
		})
	}
}

func TestInvalidGetDelimiterFromString(t *testing.T) {
	testCases := []struct {
		delimiter string
	}{
		{""},
		{"lolcats"},
	}
	for _, tt := range testCases {
		t.Run(tt.delimiter, func(t *testing.T) {
			delimiterRune, err := GetDelimiterFromString(tt.delimiter)
			if err == nil {
				t.Errorf("Expected an error for delimiter \"%s\" but instead got rune \"%#U\"", tt.delimiter, delimiterRune)
			}
		})
	}
}
