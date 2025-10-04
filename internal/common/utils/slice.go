package utils

func Int64SliceToInterface(slice []int64) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

// UniqueStrings 去重字符串slice（保持顺序）
func UniqueStrings(slice []string) []string {
	if len(slice) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))
	for i := range slice {
		if !seen[slice[i]] {
			seen[slice[i]] = true
			result = append(result, slice[i])
		}
	}
	return result
}

func StringSliceToInterface(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
