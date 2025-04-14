package cmd

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type ColumnType int

const (
	NULL_TYPE ColumnType = iota
	INT_TYPE
	FLOAT_TYPE
	BOOLEAN_TYPE
	DATE_TYPE
	DATETIME_TYPE
	STRING_TYPE
)

var datetimeFormats = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.DateTime,
}

var dateFormats = []string{
	"2006-01-02",
	"2006-1-2",
	"1/2/2006",
	"01/02/2006",
}

func getCommonType(a, b ColumnType) ColumnType {
	// swap a and b so a <= b (because this function is symmetrical):
	if a > b {
		a, b = b, a
	}
	if a == b {
		return a
	}
	// At this point, b > a, so we don't need to consider any
	// case where b <= a below.
	switch a {
	case NULL_TYPE:
		return b
	case INT_TYPE:
		if b == FLOAT_TYPE {
			return FLOAT_TYPE
		}
	case FLOAT_TYPE:
	case BOOLEAN_TYPE:
	case DATE_TYPE:
		if b == DATETIME_TYPE {
			return DATETIME_TYPE
		}
	case DATETIME_TYPE:
	case STRING_TYPE:
	}
	return STRING_TYPE
}

func ColumnTypeToString(columnType ColumnType) string {
	if columnType == NULL_TYPE {
		return "null"
	} else if columnType == INT_TYPE {
		return "int"
	} else if columnType == FLOAT_TYPE {
		return "float"
	} else if columnType == BOOLEAN_TYPE {
		return "boolean"
	} else if columnType == DATETIME_TYPE {
		return "datetime"
	} else if columnType == DATE_TYPE {
		return "date"
	} else if columnType == STRING_TYPE {
		return "string"
	} else {
		return ""
	}
}

func ColumnTypeToSqliteType(columnType ColumnType) string {
	if columnType == NULL_TYPE {
		return "TEXT"
	} else if columnType == INT_TYPE {
		return "INTEGER"
	} else if columnType == FLOAT_TYPE {
		return "REAL"
	} else if columnType == BOOLEAN_TYPE {
		return "TEXT"
	} else if columnType == DATETIME_TYPE {
		return "TEXT"
	} else if columnType == DATE_TYPE {
		return "TEXT"
	} else if columnType == STRING_TYPE {
		return "TEXT"
	} else {
		return "TEXT"
	}
}

type StringIterator interface {
	// Next returns the next string and whether more values are available
	Next() (string, bool)
}

func InferTypeFromStringIterator(valuesIter StringIterator) ColumnType {
	curType := NULL_TYPE

	for {
		value, ok := valuesIter.Next()
		if !ok {
			break
		}
		thisType := InferTypeWithRunningType(value, curType)
		if thisType > curType {
			curType = thisType
		}

		// Early termination if we already know it's a string
		if curType == STRING_TYPE {
			return curType
		}
	}

	return curType
}

func GetType(value string) ColumnType {
	if IsNullType(value) {
		return NULL_TYPE
	}
	if IsIntType(value) {
		return INT_TYPE
	}
	if IsFloatType(value) {
		return FLOAT_TYPE
	}
	if IsBooleanType(value) {
		return BOOLEAN_TYPE
	}
	if IsDateType(value) {
		return DATE_TYPE
	}
	if IsDatetimeType(value) {
		return DATETIME_TYPE
	}
	return STRING_TYPE
}

func InferTypeWithRunningType(elem string, runningType ColumnType) ColumnType {
	elemType := GetType(elem)
	return getCommonType(elemType, runningType)
}

func IsNullType(elem string) bool {
	return elem == ""
}

func IsIntType(elem string) bool {
	_, err := strconv.ParseInt(elem, 0, 0)
	return err == nil
}

func IsFloatType(elem string) bool {
	_, err := strconv.ParseFloat(elem, 64)
	return err == nil
}

func IsBooleanType(elem string) bool {
	_, err := ParseBoolean(elem)
	return err == nil
}

func ParseBoolean(elem string) (bool, error) {
	strLower := strings.ToLower(elem)
	if strLower == "t" || strLower == "true" {
		return true, nil
	} else if strLower == "f" || strLower == "false" {
		return false, nil
	}
	return false, errors.New("invalid boolean string")
}

func ParseBooleanOrPanic(elem string) bool {
	b, err := ParseBoolean(elem)
	if err != nil {
		ExitWithError(err)
	}
	return b
}

func IsDatetimeType(elem string) bool {
	_, err := ParseDatetime(elem)
	return err == nil
}

func IsDateType(elem string) bool {
	_, err := ParseDate(elem)
	return err == nil
}

func ParseDatetimeOrPanic(elem string) time.Time {
	t, err := ParseDatetime(elem)
	if err != nil {
		ExitWithError(err)
	}
	return t
}

func ParseDatetime(elem string) (time.Time, error) {
	for _, format := range datetimeFormats {
		t, err := time.Parse(format, elem)
		if err == nil {
			return t, nil
		}
	}
	// Fall back to parsing as Date (Date is a subset of Datetime)
	return ParseDate(elem)
}

func ParseDateOrPanic(elem string) time.Time {
	t, err := ParseDate(elem)
	if err != nil {
		ExitWithError(err)
	}
	return t
}

func ParseDate(elem string) (time.Time, error) {
	for _, format := range dateFormats {
		t, err := time.Parse(format, elem)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("invalid Date string")
}

func ParseFloat64OrPanic(strVal string) float64 {
	floatVal, err := ParseFloat64(strVal)
	if err != nil {
		ExitWithError(err)
	}
	return floatVal
}

func ParseFloat64(strVal string) (float64, error) {
	return strconv.ParseFloat(strVal, 64)
}

func ParseInt64OrPanic(strVal string) int64 {
	intVal, err := ParseInt64(strVal)
	if err != nil {
		ExitWithError(err)
	}
	return intVal
}

func ParseInt64(strVal string) (int64, error) {
	return strconv.ParseInt(strVal, 0, 0)
}
