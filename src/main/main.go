package main

import (
	"math/rand"
	"fmt"
	"time"
	"os"
	"estimator"
)

const eval_method_count = 2

// Stimuate entropy pools in lrng, i.e. input/output pool 
// According user configuration, it would have to produce different level of "randomness" 
// Also, there could different method to construct a "not-that-random" number
func entropy_pool(request_ch chan int, entropy_chs [eval_method_count]chan int) {
	// As for common practice, most RNG requires a true random seed to achieve unpredictability 
	// However, here we actually don't care about such attribute 
	// Also, using a fixed random seed help us reproduce the experiment results if desired 
	rand.Seed(80)
	for {
		req := <-request_ch
		// if the main process want to shut down the entropy collection pool 
		if req < 0 {
			// it should also propagate such signal to evaluators before shut down
			for i:=0; i<eval_method_count; i++ {
				entropy_chs[i] <- -1
			}
			return
		} else {
			rng_num :=  rand.Intn(req)
			for i:=0; i<eval_method_count; i++ {
				entropy_chs[i] <- rng_num
			}
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
	var entropy_chs [eval_method_count]chan int
	for i := range entropy_chs {
		entropy_chs[i] = make(chan int, 10)
	}
	// request_ch: channel between "evaluator" and "main"
	// transmit the estimated randomness/entropy from the random number sequence
	response_ch := make(chan bool, 10)

	var demo_chs [eval_method_count]chan bool
	for i := range demo_chs {
		demo_chs[i] = make(chan bool)
	}
	go entropy_pool(request_ch, entropy_chs)
	go estimator.Lrng_eval_3(entropy_chs[0], response_ch, demo_chs[0])
	go estimator.Differential_eval(entropy_chs[1], response_ch, demo_chs[1])

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
	demo_timer := time.NewTimer(time.Duration(100 * time.Millisecond))
	for k:=0; k<eval_method_count; k++ {
		demo_chs[k] <- true
		select {
		case <- demo_chs[k]:
		case <- demo_timer.C:
			fmt.Fprintf(os.Stderr, "error: timeout in main at demoing LRNG\n")
			os.Exit(1)
		}
	}
	

	// Cleanup goroutines (Potentially useless)
	request_ch <- -1
}