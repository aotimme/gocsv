# GoCSV

Command line CSV processing tool based on [csvkit](https://csvkit.readthedocs.io). But faster and less memory intensive.

The tool is built for pipelining, so every command (other than [headers](#headers)) accepts a CSV from standard input, and every command outputs to standard out.

Subcommands:

- [headers](#headers) - View the headers from a CSV.
- [behead](#behead) - Remove the header from a CSV.
- [stack](#stack) - Stack multiple CSVs into one CSV.
- [sort](#sort) - Sort a CSV based on one or more columns.
- [filter](#filter) - Extract rows whose column matches a regular expression.
- [select](#select) - Extract specified columns.
- [join](#join) - Join two CSVs based on equality of elements in a column.


## Subcommands

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

Because all of the subcommands (other than [headers](#headers)) support receiving a CSV from standard input, you can easily pipeline:

```shell
cat test-files/left-table.csv \
  | gocsv join --left --columns LID,RID test-files/right-table.csv \
  | gocsv filter --column XYZ --regex "[en]e" \
  | gocsv select --columns LID,XYZ \
  | gocsv sort --columns LID,XYZ
```

TODO
----

- [x] Enable column specification by index.


- [x] Guard against attempting to select element of row if row length is too short.


- [x] Implement right join.
- [x] Implement outer join.
- [x] Implement `stack` for stacking two CSVs on top of one another (assuming the same headers.)


- [x] Implement `sort` for sorting a CSV based on one or more columns (TBD whether it supports type inference to know when to sort a column as a string or as an int or float -- or a date even??).
- [x] Implement filtering _out_ columns via the `filter` subcommand.
- [ ] Support other delimiters (not just `,`). For reading and writing.


- [ ] Package it up (with cross-compilation?)
- [ ] Add subcommand autocomplete (for `zshell` at least)