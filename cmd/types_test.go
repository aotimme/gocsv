package cmd

import (
	"fmt"
	"testing"
)

type SliceStringIterator struct {
	data []string
	pos  int
}

func NewSliceStringIterator(data []string) *SliceStringIterator {
	return &SliceStringIterator{
		data: data,
	}
}
func (it *SliceStringIterator) Next() (string, bool) {
	if it.pos >= len(it.data) {
		return "", false
	}
	val := it.data[it.pos]
	it.pos++
	return val, true
}

func TestInferTypeFromStringIterator(t *testing.T) {
	testCases := []struct {
		values     []string
		columnType ColumnType
	}{
		{[]string{"", ""}, NULL_TYPE},
		// all types are nullable
		{[]string{"", "1"}, INT_TYPE},
		{[]string{"", "1.0"}, FLOAT_TYPE},
		{[]string{"", "true"}, BOOLEAN_TYPE},
		{[]string{"", "2023-01-02"}, DATE_TYPE},
		{[]string{"", "2023-01-02 12:00:00"}, DATETIME_TYPE},
		{[]string{"", "hello"}, STRING_TYPE},
		// all types are the same
		{[]string{"1", "2", "3"}, INT_TYPE},
		{[]string{"1.0", "2.0", "3.0"}, FLOAT_TYPE},
		{[]string{"true", "false"}, BOOLEAN_TYPE},
		{[]string{"2023-01-01", "2023-01-02"}, DATE_TYPE},
		{[]string{"2023-01-01 12:00:00", "2023-01-02 12:00:00"}, DATETIME_TYPE},
		{[]string{"hello", "world"}, STRING_TYPE},
		// date + datetime => datetime
		{[]string{"2023-01-01", "2023-01-02 12:00:00"}, DATETIME_TYPE},
		// int + float => float
		{[]string{"1", "2", "3.0"}, FLOAT_TYPE},
		// int + string => string
		{[]string{"1", "2", "hello"}, STRING_TYPE},
		// int + boolean => string
		{[]string{"1", "2", "true"}, STRING_TYPE},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			siter := NewSliceStringIterator(tt.values)
			columnType := InferTypeFromStringIterator(siter)
			if columnType != tt.columnType {
				t.Errorf("got %s; want %s", ColumnTypeToString(columnType), ColumnTypeToString(tt.columnType))
			}
		})
	}
}

func TestIsDateType(t *testing.T) {
	testCases := []struct {
		value  string
		isDate bool
	}{
		{"2023-1-2", true},
		{"1/2/2023", true},
		{"2023-01-01", true},
		{"2023-01-01 12:00:00", false},
		{"2023-01-01T12:00:00Z", false},
		{"2023-01-01 12:00:00+00:00", false},
		{"2023-01-01 12:00:00-07:00", false},
		{"2023-01-01 12:00:00.123456", false},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			isDate := IsDateType(tt.value)
			if isDate != tt.isDate {
				t.Errorf("got %t; want %t", isDate, tt.isDate)
			}
		})
	}
}

func TestIsDatetimeType(t *testing.T) {
	testCases := []struct {
		value      string
		isDatetime bool
	}{
		{"2023-01-01 12:00:00", true},
		{"2023-01-01T12:00:00Z", true},
		{"2023-01-01 12:00:00.123456", true},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			isDatetime := IsDatetimeType(tt.value)
			if isDatetime != tt.isDatetime {
				t.Errorf("got %t; want %t", isDatetime, tt.isDatetime)
			}
		})
	}
}
