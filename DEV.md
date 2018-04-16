# Dev README

## Development

After cloning this repository, set up the pre-commit hook to ensure proper formatting of the Go code:
```shell
ln -s ../../git-hooks/pre-commit .git/hooks/pre-commit
```

When developing, you can build a local version of the `gocsv` binary via running `make`. This will create a newly compiled `gocsv` binary in `bin/`.

## Releasing

To release an update to `gocsv`, make sure you have committed and pushed the most recent commit on master. Then:

1. Tag the latest commit as "latest":

   ```shell
   make tag
   ```

2. Create binaries for distribution:

   ```shell
   make dist
   ```

   This will create zip files in the `dist` directory holding the `gocsv` binaries for Linux, Unix, and Windows.

3. Upload the newly created distribution binaries to the [Latest Release](https://github.com/DataFoxCo/gocsv/releases/tag/latest) page. You will need to [edit](https://github.com/DataFoxCo/gocsv/releases/edit/latest) the release, remove the existing zip files, and upload the recently created zip files in `dist/`.
