package internal

import "sort"

// IsStringInSlice checks if a string is present in a slice of strings
// Optimized version that uses binary search for sorted slices
//
//	letters := []string{"a", "b", "c", "d"}
//  IsStringInSlice("a", letters) // true
//  IsStringInSlice("e", letters) // false
//
func IsStringInSlice(a string, list []string) bool {
	// For small slices, linear search is faster due to CPU cache locality
	if len(list) < 8 {
		for _, b := range list {
			if b == a {
				return true
			}
		}
		return false
	}

	// For larger slices, check if sorted and use binary search
	if isSorted(list) {
		return binarySearchString(a, list)
	}

	// Fall back to linear search for unsorted larger slices
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// isSorted checks if a string slice is sorted (used to determine if binary search is applicable)
func isSorted(list []string) bool {
	for i := 1; i < len(list); i++ {
		if list[i-1] > list[i] {
			return false
		}
	}
	return true
}

// binarySearchString performs binary search on a sorted string slice
func binarySearchString(target string, list []string) bool {
	return sort.SearchStrings(list, target) < len(list) && list[sort.SearchStrings(list, target)] == target
}

// StringSliceToSet converts a string slice to a map[string]bool for O(1) lookups
// Use this when you need to perform multiple lookups against the same slice
func StringSliceToSet(slice []string) map[string]bool {
	set := make(map[string]bool, len(slice))
	for _, s := range slice {
		set[s] = true
	}
	return set
}

// IsStringInSet checks if a string exists in a string set (map[string]bool)
// This is O(1) vs O(n) for slice lookups
func IsStringInSet(a string, set map[string]bool) bool {
	return set[a]
}
