# History

This section preserves the project's history: design notes, optimization sessions, and the original analysis that motivated `gitsweeper`'s squash-merge detection. **These pages are archived as-is.** They reflect the state of the project at the time they were written and are kept for context rather than as living reference material.

If you're looking for current docs, see [Architecture](../explanation/architecture.md) or [Detection Strategy](../explanation/detection-strategy.md) instead.

## What's here

### [Detection Advice (original)](advice-original.md)

The original investigation that drove `gitsweeper`'s two-pass detection design. Captures findings from a manual branch-cleanup session on a personal blog repository where multiple detection strategies were needed to identify every safely-deletable branch. Contains the per-branch breakdown that motivated the "~24% of deletable branches need deep-check" claim.

### Code-review reports

Seven artifacts from earlier optimization and refactoring sessions. They document the steps and decisions behind:

- The migration from a hand-rolled release pipeline to [goreleaser](https://goreleaser.com/).
- Performance work on the merge-detection passes (concurrent workers, batch sizing, hash-map lookups).
- A handful of code-quality improvements and audit reports.

| Report | Topic |
| --- | --- |
| [GoReleaser Migration](code-reviews/goreleaser-migration.md) | Moving the release pipeline to goreleaser |
| [Optimization Results](code-reviews/optimization-results.md) | Measured impact of the optimization pass |
| [Final Optimization Report](code-reviews/final-optimization-report.md) | Wrap-up of the optimization series |
| [Optimization Summary](code-reviews/optimization-summary.md) | Higher-level summary of the same work |
| [Improvements](code-reviews/improvements.md) | Miscellaneous code-quality improvements |
| [Advanced Optimization Analysis](code-reviews/advanced-optimization-analysis.md) | Deeper analysis of perf-critical paths |
| [Performance Analysis](code-reviews/performance-analysis.md) | Profiling and benchmark notes |

## Why keep this around?

Two reasons:

1. **Context decays.** Reading "why is this concurrent?" against the current code with no history is harder than reading the original analysis that justified it.
2. **The advice doc isn't redundant.** The user-facing [Detection Strategy](../explanation/detection-strategy.md) page is a summary; the [archived advice doc](advice-original.md) has the full per-branch breakdown and method comparison that won't fit cleanly in user docs.

If a future refactor makes any of these obsolete, they should be deleted — not silently rewritten.
