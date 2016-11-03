# TODO

## Definitely

- [ ] Refactor to pass in `*csv.Reader` to the main processing functions rather than `io.Reader`. That way we can specify options to the `*csv.Reader` before passing into the processing function.
- [ ] Subcommand to give basic stats on the CSV (number of rows and columns).
- [ ] Implement `format` subcommand to output with different delimiters, etc. This might replace the `tsv` subcommand.
- [ ] Support `-` as a filename specifying `stdin` like csvkit does.
- [ ] Support other delimiters (not just `,`) for both reading and writing.
- [ ] Implement filtering by numeric types (`--gt`, `--gte`, etc.).
- [ ] Add table to README detailing pipelining support (Y/N stdin, Y/N stdout, notes).

## Maybe

- [ ] Implement `view` subcommand with pretty printing.
- [ ] Add subcommand autocomplete (for `zshell` at least).