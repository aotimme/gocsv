package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type ViewSubcommand struct {
	maxWidth int
	maxLines int
	maxRows  int
}

func (sub *ViewSubcommand) Name() string {
	return "view"
}
func (sub *ViewSubcommand) Aliases() []string {
	return []string{}
}
func (sub *ViewSubcommand) Description() string {
	return "Display a CSV in a pretty tabular format."
}
func (sub *ViewSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.IntVar(&sub.maxWidth, "max-width", 0, "Maximum width per column")
	fs.IntVar(&sub.maxWidth, "w", 0, "Maximum width per column (shorthand)")
	fs.IntVar(&sub.maxLines, "max-lines", 0, "Maximum number of lines per cell")
	fs.IntVar(&sub.maxLines, "l", 0, "Maximum number of lines per cell (shorthand)")
	fs.IntVar(&sub.maxRows, "n", 0, "Number of rows to display")
}

func (sub *ViewSubcommand) Run(args []string) {
	if sub.maxWidth < 0 {
		fmt.Fprintln(os.Stderr, "Invalid argument --max-width")
		os.Exit(1)
	}
	if sub.maxLines < 0 {
		fmt.Fprintln(os.Stderr, "Invalid argument --max-lines")
		os.Exit(1)
	}
	if sub.maxRows < 0 {
		sub.maxRows = 0
	}

	inputCsvs := GetInputCsvsOrPanic(args, 1)
	View(inputCsvs[0], sub.maxWidth, sub.maxLines, sub.maxRows)
}

func View(inputCsv *InputCsv, maxWidth, maxLines, maxRows int) {

	imc := NewInMemoryCsvFromInputCsv(inputCsv)

	// Default to 0
	columnWidths := make([]int, imc.NumColumns())
	for j, cell := range imc.header {
		cellLength := getCellWidth(cell)
		if cellLength > columnWidths[j] {
			if maxWidth > 0 && cellLength > maxWidth {
				columnWidths[j] = maxWidth
			} else {
				columnWidths[j] = cellLength
			}
		}
	}

	// Get the actual number of rows to display
	numRowsToView := imc.NumRows()
	if maxRows > 0 && maxRows < numRowsToView {
		numRowsToView = maxRows
	}

	for i := 0; i < numRowsToView; i++ {
		row := imc.rows[i]
		for j, cell := range row {
			if columnWidths[j] == maxWidth {
				continue
			}
			cellLength := getCellWidth(cell)
			if cellLength > columnWidths[j] {
				if maxWidth > 0 && cellLength > maxWidth {
					columnWidths[j] = maxWidth
				} else {
					columnWidths[j] = cellLength
				}
			}
		}
	}

	rowSeparator := getRowSeparator(columnWidths)

	// Top of table
	fmt.Println(rowSeparator)

	// Print header
	printRow(imc.header, columnWidths, maxLines)
	fmt.Println(rowSeparator)

	// Print rows
	for i := 0; i < numRowsToView; i++ {
		row := imc.rows[i]
		printRow(row, columnWidths, maxLines)
		fmt.Println(rowSeparator)
	}
}

func getRowSeparator(widths []int) string {
	cells := make([]string, len(widths))
	for i, width := range widths {
		cells[i] = strings.Repeat("-", width)
	}
	return fmt.Sprintf("+-%s-+", strings.Join(cells, "-+-"))
}

func getCellWidth(cell string) int {
	indexOfNewline := strings.Index(cell, "\n")
	if indexOfNewline > -1 {
		return indexOfNewline + 1
	} else {
		return len(cell)
	}
}

func printRow(row []string, columnWidths []int, maxLines int) {
	rowHeight := getRowHeight(row, maxLines)
	outrowLines := make([][]string, rowHeight)
	for i, _ := range outrowLines {
		outrowLines[i] = make([]string, len(row))
	}
	copyTruncatedAndPaddedCellToOutputRow(outrowLines, row, columnWidths)
	for _, line := range outrowLines {
		fmt.Printf("| %s |\n", strings.Join(line, " | "))
	}
}

func getRowHeight(row []string, maxLines int) int {
	rowHeight := 1
	for _, cell := range row {
		numLines := strings.Count(cell, "\n") + 1
		if maxLines > 0 && numLines >= maxLines {
			return maxLines
		}
		if numLines > rowHeight {
			rowHeight = numLines
		}
	}
	return rowHeight
}

func copyTruncatedAndPaddedCellToOutputRow(outrowLines [][]string, row []string, widths []int) {
	for j, cell := range row {
		cellLines := strings.Split(cell, "\n")
		width := widths[j]
		for i, cell := range outrowLines {
			if len(cellLines) > i {
				cell[j] = getTruncatedLine(cellLines[i], width)
			} else {
				cell[j] = strings.Repeat(" ", width)
			}
		}
	}
}

func getTruncatedLine(line string, width int) string {
	if len(line) == width {
		return line
	} else if len(line) < width {
		return line + strings.Repeat(" ", width-len(line))
	} else {
		return line[:width-3] + "..."
	}
}
