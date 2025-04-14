package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

type InMemoryCsv struct {
	header []string
	rows   [][]string

	// index of column
	isIndexed bool
	index     map[string][]int
}

func NewInMemoryCsvFromInputCsv(inputCsv *InputCsv) *InMemoryCsv {
	imc := new(InMemoryCsv)
	rows, err := inputCsv.ReadAll()
	if err != nil {
		ExitWithError(err)
	}
	imc.header = rows[0]
	imc.rows = rows[1:]
	imc.isIndexed = false
	return imc
}

func (imc *InMemoryCsv) Index(columnIndex int) {
	imc.index = make(map[string][]int)
	for i, row := range imc.rows {
		rowval := row[columnIndex]
		group, ok := imc.index[rowval]
		if ok {
			group = append(group, i)
		} else {
			group = make([]int, 1)
			group[0] = i
		}
		imc.index[rowval] = group
	}
}

func (imc *InMemoryCsv) NumRows() int {
	return len(imc.rows)
}

func (imc *InMemoryCsv) NumColumns() int {
	return len(imc.header)
}

func (imc *InMemoryCsv) GetRowIndicesMatchingIndexedColumn(value string) []int {
	indices, ok := imc.index[value]
	if ok {
		return indices
	} else {
		return make([]int, 0)
	}
}

func (imc *InMemoryCsv) GetRowsMatchingIndexedColumn(value string) [][]string {
	indices := imc.GetRowIndicesMatchingIndexedColumn(value)
	rows := make([][]string, 0)
	for _, idx := range indices {
		rows = append(rows, imc.rows[idx])
	}
	return rows
}

type ColumnValueIterator struct {
	imc         *InMemoryCsv
	columnIndex int
	rowIndex    int
}

func NewColumnValueIterator(imc *InMemoryCsv, columnIndex int) *ColumnValueIterator {
	return &ColumnValueIterator{
		imc:         imc,
		columnIndex: columnIndex,
		rowIndex:    0,
	}
}
func (cvi *ColumnValueIterator) Next() (string, bool) {
	if cvi.rowIndex >= len(cvi.imc.rows) {
		return "", false
	}
	retval := cvi.imc.rows[cvi.rowIndex][cvi.columnIndex]
	cvi.rowIndex++
	return retval, true
}

func (imc *InMemoryCsv) InferType(columnIndex int) ColumnType {
	cvi := NewColumnValueIterator(imc, columnIndex)
	return InferTypeFromStringIterator(cvi)
}

func (imc *InMemoryCsv) SortRows(columnIndices []int, columnTypes []ColumnType, stable bool, reverse bool) {
	isLessFunc := func(row1Ptr, row2Ptr *[]string) bool {
		row1 := *row1Ptr
		row2 := *row2Ptr
		for i, columnIndex := range columnIndices {
			isElem1Null := IsNullType(row1[columnIndex])
			isElem2Null := IsNullType(row2[columnIndex])
			if isElem1Null && isElem2Null {
				continue
			}
			if isElem1Null && !isElem2Null {
				return true
			}
			if !isElem1Null && isElem2Null {
				return false
			}
			columnType := columnTypes[i]
			if columnType == FLOAT_TYPE {
				row1Val := ParseFloat64OrPanic(row1[columnIndex])
				row2Val := ParseFloat64OrPanic(row2[columnIndex])
				if row1Val < row2Val {
					return true
				} else if row1Val > row2Val {
					return false
				}
			} else if columnType == INT_TYPE {
				row1Val := ParseInt64OrPanic(row1[columnIndex])
				row2Val := ParseInt64OrPanic(row2[columnIndex])
				if row1Val < row2Val {
					return true
				} else if row1Val > row2Val {
					return false
				}
			} else if columnType == DATETIME_TYPE {
				row1Val := ParseDatetimeOrPanic(row1[columnIndex])
				row2Val := ParseDatetimeOrPanic(row2[columnIndex])
				if row1Val.Before(row2Val) {
					return true
				} else if row1Val.After(row2Val) {
					return false
				}
			} else if columnType == DATE_TYPE {
				row1Val := ParseDateOrPanic(row1[columnIndex])
				row2Val := ParseDateOrPanic(row2[columnIndex])
				if row1Val.Before(row2Val) {
					return true
				} else if row1Val.After(row2Val) {
					return false
				}
			} else {
				row1Val := row1[columnIndex]
				row2Val := row2[columnIndex]
				if row1Val < row2Val {
					return true
				} else if row1Val > row2Val {
					return false
				}
			}
		}
		return false // all values to sort by are equal (hence, not 'less')
	}

	SortRowsBy(isLessFunc).Sort(imc.rows, stable, reverse)
}

