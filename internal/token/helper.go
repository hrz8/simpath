package token

import "strings"

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func spaceDelimitedStringNotGreater(first, second string) bool {
	if first == "" {
		return true
	}
	secondParts := strings.Split(second, " ")

	for _, firstPart := range strings.Split(first, " ") {
		if !stringInSlice(firstPart, secondParts) {
			return false
		}
	}

	return true
}
