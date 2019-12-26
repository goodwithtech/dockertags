package utils

type StrByLen []string

func (a StrByLen) Len() int {
	return len(a)
}

func (a StrByLen) Less(i, j int) bool {
	return len(a[i]) < len(a[j])
}

func (a StrByLen) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
