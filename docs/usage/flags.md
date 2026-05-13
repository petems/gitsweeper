# Flags

Every flag works with both `gitsweeper preview` and `gitsweeper cleanup`.

| Flag | Default | Description |
| --- | --- | --- |
| `--origin` | `origin` | Name of the remote to clean up. |
| `--master` | `master` | Name of the main branch. Set to `main` for most modern GitHub repos. |
| `--skip` | _(empty)_ | Comma-separated list of branches to exclude from detection and deletion. |
| `--force` | `false` | Skip the confirmation prompt during `cleanup`. |
| `--max-commits` | `10000` | Maximum commits to scan in the main branch's history. |
| `--no-deep-check` | `false` | Disable squash-merge detection (`git cherry` / `git patch-id`). Faster but misses squash merges. |
| `--debug` | `false` | Enable debug logging to stderr. |

## Examples

### Use `main` instead of `master`

```bash
gitsweeper preview --master main
```

### Clean up against a fork's upstream

```bash
gitsweeper preview --origin upstream
```

### Skip release and hotfix branches

```bash
gitsweeper cleanup --skip "release-2024,hotfix-prod"
```

The `--skip` flag matches the short branch name (the part after `origin/`). Glob-style matching is **not** currently supported — each entry is an exact match.

### Run unattended in CI

```bash
gitsweeper cleanup --force --master main
```

Combine with `--no-deep-check` if CI time matters more than thoroughness:

```bash
gitsweeper cleanup --force --master main --no-deep-check
```

### Search deeper history on long-lived repos

The default `--max-commits 10000` is enough for most projects. Bump it on monorepos or long-lived repositories where merges may be older than that:

```bash
gitsweeper preview --max-commits 100000
```

### Debug a detection mismatch

```bash
gitsweeper preview --debug
```

Debug mode prints, for each branch, which detection method (hash match or cherry/patch-id) matched it — useful when you're investigating why a particular branch did or did not show up.
