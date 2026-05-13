# Commands

`gitsweeper` exposes two subcommands. Both must be run from inside a git working tree.

## `gitsweeper preview`

Lists every remote branch on the configured origin that has been merged into the main branch. Read-only — nothing is deleted.

```bash
$ gitsweeper preview
Fetching from the remote...

These branches have been merged into master:
  origin/feat/auth-cleanup
  origin/fix/typo-in-readme

To delete them, run again with `gitsweeper cleanup`
```

Exit codes:

| Code | Meaning |
| --- | --- |
| 0 | Ran successfully (whether or not branches were found) |
| 1 | An error occurred (no git repo, remote not found, etc.) |

## `gitsweeper cleanup`

Same detection as `preview`, but then actually deletes the branches with `git push <remote> --delete <branch>`. Prompts for confirmation before deleting.

```bash
$ gitsweeper cleanup
Fetching from the remote...

These branches have been merged into master:
  origin/feat/auth-cleanup
  origin/fix/typo-in-readme

Delete these branches? (y/n) y
Deleting origin/feat/auth-cleanup... done
Deleting origin/fix/typo-in-readme... done
```

Pass `--force` to skip the prompt — useful for cron jobs and CI:

```bash
gitsweeper cleanup --force
```

### How deletion works

Deletion shells out to your system's `git` command:

```text
git push <origin> --delete <branch>
```

This means `gitsweeper` automatically uses whatever authentication you already have configured — SSH agents, credential helpers, deploy tokens, etc. See [Architecture](../explanation/architecture.md#why-shell-out-for-deletion) for why this design choice was made.

Each deletion has a 30-second timeout. If a remote is unreachable or hangs, `gitsweeper` reports the failure and moves on to the next branch.

## Behaviour shared by both commands

- The remote is fetched first so your local view of merged branches is up-to-date.
- The main branch defaults to `master`. Override with `--master main` (or any name).
- The remote defaults to `origin`. Override with `--origin upstream`.
- Branches matching `--skip` are excluded from both the preview list and the delete list.
- Up to `--max-commits` (default `10000`) of main-branch history are scanned. Increase this if you've enabled deep history retention.
- Squash-merge detection (`git cherry` / `git patch-id`) runs by default. Disable with `--no-deep-check` for a faster but less thorough scan.

See [Flags](flags.md) for the complete option list.
