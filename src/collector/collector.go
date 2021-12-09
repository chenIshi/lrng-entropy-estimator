package collector

import (
	"fmt"
	"encoding/csv"
	"os"
	"log"
	"strconv"
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
				for i, entropies := range entropies_from_sources {
					sum := 0.0
					for j:=0;j<len(entropies);j++ {
						sum += entropies[j]
					}
					avg := sum / float64(len(entropies))
					if i == util.EVAL_LRNG3 {
						fmt.Println("LRNG entropy estimation avg = ", avg)
					}
				}
				dump_csv("eval.csv", entropies_from_sources)
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
	checkError("Can't creat file", err)

	writer := csv.NewWriter(f)
	defer writer.Flush()

	err = writer.Write([]string{"idx", "val", "type"})
	checkError("Can't write csv header", err)
	for i, entropies := range entropies_from_sources {
		for j, entropy := range entropies {
			entropy_in_str := fmt.Sprintf("%.2f", entropy)
			err = writer.Write([]string{strconv.Itoa(j), entropy_in_str, strconv.Itoa(i)})
			checkError("Can't write csv content", err)
		}
	}
}

func checkError(message string, err error) {
    if err != nil {
        log.Fatalln(message, err)
		os.Exit(1)
    }
}