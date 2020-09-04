package util

//InStringArray if str in arr, return index. or return -1
func InStringArray(arr []string, str string) int {
	for idx, item := range arr {
		if item == str {
			return idx
		}
	}

	return -1
}
