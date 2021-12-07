package util

func Min (nums []float64) float64 {
	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}
	return min
}