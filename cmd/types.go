package cmd

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type ColumnType int

// NOTE: Order matters here. Ordered by strictness descending
const (
	NULL_TYPE ColumnType = iota
	INT_TYPE
	FLOAT_TYPE
	BOOLEAN_TYPE
	DATETIME_TYPE
	DATE_TYPE
	STRING_TYPE
)

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

func InferTypeWithHint(elem string, hint ColumnType) ColumnType {
	if IsNullType(elem) {
		return NULL_TYPE
	}
	if INT_TYPE >= hint && IsIntType(elem) {
		return INT_TYPE
	}
	if FLOAT_TYPE >= hint && IsFloatType(elem) {
		return FLOAT_TYPE
	}
	if BOOLEAN_TYPE >= hint && IsBooleanType(elem) {
		return BOOLEAN_TYPE
	}
	if DATETIME_TYPE >= hint && IsDatetimeType(elem) {
		return DATETIME_TYPE
	}
	if DATE_TYPE >= hint && IsDateType(elem) {
		return DATE_TYPE
	}
	return STRING_TYPE
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
	strLower := strings.ToLower(elem)
	return strLower == "t" || strLower == "true" || strLower == "f" || strLower == "false"
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
	patterns := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
	}
	for _, pattern := range patterns {
		t, err := time.Parse(pattern, elem)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("Invalid Date string")
}

func ParseDateOrPanic(elem string) time.Time {
	t, err := ParseDate(elem)
	if err != nil {
		ExitWithError(err)
	}
	return t
}

func ParseDate(elem string) (time.Time, error) {
	patterns := []string{
		"2006-01-02",
		"2006-1-2",
		"1/2/2006",
		"01/02/2006",
	}
	for _, pattern := range patterns {
		t, err := time.Parse(pattern, elem)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("Invalid Date string")
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
