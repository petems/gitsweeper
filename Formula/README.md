# Homebrew Formula for gitsweeper

This directory contains the Homebrew formula for gitsweeper, allowing users to install the tool via Homebrew without needing a separate tap repository.

## Installation

Users can install gitsweeper directly from this repository:

```bash
# Install directly from the formula
brew install --formula Formula/gitsweeper.rb

# Force compilation from source
brew install --build-from-source --formula Formula/gitsweeper.rb
```

## How the Formula Works

The formula intelligently chooses the installation method based on the platform:

1. **Intel Macs**: Downloads pre-built `darwin-amd64` binary for speed
2. **Apple Silicon Macs**: Downloads pre-built `darwin-arm64` binary for speed  
3. **Other platforms**: Falls back to source compilation using Go

## Automatic Updates

The `.github/workflows/homebrew.yml` workflow automatically updates the formula during releases:

1. Calculates SHA256 checksums for source archive and both binary archives
2. Updates version numbers and checksums in the formula
3. The updated formula is committed to the repository

## Placeholders

The formula contains placeholders that are replaced during CI:

- `PLACEHOLDER_SOURCE_SHA256` - SHA256 of the source archive
- `PLACEHOLDER_BINARY_SHA256` - SHA256 of the Intel binary archive  
- `PLACEHOLDER_BINARY_ARM64_SHA256` - SHA256 of the ARM64 binary archive

## Testing

The formula includes basic tests to verify:
- The binary was installed correctly
- Help output contains expected content
- Version output works

## Maintenance

When making changes to the formula:

1. Test locally with `brew install --formula Formula/gitsweeper.rb`
2. Test uninstall with `brew uninstall gitsweeper`
3. Ensure placeholders are preserved for CI replacement
4. Update this documentation if behavior changes