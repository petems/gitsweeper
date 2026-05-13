# Quickstart

Once `gitsweeper` is [installed](installation.md), there are two commands you'll use day-to-day: `preview` and `cleanup`.

## See what would be deleted

`preview` is read-only — it lists every remote branch that's already been merged into the main branch, but doesn't delete anything:

```bash
$ gitsweeper preview
Fetching from the remote...

These branches have been merged into master:
  origin/merged_already_to_master

To delete them, run again with `gitsweeper cleanup`
```

## Delete the merged branches

Once you're happy with the preview, run `cleanup`:

```bash
$ gitsweeper cleanup
Fetching from the remote...

These branches have been merged into master:
  origin/merged_already_to_master

Delete these branches? (y/n)
```

Confirm with `y` to delete. To skip the prompt entirely (useful in scripts and CI), pass `--force`:

```bash
gitsweeper cleanup --force
```

## Common variations

```bash
# Use a different main branch
gitsweeper preview --master main

# Use a different remote
gitsweeper preview --origin upstream

# Skip specific branches (comma-separated)
gitsweeper preview --skip "release/*,hotfix-prod"

# Faster but less thorough — skips squash-merge detection
gitsweeper preview --no-deep-check

# Search deeper history (default is 10,000 commits)
gitsweeper preview --max-commits 50000

# Show debug logging
gitsweeper preview --debug
```

See the [flags reference](../usage/flags.md) for every option.

## What "merged" means here

`gitsweeper` considers a branch merged if **any** of these are true:

1. Its tip commit appears in the main branch's history (regular merge / fast-forward).
2. Its content matches a commit in the main branch via `git cherry` / `git patch-id` (squash merge, rebase, cherry-pick).

The second check is what makes `gitsweeper` reliable on repositories that use GitHub's *Squash and merge* button — see [Squash-Merge Detection](../usage/squash-merge-detection.md) for the details.
