package util

import (
	"os"
	"syscall"
	"log"
)

const EVAL_METHEOD_COUNT = 2

// Different types of RNGs
const (
	RNG_UNI int = iota
)

type Num_msg struct {
	Idx int
	Val int
	Rng_t int
}

// Different types of evalulators
const (
	EVAL_LRNG3 int = iota
	EVAL_DIFF
)

type Eval_msg struct {
	Idx int
	Val int
	Eval_t int
}

// Entropy to collectors
type Entropy_msg struct {
	Idx int
	Val float64
	Eval_t int
}

// Controller signals
const (
	CTRL_COMM_RESP int = iota
	CTRL_OUT_REQ
	CTRL_OUT_RESP
)

type Ctrl_msg struct {
	Idx int
	Signal int
}

func Min (nums []float64) float64 {
	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}
	return min
}

func Mkdir(filename string) {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)

	err := os.MkdirAll("eval", 0766)
	CheckError("Can't create directory in collector", err)
}

func CheckError(message string, err error) {
    if err != nil {
        log.Fatalln(message, err)
		os.Exit(1)
    }
}