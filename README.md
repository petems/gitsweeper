# `gitsweeper`

A CLI tool for cleaning up merged git branches.

`gitsweeper` finds remote branches that have already been merged into your main branch — including squash-merged branches that simpler tools miss — and helps you delete them safely.

![gitsweeper demo](docs/assets/demo.gif)

## Quick start

```bash
# Install (Linux / macOS)
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh

# Preview merged branches
gitsweeper preview

# Delete them
gitsweeper cleanup
```

## Documentation

Full documentation lives at **[gitsweeper.readthedocs.io](https://gitsweeper.readthedocs.io/)** and in the [`docs/`](docs/) folder.

Quick links:

- [Installation](docs/getting-started/installation.md) — install script, Homebrew, pre-built binaries, build from source
- [Quickstart](docs/getting-started/quickstart.md) — your first `preview` and `cleanup`
- [Commands & flags](docs/usage/commands.md) — every CLI option with examples
- [Squash-merge detection](docs/usage/squash-merge-detection.md) — how `gitsweeper` finds branches other tools miss
- [Architecture](docs/explanation/architecture.md) — why `gitsweeper` mixes `go-git` with shell-outs
- [Contributing](docs/contributing/development.md) — local dev setup, testing, release process

## License

[MIT](LICENSE.md)
