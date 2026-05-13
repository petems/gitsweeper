# Detection Strategy

This page explains how `gitsweeper` decides whether a branch is "merged", and the trade-offs between the different approaches that were considered. For the user-facing summary, see [Squash-Merge Detection](../usage/squash-merge-detection.md). For the full historical analysis that motivated this design, see the [archived advice doc](../history/advice-original.md).

## The four detection methods

When investigating ways to identify safely-deletable branches, four methods were considered:

### 1. `git branch --merged master` (commit-hash reachability)

**How it works.** Lists branches whose tip commit is reachable from the main branch's HEAD. This is what `gitsweeper` does internally in Pass 1: it walks the main branch's commit history and checks whether each branch's HEAD commit appears in that history.

**Catches.** Regular merge commits, fast-forward merges.

**Misses.**

- Squash merges (the squash creates a new commit with a different hash on the main branch).
- Rebase merges (rebased commits get new hashes).
- Anything older than `--max-commits` (default 10 000) commits back.

This is fast and uses [`go-git`](https://github.com/go-git/go-git) entirely in-process. It's also what the original `git-sweep` did, and what most "merged branch" tools do.

### 2. GitHub API — check for merged PRs by branch name

**How it works.** Query the GitHub API for PRs with `head=<branch>` and `state=merged`:

```bash
gh pr list --head "$branch" --state merged --json number,title,mergedAt
```

**Catches.** Every branch that was merged through a GitHub PR, regardless of merge strategy (merge commit, squash, rebase).

**Trade-offs.**

- Forge-specific — works on GitHub, but not GitLab / Bitbucket / Gitea / self-hosted without rewriting per provider.
- Requires auth tokens for private repos.
- Subject to API rate limits.
- Doesn't catch branches merged outside the PR flow (direct pushes, command-line merges).

**Why `gitsweeper` doesn't use this.** Locking the tool to a specific forge would be a significant ergonomics regression. Pattern 3 below gives equivalent coverage with no auth and no rate limits.

### 3. `git cherry` / `git patch-id` — diff content matching

**How it works.** Compare the *content* of commits, not their hashes.

`git cherry master branch` checks which commits on `branch` have not been applied to master, using patch-id comparison under the hood. If all output lines start with `-`, every commit's diff has already been applied — meaning the branch was squash-merged, rebased, or cherry-picked into master.

`git patch-id` produces a stable hash of a diff. By comparing the patch-id of `master...branch` against the patch-ids of recent master commits, you can detect squash merges where the entire branch was combined into a single new commit.

**Catches.** Squash merges, rebase merges, cherry-picks — anything where the content reached the main branch even though the hashes don't match.

**Trade-offs.**

- Slower than hash comparison on large histories.
- Can be confused by conflicts that were resolved differently during merge.
- Requires `git` on `PATH` and a real filesystem worktree.

**This is `gitsweeper`'s Pass 2.** It runs only on branches that Pass 1 didn't catch, so the cost is bounded.

### 4. `git branch --merged` via reflog or full ancestry walk

**How it works.** Walk every commit reachable from master and check ancestry against every branch tip. Variants include using the reflog to catch recently-deleted refs.

**Trade-offs.** Equivalent to Method 1 in coverage but slower. Provides no benefit unless you explicitly want to consider reflog state.

**Why `gitsweeper` doesn't use this.** No meaningful upside over Method 1 + Method 3 combined.

## Why two passes?

| Pass | Method | Speed | What it catches |
| --- | --- | --- | --- |
| 1 | Hash reachability via `go-git` | Fast | Regular merges, fast-forwards |
| 2 | `git cherry` + `git patch-id` via shell-out | Slower (bounded by Pass 1) | Squash merges, rebases, cherry-picks |

This split lets the fast pass handle the common case and only spends extra time on branches that genuinely need the deeper check. Pass 2 can be disabled with `--no-deep-check` if you want maximum speed and accept missing squash-merged branches.

## Why this matters

In real-world testing on a personal-blog-sized repo, 12 of 49 deletable branches (**~24%**) were detectable *only* via Method 3. All twelve had been merged via GitHub's *Squash and merge* button. Without Pass 2, `gitsweeper` would have left them sitting in the remote forever.

The full numbers and per-branch breakdown are in the [archived advice doc](../history/advice-original.md).
