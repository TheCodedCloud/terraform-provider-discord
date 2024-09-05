package common

// FetchKeyByValue returns the key of a map by its value.
func FetchKeyByValue[K comparable](m map[K]string, value string) (interface{}, bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}

	return nil, false
}
