package utils

import "strings"

func Like(a, b string) bool {
	return strings.Contains(strings.ToLower(b), strings.ToLower(a))
}

// Includes checks if list of strings includes a string
func Includes(list []string, s string) bool {
	for _, elem := range list {
		if s == elem {
			return true
		}
	}
	return false
}

// Exclude checks if list of strings excludes a string
func Excludes(list []string, s string) bool {
	for _, elem := range list {
		if s == elem {
			return false
		}
	}
	return true
}

// LikeAny returns true if any of the list elements is a substring of s
func LikeAny(s string, list ...string) bool {
	for _, elem := range list {
		if strings.Contains(strings.ToLower(s), strings.ToLower(elem)) {
			return true
		}
	}
	return false
}

// HasCommon checks if two list of strings has at least one common element
func HasCommon(listA, listB []string) bool {
	for _, elemA := range listA {
		for _, elemB := range listB {
			if elemA == elemB {
				return true
			}
		}
	}
	return false
}

// NoCommon checks if two list of strings does not have a single common element
func NoCommon(listA, listB []string) bool {
	for _, elemA := range listA {
		for _, elemB := range listB {
			if elemA == elemB {
				return false
			}
		}
	}
	return true
}

// TidySplit is like strings.Split, but the empty strings are not included in the output
func TidySplit(base string, separator string) (out []string) {
	split := strings.Split(base, separator)
	for i := 0; i < len(split); i++ {
		if split[i] != "" {
			out = append(out, split[i])
		}
	}
	return
}

// Chain joins its args with dot "."
func Chain(args ...string) string {
	return strings.Join(args, ".")
}
