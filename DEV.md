# Dev README

## Development

After cloning this repository, set up the pre-commit hook to ensure proper formatting of the Go code and the `go.mod` and `go.sum` files:
```shell
ln -s ../../git-hooks/pre-commit .git/hooks/pre-commit
```

When developing, you can build a local version of the `gocsv` binary via running `make`. This will create a newly compiled `gocsv` binary in `bin/`.

## Releasing

To release an update to `gocsv`, make sure you have committed and pushed the most recent commit on master. Then:

1. Tag the latest commit as "latest".

   ```shell
   make tag
   ```


2. Create cross-compiled binaries for distribution. Cross-compilation uses [xgo](https://github.com/karalabe/xgo) to handle issues with CGO packages in other platforms and architectures. Because `xgo` requires `docker`, you will need `docker` installed.

   ```shell
   go get -u github.com/karalabe/xgo
   make dist
   ```

   This will create zip files in the `dist` directory holding the `gocsv` binaries for various platforms and architectures.

3. Upload the newly created distribution binaries to the [Latest Release](https://github.com/aotimme/gocsv/releases/tag/latest) page. You will need to [edit](https://github.com/aotimme/gocsv/releases/edit/latest) the release, remove the existing zip files, and upload the recently created zip files in `dist/`.
