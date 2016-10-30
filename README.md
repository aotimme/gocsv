# GoCSV

Command line CSV processing tool based on [csvkit](https://csvkit.readthedocs.io). But faster and less memory intensive.

The tool is built for pipelining, so every command (other than [stack](#stack)) accepts a CSV from standard input, and every command outputs to standard out.

Subcommands:

- [clean](#clean) - Clean a CSV of common formatting issues.
- [headers](#headers) - View the headers from a CSV.
- [behead](#behead) - Remove the header from a CSV.
- [autoincrement](#autoincrement) - Add a column of incrementing integers to a CSV.
- [stack](#stack) - Stack multiple CSVs into one CSV.
- [sort](#sort) - Sort a CSV based on one or more columns.
- [filter](#filter) - Extract rows whose column matches a regular expression.
- [select](#select) - Extract specified columns.
- [join](#join) - Join two CSVs based on equality of elements in a column.


## Subcommands

### clean

Clean a CSV of common formatting issues. Currently this consists of making sure all rows are the same length (padding short rows and trimming long ones) and removing empty lines at the end.

Usage:

```shell
gocsv clean [--no-trim] FILE
```

Arguments:

- `--no-trim` (optional) Do not remove trailing rows that are empty.

### headers

View the headers of a CSV along with the index of each header.

Usage:

```bash
gocsv headers FILE
```

### behead

Remove the header from a CSV

Usage:

```shell
gocsv behead FILE
```

### autoincrement

Append (or prepend) a column of incrementing integers to each row. This can be helpful to be able to map back to the original row after a number of transformations.

Alias: `autoinc`

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

### sort

Sort a CSV by multiple columns, with or without type inference. The currently supported types are float, int, and string.

Usage:

```shell
gocsv sort --columns COLUMNS [--reverse] [--no-inference] FILE
```

Arguments:

- `--columns` A comma-separated list (in order) of the columns to sort against. See [Specifying Columns](#specifying-columns) for more details.
- `--reverse` (optional) Reverse the order of sorting. By default the sort order is ascending.
- `--no-inference` (optional) Skip type inference when sorting.

### filter

Filter a CSV by rows whose column matches a regular expression

Usage:

```shell
gocsv filter --column COLUMN --regex REGEX FILE
```

Arguments:

- `--column` Column to filter on. See [Specifying Columns](#specifying-columns) for more details.
- `--regex` Regular expression to use to match against.

### select

Select (or exclude) columns from a CSV

Usage:

```shell
gocsv select --columns COLUMNS [--exclude] FILE
```

Arguments:

- `--columns` A comma-separated list (in order) of the columns to select. If you want to select a column multiple times, you can! See [Specifying Columns](#specifying-columns) for more details.
- `--exclude` (optional) Exclude the specified columns (default is to include).

### join

Join two CSVs using an inner (default), left, right, or outer join.

Usage:

```shell
gocsv join --columns COLUMNS [--left] [--right] [--outer] LEFT_FILE RIGHT_FILE
```

Arguments:

- `--columns` A comma-separated list (in order) of the columns to use for joining. You must specify either 1 or 2 columns. When 1 is specified, it will join the CSVs using that column in both the left and right CSV. When 2 are specified, it will join using the first column on the left CSV and the second column on the right CSV. See [Specifying Columns](#specifying-columns) for more details.
- `--left` (optional) Perform a left join.
- `--right` (optional) Perform a right join.
- `--outer` (optional) Perform an outer join.

Note that by default it will perform an inner join. It will exit if you specify multiple types of join.

## N.B.

### Specifying Columns

When specifying a column on the command line, you can specify either than name or the index of the column. The tool will _always_ try to interpret the column by index, and then by name. The tool uses 1-based indexing (as in the output of the [headers](#headers) subcommand). When specifying the name, it will use the the _first_ column that matches.

## Pipelining

Because all of the subcommands (other than [stack](#stack)) support receiving a CSV from standard input, you can easily pipeline:

```shell
cat test-files/left-table.csv \
  | gocsv join --left --columns LID,RID test-files/right-table.csv \
  | gocsv filter --column XYZ --regex "[en]e" \
  | gocsv select --columns LID,XYZ \
  | gocsv sort --columns LID,XYZ
```

## Installation

For the latest pre-built binaries, see the [Latest Release](https://github.com/DataFoxCo/gocsv/releases/tag/latest) page.

### Apple OS X

To install the pre-built binary for Apple OS X, download the `gocsv-darwin-amd64.zip` file. It should download into your `~/Downloads` directory. To install it, open a Terminal window and do the following:

```shell
cd ~/Downloads
unzip gocsv-darwin-amd64.zip
mv gocsv-darwin-amd64/gocsv /usr/local/bin
rmdir gocsv-darwin-amd64
```

To verify that it has installed, open a _new_ Terminal windo and run

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

To verify that it has installed, open a _new_ Terminal windo and run

```shell
gocsv help
```

You should see the `gocsv` help message.

### Windows

Download `gocsv-windows-amd64.zip`. Then good luck.

TODO
----

- [ ] Support `-` as a filename specifying `stdin` like csvkit does.
- [ ] Support other delimiters (not just `,`) for both reading and writing.
- [ ] Implement filtering by numeric types.
- [ ] Add subcommand autocomplete (for `zshell` at least).

