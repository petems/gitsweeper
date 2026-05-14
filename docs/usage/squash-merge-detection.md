# Squash-Merge Detection

The original `git-sweep` tool — and `gitsweeper`'s initial implementation — detected merged branches by checking whether a branch's tip commit appeared in the main branch's history. That works for regular merges and fast-forward merges, but it **misses squash merges entirely**.

Squash merging is the default merge strategy on many GitHub repositories (the *Squash and merge* button on PRs). When you squash, GitHub combines all of the branch's commits into a single new commit on the main branch with a different hash. Since the original branch commits never appear on the main branch, hash-based detection can't find them.

In real-world testing, **~24% of deletable branches were only detectable via squash-merge detection**.

## How `gitsweeper` handles it

`gitsweeper` runs a two-pass detection strategy:

### Pass 1 — Hash matching (fast)

Walks the main branch's commit history and checks whether each branch's HEAD commit appears in it. Catches:

- Regular merge commits
- Fast-forward merges

This pass is fast and uses [`go-git`](https://github.com/go-git/go-git) entirely in-process.

### Pass 2 — Cherry / patch-id (thorough)

For every branch *not* matched by Pass 1, `gitsweeper` shells out to:

- `git cherry` — compares individual commit diffs. Handles single-commit squash merges, rebases, and cherry-picks.
- `git patch-id` — compares the combined branch diff against upstream commits. Handles multi-commit squash merges.

This pass works with **any git hosting provider** — it uses vanilla git commands rather than the GitHub or GitLab API. That means no tokens to configure and no rate limits to worry about.

## Controlling detection

Both passes respect `--max-commits` (default `10000`). You can disable Pass 2 entirely if you want a faster scan and don't care about squash merges:

```bash
# Disable squash-merge detection
gitsweeper preview --no-deep-check

# Search more history (both passes respect this limit)
gitsweeper preview --max-commits 50000
```

## When Pass 2 is skipped

Pass 2 needs a real filesystem worktree (it shells out to `git`). It is automatically skipped — with an info-level log message — when:

- The repository is bare.
- The repository is in-memory (test fixtures).
- `git` is not on the user's `PATH`.

In those cases `gitsweeper` falls back to Pass 1 results only.

## Why not the GitHub API?

The simpler alternative would be to query the GitHub API for merged PRs by branch name. We deliberately avoid this because:

- It would tie `gitsweeper` to a single forge (GitHub vs. GitLab vs. Bitbucket vs. Gitea vs. self-hosted).
- It would require auth tokens for private repos.
- It would hit API rate limits on busy repositories.

The patch-id approach gives us the same coverage with zero configuration. The trade-off is that it's slower — typically a few seconds extra on a typical repo, which is fine for a tool you run once a week.

For the history of how this two-pass design was chosen, see the [original detection-strategy advice](../explanation/detection-strategy.md).
