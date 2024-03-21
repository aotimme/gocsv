# GoCSV

Command line CSV processing tool inspired by [csvkit](https://csvkit.readthedocs.io). But faster and less memory intensive. And written in Go.

To install on Apple OS X, open a Terminal window and run

```shell
/bin/bash <(curl -s https://raw.githubusercontent.com/aotimme/gocsv/master/scripts/install-latest-darwin-amd64.sh)
```

For other platforms, see the [Installation](#installation) section.

### Table of Contents

- [Introduction](#introduction)
- [Subcommands](#subcommands)
- [Specifying Columns](#specifying-columns)
- [Regular Expression Syntax](#regular-expression-syntax)
- [Pipelining](#pipelining)
- [Changing the Default Delimiter](#changing-the-default-delimiter)
- [Inference](#inference)
- [Examples](#examples)
- [Debugging](#debugging)
- [Installation](#installation)

## Introduction

The tool is built for [pipelining](#pipelining), so most commands accept a CSV from standard input and output to standard output.

Subcommands:

- [add](#add) (aliases: `template`, `tmpl`) - Add a column to a CSV.
- [autoincrement](#autoincrement) (alias: `autoinc`) - Add a column of incrementing integers to a CSV.
- [behead](#behead) - Remove header row(s) from a CSV.
- [cap](#cap) - Add a header row to a CSV.
- [clean](#clean) - Clean a CSV of common formatting issues.
- [delimiter](#delimiter) (alias: `delim`) - Change the delimiter being used for a CSV.
- [describe](#describe) - Get basic information about a CSV.
- [dimensions](#dimensions) (alias: `dims`) - Get the dimensions of a CSV.
- [filter](#filter) - Extract rows whose column match some criterion.
- [head](#head) - Extract the first _N_ rows from a CSV.
- [headers](#headers) - View the headers from a CSV.
- [join](#join) - Join two CSVs based on equality of elements in a column.
- [ncol](#ncol) - Get the number of columns in a CSV.
- [nrow](#nrow) - Get the number of rows in a CSV.
- [rename](#rename) - Rename the headers of a CSV.
- [replace](#replace) - Replace values in cells by regular expression.
- [sample](#sample) - Sample rows.
- [select](#select) - Extract specified columns.
- [sort](#sort) - Sort a CSV based on one or more columns.
- [split](#split) - Split a CSV into multiple files.
- [sql](#sql) - Run SQL queries on CSVs.
- [stack](#stack) - Stack multiple CSVs into one CSV.
- [stats](#stats) - Get some basic statistics on a CSV.
- [tail](#tail) - Extract the last _N_ rows from a CSV.
- [transpose](#transpose) - Transpose a CSV.
- [tsv](#tsv) - Transform a CSV into a TSV.
- [unique](#unique) (alias: `uniq`) - Extract unique rows based upon certain columns.
- [view](#view) - Display a CSV in a pretty tabular format.
- [xlsx](#xlsx) - Convert sheets of a XLSX file to CSV.
- [zip](#zip) - Zip multiple CSVs into one CSV.

To view the usage of `gocsv` at the command line, use the `help` subcommand (i.e. `gocsv help`). This will also print out the version of the `gocsv` binary as well as the hash of the git commit of this repository on which the binary was built. To view only the version and git hash, use the `version` subcommand (i.e. `gocsv version`).

## Subcommands

### add

_Aliases:_ `template`, `tmpl`

Add a column to a CSV.

Usage:

```shell
gocsv add [--prepend] [--name NAME] [--template TEMPLATE] FILE
```

Arguments:

- `--prepend` (optional) Prepend the new column rather than the default append.
- `--name` (shorthand `-n`, optional) Specify a name for the new column. Defaults to the empty string.
- `--template` (shorthand `-t`, optional) Template for column.

Note that the `--template` argument for this subcommand is a string providing a template for the new column. Templates are parsed using the [html/template](https://golang.org/pkg/html/template/) package provided by Go and can reference any column by the _name_ of the column, along with a special variable `index` that represents the row number (starting at `1`).

For example, if your CSV has a column named `Name`, you can do

```shell
gocsv add -t "Hello, {{.Name}}! You are number {{.index}} in line."
```

For multi-word columns there is a slightly different syntax. Say you have a column called `Full Name`. Then the following template would work:

```shell
gocsv add -t 'Hello {{index . "Full Name"}}! You are number {{.index} in line.'
```

GoCSV has been loaded with utility functions from [Sprig](https://github.com/Masterminds/sprig). This will help you to perform wide range of text manipulation on your template on top of built-in Go template functionalities.

Here is an example of how to add a new column of extracting hashtags using [RegEx (RE2)](https://github.com/google/re2/wiki/Syntax), sort, remove duplicates and then join with comma seperated from an existing column:

```shell
gocsv add -t '{{ regexFindAll "#[\\w\\-]+" .Comments -1 | sortAlpha | uniq | join ", " }}' --name 'Hashtags'
```

For further reference on the options available for text manipulation see [Sprig documentation](http://masterminds.github.io/sprig/).

> **TIP:** You will have to take note that Regular Expressions need to be escaped.

### autoincrement

_Alias:_ `autoinc`

Append (or prepend) a column of incrementing integers to each row. This can be helpful to be able to map back to the original row after a number of transformations.

Usage:

```shell
gocsv autoincrement [--prepend] [--name NAME] [--seed SEED] FILE
```

Arguments:

- `--prepend` (optional) Prepend the new column rather than the default append.
- `--name` (optional) Specify a name for the autoincrementing column. Defaults to `ID`.
- `--seed` (optional) Specify the integer to begin incrementing from. Defaults to `1`.

### behead

Remove the header from a CSV

Usage:

```shell
gocsv behead [-n N] FILE
```

Arguments:

- `-n` (optional) Number of header rows to remove. Defaults to 1.

### cap

Add a header row to a CSV

Usage:

```shell
gocsv cap --names NAMES [--truncate-names] [--default-name DEFAULT_NAME] FILE
```

Arguments:

- `--names` [optional] A comma-separated list of names to add as headers to each column.
- `--truncate-names` [optional] If there are fewer columns than the number of names provided by `--names`, drop the extra column names.
- `--default-name` [optional] Fill out the balance of columns not specified by `--names` by using this as the base of the default names for the additional columns. If the default name is `Extra Col`, then the first additional column will be named "Extra Col", the second "Extra Col 1", the third "Extra Col 2", etc.

The subcommand will error if:

- Both `--names` and `--default-name` are not specified, or
- The number of names provided by `--names` is greater than the number of columns and `--truncate-names` is not specified, or
- the number of names provided by `--names` is less than the number of columns and `--default-name` is not specified.

For example:

```shell
echo Jamie,52,Purple | gocsv cap --names 'Name,Age,Favorite color'
Name,Age,Favorite color
Jamie,52,Purple
```

```shell
echo Jamie,52,Purple | gocsv cap --names 'Name' --default-name 'Col'
Name,Col,Col 1
Jamie,52,Purple
```

```shell
echo Jamie,52,Purple | gocsv cap --default-name 'Col'
Col,Col 1,Col 2
Jamie,52,Purple
```

### clean

Clean a CSV of common formatting issues. Currently this consists of making sure all rows are the same length (padding short rows and trimming long ones) and removing empty rows at the end.

Note that this subcommand, along with other subcommands, will include a newline at the end of the last line of the outputted CSV. This is because `gocsv` assumes that every row in a CSV (or other-delimited text file) will end in a new line.

Usage:

```shell
gocsv clean [--verbose] [--no-trim] [--strip-bom] [--excel] [--numbers] FILE
```

Arguments:

- `--verbose` (optional) Print out to stderr when cleaning the CSV.
- `--no-trim` (optional) Do not remove trailing rows that are empty.
- `--add-bom` (optional) Ensure that a BOM that exists at the beginning of the CSV.
- `--strip-bom` (optional) Remove any BOM that exists at the beginning of the CSV.
- `--excel` (optional) Clean the CSV for issues that will cause problems with Excel. See [Excel specifications and limitations](https://support.office.com/en-us/article/Excel-specifications-and-limits-16c69c74-3d6a-4aaf-ba35-e6eb276e8eaa).
  - Truncate any cells that exceed the maximum character limit of 32767.
- `--numbers` (optional) Clean the CSV for issues that will cause problems with Numbers.
  - Truncate the number of rows in the CSV at 65535, the maximum amount of rows that Numbers displays.

Note that only one of `--add-bom` or `--strip-bom` can be specified.

### delimiter

_Alias_: `delim`

Change the delimiter being used for a CSV.

Usage:

```bash
gocsv delim [--input INPUT_DELIMITER] [--output OUTPUT_DELIMITER] FILE
```

Arguments:

- `--input` (shorthand `-i`, optional) The delimiter used in the input. Defaults to `,`.
- `--output` (shorthand `-o`, optional) The delimiter used in the output. Defaults to `,`.

### describe

Get basic information about a CSV. This will output the number of rows and columns in the CSV, the column headers in the CSV, and the [inferred](#inference) type of each column.

Usage

```shell
gocsv describe FILE
```

### dimensions

_Alias:_ `dims`

Get the dimensions of a CSV.

Usage

```shell
gocsv dimensions [--csv] FILE
```

Arguments:

- `--csv` (optional) Output the results as a CSV.

### filter

Filter a CSV by rows whose columns match some criterion.

Usage:

```shell
gocsv filter [--columns COLUMNS] [--equals STR] [--regex REGEX] [--gt N] [--gte N] [--lt N] [--lte N] [--exclude] FILE
```

Arguments:

- `--columns` (optional, shorthand `-c`) A comma-separated list of the columns to filter against. If no columns are specified, then filter checks every column on a row. If a row matches on any of the columns, the row is considered a match. See [Specifying Columns](#specifying-columns) for more details.
- `--equals` (optional, shorthand `-eq`) String to match against.
- `--regex` (optional) Regular expression to use to match against. See [Regular Expression Syntax](#regular-expression-syntax) for the syntax.
- `--case-insensitive` (optional, shorthand `-i`) When using the `--regex` flag, use this flag to specify a case insensitive match rather than the default case sensitive match.
- `--gt` , `--gte`, `--lt`, `--lte` (optional) Compare against the [inferred types int, float, datetime](#inference); booleans and strings cannot be compared.
- `--exclude` (optional) Exclude rows that match. Default is to include.

Note that one of `--regex`, `--equals` (`-eq`), `--gt` , `--gte`, `--lt`, or `--lte` must be specified.

### head

Extract the first _N_ rows from a CSV.

Usage:

```shell
gocsv head [-n N] FILE
```

Arguments:

- `-n` (optional) The number of rows to extract. If `N` is an integer, it will extract the first _N_ rows. If `N` is prepended with `+`, it will extract all except the last _N_ rows.

### headers

View the headers of a CSV along with the index of each header.

Usage:

```bash
gocsv headers [--csv] FILE
```

Arguments:

- `--csv` (optional) Output the results as a CSV.

### join

Join two CSVs using an inner (default), left, right, or outer join.

Usage:

```shell
gocsv join --columns COLUMNS [--left] [--right] [--outer] LEFT_FILE RIGHT_FILE
```

Arguments:

- `--columns` (shorthand `-c`) A comma-separated list (in order) of the columns to use for joining. You must specify either 1 or 2 columns. When 1 is specified, it will join the CSVs using that column in both the left and right CSV. When 2 are specified, it will join using the first column on the left CSV and the second column on the right CSV. See [Specifying Columns](#specifying-columns) for more details.
- `--left` (optional) Perform a left join (i.e. left outer join).
- `--right` (optional) Perform a right join (i.e. right outer join).
- `--outer` (optional) Perform an outer join (i.e. full outer join).

Note that by default it will perform an inner join. It will exit if you specify multiple types of join.

### ncol

Get the number of columns in a CSV.

Usage:

```shell
gocsv ncol FILE
```

### nrow

Get the number of rows in a CSV.

Usage:

```shell
gocsv nrow FILE
```

### rename

Rename the headers of a CSV.

Usage:

```shell
gocsv rename --columns COLUMNS --names NAMES FILE
```

Arguments:

- `--columns` (shorthand `-c`) A comma-separated list of the columns to rename. See [Specifying Columns](#specifying-columns) for more details.
- `--names` A comma-separated list of names to change each column to. This must be the same length as and match the order of the `columns` argument.

### replace

Replace values in cells by regular expression.

Usage:

```shell
gocsv replace [--columns COLUMNS] --regex REGEX --repl REPLACEMENT FILE
```

Arguments:

- `--columns` (optional, shorthand `-c`) A comma-separated list of the columns to run replacements on. If no columns are specified, then replace runs the replacement operation on cells in every column. See [Specifying Columns](#specifying-columns) for more details.
- `--regex` Regular expression to use to match against for replacement. See [Regular Expression Syntax](#regular-expression-syntax) for the syntax.
- `--case-insensitive` (optional, shorthand `-i`) Use this flag to specify a case insensitive match for replacement rather than the default case sensitive match.
- `--repl` String to use for replacement.

Note that if you have a capture group in the `--regex` argument you can reference that in the replacement argument using `"\$1"` for the first capture group, `"\$2"` for the second capture group, etc.

### sample

Sample rows from a CSV

Usage

```shell
gocsv sample -n NUM_ROWS [--replace] [--seed SEED] FILE
```

Arguments:

- `-n` The number of rows to sample.
- `--replace` (optional) Whether to sample with replacement. Defaults to `false`.
- `--seed` (optional) Integer seed to use for generating pseudorandom numbers for sampling.

### select

Select (or exclude) columns from a CSV

Usage:

```shell
gocsv select --columns COLUMNS [--exclude] FILE
```

Arguments:

- `--columns` (shorthand `-c`) A comma-separated list (in order) of the columns to select. If you want to select a column multiple times, you can! See [Specifying Columns](#specifying-columns) for more details.
- `--exclude` (optional) Exclude the specified columns (default is to include).

### sort

Sort a CSV by multiple columns, with or without type [inference](#inference). The currently supported types are float, int, datetime, and string.

Usage:

```shell
gocsv sort --columns COLUMNS [--stable] [--reverse] [--no-inference] FILE
```

Arguments:

- `--columns` (shorthand `-c`) A comma-separated list (in order) of the columns to sort against. See [Specifying Columns](#specifying-columns) for more details.
- `--stable` (optional) Keep the original order of equal rows while sorting.
- `--reverse` (optional) Reverse the order of sorting. By default the sort order is ascending.
- `--no-inference` (optional) Skip type inference when sorting.

When `--stable` and `--reverse` are both specified, the original order of equal rows is preserved (and not reversed).

### split

Split a CSV into multiple files.

Usage:

```shell
gocsv split --max-rows N [--filename-base FILENAME] FILE
```

Arguments:

- `--max-rows` Maximum number of rows per final CSV.
- `--filename-base` (optional) Prefix of the resulting files. The file outputs will be appended with `"-1.csv"`,`"-2.csv"`, etc. If not specified, the base filename will be the same as the base of the input filename, unless the input is specified by standard input. If so, then the base filename will be `out`.

### sql

Run SQL queries on CSVs.

Usage:

```shell
gocsv sql --query QUERY FILE [FILES]
```

Arguments:

- `--query` (shorthand `-q`) The SQL query to run.

When passing in files, you may read from standard input by specifying the filename `-`.

Table names are derived from the CSV filenames by taking the base filename without the file extension. For example, `test-files/stats.csv` is referenced as a table with the name `stats`. The table from standard input `-` should be referenced as the table `stdin`.

This subcommand uses SQLite3 under the hood. It attempts to infer column types for defining the SQL tables, but all the rules of dynamic typing and type affinity in SQLite3 still pertain.

See [Datatypes In SQLite Version 3](https://www.sqlite.org/datatype3.html) for more information.

Also note that this subcommand makes no attempts to prevent SQL injection (either via the input CSVs or via the query).

### stack

Stack multiple CSVs to create a larger CSV. Optionally include an indication of which file a row came from in the final CSV.

Usage:

```shell
gocsv stack [--filenames] [--groups GROUPS] [--group-name GROUP_NAME] FILE [FILES]
```

Arguments:

- `--filenames` (optional) Use the names of each file as the group variable. By default the column will be named "File".
- `--groups` (optional) Comma-separated list to use as the names of the groups for each row. There must be as many groups as there are files. By default the column will be named "Group".
- `--group-name` (optional) Name of the grouping column in the final CSV.

Note that `--groups` and `--filenames` are mutually exclusive.

Specifying a file by name `-` will read a CSV from standard input.

### stats

Get some basic statistics on a CSV, uses [inference](#inference).

Usage:

```shell
gocsv stats FILE
```

### tail

Extract the last _N_ rows from a CSV.

Usage:

```shell
gocsv tail [-n N] FILE
```

Arguments:

- `-n` (optional) The number of rows to extract. If `N` is an integer, it will extract the last _N_ rows. If `N` is prepended with `+`, it will extract all except the first _N_ rows.

### transpose

Transpose a CSV.

Usage:

```shell
gocsv tranpose FILE
```

### tsv

Transform a CSV into a TSV. It is shortand for `gocsv delim -o "\t" FILE`. This can very useful if you want to pipe the result to `pbcopy` (OS X) in order to paste it into a spreadsheet tool.

Usage:

```shell
gocsv tsv FILE
```

### unique

_Alias:_ `uniq`

Extract unique rows based upon certain columns.

Usage:

```shell
gocsv unique [--columns COLUMNS] [--sorted] [--count] FILE
```

Arguments

- `--columns` (optional, shorthand `-c`) A comma-separated list (in order) of the columns to use to define uniqueness. If no columns are specified, it will perform uniqueness across the entire row. See [Specifying Columns](#specifying-columns) for more details.
- `--sorted` (optional) Specify whether the input is sorted. If the input is sorted, the unique subcommand will run more efficiently.
- `--count` (optional) Append a column with the header "Count" to keep track of how many times that unique row occurred in the input.

### view

Display a CSV in a pretty tabular format.

Usage:

```shell
gocsv view [-n N] [--max-width N] FILE
```

Arguments:

- `-n` (optional) Display only the first _N_ rows of the CSV.
- `--max-width` (optional, shorthand `-w`) The maximum width of each cell for display. If a cell exceeds the maximum width, it will be truncated in the display.
- `--max-lines` (optional, shorthand `-l`) The maximum number of lines to display per cell.

If the length of a cell exceeds `--max-width` it will be truncated with an ellipsis.

### xlsx

Convert sheets of a XLSX file to CSV.

Usage:

```shell
gocsv xlsx [--list-sheets] [--dirname DIRNAME] [--sheet SHEET] FILE
```

Arguments:

- `--list-sheets` (optional) List the sheets in the XLSX file.
- `--sheet` (optional) Specify the sheet (by index or name) of the sheet to convert.
- `--dirname` (optional) Name of directory to output CSV conversions of sheets from `FILE`. If this is not specified, the command will output the CSV files to a directory with the same name as `FILE` (without the `.xlsx` extension).

By default the `xlsx` subcommand will convert all the sheets in `FILE` to CSVs to a directory with the same name as `FILE`.

### zip

Zip multiple CSVs into one CSV.

Usage:

```shell
gocsv zip FILE [FILES]
```

Specifying a file by name `-` will read a CSV from standard input.

## Specifying Columns

When specifying columns on the command line (i.e. with the `--columns` or `-c` flags), you can specify either the indices or the names of the columns. The tool will always try to interpret columns first by index and then by name.

#### Specifying Columns by Index

The tool uses 1-based indexing (as in the output of the [headers](#headers) subcommand).

The tool also allows for specification of ranges with indices (e.g. `2-4`) including reverse ranges (e.g. `4-2`). It also allows for open-ended ranges on indexes (e.g. `2-` or `-4`). In the former case (a-) it will include all columns from `a` on. In the latter case (`-b`) it will include all columns before `b` and `b` itself.

### Specifying Columns by Name

When specifying the name of a column, it will match all columns that are exact case-sensitive matches.

When referencing a column name that has whitespace, either escape the whitespace with `\` or use quotes (`"`) around the column name.

For example, if you have a column named `Hello World`,

```shell
gocsv select -c "Hello World" test.csv
```

or

```shell
gocsv select -c Hello\ World test.csv
```

When referencing multiple columns, specify column names as a comma-delimited list with no spaces between the columns. If any of the column names have whitespace, enclose the entire list in a single set of quotes.

```shell
gocsv select -c "Hello World,Foo Bar" test.csv
```

## Regular Expression Syntax

A few of the subcommands allow the ability to pass in regular expressions via a `--regex` flag (e.g. [filter](#filter) and [replace](#replace)).

Because the regular expressions passed in to the `--regex` flag are parsed by the underlying [regexp](https://golang.org/pkg/regexp/) Go package, see the [regexp/syntax](https://golang.org/pkg/regexp/syntax/) documentation for more details on the syntax. It is based on the syntax accepted by [RE2](https://github.com/google/re2/wiki/Syntax).

## Pipelining

Because all of the subcommands support receiving a CSV from standard input, you can easily pipeline:

```shell
cat test-files/left-table.csv \
  | gocsv join --left --columns LID,RID test-files/right-table.csv \
  | gocsv filter --columns XYZ --regex "[ev]e-\d$" \
  | gocsv select --columns LID,XYZ \
  | gocsv sort --columns LID,XYZ
```

### Pipelining Support

| Subcommand    |    Input            |  Output  |
| ------------- | :-----------------: | :------: |
| add           |  &#x2714;           | &#x2714; |
| autoincrement |  &#x2714;           | &#x2714; |
| behead        |  &#x2714;           | &#x2714; |
| clean         |  &#x2714;           | &#x2714; |
| delimiter     |  &#x2714;           | &#x2714; |
| describe      |  &#x2714;           |   N/A    |
| dimensions    |  &#x2714;           | &#x2714;<sup>*</sup> |
| filter        |  &#x2714;           | &#x2714; |
| head          |  &#x2714;           | &#x2714; |
| headers       |  &#x2714;           | &#x2714;<sup>*</sup> |
| join          |  &#x2714;           | &#x2714; |
| ncol          |  &#x2714;           |   N/A    |
| nrow          |  &#x2714;           |   N/A    |
| rename        |  &#x2714;           | &#x2714; |
| replace       |  &#x2714;           | &#x2714; |
| sample        |  &#x2714;           | &#x2714; |
| select        |  &#x2714;           | &#x2714; |
| sort          |  &#x2714;           | &#x2714; |
| split         |  &#x2714;           |   N/A    |
| sql           |  &#x2714;<sup>&#x2020;</sup>   | &#x2714; |
| stack         |  &#x2714;<sup>&#x2020;</sup>   | &#x2714; |
| stats         |  &#x2714;           |   N/A    |
| tail          |  &#x2714;           | &#x2714; |
| transpose     |  &#x2714;           | &#x2714; |
| tsv           |  &#x2714;           | &#x2714; |
| unique        |  &#x2714;           | &#x2714; |
| view          |  &#x2714;           |   N/A    |
| xlsx          |     N/A             | &#x2021; |

\* `dimensions` and `headers` write to CSV format when using the `--csv` argument.

&#x2020; `stack` and `sql` read from standard input when specifying the filename as `-`.

&#x2021; `xlsx` sends output to standard out when using the `--sheet` flag.

## Changing the Default Delimiter

While `gocsv` generally assumes standard CSVs (per [RFC 4180](https://tools.ietf.org/html/rfc4180)), you can specify a default delimiter other than `,` using the `GOCSV_DELIMITER` environment variable. The delimiter _must_ evaluate to exactly 1 ["rune"](https://go.dev/doc/go1#rune). If it does not, `gocsv` will error.

For example, to use semicolon-delimited files:

```shell
export GOCSV_DELIMITER=";"
gocsv select -c 1 semicolon-delimited.scsv
```

Or, to use tab-delimited files (TSVs):

```shell
export GOCSV_DELIMITER="\t"
gocsv select -c 1 tab-delimited.tsv
```

Or, for more exotic delimiters you can use hexadecimal or unicode (e.g. `\x01` or `\u0001` for the SOH delimiter):

```shell
export GOCSV_DELIMITER="\x01"
gocsv select -c 1 soh-delimited.tsv
```

## Inference

GoCSV can infer values (in increasing order of precedence) for int, float, boolean, and datetime.

- Ints and floats are any number that can be parsed by strconv's ParseInt and ParseFloat functions.

  A column with ints and floats will be inferred as floats by the [describe](#describe) and [stats](#stats) subcommands—floats have greater precedence than ints.

- Booleans are any of "T", "True", "F", or "False" (case invariant).

- Datetimes are values that match any of [pkg time's](https://pkg.go.dev/time@go1.14#pkg-constants) predefined layouts, ANSIC, UnixDate, RubyDate, RFC822, RFC822Z, RFC850, RFC1123, RFC1123Z, and RFC3339 , as well as the date-only layouts, "2006-01-02", "2006-1-2", "1/2/2006", "01/02/2006".

  A datetime column can have values with different layouts that match any of the previously mentioned layouts.

  A custom layout can be supplied with the environment variable `GOCSV_TIMELAYOUT` (like in the GOCSV_DELIMITER examples above). Using the env var puts that layout of the "front of the line" of the precviously mentioned layouts to be tried, so the other layouts can still match values in any of the input columns.

## Examples

##### Copy Values

```shell
gocsv tsv test-files/left-table.csv | pbcopy
```

##### Reorder Columns

```shell
gocsv select --columns 2,1 test-files/left-table.csv
```

##### Duplicate Columns

```shell
gocsv select --columns 1,1,2,2 test-files/left-table.csv
```

##### VLOOKUP aka Join

```shell
gocsv join --left --columns LID,RID test-files/left-table.csv test-files/right-table.csv
```

##### Distinct Column Values

```shell
gocsv select --columns LID test-files/left-table.csv | gocsv behead | sort | uniq
```

##### Count of Distinct Column Values

```shell
gocsv select --columns LID test-files/left-table.csv | gocsv behead | sort | uniq -c | sort -nr
```

##### Extract Rows Matching Regular Expression

```shell
gocsv filter --columns ABC --regex "-1$" test-files/left-table.csv
```

##### Extract Rows with Blank Values in Column

```shell
gocsv filter --columns Stringer --regex "^$" test-files/stats.csv
```

If you also want to match on cells that have only whitespace, you can use a regular expression like `"^s*$"`.

##### Replace Content in Cells By Regular Expression

```shell
gocsv replace --columns ABC --regex "^(.*)-(\d)$" -i --repl "\$2-\$1" test-files/left-table.csv
```

##### Sort by Multiple Columns

```shell
gocsv sort --columns LID,ABC --reverse test-files/left-table.csv
```

##### Combine Multiple CSVs

```shell
gocsv stack --groups "Primer Archivo,Segundo Archivo,Tercer Archivo" --group-name "Orden de Archivo" test-files/stack-1.csv test-files/stack-2.csv test-files/stack-3.csv
```

To do the same via pipelining through standard input,

```shell
cat test-files/stack-1.csv | gocsv stack --groups "Primer Archivo,Segundo Archivo,Tercer Archivo" --group-name "Orden de Archivo" - test-files/stack-2.csv test-files/stack-3.csv
```

##### Create a Column from a Template

```shell
cat test-files/stats.csv | gocsv add -t "Row {{.index}}: {{if eq .Boolean \"T\"}}{{.Floater}}{{else}}{{.Integer}}{{end}}" -name "Integer or Floater"
```

## Debugging

To enable debugging mode when running a `gocsv` command, specify the `--debug` command line argument to any subcommand (other than `gocsv help` and `gocsv version`). Any errors will then also print out a stack trace.

## Installation

For the latest pre-built binaries, cross-compiled using [xgo](https://github.com/crazy-max/xgo), see the [Latest Release](https://github.com/aotimme/gocsv/releases/latest) page.

### Apple OS X

#### Simple Version

Open a Terminal window and paste the following command:

```shell
/bin/bash <(curl -s https://raw.githubusercontent.com/aotimme/gocsv/master/scripts/install-latest-darwin-amd64.sh)
```

This will install `gocsv` at `/usr/local/bin/gocsv`.

#### Detailed Version

To install the pre-built binary for Apple OS X, download the `gocsv-darwin-amd64.zip` file. It should download into your `~/Downloads` directory. To install it, open a Terminal window and do the following:

```shell
cd ~/Downloads
unzip gocsv-darwin-amd64.zip
mv gocsv-darwin-amd64/gocsv /usr/local/bin
rmdir gocsv-darwin-amd64
```

To verify that it has installed, open a _new_ Terminal window and run

```shell
gocsv help
```

You should see the `gocsv` help message.

### Linux

Installing the pre-built binary for Linux is very similar to installing the binary for Apple OS X. First, download `gocsv-linux-amd64.zip`. Assuming this downloads to your `~/Downloads` directory, open a Terminal window and run the following commands:

```shell
cd ~/Downloads
unzip gocsv-linux-amd64.zip
mv gocsv-linux-amd64/gocsv /usr/local/bin
rmdir gocsv-linux-amd64
```

To verify that it has installed, open a _new_ Terminal window and run

```shell
gocsv help
```

You should see the `gocsv` help message.

### Windows

Download `gocsv-windows-amd64.zip`. Unzip the file and you should see a file `gocsv.exe`. Put that executable in the appropriate location and it should work.
