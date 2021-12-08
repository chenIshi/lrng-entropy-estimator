package util

import (
	"encoding/csv"
	"os"
	"log"
)

func Min (nums []float64) float64 {
	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}
	return min
}

func Collector() {
	
}

func dump_csv(filename string, nums []int) {
	f, err := os.Create(filename)
	defer f.Close()

	if err != nil {
		log.Fatalln("error writing ", filename, " due to ", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
}