package util

// 获取map的key，返回所有key组成的slice
func MapKeys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for key, _ := range data {
		keys = append(keys, key)
	}
	return keys
}

// 获取map的key，返回所有key组成的slice
func MapIntKeys(data map[int]int) []int {
	keys := make([]int, 0, len(data))
	for key, _ := range data {
		keys = append(keys, key)
	}
	return keys
}
