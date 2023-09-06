package scrub

const (
	mask = "********"
)

// Replace all middle characters with asterisks.
func Middle(s string) string {
	l := len(s)
	if l <= 8 {
		return mask
	}
	if l > 8 && l <= 10 {
		return s[:2] + mask + s[l-2:]
	}
	if l > 10 && l < 20 {
		return s[:3] + mask + s[l-3:]
	}
	if l >= 20 && l < 30 {
		return s[:5] + mask + s[l-5:]
	}
	return s[:10] + mask + s[l-10:]
}
