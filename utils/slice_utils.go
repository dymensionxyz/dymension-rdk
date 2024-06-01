package utils

import (
	"slices"
)

// RemoveDuplicates removes duplicates from a slice of any comparable type.
func RemoveDuplicates[T comparable](input []T) []T {
	uniqueMap := make(map[T]bool)
	var uniqueList []T
	for _, v := range input {
		if !uniqueMap[v] {
			uniqueMap[v] = true
			uniqueList = append(uniqueList, v)
		}
	}
	return uniqueList
}

// RemoveValueFromSlice removes all occurrences of a specific value from a slice of any comparable type.
func RemoveValueFromSlice[T comparable](list []T, x T) []T {
	return slices.DeleteFunc(list, func(v T) bool { return v == x })
}

// ShiftValueToEnd takes a value and slice of any comparable type
// Further removes the duplicate from slice and moves the given value at the end of slice.
func ShiftValueToEnd[T comparable](list []T, x T) []T {
	list = RemoveDuplicates(list)
	var index int = -1
	for i, val := range list {
		if val == x {
			index = i
			break
		}
	}
	if index == -1 {
		return list
	}
	list = append(list[:index], list[index+1:]...)
	list = append(list, x)
	return list
}
