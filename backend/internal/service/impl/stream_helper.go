package impl

func splitContent(text string, maxLen int) []string {
	if len([]rune(text)) <= maxLen {
		return []string{text}
	}

	runes := []rune(text)
	var chunks []string
	for len(runes) > 0 {
		end := maxLen
		if end > len(runes) {
			end = len(runes)
		}
		if end < len(runes) {
			if idx := lastIndexRune(runes[:end], '\n'); idx > end/2 {
				end = idx + 1
			}
		}
		chunks = append(chunks, string(runes[:end]))
		runes = runes[end:]
	}
	return chunks
}

func lastIndexRune(s []rune, c rune) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}
