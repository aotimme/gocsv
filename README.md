# GoCSV

Command line CSV processing tool based on [csvkit](https://csvkit.readthedocs.io). But faster and less memory intensive.

Subcommands:
- [headers](#headers) - View the headers from a CSV.
- [behead](#behead) - Remove the header from a CSV.
- [stack](#stack) - Stack multiple CSVs into one CSV.
- [sort](#sort) - Sort a CSV based on one or more columns.
- [filter](#filter) - Extract rows whose column matches a regular expression.
- [select](#select) - Extract specified columns.
- [join](#join) - Join two CSVs based on equality of elements in a column.


Every subcommand (other than `headers`) is designed to accept either CSV files as arguments our via standard input for pipelining.

## Subcommands

### headers

View the headers of a CSV.

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

Stack multiple CSVs to create a larger CSV

Usage:

```shell
gocsv stack [--filenames] [--groups GROUPS] [--group-name GROUP_NAME] FILE [FILES]
```

### sort

Sort a CSV by multiple columns, with or without type inference.

Usage:

```shell
gocsv sort --columns COLUMNS [--reverse] [--no-inference] FILE
```

### filter

Filter a CSV by rows whose column matches a regular expression

Usage:

```shell
gocsv filter --column COLUMN --regex REGEX FILE
```

### select

Select (or exclude) columns from a CSV

Usage:

```shell
gocsv select --columns COLUMNS [--exclude] FILE
```

### join

Join two CSVs using an inner (default), left, right, or outer join.

Usage:

```shell
gocsv join --columns COLUMNS [--left] [--right] [--outer] LEFT_FILE RIGHT_FILE
```

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