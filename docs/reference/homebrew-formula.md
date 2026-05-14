# Homebrew Formula

`gitsweeper` ships its Homebrew formula directly in the repository at [`Formula/gitsweeper.rb`](https://github.com/petems/gitsweeper/blob/master/Formula/gitsweeper.rb) so users can install via `brew` without needing a separate tap repository.

## Installation

```bash
# Install directly from the formula in this repo
brew install --formula Formula/gitsweeper.rb

# Force compilation from source
brew install --build-from-source --formula Formula/gitsweeper.rb
```

## How the formula picks an install method

The formula picks the fastest install method that works on the current platform:

| Platform | Default behaviour |
| --- | --- |
| Intel Macs | Downloads the pre-built `darwin-amd64` binary |
| Apple Silicon Macs | Downloads the pre-built `darwin-arm64` binary |
| Other platforms | Compiles from source using Go |

Users can override and force a source build with `--build-from-source`.

## Automatic updates on release

Every tagged release runs the [release workflow](https://github.com/petems/gitsweeper/blob/master/.github/workflows/release.yml), which:

1. Calculates SHA256 checksums for the source archive and both binary archives.
2. Substitutes those into the formula's three placeholders.
3. Commits the updated formula back to the repository.

### Placeholders

The committed formula in `master` contains three placeholder strings that CI replaces during release:

| Placeholder | Replaced with |
| --- | --- |
| `PLACEHOLDER_SOURCE_SHA256` | SHA256 of the source tarball |
| `PLACEHOLDER_BINARY_SHA256` | SHA256 of the `darwin-amd64` archive |
| `PLACEHOLDER_BINARY_ARM64_SHA256` | SHA256 of the `darwin-arm64` archive |

**Do not** commit a formula with real SHA256 values into `master`. The placeholders must be present for the release workflow to substitute correctly.

## Testing the formula

The formula includes basic checks that verify:

- The binary was installed correctly.
- Help output contains expected content.
- `gitsweeper --version` returns a sensible value.

To exercise the install + uninstall locally:

```bash
brew install --formula Formula/gitsweeper.rb
gitsweeper --help
brew uninstall gitsweeper
```

## Maintenance checklist

When making changes to the formula:

1. Test locally with `brew install --formula Formula/gitsweeper.rb`.
2. Test the uninstall path with `brew uninstall gitsweeper`.
3. Confirm the three `PLACEHOLDER_*` strings are still present and unmodified.
4. Update this page if formula behaviour changes.
