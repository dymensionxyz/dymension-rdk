package utils

import "strings"

// ConvertStringToStringMap converts a string to a map[string]string.
// The input string is expected to be for example "key1:value1,key2:value2" where seperators can be specified.
func ConvertStringToStringMap(input string, mapSeperator string, kvSeperator string) map[string]string {
	resultMap := make(map[string]string)
	if input != "" {
		for _, keyValue := range strings.Split(input, mapSeperator) {
			kv := strings.Split(keyValue, kvSeperator)
			resultMap[kv[0]] = kv[1]
		}
	}
	return resultMap
}
