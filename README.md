# GoCSV

Command line CSV processing tool based on [csvkit](https://csvkit.readthedocs.io). But faster and less memory intensive. And written in Go.

To install on Apple OS X, open a Terminal window and run

```shell
/bin/bash <(curl -s https://raw.githubusercontent.com/DataFoxCo/gocsv/latest/scripts/install-latest-darwin-amd64.sh)
```

### Table of Contents

- [Introduction](#introduction)
- [Subcommands](#subcommands)
- [Specifying Columns](#specifying-columns)
- [Pipelining](#pipelining)
- [Examples](#examples)
- [Installation](#installation)

## Introduction

The tool is built for [pipelining](#pipelining), so most commands accept a CSV from standard output and output to standard output.

Subcommands:

- [describe](#describe) - Get basic information about a CSV.
- [dimensions](#dimensions) (alias: `dims`) - Get the dimensions of a CSV.
- [clean](#clean) - Clean a CSV of common formatting issues.
- [tsv](#tsv) - Transform a CSV into a TSV.
- [head](#head) - Extract the first _N_ rows from a CSV.
- [tail](#tail) - Extract the last _N_ rows from a CSV.
- [headers](#headers) - View the headers from a CSV.
- [view](#view) - Display a CSV in a pretty tabular format.
- [rename](#rename) - Rename the headers of a CSV.
- [behead](#behead) - Remove header row(s) from a CSV.
- [autoincrement](#autoincrement) (alias: `autoinc`) - Add a column of incrementing integers to a CSV.
- [stack](#stack) - Stack multiple CSVs into one CSV.
- [split](#split) - Split a CSV into multiple files.
- [sort](#sort) - Sort a CSV based on one or more columns.
- [filter](#filter) - Extract rows whose columns match a regular expression.
- [select](#select) - Extract specified columns.
- [sample](#sample) - Sample rows.
- [unique](#unique) (alias: `uniq`) - Extract unique rows based upon certain columns.
- [join](#join) - Join two CSVs based on equality of elements in a column.
- [xlsx](#xlsx) - Convert sheets of a XLSX file to CSV.


## Subcommands

### describe

Get basic information about a CSV. This will output the number of rows and columns in the CSV, the column headers in the CSV, and the inferred type of each column.

Usage

```shell
gocsv describe FILE
```

### dimensions

_Alias:_ `dims`

Get the dimensions of a CSV.

Usage

```shell
gocsv dimensions FILE
```

### clean

Clean a CSV of common formatting issues. Currently this consists of making sure all rows are the same length (padding short rows and trimming long ones) and removing empty rows at the end.

Usage:

```shell
gocsv clean [--verbose] [--no-trim] [--excel] [--numbers] FILE
```

Arguments:

- `--verbose` (optional) Print out to stderr when cleaning the CSV. 
- `--no-trim` (optional) Do not remove trailing rows that are empty.
- `--excel` (optional) Clean the CSV for issues that will cause problems with Excel. See [Excel specifications and limitations](https://support.office.com/en-us/article/Excel-specifications-and-limits-16c69c74-3d6a-4aaf-ba35-e6eb276e8eaa).
  - Truncate any cells that exceed the maximum character limit of 32767.
- `--numbers` (optional) Clean the CSV for issues that will cause problems with Numbers.
  - Truncate the number of rows in the CSV at 65535, the maximum amount of rows that Numbers displays.

### tsv

Transform a CSV into a TSV. This can very useful if you want to pipe the result to `pbcopy` (OS X) in order to paste it into a spreadsheet tool.

Usage:

```shell
gocsv tsv FILE
```

### head

Extract the first _N_ rows from a CSV.

Usage:

```shell
gocsv head [-n N] FILE
```

Arguments:

- `-n` (optional) The number of rows to extract. If `N` is an integer, it will extract the first _N_ rows. If `N` is prepended with `+`, it will extract all except the last _N_ rows.

### tail

Extract the last _N_ rows from a CSV.

Usage:

```shell
gocsv tail [-n N] FILE
```

Arguments:

- `-n` (optional) The number of rows to extract. If `N` is an integer, it will extract the last _N_ rows. If `N` is prepended with `+`, it will extract all except the first _N_ rows.

### headers

View the headers of a CSV along with the index of each header.

Usage:

```bash
gocsv headers FILE
```

### view

Display a CSV in a pretty tabular format.

Usage:

```shell
gocsv view [-n N] [--max-width N] FILE
```

Arguments:

- `-n` (optional) Display only the first _N_ rows of the CSV.
- `--max-width` (optional, default 20) The maximum width of each cell for display. If a cell exceeds the maximum width, it will be truncated in the display.

If the length of a cell exceeds `--max-width` it will be truncated with an ellipsis. If a cell contains a new-line character, only the first line will be displayed.

### rename

Rename the headers of a CSV.

Usage:

```shell
gocsv rename --columns COLUMNS --names NAMES FILE
```

Arguments:

- `--columns` (shorthand `-c`) A comma-separated list of the columns to rename. See [Specifying Columns](#specifying-columns) for more details.
- `--names` A comma-separated list of names to change each column to. This must be the same length as and match the order of the `columns` argument.

### behead

Remove the header from a CSV

Usage:

```shell
gocsv behead [-n N] FILE
```

Arguments:

- `-n` (optional) Number of header rows to remove. Defaults to 1.

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

Also note that the `stack` subcommand does not support piping from standard input.

### split

Split a CSV into multiple files.

Usage:

```shell
gocsv split --max-rows N [--filename-base FILENAME] FILE
```

Arguments:

- `--max-rows` Maximum number of rows per final CSV.
- `--filename-base` (optional) Prefix of the resulting files. The file outputs will be appended with `"-1.csv"`,`"-2.csv"`, etc. If not specified, the base filename will be the same as the base of the input filename, unless the input is specified by standard input. If so, then the base filename will be `out`.

### sort

Sort a CSV by multiple columns, with or without type inference. The currently supported types are float, int, and string.

Usage:

```shell
gocsv sort --columns COLUMNS [--reverse] [--no-inference] FILE
```

Arguments:

- `--columns` (shorthand `-c`) A comma-separated list (in order) of the columns to sort against. See [Specifying Columns](#specifying-columns) for more details.
- `--reverse` (optional) Reverse the order of sorting. By default the sort order is ascending.
- `--no-inference` (optional) Skip type inference when sorting.

### filter

Filter a CSV by rows whose columns match a regular expression

Usage:

```shell
gocsv filter [--columns COLUMNS] [--regex REGEX] [--gt N] [--gte N] [--lt N] [--lte N] [--exclude] FILE
```

Arguments:

- `--columns` (optional, shorthand `-c`) A comma-separated list of the columns to filter against. If no columns are specified, then filter checks every column on a row. If a row matches on any of the columns, the row is considered a match. See [Specifying Columns](#specifying-columns) for more details.
- `--regex` (optional) Regular expression to use to match against.
- `--case-insensitive` (optional, shorthand `-i`) When using the `--regex` flag, use this flag to specify a case insensitive match rather than the default case sensitive match.
- `--gt` , `--gte`, `--lt`, `--lte` (optional) Compare against a number.
- `--exclude` (optional) Exclude rows that match. Default is to include.

Note that one of `--regex`, `--gt` , `--gte`, `--lt`, or `--lte` must be specified.

### select

Select (or exclude) columns from a CSV

Usage:

```shell
gocsv select --columns COLUMNS [--exclude] FILE
```

Arguments:

- `--columns` (shorthand `-c`) A comma-separated list (in order) of the columns to select. If you want to select a column multiple times, you can! See [Specifying Columns](#specifying-columns) for more details.
- `--exclude` (optional) Exclude the specified columns (default is to include).

### unique

_Alias:_ `uniq`

Extract unique rows based upon certain columns.

Usage:

```shell
gocsv unique [--columns COLUMNS] [--sorted] FILE
```

Arguments

- `--columns` (optional, shorthand `-c`) A comma-separated list (in order) of the columns to use to define uniqueness. If no columns are specified, it will perform uniqueness across the entire row. See [Specifying Columns](#specifying-columns) for more details.
- `--sorted` (optional) Specify whether the input is sorted. If the input is sorted, the unique subcommand will run more efficiently.

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

### join

Join two CSVs using an inner (default), left, right, or outer join.

Usage:

```shell
gocsv join --columns COLUMNS [--left] [--right] [--outer] LEFT_FILE RIGHT_FILE
```

Arguments:

- `--columns` (shorthand `-c`) A comma-separated list (in order) of the columns to use for joining. You must specify either 1 or 2 columns. When 1 is specified, it will join the CSVs using that column in both the left and right CSV. When 2 are specified, it will join using the first column on the left CSV and the second column on the right CSV. See [Specifying Columns](#specifying-columns) for more details.
- `--left` (optional) Perform a left join.
- `--right` (optional) Perform a right join.
- `--outer` (optional) Perform an outer join.

Note that by default it will perform an inner join. It will exit if you specify multiple types of join.

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

## Specifying Columns

When specifying a column on the command line, you can specify either the index or the name of the column. The tool will always try to interpret the column first by index and then by name. The tool uses 1-based indexing (as in the output of the [headers](#headers) subcommand). When specifying the name, it will use the first column that is an exact case-sensitive match.

## Pipelining

Because all of the subcommands (other than [stack](#stack)) support receiving a CSV from standard input, you can easily pipeline:

```shell
cat test-files/left-table.csv \
  | gocsv join --left --columns LID,RID test-files/right-table.csv \
  | gocsv filter --columns XYZ --regex "[ev]e-\d$" \
  | gocsv select --columns LID,XYZ \
  | gocsv sort --columns LID,XYZ
```

### Pipelining Support

| Subcommand    |    Input    |  Output  |
| ------------- | :---------: | :------: |
| describe      |  &#x2714;   |   N/A    |
| dimensions    |  &#x2714;   |   N/A    |
| clean         |  &#x2714;   | &#x2714; |
| tsv           |  &#x2714;   | &#x2714; |
| head          |  &#x2714;   | &#x2714; |
| tail          |  &#x2714;   | &#x2714; |
| headers       |  &#x2714;   |   N/A    |
| view          |  &#x2714;   |   N/A    |
| rename        |  &#x2714;   | &#x2714; |
| behead        |  &#x2714;   | &#x2714; |
| autoincrement |  &#x2714;   | &#x2714; |
| stack         | Coming Soon | &#x2714; |
| split         |  &#x2714;   |   N/A    |
| sort          |  &#x2714;   | &#x2714; |
| filter        |  &#x2714;   | &#x2714; |
| select        |  &#x2714;   | &#x2714; |
| sample        |  &#x2714;   | &#x2714; |
| unique        |  &#x2714;   | &#x2714; |
| join          |  &#x2714;   | &#x2714; |
| xlsx          |     N/A     |    *     |

\* `xlsx` sends output to standard out when using the `--sheet` flag.

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

##### Sort by Multiple Columns

```shell
gocsv sort --columns LID,ABC --reverse test-files/left-table.csv
```

##### Combine Multiple CSVs

```shell
gocsv stack --groups "Primer Archivo,Segundo Archivo,Tercer Archivo" --group-name "Orden de Archivo" test-files/stack-1.csv test-files/stack-2.csv test-files/stack-3.csv
```

## Installation

For the latest pre-built binaries, see the [Latest Release](https://github.com/DataFoxCo/gocsv/releases/tag/latest) page.

### Apple OS X

#### Simple Version

Open a Terminal window and paste the following command:

```shell
/bin/bash <(curl -s https://raw.githubusercontent.com/DataFoxCo/gocsv/latest/scripts/install-latest-darwin-amd64.sh)
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

Download `gocsv-windows-amd64.zip`. Then good luck.
