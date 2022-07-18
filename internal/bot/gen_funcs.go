package bot

func isContain[T comparable](val T, values []T) bool {
	for _, v := range values {
		if val == v {
			return true
		}
	}

	return false
}
