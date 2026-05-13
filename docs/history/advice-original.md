# Gitsweeper: Branch Detection Improvement Advice

This document captures findings from a manual branch cleanup session on `petersouter.xyz`, where multiple detection strategies were needed to identify all safely-deletable branches. The goal is to inform improvements to [gitsweeper](https://github.com/petems/gitsweeper).

## Summary

Out of 80 local branches, **49 were safe to delete**. gitsweeper's current approach (commit hash matching against master's history) would have caught **37 of them**. The remaining **12** required checking the GitHub API for squash-merged PRs — a detection method gitsweeper does not currently implement.

## Detection Methods Used

### Method 1: `git branch --merged master` (what gitsweeper approximates)

**How it works:** Lists branches whose tip commit is reachable from master's HEAD. This is essentially what gitsweeper does internally — it walks master's commit history and checks if each branch's HEAD commit appears in that history.

**What it found:** 37 branches.

**Why it works:** When a branch is merged via a regular merge commit or fast-forward merge, the branch's commits become ancestors of master. Git can trivially determine reachability.

**Limitations:**
- Does not detect squash merges (the squash creates a new commit on master with a different hash)
- Does not detect rebase merges (rebased commits get new hashes)
- gitsweeper additionally limits history to 10,000 commits, which could miss older merges

### Method 2: GitHub API — check for merged PRs by branch name

**How it works:** For each branch that was _not_ detected by Method 1, query the GitHub API for PRs with `head` matching the branch name and `state=merged`.

**Command used:**
```bash
gh pr list --head "$branch" --state merged --json number,title,mergedAt
```

**What it found:** 12 additional branches that had been squash-merged via GitHub PRs.

| Branch | PR | Merged Date |
|---|---|---|
| `chore/add-claude-settings` | #163 | 2026-04-09 |
| `chore/fix-claude-and-agent-setup` | #162 | 2026-04-03 |
| `claude/test-blog-editor-skill-fdyeQ` | #161 | 2026-04-03 |
| `docs/add-url-permalink-convention` | #164 | 2026-04-09 |
| `feat/improve-blog-writing-skills` | #146 | 2026-04-03 |
| `feat/improve-conference-talks-section` | #96 | 2026-01-28 |
| `feat/quick-april-blog-post` | #160 | 2026-04-08 |
| `petems/add-garden-section` | #149 | 2026-03-20 |
| `petems/fix-hugo-new-docs` | #142 | 2026-03-02 |
| `petems/migrate-to-papermod` | #147 | 2026-03-18 |
| `refactor/fix-theme-issues-and-fixes` | #110 | 2026-01-28 |
| `refactor/garden-topic-subdirectories` | #155 | 2026-03-21 |

**Why Method 1 missed these:** All 12 were squash-merged on GitHub. A squash merge creates a single new commit on master that combines all the branch's commits. The resulting commit hash on master is different from any commit on the branch, so commit-hash-based reachability checks fail.

### Method 3 (not used, but worth considering): Patch-ID / diff comparison

**How it works:** Compare the `git patch-id` of commits on a branch against patch-IDs of commits on master. Patch-IDs are based on the diff content, not the commit hash, so they can match squash-merged or cherry-picked content.

```bash
# Get patch-id for the branch's combined diff
git diff master...branch | git patch-id --stable

# Compare against patch-ids of recent master commits
git log master --format=%H | while read hash; do
  git diff "$hash^..$hash" | git patch-id --stable
done
```

**Trade-offs:**
- Can detect squash merges and cherry-picks without needing a forge API
- Works offline / with any Git hosting provider
- More expensive to compute (requires diffing every candidate commit)
- May produce false positives if the same change was applied differently

### Method 4 (not used, but worth considering): `git cherry` / symmetric diff

**How it works:** `git cherry master branch` checks which commits on `branch` have not been applied to master, using patch-ID comparison under the hood. If all commits show `-` (already applied), the branch is fully merged.

```bash
# If output is empty or all lines start with '-', branch is merged
git cherry master branch
```

**Trade-offs:**
- Built-in Git command, no API needed
- Handles rebases and cherry-picks
- Slower than hash comparison on large histories
- Can be confused by conflicts resolved differently during merge

## Recommendations for gitsweeper

### High Impact

1. **Add squash-merge detection via forge API (GitHub/GitLab/Bitbucket)**
   - This is the biggest gap. Squash merges are the default merge strategy on many GitHub repos.
   - Use the GitHub API (`GET /repos/{owner}/{repo}/pulls?head={branch}&state=closed`) to check if a branch's PR was merged.
   - Could be behind a `--check-forge` or `--github` flag.
   - The `gh` CLI or GitHub REST API both support this.

2. **Add `git cherry`-based detection as a fallback**
   - For cases where no forge API is available, `git cherry` can detect squash merges and rebases.
   - Run this as a second pass on branches that weren't caught by hash matching.
   - Could be behind a `--deep-check` flag since it's slower.

### Medium Impact

3. **Make `MaxCommitsToCheck` configurable**
   - The current hard limit of 10,000 commits could miss merges in large repos.
   - Add a `--max-commits` flag.

4. **Auto-detect default branch**
   - Instead of defaulting to "master", detect the repo's default branch from `refs/remotes/origin/HEAD` or the GitHub API.
   - Fall back to "master" or "main" if detection fails.

5. **Support glob patterns in `--skip`**
   - Allow `--skip="renovate/*"` to skip all Renovate branches, for example.

### Lower Impact

6. **Report detection method in output**
   - When listing merged branches, indicate _how_ each was detected (hash match, API check, patch-ID match). This builds user trust and helps debugging.

7. **Dry-run diff summary**
   - In preview mode, show how many unique commits each branch has vs master. A branch with 0 unique commits is almost certainly safe to delete regardless of merge method.

## Why This Matters

In the `petersouter.xyz` repo, **24% of deletable branches (12 out of 49)** were only detectable via the GitHub API. This is not an unusual ratio — squash merging is the default strategy for many teams and projects. Without forge-aware detection, gitsweeper will consistently undercount merged branches on repos that use squash merges.
