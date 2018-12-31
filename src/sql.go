package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type SqlSubcommand struct {
	queryString string
}

func (sub *SqlSubcommand) Name() string {
	return "sql"
}
func (sub *SqlSubcommand) Aliases() []string {
	return []string{}
}
func (sub *SqlSubcommand) Description() string {
	return "Run SQL queries on CSVs."
}
func (sub *SqlSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.queryString, "query", "", "SQL query")
	fs.StringVar(&sub.queryString, "q", "", "SQL query (shorthand)")
}

func (sub *SqlSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, -1)
	outputCsv := NewOutputCsvFromInputCsvs(inputCsvs)
	sub.RunSql(inputCsvs, outputCsv)
}

func (sub *SqlSubcommand) RunSql(inputCsvs []*InputCsv, outputCsvWriter OutputCsvWriter) {
	query := sub.queryString

	// 1. Create the SQLite DB
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		ExitWithError(err)
	}
	defer db.Close()

	// 2. Create and populate the tables in the SQL DB
	for _, inputCsv := range inputCsvs {
		PopulateSqlTable(db, inputCsv)
	}
	// 3. Run the query
	rows, err := db.Query(query)
	if err != nil {
		ExitWithError(err)
	}
	defer rows.Close()

	// 4. Write the results
	columns, err := rows.Columns()
	if err != nil {
		ExitWithError(err)
	}
	outputCsvWriter.Write(columns)

	// See: https://stackoverflow.com/a/14500756
	readRow := make([]interface{}, len(columns))
	writeRow := make([]sql.NullString, len(columns))
	for i := range writeRow {
		readRow[i] = &writeRow[i]
	}
	csvRow := make([]string, len(columns))

	for rows.Next() {
		err := rows.Scan(readRow...)
		if err != nil {
			ExitWithError(err)
		}
		for i, elem := range writeRow {
			if elem.Valid {
				csvRow[i] = elem.String
			} else {
				csvRow[i] = ""
			}
		}
		outputCsvWriter.Write(csvRow)
	}
}

func PopulateSqlTable(db *sql.DB, inputCsv *InputCsv) {
	tableName := inputCsv.Name()
	imc := NewInMemoryCsvFromInputCsv(inputCsv)
	allVariables := make([]interface{}, 2*len(imc.header)+1)
	allVariables[0] = tableName
	createStatement := "CREATE TABLE [%s]("
	for i, headerName := range imc.header {
		allVariables[2*i+1] = headerName
		columnType := imc.InferType(i)
		allVariables[2*i+2] = ColumnTypeToSqliteType(columnType)
		if i > 0 {
			createStatement += ", "
		}
		createStatement += "[%s] %s NULL"
	}
	createStatement += ");"
	// Unfortunately using `db.Prepare` with `?` variables wouldn't work
	preparedStatement := fmt.Sprintf(createStatement, allVariables...)
	_, err := db.Exec(preparedStatement)
	if err != nil {
		ExitWithError(err)
	}

	escapedHeaders := make([]string, len(imc.header))
	for i, headerName := range imc.header {
		escapedHeaders[i] = fmt.Sprintf("[%s]", headerName)
	}
	tableColumns := fmt.Sprintf("[%s](%s)", tableName, strings.Join(escapedHeaders, ", "))
	valuesQuestions := make([]string, len(imc.header))
	for i := range valuesQuestions {
		valuesQuestions[i] = "?"
	}
	tableValues := fmt.Sprintf("values(%s)", strings.Join(valuesQuestions, ", "))
	insertStatement := fmt.Sprintf("INSERT INTO %s %s", tableColumns, tableValues)
	preparedInsert, err := db.Prepare(insertStatement)
	if err != nil {
		ExitWithError(err)
	}
	valuesRow := make([]interface{}, len(imc.header))
	for _, row := range imc.rows {
		for i, elem := range row {
			valuesRow[i] = elem
		}
		_, err = preparedInsert.Exec(valuesRow...)
		if err != nil {
			ExitWithError(err)
		}
	}
}
