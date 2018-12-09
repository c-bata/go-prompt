package bisect

import "sort"

// Right to locate the insertion point for v in a to maintain sorted order.
func Right(a []int, v int) int {
	return bisectRightRange(a, v, 0, len(a))
}

func bisectRightRange(a []int, v int, lo, hi int) int {
	s := a[lo:hi]
	return sort.Search(len(s), func(i int) bool {
		return s[i] > v
	})
}