func (imc *InMemoryCsv) SampleRowIndicesWithReplacement(numRows, seed int) []int {
	totalRows := imc.NumRows()
	retval := make([]int, numRows)
	for i := 0; i < numRows; i++ {
		retval[i] = rand.Intn(totalRows)
	}
	return retval
}

func (imc *InMemoryCsv) SampleRowIndicesWithoutReplacement(numRows, seed int) []int {
	totalRows := imc.NumRows()
	permuted := rand.Perm(totalRows)
	retval := make([]int, numRows)
	for i := 0; i < numRows; i++ {
		retval[i] = permuted[i]
	}
	return retval
}

func (imc *InMemoryCsv) SampleRowIndices(numRows int, replace bool, seed int) []int {
	// NOTE: Updating global `rand` variable for the life of the proces...
	if seed != 0 {
		rand.Seed(int64(seed))
	} else {
		rand.Seed(time.Now().UTC().UnixNano())
	}
	if replace {
		return imc.SampleRowIndicesWithReplacement(numRows, seed)
	} else {
		return imc.SampleRowIndicesWithoutReplacement(numRows, seed)
	}
}

func (imc *InMemoryCsv) PrintStats() {
	numColumns := imc.NumColumns()
	buffers := make([]bytes.Buffer, numColumns)

	// Use a WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup
	wg.Add(numColumns)

	// Process each column concurrently
	for i := 0; i < numColumns; i++ {
		go func(columnIndex int) {
			defer wg.Done()
			bufferedWriter := bufio.NewWriter(&buffers[columnIndex])
			imc.FprintStatsForColumn(bufferedWriter, columnIndex)
			bufferedWriter.Flush()
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Print all column stats in order
	for i := 0; i < numColumns; i++ {
		fmt.Print(buffers[i].String())
	}

	fmt.Printf("Number of rows: %d\n", imc.NumRows())
}

func (imc *InMemoryCsv) GetPrintStatsForColumn(columnIndex int) bytes.Buffer {
	var buf bytes.Buffer
	bufferedWriter := bufio.NewWriter(&buf)
	imc.FprintStatsForColumn(bufferedWriter, columnIndex)
	bufferedWriter.Flush()
	return buf
}

func (imc *InMemoryCsv) FprintStatsForColumn(w io.Writer, columnIndex int) {
	fmt.Fprintf(w, "%d. %s\n", columnIndex+1, imc.header[columnIndex])
	columnType := imc.InferType(columnIndex)
	fmt.Fprintf(w, "  Type: %s\n", ColumnTypeToString(columnType))
	imc.FprintColumnNumberNulls(w, columnIndex)
	if columnType == NULL_TYPE {
		// continue
	} else if columnType == INT_TYPE {
		imc.FprintStatsForColumnAsInt(w, columnIndex)
	} else if columnType == FLOAT_TYPE {
		imc.FprintStatsForColumnAsFloat(w, columnIndex)
	} else if columnType == BOOLEAN_TYPE {
		imc.FprintStatsForColumnAsBoolean(w, columnIndex)
	} else if columnType == DATE_TYPE {
		imc.FprintStatsForColumnAsDatetimeWithFormat(w, columnIndex, time.DateOnly)
	} else if columnType == DATETIME_TYPE {
		imc.FprintStatsForColumnAsDatetimeWithFormat(w, columnIndex, time.RFC3339)
	} else if columnType == STRING_TYPE {
		imc.FprintStatsForColumnAsString(w, columnIndex)
	}
}

func (imc *InMemoryCsv) FprintColumnNumberNulls(w io.Writer, columnIndex int) {
	numNulls := imc.CountNullsInColumn(columnIndex)
	fmt.Fprintf(w, "  Number NULL: %d\n", numNulls)
}

func (imc *InMemoryCsv) CountNullsInColumn(columnIndex int) int {
	numNulls := 0
	for _, row := range imc.rows {
		cell := row[columnIndex]
		if IsNullType(cell) {
			numNulls += 1
		}
	}
	return numNulls
}

func (imc *InMemoryCsv) FprintStatsForColumnAsInt(w io.Writer, columnIndex int) {
	numNulls := imc.CountNullsInColumn(columnIndex)
	intArray := make([]int64, imc.NumRows()-numNulls)
	i := 0
	for _, row := range imc.rows {
		if !IsNullType(row[columnIndex]) {
			intArray[i] = ParseInt64OrPanic(row[columnIndex])
			i++
		}
	}
	ics := NewIntColumnsStats(intArray)
	ics.CalculateAllStats()

	fmt.Fprintf(w, "  Min: %d\n", ics.min)
	fmt.Fprintf(w, "  Max: %d\n", ics.max)
	fmt.Fprintf(w, "  Sum: %d\n", ics.sum)
	fmt.Fprintf(w, "  Mean: %f\n", ics.mean)
	fmt.Fprintf(w, "  Median: %f\n", ics.median)
	fmt.Fprintf(w, "  Standard Deviation: %f\n", ics.stdev)
	fmt.Fprintf(w, "  Unique values: %d\n", len(ics.valueCounts))
	numFrequent := 5
	if numFrequent > len(ics.valueCounts) {
		numFrequent = len(ics.valueCounts)
	}
	fmt.Fprintf(w, "  %d most frequent values:\n", numFrequent)
	for i := 0; i < numFrequent; i++ {
		fmt.Fprintf(w, "      %d: %d\n", ics.valueCounts[i].value, ics.valueCounts[i].count)
	}
}

type IntColumnStats struct {
	array               []int64
	min, max, sum       int64
	mean, median, stdev float64
	valueCounts         []IntValueCount
}

func NewIntColumnsStats(intArray []int64) *IntColumnStats {
	ics := new(IntColumnStats)
	ics.array = intArray
	return ics
}

func (ics *IntColumnStats) CalculateAllStats() {
	ics.CalculateMin()
	ics.CalculateMax()
	ics.CalculateSum()
	ics.CalculateMean()
	ics.CalculateMedian()
	ics.CalculateStdDev()
	ics.CalculateValueCounts()
}

func (ics *IntColumnStats) CalculateMin() {
	ics.min = math.MaxInt64
	for _, intVal := range ics.array {
		if intVal < ics.min {
			ics.min = intVal
		}
	}
}

func (ics *IntColumnStats) CalculateMax() {
	ics.max = math.MinInt64
	for _, intVal := range ics.array {
		if intVal > ics.max {
			ics.max = intVal
		}
	}
}

func (ics *IntColumnStats) CalculateSum() {
	ics.sum = 0
	for _, intVal := range ics.array {
		ics.sum += intVal
	}
}

func (ics *IntColumnStats) CalculateMean() {
	ics.mean = float64(ics.sum) / float64(len(ics.array))
}

type Int64Array []int64

func (a Int64Array) Len() int           { return len(a) }
func (a Int64Array) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Int64Array) Less(i, j int) bool { return a[i] < a[j] }

func (ics *IntColumnStats) CalculateMedian() {
	arrayLen := len(ics.array)
	sortedArray := make([]int64, arrayLen)
	copy(sortedArray, ics.array)
	sort.Sort(Int64Array(sortedArray))
	if len(sortedArray)%2 == 0 {
		ics.median = (float64(sortedArray[arrayLen/2-1]) + float64(sortedArray[arrayLen/2-1])) / 2.0
	} else {
		ics.median = float64(sortedArray[arrayLen/2])
	}
}

func (ics *IntColumnStats) CalculateStdDev() {
	sum := 0.0
	for _, intVal := range ics.array {
		elem := float64(intVal) - ics.mean
		sum += elem * elem
	}
	ics.stdev = math.Sqrt(sum / float64(len(ics.array)-1))
}

func (ics *IntColumnStats) CalculateValueCounts() {
	valueCountsMap := make(map[int64]int)
	for _, intVal := range ics.array {
		count, ok := valueCountsMap[intVal]
		if ok {
			valueCountsMap[intVal] = count + 1
		} else {
			valueCountsMap[intVal] = 1
		}
	}
	ics.valueCounts = make([]IntValueCount, len(valueCountsMap))
	i := 0
	for value, count := range valueCountsMap {
		ics.valueCounts[i] = IntValueCount{value, count}
		i++
	}
	sort.Sort(sort.Reverse(IntValueCountByCount(ics.valueCounts)))
}

type IntValueCount struct {
	value int64
	count int
}

type IntValueCountByCount []IntValueCount

func (a IntValueCountByCount) Len() int           { return len(a) }
func (a IntValueCountByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IntValueCountByCount) Less(i, j int) bool { return a[i].count < a[j].count }

func (imc *InMemoryCsv) FprintStatsForColumnAsFloat(w io.Writer, columnIndex int) {
	numNulls := imc.CountNullsInColumn(columnIndex)
	floatArray := make([]float64, imc.NumRows()-numNulls)
	i := 0
	for _, row := range imc.rows {
		if !IsNullType(row[columnIndex]) {
			floatArray[i] = ParseFloat64OrPanic(row[columnIndex])
			i++
		}
	}
	fcs := NewFloatColumnsStats(floatArray)
	fcs.CalculateAllStats()

	fmt.Fprintf(w, "  Min: %f\n", fcs.min)
	fmt.Fprintf(w, "  Max: %f\n", fcs.max)
	fmt.Fprintf(w, "  Sum: %f\n", fcs.sum)
	fmt.Fprintf(w, "  Mean: %f\n", fcs.mean)
	fmt.Fprintf(w, "  Median: %f\n", fcs.median)
	fmt.Fprintf(w, "  Standard Deviation: %f\n", fcs.stdev)
	fmt.Fprintf(w, "  Unique values: %d\n", len(fcs.valueCounts))
	numFrequent := 5
	if numFrequent > len(fcs.valueCounts) {
		numFrequent = len(fcs.valueCounts)
	}
	fmt.Fprintf(w, "  %d most frequent values:\n", numFrequent)
	for i := 0; i < numFrequent; i++ {
		fmt.Fprintf(w, "      %f: %d\n", fcs.valueCounts[i].value, fcs.valueCounts[i].count)
	}
}

type FloatColumnStats struct {
	array               []float64
	min, max, sum       float64
	mean, median, stdev float64
	valueCounts         []FloatValueCount
}

func NewFloatColumnsStats(floatArray []float64) *FloatColumnStats {
	fcs := new(FloatColumnStats)
	fcs.array = floatArray
	return fcs
}

func (fcs *FloatColumnStats) CalculateAllStats() {
	fcs.CalculateMin()
	fcs.CalculateMax()
	fcs.CalculateSum()
	fcs.CalculateMean()
	fcs.CalculateMedian()
	fcs.CalculateStdDev()
	fcs.CalculateValueCounts()
}

func (fcs *FloatColumnStats) CalculateMin() {
	fcs.min = math.Inf(1)
	for _, floatVal := range fcs.array {
		if floatVal < fcs.min {
			fcs.min = floatVal
		}
	}
}

func (fcs *FloatColumnStats) CalculateMax() {
	fcs.max = math.Inf(-1)
	for _, floatVal := range fcs.array {
		if floatVal > fcs.max {
			fcs.max = floatVal
		}
	}
}

func (fcs *FloatColumnStats) CalculateSum() {
	fcs.sum = 0.0
	for _, floatVal := range fcs.array {
		fcs.sum += floatVal
	}
}

func (fcs *FloatColumnStats) CalculateMean() {
	fcs.mean = fcs.sum / float64(len(fcs.array))
}

func (fcs *FloatColumnStats) CalculateMedian() {
	arrayLen := len(fcs.array)
	sortedArray := make([]float64, arrayLen)
	copy(sortedArray, fcs.array)
	sort.Float64s(sortedArray)
	if len(sortedArray)%2 == 0 {
		fcs.median = (float64(sortedArray[arrayLen/2-1]) + float64(sortedArray[arrayLen/2-1])) / 2.0
	} else {
		fcs.median = float64(sortedArray[arrayLen/2])
	}
}

func (fcs *FloatColumnStats) CalculateStdDev() {
	sum := 0.0
	for _, floatVal := range fcs.array {
		elem := float64(floatVal) - fcs.mean
		sum += elem * elem
	}
	fcs.stdev = math.Sqrt(sum / float64(len(fcs.array)-1))
}

func (fcs *FloatColumnStats) CalculateValueCounts() {
	valueCountsMap := make(map[float64]int)
	for _, floatVal := range fcs.array {
		count, ok := valueCountsMap[floatVal]
		if ok {
			valueCountsMap[floatVal] = count + 1
		} else {
			valueCountsMap[floatVal] = 1
		}
	}
	fcs.valueCounts = make([]FloatValueCount, len(valueCountsMap))
	i := 0
	for value, count := range valueCountsMap {
		fcs.valueCounts[i] = FloatValueCount{value, count}
		i++
	}
	sort.Sort(sort.Reverse(FloatValueCountByCount(fcs.valueCounts)))
}

type FloatValueCount struct {
	value float64
	count int
}

type FloatValueCountByCount []FloatValueCount

func (a FloatValueCountByCount) Len() int           { return len(a) }
func (a FloatValueCountByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FloatValueCountByCount) Less(i, j int) bool { return a[i].count < a[j].count }

func (imc *InMemoryCsv) FprintStatsForColumnAsBoolean(w io.Writer, columnIndex int) {
	numTrue := 0
	numFalse := 0
	for _, row := range imc.rows {
		value := strings.Trim(row[columnIndex], " ")
		if !IsNullType(value) {
			bval := ParseBooleanOrPanic(value)
			if bval {
				numTrue++
			} else {
				numFalse++
			}
		}
	}
	fmt.Fprintf(w, "  Number TRUE: %d\n", numTrue)
	fmt.Fprintf(w, "  Number FALSE: %d\n", numFalse)
}

func (imc *InMemoryCsv) FprintStatsForColumnAsDatetimeWithFormat(w io.Writer, columnIndex int, format string) {
	numNulls := imc.CountNullsInColumn(columnIndex)
	dateArray := make([]time.Time, imc.NumRows()-numNulls)
	i := 0
	for _, row := range imc.rows {
		if !IsNullType(row[columnIndex]) {
			dateArray[i] = ParseDatetimeOrPanic(row[columnIndex])
			i++
		}
	}
	dcs := NewDateColumnsStats(dateArray, format)
	dcs.CalculateAllStats()

	fmt.Fprintf(w, "  Min: %s\n", dcs.min.Format(format))
	fmt.Fprintf(w, "  Max: %s\n", dcs.max.Format(format))
	fmt.Fprintf(w, "  Unique values: %d\n", len(dcs.valueCounts))
	numFrequent := 5
	if numFrequent > len(dcs.valueCounts) {
		numFrequent = len(dcs.valueCounts)
	}
	fmt.Fprintf(w, "  %d most frequent values:\n", numFrequent)
	for i := 0; i < numFrequent; i++ {
		fmt.Fprintf(w, "      %s: %d\n", dcs.valueCounts[i].value, dcs.valueCounts[i].count)
	}
}

type DateColumnStats struct {
	array       []time.Time
	min, max    time.Time
	valueCounts []StringValueCount
	format      string
}

func NewDateColumnsStats(dateArray []time.Time, format string) *DateColumnStats {
	dcs := new(DateColumnStats)
	dcs.array = dateArray
	dcs.format = format
	return dcs
}

func (dcs *DateColumnStats) CalculateAllStats() {
	dcs.CalculateMin()
	dcs.CalculateMax()
	dcs.CalculateValueCounts()
}

func (dcs *DateColumnStats) CalculateMin() {
	for i, dateVal := range dcs.array {
		if i == 0 || dateVal.Before(dcs.min) {
			dcs.min = dateVal
		}
	}
}

func (dcs *DateColumnStats) CalculateMax() {
	for i, dateVal := range dcs.array {
		if i == 0 || dateVal.After(dcs.max) {
			dcs.max = dateVal
		}
	}
}

func (dcs *DateColumnStats) CalculateValueCounts() {
	valueCountsMap := make(map[string]int)
	for _, dateVal := range dcs.array {
		dateStr := dateVal.Format(dcs.format)
		count, ok := valueCountsMap[dateStr]
		if ok {
			valueCountsMap[dateStr] = count + 1
		} else {
			valueCountsMap[dateStr] = 1
		}
	}
	dcs.valueCounts = make([]StringValueCount, len(valueCountsMap))
	i := 0
	for value, count := range valueCountsMap {
		dcs.valueCounts[i] = StringValueCount{value, count}
		i++
	}
	sort.Sort(sort.Reverse(StringValueCountByCount(dcs.valueCounts)))
}

func (imc *InMemoryCsv) FprintStatsForColumnAsString(w io.Writer, columnIndex int) {
	numNulls := imc.CountNullsInColumn(columnIndex)
	stringArray := make([]string, imc.NumRows()-numNulls)
	i := 0
	for _, row := range imc.rows {
		if !IsNullType(row[columnIndex]) {
			stringArray[i] = row[columnIndex]
			i++
		}
	}
	scs := NewStringColumnsStats(stringArray)
	scs.CalculateAllStats()

	fmt.Fprintf(w, "  Unique values: %d\n", len(scs.valueCounts))
	fmt.Fprintf(w, "  Max length: %d\n", scs.maxLength)
	numFrequent := 5
	if numFrequent > len(scs.valueCounts) {
		numFrequent = len(scs.valueCounts)
	}
	fmt.Fprintf(w, "  %d most frequent values:\n", numFrequent)
	for i := 0; i < numFrequent; i++ {
		fmt.Fprintf(w, "      %s: %d\n", scs.valueCounts[i].value, scs.valueCounts[i].count)
	}
}

type StringColumnStats struct {
	array       []string
	valueCounts []StringValueCount
	maxLength   int
}

func NewStringColumnsStats(stringArray []string) *StringColumnStats {
	scs := new(StringColumnStats)
	scs.array = stringArray
	return scs
}

func (scs *StringColumnStats) CalculateAllStats() {
	scs.CalculateMaxLength()
	scs.CalculateValueCounts()
}

func (scs *StringColumnStats) CalculateMaxLength() {
	scs.maxLength = -1
	for _, elem := range scs.array {
		if len(elem) > scs.maxLength {
			scs.maxLength = len(elem)
		}
	}
}

func (scs *StringColumnStats) CalculateValueCounts() {
	valueCountsMap := make(map[string]int)
	for _, stringVal := range scs.array {
		count, ok := valueCountsMap[stringVal]
		if ok {
			valueCountsMap[stringVal] = count + 1
		} else {
			valueCountsMap[stringVal] = 1
		}
	}
	scs.valueCounts = make([]StringValueCount, len(valueCountsMap))
	i := 0
	for value, count := range valueCountsMap {
		scs.valueCounts[i] = StringValueCount{value, count}
		i++
	}
	sort.Sort(sort.Reverse(StringValueCountByCount(scs.valueCounts)))
}

type StringValueCount struct {
	value string
	count int
}

type StringValueCountByCount []StringValueCount

func (a StringValueCountByCount) Len() int           { return len(a) }
func (a StringValueCountByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StringValueCountByCount) Less(i, j int) bool { return a[i].count < a[j].count }
