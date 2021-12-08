package main

import (
	"fmt"
	"time"
	"os"
	"flag"
	"estimator"
	"rng"
)

const eval_method_count = 2

var testscale = flag.Int("n", 5000, "Amount of evaulated random number.")
var maxrng = flag.Int("max", 50, "Maximum of random generated numbers.")

func main() {
	// User configuration
	flag.Parse()
	test_scale := *testscale
	max_rng := *maxrng

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
	go rng.Entropy_pool(request_ch, entropy_chs)
	go estimator.Lrng_eval_3(entropy_chs[0], response_ch, demo_chs[0])
	go estimator.Differential_eval(entropy_chs[1], response_ch, demo_chs[1])

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