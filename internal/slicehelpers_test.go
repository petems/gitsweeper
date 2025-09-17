package internal

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStringInSlice_SmallSlice(t *testing.T) {
	slice := []string{"a", "b", "c"}

	assert.True(t, IsStringInSlice("a", slice))
	assert.True(t, IsStringInSlice("b", slice))
	assert.True(t, IsStringInSlice("c", slice))
	assert.False(t, IsStringInSlice("d", slice))
	assert.False(t, IsStringInSlice("", slice))
}

func TestIsStringInSlice_LargeSlice(t *testing.T) {
	// Create a large slice
	slice := make([]string, 20)
	for i := 0; i < 20; i++ {
		slice[i] = string(rune('a' + i))
	}

	assert.True(t, IsStringInSlice("a", slice))
	assert.True(t, IsStringInSlice("t", slice))
	assert.False(t, IsStringInSlice("z", slice))
}

func TestIsStringInSlice_SortedSlice(t *testing.T) {
	// Create a sorted slice for binary search
	slice := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "kiwi", "lemon", "mango"}

	assert.True(t, IsStringInSlice("apple", slice))
	assert.True(t, IsStringInSlice("kiwi", slice))
	assert.True(t, IsStringInSlice("mango", slice))
	assert.False(t, IsStringInSlice("orange", slice))
	assert.False(t, IsStringInSlice("zzz", slice))
}

func TestIsStringInSlice_UnsortedLargeSlice(t *testing.T) {
	// Create an unsorted large slice
	slice := []string{"zebra", "apple", "banana", "yak", "cherry", "date", "xray", "elderberry", "fig", "grape"}

	assert.True(t, IsStringInSlice("zebra", slice))
	assert.True(t, IsStringInSlice("apple", slice))
	assert.True(t, IsStringInSlice("grape", slice))
	assert.False(t, IsStringInSlice("orange", slice))
}

func TestIsSorted(t *testing.T) {
	assert.True(t, isSorted([]string{"a", "b", "c", "d"}))
	assert.True(t, isSorted([]string{"apple", "banana", "cherry"}))
	assert.False(t, isSorted([]string{"c", "a", "b"}))
	assert.False(t, isSorted([]string{"zebra", "apple"}))
	assert.True(t, isSorted([]string{}))         // empty slice is sorted
	assert.True(t, isSorted([]string{"single"})) // single element is sorted
}

func TestBinarySearchString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry", "date", "elderberry"}

	assert.True(t, binarySearchString("apple", slice))
	assert.True(t, binarySearchString("cherry", slice))
	assert.True(t, binarySearchString("elderberry", slice))
	assert.False(t, binarySearchString("orange", slice))
	assert.False(t, binarySearchString("aaa", slice))
	assert.False(t, binarySearchString("zzz", slice))
}

func TestStringSliceToSet(t *testing.T) {
	slice := []string{"a", "b", "c", "a"} // duplicate "a"
	set := StringSliceToSet(slice)

	assert.True(t, set["a"])
	assert.True(t, set["b"])
	assert.True(t, set["c"])
	assert.False(t, set["d"])
	assert.Len(t, set, 3) // duplicates removed
}

func TestStringSliceToSetTrimsWhitespace(t *testing.T) {
	slice := []string{" feature/foo", "bar ", "  baz  ", ""}
	set := StringSliceToSet(slice)

	assert.True(t, set["feature/foo"])
	assert.True(t, set["bar"])
	assert.True(t, set["baz"])
	assert.False(t, set[" feature/foo"])
	assert.False(t, set["bar "])
	assert.False(t, set["  baz  "])
	assert.Len(t, set, 3)
}

func TestIsStringInSet(t *testing.T) {
	set := map[string]bool{
		"apple":  true,
		"banana": true,
		"cherry": true,
	}

	assert.True(t, IsStringInSet("apple", set))
	assert.True(t, IsStringInSet("banana", set))
	assert.True(t, IsStringInSet("cherry", set))
	assert.False(t, IsStringInSet("date", set))
	assert.False(t, IsStringInSet("", set))
}

// Benchmark tests to verify performance improvements.
func BenchmarkIsStringInSlice_Small(b *testing.B) {
	slice := []string{"a", "b", "c", "d", "e"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsStringInSlice("c", slice)
	}
}

func BenchmarkIsStringInSlice_Large_Sorted(b *testing.B) {
	slice := make([]string, 100)
	for i := 0; i < 100; i++ {
		slice[i] = string(rune('a'+(i%26))) + string(rune('a'+(i/26)))
	}
	sort.Strings(slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsStringInSlice("ba", slice)
	}
}

func BenchmarkIsStringInSlice_Large_Unsorted(b *testing.B) {
	slice := make([]string, 100)
	for i := 0; i < 100; i++ {
		slice[i] = string(rune('z'-(i%26))) + string(rune('z'-(i/26)))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsStringInSlice("ba", slice)
	}
}

func BenchmarkStringSliceToSet(b *testing.B) {
	slice := make([]string, 100)
	for i := 0; i < 100; i++ {
		slice[i] = string(rune('a' + (i % 26)))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringSliceToSet(slice)
	}
}

func BenchmarkIsStringInSet(b *testing.B) {
	slice := make([]string, 100)
	for i := 0; i < 100; i++ {
		slice[i] = string(rune('a' + (i % 26)))
	}
	set := StringSliceToSet(slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsStringInSet("m", set)
	}
}
