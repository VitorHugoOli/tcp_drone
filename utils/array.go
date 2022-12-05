package utils

func Contains(slice []int, element int) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func Difference(a, b []int) []int {
	mb := make(map[int]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []int
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func Remove(slice []int, s int) []int {
	var i int
	for i = range slice {
		if slice[i] == s {
			break
		}
	}
	return append(slice[:i], slice[i+1:]...)
}

// IndexOf returns the index of the first instance of x in a, or -1 if x is not present in a.
func IndexOf(a []int, x int) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

// InsertAt inserts an element at a given index
func InsertAt(slice []int, index int, element int) []int {
	return append(slice[:index], append([]int{element}, slice[index:]...)...)
}
