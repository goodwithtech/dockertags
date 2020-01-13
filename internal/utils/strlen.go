package utils

// StrByLen is alias string slice
type StrByLen []string

// Len is
func (a StrByLen) Len() int { return len(a) }

// Less is
func (a StrByLen) Less(i, j int) bool { return len(a[i]) < len(a[j]) }

// Swap is
func (a StrByLen) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
