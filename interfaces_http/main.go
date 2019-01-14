package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type logWriter struct {
}

func main() {
	resp, err := http.Get("http://google.com")
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	lw := logWriter{}
	io.Copy(lw, resp.Body)

}
func (logWriter) Write(bs []byte) (int, error) {
	fmt.Println(string(bs))
	fmt.Println("size:", len(bs))
	return len(bs), nil

}
