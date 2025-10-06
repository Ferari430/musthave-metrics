package main

import "fmt"

func merge(nums1 []int, nums2 []int) []int {

	var result []int
	i := 0
	j := 0

	for i < len(nums1) && j < len(nums2) {

		if nums1[i] < nums2[j] {
			result = append(result, nums1[i])
			i++
		} else {
			result = append(result, nums2[j])
			j++
		}
	}

	for j < len(nums2)-1 {
		result = append(result, nums2[j])
		j++
	}
	return result
}

func main() {
	nums1 := []int{1, 2, 3}
	nums2 := []int{4, 5, 6}

	fmt.Println(merge(nums1, nums2)) // [1 2 3 4 5 6]
}
