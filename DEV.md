# Dev README

## Development

After cloning this repository, set up the pre-commit hook to ensure proper formatting of the Go code and the `go.mod` and `go.sum` files:
```shell
ln -s ../../git-hooks/pre-commit .git/hooks/pre-commit
```

When developing, you can build a local version of the `gocsv` binary via running `make`. This will create a newly compiled `gocsv` binary in `bin/`.

## Releasing

To release an update to `gocsv`, make sure you have committed and pushed the most recent commit on master. Then:

1. Tag the latest commit following [semantic versioning](https://semver.org) and push the tag.

   ```shell
   git tag -a v1.2.3 -m "Release v1.2.3"
   git push origin v1.2.3
   ```


2. Create cross-compiled binaries for distribution.

   ```shell
   make dist
   ```

   This will create zip files in the `dist` directory holding the `gocsv` binaries for various platforms and architectures. The version to be returned by `gocsv version` is derived from the Git tag, so make sure to only run this after having carried out step 1.

3. Navigate to the [Releases page](https://github.com/aotimme/gocsv/releases) and draft a new release. Upload the newly created distribution binaries (zip files in `dist/`) by dropping them on the page.
