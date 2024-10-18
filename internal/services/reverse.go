package services

// GetReverseMessage returns the input string reversed.
func GetReverseMessage(s string) string {
	runes := []rune(s)
	reversedRunes := make([]rune, len(runes))

	for i := len(runes) - 1; i >= 0; i-- {
		reversedRunes[len(runes)-1-i] = runes[i]
	}

	return string(reversedRunes)
}
