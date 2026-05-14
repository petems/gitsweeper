# Release Process

`gitsweeper` releases are driven by git tags and built by GitHub Actions using [goreleaser](https://goreleaser.com/). Pushing a tag matching `v*` triggers the release workflow, which builds cross-platform binaries and creates a GitHub release.

## Supported platforms

| OS | Architectures |
| --- | --- |
| Linux | amd64, arm64 |
| macOS | amd64 (Intel), arm64 (Apple Silicon) |
| Windows | amd64 |

## Cutting a release

1. **Update the version constant** in `main.go`:

    ```go
    const Version = "x.y.z"
    ```

2. **Commit the version bump**:

    ```bash
    git add main.go
    git commit -m "chore: bump version to x.y.z"
    ```

3. **Tag and push**:

    ```bash
    git tag vx.y.z
    git push origin vx.y.z
    ```

4. **GitHub Actions takes over**, automatically:

    - Building binaries for every supported platform.
    - Creating compressed archives (`.tar.gz` for Unix, `.zip` for Windows).
    - Creating a GitHub release tied to the tag.
    - Uploading every archive plus a `checksums.txt` file as release assets.
    - Updating the Homebrew formula (in `Formula/gitsweeper.rb`) with new SHA256 checksums.

The full release workflow lives in [`.github/workflows/release.yml`](https://github.com/petems/gitsweeper/blob/master/.github/workflows/release.yml).

## Local builds

### Build for every platform

```bash
make build-all
```

Output lands in `bin/` (one binary per platform/arch).

### Build release archives locally

```bash
make release-snapshot
```

This runs `goreleaser` in snapshot mode — same artifacts as a real release, just without publishing. Output lands in `dist/`.

You can also use the bundled smoke test:

```bash
./test-goreleaser.sh
```

## Release artifacts

Each release publishes:

```text
gitsweeper-vx.y.z-linux-amd64.tar.gz
gitsweeper-vx.y.z-linux-arm64.tar.gz
gitsweeper-vx.y.z-darwin-amd64.tar.gz
gitsweeper-vx.y.z-darwin-arm64.tar.gz
gitsweeper-vx.y.z-windows-amd64.zip
gitsweeper_x.y.z_checksums.txt
```

Each archive contains:

- The `gitsweeper` binary (or `gitsweeper.exe` on Windows).
- A copy of `README.md`.

## Installation paths after release

Users can install the new release immediately via:

- The install script (`install.sh` defaults to the latest GitHub release).
- Homebrew (`brew install --formula Formula/gitsweeper.rb` — the formula is auto-updated by the release workflow).
- Direct download from the [GitHub releases page](https://github.com/petems/gitsweeper/releases).

See [Installation](../getting-started/installation.md) for the user-facing instructions.

## Troubleshooting

??? failure "Release workflow didn't fire"
    - Confirm the tag matches `v*` (e.g. `v0.1.2`, not `0.1.2` or `release-0.1.2`).
    - Check that the repository's GitHub Actions permissions allow workflow runs.
    - Look at the **Actions** tab for the detailed error log.

??? failure "Binary doesn't run on the target platform"
    - Verify the user downloaded the right architecture (`uname -m`).
    - Check execute permissions on Linux/macOS: `chmod +x gitsweeper`.
    - Confirm the binary isn't quarantined by Gatekeeper on macOS (`xattr -d com.apple.quarantine gitsweeper`).

??? failure "Homebrew formula update step failed"
    - The formula expects three SHA256 placeholders — see the [Homebrew formula reference](../reference/homebrew-formula.md).
    - Re-running the release workflow from the **Actions** tab is safe; it's idempotent.
