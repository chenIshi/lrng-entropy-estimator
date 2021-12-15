package rng

import (
	"math/rand"
	"fmt"
	"time"
	"os"
	"util"
)

type RandFunc func(num_chan chan int)

// Stimuate entropy pools in lrng, i.e. input/output pool 
// According user configuration, it would have to produce different level of "randomness" 
// Also, there could different method to construct a "not-that-random" number
func Entropy_pool(request_ch chan util.Num_msg, entropy_chs [util.EVAL_METHEOD_COUNT]chan util.Eval_msg, rng_type util.Rng_t) {
	if int(rng_type) >= util.RNG_METHOD_COUNT {
		fmt.Print("Using not supported RNG method")
		os.Exit(1)
	}
	rand_funcs := register_rngs()

	num_chan := make(chan int)
	go rand_funcs[int(rng_type)](num_chan)
	for {
		req := <-request_ch
		// if the main process want to shut down the entropy collection pool 
		if req.Val < 0 {
			// it should also propagate such signal to evaluators before shut down
			for i:=0; i<util.EVAL_METHEOD_COUNT; i++ {
				entropy_chs[i] <- util.Eval_msg{Idx: req.Idx, Val: -1, Eval: util.EVAL_LRNG3}
			}
			return
		} else {
			if int(req.Rng) >= util.RNG_METHOD_COUNT {
				fmt.Print("Using not supported RNG method")
				continue
			}
			num_chan <- req.Val

			rng_timer := time.NewTimer(time.Duration(100 * time.Millisecond))
			var rand_val int
			select {
			case num := <- num_chan:
				rand_val = num
			case <- rng_timer.C:
				fmt.Fprintf(os.Stderr, "error: timeout in main at rng\n")
				os.Exit(1)
			}
			
			for i:=0; i<util.EVAL_METHEOD_COUNT; i++ {
				entropy_chs[i] <- util.Eval_msg{Idx: req.Idx, Val: rand_val, Eval: util.EVAL_LRNG3}
			}
		}
	}
}

func register_rngs() []RandFunc {
	var rng_handler []RandFunc
	rng_handler = append(rng_handler, Rand_buildin_pseudo)
	rng_handler = append(rng_handler, Rand_rule30)
	return rng_handler
}

func Rand_buildin_pseudo(num_chan chan int) {
	// As for common practice, most RNG requires a true random seed to achieve unpredictability 
	// However, here we actually don't care about such attribute 
	// Also, using a fixed random seed help us reproduce the experiment results if desired 
	rand.Seed(80)

	for {
		timer := time.NewTimer(time.Duration(80 * time.Millisecond))
		select {
		case req_max := <- num_chan:
			num_chan <- rand.Intn(req_max)
		case <- timer.C:
			fmt.Fprintf(os.Stderr, "error: timeout in main at build-in rng\n")
			os.Exit(1)
		}
	}
}

func Rand_rule30(num_chan chan int) {
	for {
		timer := time.NewTimer(time.Duration(80 * time.Millisecond))
		select {
		case req_max := <- num_chan:
			num_chan <- req_max
		case <- timer.C:
			fmt.Fprintf(os.Stderr, "error: timeout in main at build-in rng\n")
			os.Exit(1)
		}
	}
}