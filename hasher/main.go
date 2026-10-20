package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
)

func hashFile(r io.Reader) []byte {
	buf := make([]byte, 64*1024)
	h := sha256.New()
	for {
		n, err := r.Read(buf)
		h.Write(buf[:n])

		if err == io.EOF {
			break
		}

		if err != nil {
			return []byte(err.Error())
		}
	}

	return h.Sum(nil)
}

func main() {
	flag.Parse()

	args := flag.Args()

	for _, name := range args {
		file, err := os.Open(name)
		if err != nil {
			fmt.Println(err)
			continue
		}

		hashSum := hashFile(file)

		fmt.Printf("%x  %s\n", hashSum, name)
	}
}
