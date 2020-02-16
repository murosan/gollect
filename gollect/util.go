package gollect

func trimQuotes(s string) string {
	return s[1 : len(s)-1]
}
