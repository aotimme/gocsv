"""Offline installer tests using real ZIP extraction and mocked platform/downloads."""
import os
from pathlib import Path
import shutil
import subprocess
import tempfile
import unittest
import zipfile

SCRIPTS = Path(__file__).resolve().parents[1]


class InstallerTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="gocsv installer ")
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.bin = self.root / "bin"
        self.bin.mkdir()
        self.install = self.root / "install with spaces"
        self.scratch = self.root / "scratch"
        self.scratch.mkdir()
        self.env = dict(os.environ, PATH=str(self.bin), HOME=str(self.root),
                        TMPDIR=str(self.scratch), GOCSV_INSTALL_DIR=str(self.install),
                        TEST_OS="Linux", TEST_ARCH="x86_64",
                        TEST_ZIP=str(self.root / "release.zip"),
                        TEST_URL=str(self.root / "url"))
        self.env.pop("GOCSV_VERSION", None)
        for command in ("unzip", "install", "mkdir", "mv", "mktemp", "rm"):
            (self.bin / command).symlink_to(shutil.which(command))
        self.mock("uname", 'case "$1" in -s) echo "$TEST_OS";; -m) echo "$TEST_ARCH";; esac')
        self.downloader("curl")

    def mock(self, name, body):
        path = self.bin / name
        path.write_text("#!/bin/sh\nset -eu\n" + body + "\n")
        path.chmod(0o755)

    def downloader(self, name):
        self.mock(name, '''
[ "${TEST_FAILURE:-}" != download ] || exit 22
printf '%s' "$2" > "$TEST_URL"
[ "$3" = -o ] || [ "$3" = -O ]
/bin/cp "$TEST_ZIP" "$4"
''')

    def archive(self, os_name="linux", arch="amd64", missing=False):
        with zipfile.ZipFile(self.env["TEST_ZIP"], "w") as archive:
            archive.writestr("unexpected" if missing else f"gocsv-{os_name}-{arch}/gocsv",
                             "#!/bin/sh\necho fixture-version\n")

    def run_installer(self, piped=False):
        script = SCRIPTS / "install.sh"
        result = subprocess.run(["/bin/sh"] if piped else ["/bin/sh", str(script)],
                                input=script.read_text() if piped else None,
                                text=True, capture_output=True, env=self.env)
        self.assertEqual(list(self.scratch.iterdir()), [], result.stderr)
        return result

    def test_supported_platforms(self):
        platforms = [("Darwin", "x86_64", "darwin", "amd64"),
                     ("Darwin", "arm64", "darwin", "arm64"),
                     ("Linux", "x86_64", "linux", "amd64"),
                     ("Linux", "i686", "linux", "386"),
                     ("Linux", "aarch64", "linux", "arm64"),
                     ("Linux", "armv7l", "linux", "arm"),
                     ("Linux", "ppc64le", "linux", "ppc64le"),
                     ("Linux", "riscv64", "linux", "riscv64")]
        for kernel, machine, os_name, arch in platforms:
            with self.subTest(kernel=kernel, machine=machine):
                self.env.update(TEST_OS=kernel, TEST_ARCH=machine)
                self.archive(os_name, arch)
                result = self.run_installer(piped=True)
                self.assertEqual(result.returncode, 0, result.stderr)
                self.assertEqual(Path(self.env["TEST_URL"]).read_text(),
                                 f"https://github.com/aotimme/gocsv/releases/latest/download/gocsv-{os_name}-{arch}.zip")
                binary = self.install / "gocsv"
                self.assertTrue(os.access(binary, os.X_OK))
                self.assertEqual(subprocess.check_output([str(binary)], text=True).strip(), "fixture-version")

    def test_default_directory(self):
        self.env.pop("GOCSV_INSTALL_DIR")
        self.archive()
        result = self.run_installer()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertTrue((self.root / ".local/bin/gocsv").is_file())
        self.assertIn("PATH", result.stdout)

    def test_pinned_version(self):
        self.env["GOCSV_VERSION"] = "v1.0.6"
        self.archive()
        result = self.run_installer(piped=True)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(Path(self.env["TEST_URL"]).read_text(),
                         "https://github.com/aotimme/gocsv/releases/download/v1.0.6/gocsv-linux-amd64.zip")

    def test_invalid_version(self):
        self.env["GOCSV_VERSION"] = "../latest"
        result = self.run_installer()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("GOCSV_VERSION", result.stderr)
        self.assertFalse(Path(self.env["TEST_URL"]).exists())

    def test_wget_fallback(self):
        (self.bin / "curl").unlink()
        self.downloader("wget")
        self.archive()
        result = self.run_installer()
        self.assertEqual(result.returncode, 0, result.stderr)

    def test_unsupported_platforms(self):
        for kernel, arch in [("FreeBSD", "amd64"), ("Linux", "s390x"), ("Linux", "armv6l"), ("Darwin", "i686")]:
            with self.subTest(kernel=kernel, arch=arch):
                self.env.update(TEST_OS=kernel, TEST_ARCH=arch)
                self.assertNotEqual(self.run_installer().returncode, 0)
                self.assertFalse(Path(self.env["TEST_URL"]).exists())

    def test_failures_preserve_existing_install(self):
        self.install.mkdir()
        binary = self.install / "gocsv"
        binary.write_text("existing installation")
        for failure in ("download", "corrupt", "missing"):
            with self.subTest(failure=failure):
                self.env["TEST_FAILURE"] = failure
                self.archive(missing=failure == "missing")
                if failure == "corrupt":
                    Path(self.env["TEST_ZIP"]).write_text("not a ZIP")
                result = self.run_installer()
                self.assertNotEqual(result.returncode, 0)
                self.assertNotIn("GoCSV installed", result.stdout)
                self.assertEqual(binary.read_text(), "existing installation")

    def test_failed_copy_preserves_existing_install(self):
        self.archive()
        self.install.mkdir()
        binary = self.install / "gocsv"
        binary.write_text("existing installation")
        (self.bin / "install").unlink()
        self.mock("install", 'printf partial > "$4"; exit 1')
        result = self.run_installer()
        self.assertNotEqual(result.returncode, 0)
        self.assertEqual(binary.read_text(), "existing installation")
        self.assertEqual(list(self.install.iterdir()), [binary])

    def test_missing_dependencies(self):
        for command in ("unzip", "curl"):
            path = self.bin / command
            saved = path.readlink() if path.is_symlink() else None
            contents = None if saved else path.read_text()
            path.unlink()
            result = self.run_installer()
            self.assertNotEqual(result.returncode, 0)
            self.assertIn("required", result.stderr)
            if saved:
                path.symlink_to(saved)
            else:
                path.write_text(contents)
                path.chmod(0o755)

    def test_destination_is_directory(self):
        self.archive()
        (self.install / "gocsv").mkdir(parents=True)
        result = self.run_installer()
        self.assertNotEqual(result.returncode, 0)
        self.assertEqual(list((self.install / "gocsv").iterdir()), [])

    def test_unwritable_destination(self):
        self.archive()
        self.install.write_text("a file, not a directory")
        result = self.run_installer()
        self.assertNotEqual(result.returncode, 0)
        self.assertNotIn("GoCSV installed", result.stdout)


if __name__ == "__main__":
    unittest.main()
