# Architecture

This page covers how `gitsweeper` is laid out, what libraries it leans on, and the trade-offs behind the design choices.

## Source layout

```text
gitsweeper/
├── main.go                  CLI entry point + flag parsing
├── internal/
│   ├── githelpers.go        go-git operations: branch listing, hash-match pass
│   ├── cherrycheck.go       Pass 2: shells out to `git cherry` / `git patch-id`
│   ├── prompthelpers.go     y/N confirmation prompts
│   ├── loghelpers.go        Lightweight info/fatal logging
│   └── slicehelpers.go      Slice/set utilities
└── features/                Cucumber acceptance specs (Ruby + Docker)
```

Everything CLI-related lives in `main.go`. Everything else lives under `internal/` so it can't be imported by other modules — this is a binary, not a library.

## Two-pass merge detection

`GetMergedBranches` in [`internal/githelpers.go`](https://github.com/petems/gitsweeper/blob/master/internal/githelpers.go) drives the detection. It runs in two passes:

1. **Pass 1 — hash matching (via `go-git`).** Walks the main branch's commit history and checks whether each remote branch's HEAD hash appears in it. Catches regular merges and fast-forwards.
2. **Pass 2 — cherry / patch-id (via shell-out).** For every branch *not* matched by Pass 1, `gitsweeper` shells out to `git cherry` and `git patch-id` to compare diff content rather than commit hashes. Catches squash merges and rebases.

Pass 2 can be disabled with `--no-deep-check`. Both passes respect the `--max-commits` limit.

Pass 1 is concurrent when there are more than 10 branches to check (`findMergedBranchesConcurrent`), sequential below that threshold. The worker count is capped at `runtime.NumCPU()`.

See [Squash-Merge Detection](../usage/squash-merge-detection.md) for the user-facing summary.

## `go-git` vs. shelling out

`gitsweeper` uses [`go-git`](https://github.com/go-git/go-git) for everything **except** two things:

| Operation | Implementation | Why |
| --- | --- | --- |
| List branches, walk commits | `go-git` | Pure Go, no fork, cross-platform |
| **Branch deletion** | Shell out to `git push --delete` | go-git's authentication support is incomplete |
| **`git cherry` / `git patch-id`** | Shell out | `go-git` doesn't implement these commands |

### Why shell out for deletion

`go-git`'s authentication API is famously complex. Production environments use a wide variety of auth methods — SSH keys with passphrases, SSH agents, credential helpers, OAuth tokens, deploy keys, GitHub App tokens — and `go-git` doesn't handle all of them well. See the [long-standing tracking issue](https://github.com/go-git/go-git/issues/28).

Rather than re-implement authentication, `gitsweeper` lets the user's existing `git` configuration handle it. Whatever combination of SSH agent and credential helper you already have working, `gitsweeper` inherits for free.

The deletion call (in `DeleteBranch`):

```text
git push <remote> --delete <branchShortName>
```

…runs with `GIT_TERMINAL_PROMPT=0` and a 30-second timeout. Inputs are validated against shell-injection patterns (no leading `-`, non-empty) before being passed as separate arguments to `exec.CommandContext` — no shell interpolation involved.

### Why shell out for cherry / patch-id

`go-git` simply doesn't implement `git cherry` or `git patch-id`. Re-implementing them from scratch would be a significant project, and the shell-out path works on every platform where `git` is installed (which is "everywhere `gitsweeper` would be useful").

## Go version requirement

**Go 1.23 or later.** This is enforced by `go-git` v5.14.0+, which requires `golang.org/x/crypto@v0.35.0+` to mitigate [CVE-2025-22869 / GO-2025-3487](https://pkg.go.dev/vuln/GO-2025-3487), a denial-of-service vulnerability in SSH implementations. CI workflows and local development should all use Go 1.23+.

If you're stuck on an older Go, the last compatible `gitsweeper` builds use `go-git` v5.13.x.

## Binary size optimization

The binary is built with `-s -w -trimpath` and additional `goreleaser` strip flags, producing a ~7.8 MB binary across all platforms. See [the goreleaser config](https://github.com/petems/gitsweeper/blob/master/.goreleaser.yml) and the historical [optimization notes](../history/code-reviews/final-optimization-report.md) for the full set of flags.

## Authentication assumption

`gitsweeper` assumes you already have working `git push` credentials. If `git push origin --delete <branch>` doesn't work from your shell, neither will `gitsweeper cleanup`. The error messages from `gitsweeper` deliberately propagate `git`'s own output so you can debug the underlying auth issue directly.
