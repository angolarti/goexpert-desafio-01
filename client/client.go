package main

import (
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	c := http.Client{Timeout: time.Second * 300}
	resp, err := c.Get("http://localhost:8080/cotacao")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bid, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("./cotacao.txt")

	_, err = f.Write([]byte("Dólar: " + string(bid)))
	if err != nil {
		panic(err)
	}
	f.Close()
}
