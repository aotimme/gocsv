module github.com/aotimme/gocsv

replace (
	github.com/aotimme/gocsv/cmd => ./cmd
	github.com/aotimme/gocsv/csv => ./csv
)

require (
	github.com/aotimme/gocsv/cmd v0.0.0
	github.com/aotimme/gocsv/csv v0.0.0
	github.com/tealeg/xlsx v1.0.5
)
