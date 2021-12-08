package rng

import (
	"math/rand"
	"util"
)

const eval_method_count = 2

// Stimuate entropy pools in lrng, i.e. input/output pool 
// According user configuration, it would have to produce different level of "randomness" 
// Also, there could different method to construct a "not-that-random" number
func Entropy_pool(request_ch chan util.Num_msg, entropy_chs [eval_method_count]chan util.Eval_msg) {
	// As for common practice, most RNG requires a true random seed to achieve unpredictability 
	// However, here we actually don't care about such attribute 
	// Also, using a fixed random seed help us reproduce the experiment results if desired 
	rand.Seed(80)
	for {
		req := <-request_ch
		// if the main process want to shut down the entropy collection pool 
		if req.Val < 0 {
			// it should also propagate such signal to evaluators before shut down
			for i:=0; i<eval_method_count; i++ {
				entropy_chs[i] <- util.Eval_msg{Idx: req.Idx, Val: -1, Eval_t: util.EVAL_LRNG3}
			}
			return
		} else {
			rng_num :=  rand.Intn(req.Val)
			for i:=0; i<eval_method_count; i++ {
				entropy_chs[i] <- util.Eval_msg{Idx: req.Idx, Val: rng_num, Eval_t: util.EVAL_LRNG3}
			}
		}
	}
}