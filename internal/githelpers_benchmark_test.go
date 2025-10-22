package internal

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

// generateTestHashes creates a slice of test hashes for benchmarking.
func generateTestHashes(count int) []plumbing.Hash {
	hashes := make([]plumbing.Hash, count)
	for i := 0; i < count; i++ {
		// Create deterministic but unique hashes
		data := []byte(fmt.Sprintf("commit-%d", i))
		hashes[i] = plumbing.NewHash(fmt.Sprintf("%x", sha1.Sum(data)))
	}
	return hashes
}

// generateTestBranches creates test branch info for benchmarking.
func generateTestBranches(hashes []plumbing.Hash) []BranchInfo {
	branches := make([]BranchInfo, len(hashes))
	for i, hash := range hashes {
		branches[i] = BranchInfo{
			Name:   fmt.Sprintf("origin/branch-%d", i),
			Hash:   hash,
			Remote: "origin",
			Short:  fmt.Sprintf("branch-%d", i),
		}
	}
	return branches
}

// Shared test case definitions for benchmarks.
var benchmarkTestCases = []struct {
	name      string
	branches  int
	commits   int
	matchRate float64 // percentage of commits that match branches
}{
	{"10branches_100commits_10pct", 10, 100, 0.1},
	{"100branches_1000commits_5pct", 100, 1000, 0.05},
	{"500branches_5000commits_2pct", 500, 5000, 0.02},
	{"1000branches_10000commits_1pct", 1000, 10000, 0.01},
}

// generateBenchmarkData creates branch and commit data for benchmarking.
// Returns branch info and commit hashes with the specified match rate.
func generateBenchmarkData(branches, commits int, matchRate float64) ([]BranchInfo, []plumbing.Hash) {
	branchHashes := generateTestHashes(branches)
	branchInfos := generateTestBranches(branchHashes)

	commitHashes := make([]plumbing.Hash, commits)
	matchCount := int(float64(commits) * matchRate)
	for i := 0; i < matchCount; i++ {
		// Use actual branch hashes for matches
		commitHashes[i] = branchHashes[i%len(branchHashes)]
	}
	for i := matchCount; i < commits; i++ {
		// Generate unique non-matching hashes
		data := []byte(fmt.Sprintf("non-match-%d", i))
		commitHashes[i] = plumbing.NewHash(fmt.Sprintf("%x", sha1.Sum(data)))
	}

	return branchInfos, commitHashes
}

// BenchmarkHashMapString benchmarks the old approach using string keys.
func BenchmarkHashMapString(b *testing.B) {
	for _, tc := range benchmarkTestCases {
		b.Run(tc.name, func(b *testing.B) {
			branches, commits := generateBenchmarkData(tc.branches, tc.commits, tc.matchRate)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				// OLD APPROACH: String-based map
				branchHashMap := make(map[string][]BranchInfo, len(branches))
				for _, branch := range branches {
					hashStr := branch.Hash.String() // STRING CONVERSION
					branchHashMap[hashStr] = append(branchHashMap[hashStr], branch)
				}

				// Simulate commit processing loop
				var matches int
				for _, commit := range commits {
					commitHash := commit.String() // STRING CONVERSION
					if branchInfos, exists := branchHashMap[commitHash]; exists {
						matches += len(branchInfos)
					}
				}
			}
		})
	}
}

// BenchmarkHashMapPlumbingHash benchmarks the new approach using plumbing.Hash keys.
func BenchmarkHashMapPlumbingHash(b *testing.B) {
	for _, tc := range benchmarkTestCases {
		b.Run(tc.name, func(b *testing.B) {
			branches, commits := generateBenchmarkData(tc.branches, tc.commits, tc.matchRate)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				// NEW APPROACH: plumbing.Hash-based map
				branchHashMap := make(map[plumbing.Hash][]BranchInfo, len(branches))
				for _, branch := range branches {
					branchHashMap[branch.Hash] = append(branchHashMap[branch.Hash], branch)
				}

				// Simulate commit processing loop
				var matches int
				for _, commit := range commits {
					if branchInfos, exists := branchHashMap[commit]; exists {
						matches += len(branchInfos)
					}
				}
			}
		})
	}
}

// BenchmarkHashStringConversion benchmarks the cost of hash string conversion alone.
func BenchmarkHashStringConversion(b *testing.B) {
	hashes := generateTestHashes(1000)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, hash := range hashes {
			_ = hash.String()
		}
	}
}

// BenchmarkHashComparison benchmarks the cost of comparing hashes.
func BenchmarkHashComparison(b *testing.B) {
	hashes := generateTestHashes(1000)
	target := hashes[500]

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, hash := range hashes {
			_ = hash == target
		}
	}
}
