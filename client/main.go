package main

import (
	"fmt"
	"net/http"
)

func main() {

    msg := make(chan string)

    go printer(msg)

    for i := 0; i < 100; i++ {
        res, err := http.Get("http://localhost:8080")

        defer func() { res.Body.Close()}()

        if err != nil {
            fmt.Print(err)
        }
        msg <- res.Header.Get("Date")
    }

}

func printer(msg chan string) {
    for m := range msg {
        fmt.Println(m)
    }
}
