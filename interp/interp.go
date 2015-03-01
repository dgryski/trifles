// Package interp is an interpolation search
// Code from https://en.wikipedia.org/wiki/Interpolation_search
package interp

func Search(sortedArray []int, toFind int) int {

	// Returns index of toFind in sortedArray, or -1 if not found
	low := 0
	high := len(sortedArray) - 1
	var mid int

	for sortedArray[low] <= toFind && toFind <= sortedArray[high] {
		mid = low + int((float64(toFind-sortedArray[low])*float64(high-low))/float64(sortedArray[high]-sortedArray[low])) //out of range is possible  here
		if sortedArray[mid] < toFind {
			low = mid + 1
		} else if sortedArray[mid] > toFind {
			// Repetition of the comparison code is forced by syntax limitations.
			high = mid - 1
		} else {
			return mid
		}
	}

	if sortedArray[low] == toFind {
		return low
	} else {
		return -1 // Not found
	}
}
