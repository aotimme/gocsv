package cmd

import "testing"

func TestGetCellWidth(t *testing.T) {
	testCases := []struct {
		cell     string
		maxLines int
		width    int
	}{
		{"what me worry", 0, 13},
		{"what me worry", 1, 13},
		{"what me worry", 2, 13},
		{"what\nme worry", 0, 8},
		{"what\nme worry", 1, 4},
		{"what\nme worry", 2, 8},
		{"smiley 😊 face", 0, 13},
	}
	for i, testCase := range testCases {
		cellWidth := getCellWidth(testCase.cell, testCase.maxLines)
		if cellWidth != testCase.width {
			t.Errorf("Test case %d: expected cell width of %d but got %d", i, testCase.width, cellWidth)
		}
	}
}

func TestGetTruncatedLine(t *testing.T) {
	testCases := []struct {
		line          string
		width         int
		truncatedLine string
	}{
		{"what me worry", 10, "what me..."},
		{"what me worry", 11, "what me ..."},
		{"what me worry", 12, "what me w..."},
		{"what me worry", 13, "what me worry"},
		{"what me worry", 14, "what me worry "},
		{"what me worry", 15, "what me worry  "},
		// https://github.com/aotimme/gocsv/issues/47
		{"Foobarbaz 日本のルーン", 14, "Foobarbaz 日..."},
		{"Foobarbaz 日本のルーン", 15, "Foobarbaz 日本..."},
		{"Foobarbaz 日本のルーン", 16, "Foobarbaz 日本のルーン"},
		{"Foobarbaz 日本のルーン", 17, "Foobarbaz 日本のルーン "},
		{"Foobarbaz 日本のルーン", 18, "Foobarbaz 日本のルーン  "},
	}
	for _, testCase := range testCases {
		truncatedLine := getTruncatedLine(testCase.line, testCase.width)
		if truncatedLine != testCase.truncatedLine {
			t.Errorf("getTruncatedLine(%q, %d) = %q; want %q", testCase.line, testCase.width, truncatedLine, testCase.truncatedLine)
		}
	}
}

func TestGetCellHeight(t *testing.T) {
	testCases := []struct {
		cell       string
		maxLines   int
		cellHeight int
	}{
		{"what me worry", 0, 1},
		{"what\nme worry", 0, 2},
		{"what\nme worry", 1, 1},
		{"what\nme worry", 2, 2},
		{"what\nme worry", 1, 1},
		{"what\nme worry", 2, 2},
		{"what\nme worry\n", 0, 3},
		{"what\nme worry\n", 1, 1},
		{"what\nme worry\n", 2, 2},
	}
	for i, testCase := range testCases {
		cellHeight := getCellHeight(testCase.cell, testCase.maxLines)
		if cellHeight != testCase.cellHeight {
			t.Errorf("Test case %d: expected cell height %d but got %d", i, testCase.cellHeight, cellHeight)
		}
	}
}
