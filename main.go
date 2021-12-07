package main

import (
	"math/rand"
	"fmt"
	"time"
	"os"
	"math"
)

const eval_method_count = 2

// Stimuate entropy pools in lrng, i.e. input/output pool 
// According user configuration, it would have to produce different level of "randomness" 
// Also, there could different method to construct a "not-that-random" number
func entropy_pool(request_ch chan int, entropy_1ch chan int, entropy_2ch chan int) {
	// As for common practice, most RNG requires a true random seed to achieve unpredictability 
	// However, here we actually don't care about such attribute 
	// Also, using a fixed random seed help us reproduce the experiment results if desired 
	rand.Seed(80)
	for {
		req := <-request_ch
		// if the main process want to shut down the entropy collection pool 
		if req < 0 {
			// it should also propagate such signal to evaluators before shut down
			entropy_1ch <- -1 
			entropy_2ch <- -1 
			return
		} else {
			rng_num :=  rand.Intn(req)
			entropy_1ch <- rng_num
			entropy_2ch <- rng_num
		}
	}
}

func min (nums []int) int {
	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}
	return min
}

func lrng_eval_3(entropy_ch chan int, response_ch chan bool, demo_ch chan bool) {
	var estimations []float64
	// LRNG entropy estimation requires to keep the previous states 
	// By default, it would require a 3 step of differential, that is 4 states in total
	pre_states := []int{-1, -1,-1, -1}
	// As stated above, LRNG requires to collect previous states before estimation 
	// In this case, it would start actual estimation until the fourth round
	eval_round := 0

	for {
		select {
		case rng_num := <- entropy_ch:
			if rng_num < 0 {
				return
			} else {
				// Push the generated rng number into pre_states 
				// Also kick out the eldest state (that is, FIFO)
				pre_states = append(pre_states[1:], rng_num)
				if eval_round >= 3 {
					// LRNF entropy estimation: see [The Linux Pseudorandom Number Generator Revisited](https://eprint.iacr.org/2012/251.pdf) for details 
					delta1 := pre_states[3] - pre_states[2]
					delta2 := pre_states[3] - 2 * pre_states[2] + pre_states[1]
					delta3 := pre_states[3] - 3 * pre_states[2] + 3 * pre_states[1] - pre_states[0]
					delta := float64(min([]int{delta1, delta2, delta3}))
					entropy := 0.0
					if delta > (1 << 12) {
						entropy = 11
					} else if delta > 2{
						entropy = math.Log2(delta)
					}
					estimations = append(estimations, entropy)
				}
				eval_round ++
				response_ch <- true
			}
		case <-demo_ch:
			if eval_round < 3 {
				fmt.Fprintf(os.Stderr, "LRNG entropy estimation hasn't begin\n")
				os.Exit(1)
			} else {
				sum := 0.0
				for i:=0;i<len(estimations);i++ {
					sum += estimations[i]
				}
				avg := sum / float64(len(estimations))
				fmt.Println("LRNG entropy estimation avg = ", avg)
				demo_ch <- true
			}
		}		
	}
}

func differential_eval(entropy_ch chan int, response_ch chan bool, demo_ch chan bool) {
	for {
		select {
		case rng_num := <- entropy_ch:
			if rng_num < 0 {
				return
			} else {
				response_ch <- true
			}
		case <-demo_ch:
			fmt.Println("Differential entropy estimation not implemented yet!")
			demo_ch <- false
		}		
	}
}

func main() {
	// User configuration

	// Module initialization
	// request_ch: channel between "main" and "entropy_pool"
	// transmit the maximum of random number (minimum being 0)
	// **Sending a negative signal means to shut down the goroutine**
	request_ch := make(chan int, 10)
	// requestX_ch: channel between "entropy_pool" and "evaluatorX"
	// transmit the generated "random number" 
	// **Sending a negative signal means to shut down the goroutine**
	// TODO: replace with a channel array
	entropy_1ch := make(chan int, 10)
	entropy_2ch := make(chan int, 10)
	// request_ch: channel between "evaluator" and "main"
	// transmit the estimated randomness/entropy from the random number sequence
	response_ch := make(chan bool, 10)
	demo_1ch := make(chan bool)
	demo_2ch := make(chan bool)
	go entropy_pool(request_ch, entropy_1ch, entropy_2ch)
	go lrng_eval_3(entropy_1ch, response_ch, demo_1ch)
	go differential_eval(entropy_2ch, response_ch, demo_2ch)

	var test_scale int = 1000
	var max_rng int = 50

	for i:=0; i<test_scale; i++ {
		request_ch <- max_rng
		
		lrng_timer := time.NewTimer(time.Duration(100 * time.Millisecond))
		for j:=0; j<eval_method_count; j++ {
			select {
			case <- response_ch:
				continue
			case <- lrng_timer.C:
				fmt.Fprintf(os.Stderr, "error: timeout in main at %d, %d\n", i, j)
				os.Exit(1)
			}
		}
	}

	// Output evaluation results using different estimation approaches 
	demo_1ch <- true
	demo_2ch <- true


	demo_timer := time.NewTimer(time.Duration(100 * time.Millisecond))
	select {
	case <- demo_1ch:
	case <- demo_timer.C:
		fmt.Fprintf(os.Stderr, "error: timeout in main at demoing LRNG\n")
		os.Exit(1)
	}
	select {
	case <- demo_2ch:
	case <- demo_timer.C:
		fmt.Fprintf(os.Stderr, "error: timeout inmain at demoing Differential\n")
		os.Exit(1)
	}
	

	// Cleanup goroutines (Potentially useless)
	request_ch <- -1
}