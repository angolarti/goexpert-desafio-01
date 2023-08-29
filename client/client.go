package main

import (
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	c := http.Client{Timeout: time.Millisecond * 300}
	resp, err := c.Get("http://localhost:8080/cotacao")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bid, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if string(bid) != "" {
		f, err := os.Create("cotacao.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = f.Write([]byte("Dólar: " + string(bid)))
		if err != nil {
			panic(err)
		}
	}
}
