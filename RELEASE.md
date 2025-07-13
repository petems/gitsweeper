# Release Process

This document describes the automated release process for gitsweeper.

## Overview

The project uses GitHub Actions to automatically build cross-platform binaries and create GitHub releases when tags are pushed to the repository.

## Supported Platforms

The release process builds binaries for the following platforms:

- **Linux**
  - x86_64 (amd64)
  - ARM64 (arm64)
- **macOS**
  - Intel (amd64)
  - Apple Silicon (arm64)
- **Windows**
  - x86_64 (amd64)

## How to Create a Release

1. **Update the version** in `main.go`:
   ```go
   const Version = "x.y.z"
   ```

2. **Commit the version change**:
   ```bash
   git add main.go
   git commit -m "Bump version to x.y.z"
   ```

3. **Create and push a tag**:
   ```bash
   git tag vx.y.z
   git push origin vx.y.z
   ```

4. **GitHub Actions will automatically**:
   - Build binaries for all supported platforms
   - Create compressed archives (`.tar.gz` for Unix, `.zip` for Windows)
   - Create a GitHub release with the tag
   - Upload all archives as release assets

## Local Development

### Build for all platforms locally

```bash
make build-all
```

This will create binaries in the `bin/` directory for all supported platforms.

### Create release archives locally

```bash
make release-archives
```

This will create compressed archives in the `dist/` directory.

## Release Artifacts

Each release includes the following downloadable artifacts:

- `gitsweeper-vx.y.z-linux-amd64.tar.gz`
- `gitsweeper-vx.y.z-linux-arm64.tar.gz`
- `gitsweeper-vx.y.z-darwin-amd64.tar.gz`
- `gitsweeper-vx.y.z-darwin-arm64.tar.gz`
- `gitsweeper-vx.y.z-windows-amd64.zip`

Each archive contains:
- The `gitsweeper` binary (or `gitsweeper.exe` for Windows)
- The `README.md` file

## Installation Instructions

Users can install gitsweeper by:

1. Downloading the appropriate archive for their platform from the GitHub releases page
2. Extracting the archive
3. Moving the binary to a directory in their PATH
4. Making it executable (Linux/macOS): `chmod +x gitsweeper`

## Troubleshooting

### Release failed to create

- Check that the tag follows the `vx.y.z` format
- Ensure the repository has the necessary permissions for GitHub Actions
- Check the Actions tab for detailed error logs

### Binary doesn't work on target platform

- Verify the correct architecture was downloaded
- Check that the binary has execute permissions (Linux/macOS)
- Ensure the target system meets the minimum requirements for Go binaries

## GitHub Actions Workflow

The release workflow (`.github/workflows/release.yml`) consists of two jobs:

1. **Build Job**: Builds binaries for all platforms in parallel using a matrix strategy
2. **Release Job**: Creates the GitHub release and uploads all artifacts

The workflow is triggered only on tag pushes that match the pattern `v*`.