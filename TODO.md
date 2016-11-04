# TODO

## Definitely

- [ ] Subcommand to give basic stats on the CSV (number of rows and columns).
- [ ] Implement `format` subcommand to output with different delimiters, etc. This might replace the `tsv` subcommand.
- [ ] Support `-` as a filename specifying `stdin` like csvkit does.
- [ ] Support other delimiters (not just `,`) for both reading and writing.
- [ ] Implement filtering by numeric types (`--gt`, `--gte`, etc.).

## Maybe

- [ ] Implement `view` subcommand with pretty printing.
- [ ] Add subcommand autocomplete (for `zshell` at least).