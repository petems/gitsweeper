# gitsweeper

A CLI tool for cleaning up merged git branches.

`gitsweeper` identifies remote branches that have already been merged into your main branch — including squash-merged branches that simpler tools miss — and helps you delete them safely.

![gitsweeper demo](assets/demo.gif)

## Why this tool exists

`gitsweeper` is a modern Go rewrite of the Python-based [`git-sweep`](https://github.com/arc90/git-sweep), which has been unmaintained since 2016. Beyond replacing an abandoned tool, `gitsweeper` adds **squash-merge detection** — in real-world testing, ~24% of safely-deletable branches were *only* findable via the deeper diff-based check.

See [Motivation](explanation/motivation.md) for the full story.

## Quick start

```bash
# Install (Linux / macOS)
curl -sSfL https://raw.githubusercontent.com/petems/gitsweeper/master/install.sh | sh

# See which remote branches have been merged
gitsweeper preview

# Delete them
gitsweeper cleanup
```

For other installation methods see [Installation](getting-started/installation.md).

## Where to go next

<div class="grid cards" markdown>

-   :material-rocket-launch: **[Quickstart](getting-started/quickstart.md)**

    Run your first `preview` and `cleanup` in 60 seconds.

-   :material-console-line: **[Commands & flags](usage/commands.md)**

    Reference for every CLI option, with examples.

-   :material-magnify: **[Squash-merge detection](usage/squash-merge-detection.md)**

    How `gitsweeper` finds branches that other tools miss.

-   :material-source-branch: **[Architecture](explanation/architecture.md)**

    Why `gitsweeper` mixes `go-git` with shell-outs to `git`.

</div>

## Project links

- **Source:** [github.com/petems/gitsweeper](https://github.com/petems/gitsweeper)
- **Releases:** [GitHub releases](https://github.com/petems/gitsweeper/releases)
- **License:** [MIT](https://github.com/petems/gitsweeper/blob/master/LICENSE.md)
