package internal

// IsStringInSlice checks if a string is present in a slice of strings
//
//	letters := []string{"a", "b", "c", "d"}
//  IsStringInSlice("a", letters) // true
//  IsStringInSlice("e", letters) // false
//
func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
