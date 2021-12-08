package main

import (
	"fmt"
	"time"
	"os"
	"flag"
	"estimator"
	"collector"
	"rng"
	"util"
)



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
	request_ch := make(chan util.Num_msg, 10)
	// requestX_ch: channel between "entropy_pool" and "evaluatorX"
	// transmit the generated "random number" 
	// **Sending a negative signal means to shut down the goroutine**
	var entropy_chs [util.EVAL_METHEOD_COUNT]chan util.Eval_msg
	for i := range entropy_chs {
		entropy_chs[i] = make(chan util.Eval_msg, 10)
	}
	// request_ch: channel between "evaluator" and "main"
	// transmit the estimated randomness/entropy from the random number sequence
	response_ch := make(chan util.Entropy_msg, 10)

	var demo_chs [util.EVAL_METHEOD_COUNT]chan bool
	for i := range demo_chs {
		demo_chs[i] = make(chan bool)
	}

	ctrl_ch := make(chan util.Ctrl_msg, 10)

	go rng.Entropy_pool(request_ch, entropy_chs)
	go estimator.Lrng_eval_3(entropy_chs[0], response_ch)
	go estimator.Differential_eval(entropy_chs[1], response_ch)
	go collector.Collector(response_ch, ctrl_ch)

	for i:=0; i<test_scale; i++ {
		request_ch <- util.Num_msg{Idx:i, Val:max_rng, Rng_t: util.RNG_UNI}
		
		lrng_timer := time.NewTimer(time.Duration(160 * time.Millisecond))
		select {
		case <- ctrl_ch:
		case <- lrng_timer.C:
			fmt.Fprintf(os.Stderr, "error: timeout in main\n")
			os.Exit(1)
		}
	}

	// Output evaluation results using different estimation approaches 
	demo_timer := time.NewTimer(time.Duration(100 * time.Millisecond))
	ctrl_ch <- util.Ctrl_msg{Idx: 0, Signal: util.CTRL_OUT_REQ}
	select {
	case <- ctrl_ch:
	case <- demo_timer.C:
		fmt.Fprintf(os.Stderr, "error: timeout in main at demoing LRNG\n")
		os.Exit(1)
	}

	/*
	for k:=0; k<util.EVAL_METHEOD_COUNT; k++ {
		demo_chs[k] <- true
		select {
		case <- demo_chs[k]:
		case <- demo_timer.C:
			fmt.Fprintf(os.Stderr, "error: timeout in main at demoing LRNG\n")
			os.Exit(1)
		}
	}
	*/
	

	// Cleanup goroutines (Potentially useless)
	request_ch <- util.Num_msg{Idx:1, Val:-1, Rng_t: util.RNG_UNI}
}