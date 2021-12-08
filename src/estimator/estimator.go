package estimator

import (
	"math"
	"util"
)

func Lrng_eval_3(entropy_ch chan util.Eval_msg, response_ch chan util.Entropy_msg) {
	// LRNG entropy estimation requires to keep the previous states 
	// By default, it would require a 3 step of differential, that is 4 states in total
	pre_states := []int{-1, -1,-1, -1}
	// As stated above, LRNG requires to collect previous states before estimation 
	// In this case, it would start actual estimation until the fourth round
	eval_round := 0

	for {
		rng_num := <- entropy_ch
		if rng_num.Val < 0 {
			return
		} else {
			// Push the generated rng number into pre_states 
			// Also kick out the eldest state (that is, FIFO)
			pre_states = append(pre_states[1:], rng_num.Val)
			if eval_round >= 3 {
				// LRNF entropy estimation: see [The Linux Pseudorandom Number Generator Revisited](https://eprint.iacr.org/2012/251.pdf) for details 
				delta1 := math.Abs(float64(pre_states[3] - pre_states[2]))
				delta2 := math.Abs(float64(pre_states[3] - 2 * pre_states[2] + pre_states[1]))
				delta3 := math.Abs(float64(pre_states[3] - 3 * pre_states[2] + 3 * pre_states[1] - pre_states[0]))
				delta := util.Min([]float64{delta1, delta2, delta3})
				entropy := 0.0
				if delta > (1 << 12) {
					entropy = 11
				} else if delta > 2{
					entropy = math.Log2(delta)
				}
				response_ch <- util.Entropy_msg{Idx: rng_num.Idx, Val: entropy, Eval_t: util.EVAL_LRNG3}
			} else {
				response_ch <- util.Entropy_msg{Idx: rng_num.Idx, Val: 0, Eval_t: util.EVAL_LRNG3}
			}
			eval_round ++
		}
		/*
		select {
		case rng_num := <- entropy_ch:
			if rng_num.Val < 0 {
				return
			} else {
				// Push the generated rng number into pre_states 
				// Also kick out the eldest state (that is, FIFO)
				pre_states = append(pre_states[1:], rng_num.Val)
				if eval_round >= 3 {
					// LRNF entropy estimation: see [The Linux Pseudorandom Number Generator Revisited](https://eprint.iacr.org/2012/251.pdf) for details 
					delta1 := math.Abs(float64(pre_states[3] - pre_states[2]))
					delta2 := math.Abs(float64(pre_states[3] - 2 * pre_states[2] + pre_states[1]))
					delta3 := math.Abs(float64(pre_states[3] - 3 * pre_states[2] + 3 * pre_states[1] - pre_states[0]))
					delta := util.Min([]float64{delta1, delta2, delta3})
					entropy := 0.0
					if delta > (1 << 12) {
						entropy = 11
					} else if delta > 2{
						entropy = math.Log2(delta)
					}
					response_ch <- util.Entropy_msg{Idx: rng_num.Idx, Val: entropy, Eval_t: util.EVAL_LRNG3}
				} else {
					response_ch <- util.Entropy_msg{Idx: rng_num.Idx, Val: 0, Eval_t: util.EVAL_LRNG3}
				}
				eval_round ++
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
		} */		
	}
}

func Differential_eval(entropy_ch chan util.Eval_msg, response_ch chan util.Entropy_msg) {
	for {
		rng_num := <- entropy_ch
		if rng_num.Val < 0 {
			return
		} else {
			response_ch <- util.Entropy_msg{Idx: rng_num.Idx, Val: 0, Eval_t: util.EVAL_DIFF}
		}
		/*
		select {
		case rng_num := <- entropy_ch:
			if rng_num.Val < 0 {
				return
			} else {
				response_ch <- util.Entropy_msg{Idx: rng_num.Idx, Val: 0, Eval_t: util.EVAL_DIFF}
			}
		case <-demo_ch:
			fmt.Println("Differential entropy estimation not implemented yet!")
			demo_ch <- false
		}	
		*/	
	}
}