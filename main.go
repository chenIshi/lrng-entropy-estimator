package main

import (
	"fmt"
)

// EVal for "Entropy Value"
type EVal struct{
	evaluationType string
	value float32
}

// stimuate entropy pools in lrng, i.e. input/output pool
func entropy_pool(request_ch chan bool, entropy_ch chan int) {

}

func lrng_eval_3(entropy_ch chan int, response_ch chan EVal) {

}

func differential_eval(entropy_ch chan int, response_ch chan EVal) {

}

func main() {
	// user configuration

	request_ch := make(chan bool)
	entropy_ch := make(chan int)
	response_ch := make(chan EVal)
	go entropy_pool(request_ch, entropy_ch)
	go lrng_eval_3(entropy_ch, response_ch)
	go differential_eval(entropy_ch, response_ch)
	fmt.Println("Hello")
}