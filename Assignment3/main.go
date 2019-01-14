package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

// type logWriter struct {
// }

func main() {
	fmt.Println(os.Args[1])
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// lw := logWriter{}
	// io.Copy(lw, file)
	io.Copy(os.Stdout, file)
}

// func (logWriter) Write(bs []byte) (int, error) {
// 	fmt.Println(string(bs))
// 	fmt.Println("size:", len(bs))
// 	return len(bs), nil

// }
