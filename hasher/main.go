package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
)

func hashFile(name string) ([]byte, error) {
	file, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf := make([]byte, 64*1024)
	h := sha256.New()
	for {
		n, err := file.Read(buf)
		h.Write(buf[:n])

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	return h.Sum(nil), nil
}

func main() {
	flag.Parse()

	args := flag.Args()

	for _, name := range args {
		hashSum, err := hashFile(name)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%x  %s\n", hashSum, name)
	}
}
