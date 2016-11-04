# TODO

## Definitely

- [ ] Subcommand to give basic stats on the CSV (number of rows and columns).
- [ ] Implement `format` subcommand to output with different delimiters, etc. This might replace the `tsv` subcommand.
- [ ] Support `-` as a filename specifying `stdin` like csvkit does.
- [ ] Support other delimiters (not just `,`) for both reading and writing.
- [ ] Replace `fmt.Fprintf` and `os.Exit(2)` with `log.Fatalf`.

## Maybe

- [ ] Implement `view` subcommand with pretty printing.
- [ ] Add subcommand autocomplete (for `zshell` at least).