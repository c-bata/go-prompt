package bisect

import "sort"

// Right to locate the insertion point for v in a to maintain sorted order.
func Right[T ~int](a []T, v T) T {
	return bisectRightRange(a, v, 0, len(a))
}

func bisectRightRange[T ~int](a []T, v T, lo, hi int) T {
	s := a[lo:hi]
	return T(sort.Search(len(s), func(i int) bool {
		return s[i] > v
	}))
}
