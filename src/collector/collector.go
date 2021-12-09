package collector

import (
	"fmt"
	"encoding/csv"
	"os"
	"log"
	"time"
	"util"
)

func Collector(response_ch chan util.Entropy_msg, ctrl_ch chan util.Ctrl_msg) {
	// TODO: Improve with a sliding-window to buffer response counts 
	// This is used to prevent some evaluation method is way faster than 
	// others, causing congested channel 
	comm_resp_cnt := 0

	var entropies_from_sources [util.EVAL_METHEOD_COUNT][] float64
	for {
		collector_timer := time.NewTimer(time.Duration(80 * time.Millisecond))
		select {
		case resp:= <- response_ch:
			entropies_from_sources[resp.Eval_t] = append(entropies_from_sources[resp.Eval_t], resp.Val)
			comm_resp_cnt ++
			if comm_resp_cnt >= util.EVAL_METHEOD_COUNT {
				ctrl_ch <- util.Ctrl_msg{Idx: resp.Idx, Signal: util.CTRL_COMM_RESP}
				comm_resp_cnt = 0
			}
		case ctrl_sig := <- ctrl_ch:
			// jump to info dumping stage then shut down itself
			if ctrl_sig.Signal == util.CTRL_OUT_REQ {
				for i := range entropies_from_sources {
					sum := 0.0
					for j:=0;j<len(entropies_from_sources[i]);j++ {
						sum += entropies_from_sources[i][j]
					}
					avg := sum / float64(len(entropies_from_sources[i]))
					if i == util.EVAL_LRNG3 {
						fmt.Println("LRNG entropy estimation avg = ", avg)
					}
				}
				ctrl_ch <- util.Ctrl_msg{Idx: ctrl_sig.Idx, Signal: util.CTRL_OUT_RESP}
				return
			}
		case <- collector_timer.C:
			fmt.Fprintf(os.Stderr, "error: timeout in collector\n")
			os.Exit(1)
		}
	}
}

func dump_csv(filename string, entropies_from_sources [util.EVAL_METHEOD_COUNT][]float64) {
	f, err := os.Create(filename)
	defer f.Close()

	if err != nil {
		log.Fatalln("error writing ", filename, " due to ", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
}

func plot(figname string, entropies_from_sources [util.EVAL_METHEOD_COUNT][]float64) {
	
}