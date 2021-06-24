package Utils

// Contains check that `source` array contains specific `value`
func Contains(source []string, value string) bool {
	for _, it := range source {
		if it == value {
			return true
		}
	}

	return false
}