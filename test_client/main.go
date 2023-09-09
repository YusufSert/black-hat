package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

var rf = newRequestFlag("r", 100, "number of requests")

func main() {

	msg := make(chan time.Duration)
	var wg sync.WaitGroup

	for i := 0; i < *rf; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()
			start := time.Now()
			res, err := http.Get("http://localhost:8080")
			msg <- time.Since(start) // no race condition bc channels block

			defer res.Body.Close()

			if err != nil {
				fmt.Print(err)
			}
		}()
	}

	go func() {
		wg.Wait() // waits for signal to stop blocking
		close(msg)
	}()

	responseTimes := make([]float64, 0, 0)

	for m := range msg {
		responseTimes = append(responseTimes, m.Seconds())
		fmt.Println(len(responseTimes))
	}

	// mean
	var sum float64
	for _, n := range responseTimes {
		sum += n
	}
	mean := sum / 100

	// variance
	var variance float64
	var x float64
	for _, n := range responseTimes {
		x = n - mean
		variance += math.Pow(x, 2)
	}

	// square root of variance = standard deviation
	std := math.Pow(variance, 2)
	fmt.Println(std, mean)

}

////////////////////flags

type requestFlag struct {
	n int
}

func (f *requestFlag) Set(s string) error {
	var value int
	_, err := fmt.Sscanf(s, "%d", &value)
	if err != nil {
		return err
	}
	f.n = value
	return nil
}

func (f *requestFlag) String() string {
	return fmt.Sprint(f.n)
}

func newRequestFlag(name string, value int, usage string) *int {
	f := requestFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.n
}
