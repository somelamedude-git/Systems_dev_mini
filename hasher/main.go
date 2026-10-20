package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func main() {
	name := os.Args[1]
	file, err := os.Open(name)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	h := sha256.New()
	buf := make([]byte, 64*1024)

	for {
		n, err := file.Read(buf)
		h.Write(buf[:n])

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Printf("%x  %s\n", h.Sum(nil), name)
}
