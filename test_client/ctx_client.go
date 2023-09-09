package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var duration = flag.Duration("d", 5*time.Second, "request-response duration")

//var requestNubmer = newRequestFlag("r", 5, "total request")

func main() {
	// The client has five seconds to connect to the web server, send the request, read the response
	// header, and pass the response to your code. You then have the remainder of the five seconds to
	// read the response body
	msg := make(chan time.Duration)

	for i := 0; i < 100; i++ {

		go func(<-chan time.Duration) {
			ctx, cancel := context.WithTimeout(context.Background(), *duration)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
			if err != nil {
				log.Fatal(err)
			}
			req.Close = true

			start := time.Now()
			resp, err := http.DefaultClient.Do(req)
			rd := time.Since(start)

			msg <- rd

			if err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					log.Fatal(err)
				}
				return
			}
			_ = resp.Body.Close()
		}(msg)
	}

	for m := range msg {
		fmt.Println(m)
	}
}

//By default, Go’s HTTP client maintains the underlying TCP connection to a
//web server after reading its response unless explicitly told to disconnect by
//the server. Although this is desirable behavior for most use cases because it
//allows you to use the same TCP connection for multiple requests
