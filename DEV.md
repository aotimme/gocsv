# Dev README

## Development

After cloning this repository, set up the pre-commit hook to ensure proper formatting of the Go code and the `go.mod` and `go.sum` files:

```shell
ln -s ../../git-hooks/pre-commit .git/hooks/pre-commit
```

When developing, you can build a local version of the `gocsv` binary via running `make`. This will create a newly compiled `gocsv` binary in `bin/`.

## Testing

Run `make test` to run the tests with coverage.

If this reports `go: no such tool "covdata"`, you may be hitting a [Go toolchain switching bug](https://github.com/golang/go/issues/75031). Explicitly select the toolchain reported by `go env GOVERSION`:

```shell
GOTOOLCHAIN="$(go env GOVERSION)+auto" make test
```

### Installers

Run the offline shell installer tests on Linux or macOS:

```shell
python3 scripts/tests/test_install.py
```

On Windows, run `./scripts/tests/test-install.ps1` in Windows PowerShell 5.1 or PowerShell 7. These tests use fixture ZIP files, exercise installation and user PATH updates in a temporary directory, and restore the original PATH afterward.

The `Installers` GitHub Actions workflow runs these tests and installs the latest published release on Linux, Intel and Apple Silicon macOS, and x64 and ARM64 Windows. Windows runs under both PowerShell versions. The smoke tests download actual release assets and run `gocsv version`.

## Releasing

To release an update to `gocsv`, make sure you have committed and pushed the most recent commit on main. Then:

1. Tag the latest commit following [semantic versioning](https://semver.org) and push the tag.

   ```shell
   tagName=v9.8.7
   git tag -a $tagName -m "Release $tagName"
   git push origin $tagName
   ```

2. Create cross-compiled binaries for distribution.

   ```shell
   make dist
   ```

   This will create zip files in the `dist` directory holding the `gocsv` binaries for various platforms and architectures. The version to be returned by `gocsv version` is derived from the Git tag, so make sure to only run this after having carried out step 1. The git commands in build scripts will find the most recent tag name.

3. Navigate to the [Releases page](https://github.com/aotimme/gocsv/releases) and draft a new release. Upload the newly created distribution binaries (zip files in `dist/`) by dropping them on the page.
