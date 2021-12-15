package util

import (
	"os"
	"syscall"
	"log"
	"fmt"
)

const RNG_METHOD_COUNT = 2
const EVAL_METHEOD_COUNT = 3

// Different types of RNGs
type Rng_t int
const (
	RNG_UNI Rng_t = iota
	RNG_RULE30
)

func (r Rng_t) String() string {
	switch r {
	case RNG_UNI:
		return "Uniform"
	case RNG_RULE30:
		return "Rule_30"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

type Num_msg struct {
	Idx int
	Val int
	Rng Rng_t
}

// Different types of evalulators
type Eval_t int
const (
	EVAL_LRNG2 Eval_t = iota
	EVAL_LRNG3
	EVAL_LRNG4
)

func (e Eval_t) String() string {
	switch e {
	case EVAL_LRNG2:
		return "LRNG2"
	case EVAL_LRNG3:
		return "LRNG3"
	case EVAL_LRNG4: 
		return "LRNG4"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

type Eval_msg struct {
	Idx int
	Val int
	Eval Eval_t
}

// Entropy to collectors
type Entropy_msg struct {
	Idx int
	Val float64
	Eval Eval_t
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