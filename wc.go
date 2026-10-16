package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func Count(r io.Reader) (int, int, int, error) {
	buf := make([]byte, 64*1024)
	bytes := 0
	lines := 0
	words := 0

	var inWord bool
	inWord = false
	for {
		n, err := r.Read(buf)

		bytes += n
		for _, b := range buf[:n] {
			if b == '\n' {
				lines++
			}

			if b == ' ' || b == '\t' || b == '\n' || b == '\v' || b == '\r' || b == '\f' {
				inWord = false
			} else {
				if !inWord {
					inWord = true
					words++
				}
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				return 0, 0, 0, err
			}
		}
	}
	return lines, words, bytes, nil
}

func CountFile(name string) (int, int, int, error) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return 0, 0, 0, err
	}

	return Count(file)
}

func main() {
	lines := flag.Bool("l", false, "lines")
	words := flag.Bool("w", false, "words")
	bytes := flag.Bool("b", false, "bytes")

	flag.Parse()

	args := flag.Args() // will give the file names

	if !*lines && !*bytes && !*words {
		*lines, *bytes, *words = true, true, true
	}

	if len(args) == 0 {
		l, w, b, err := Count(os.Stdin)
		if err != nil {
			fmt.Println(err)
			return
		}

		if *lines {
			fmt.Println(l)
		}

		if *words {
			fmt.Println(w)
		}

		if *bytes {
			fmt.Println(b)
		}
	}

	for _, name := range args {
		file, err := os.Open(name)

		if err != nil {
			fmt.Println(err)
			continue
		}

		defer file.Close()

		l, w, b, err := Count(file)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if *lines {
			fmt.Println(l)
		}

		if *words {
			fmt.Println(w)
		}

		if *bytes {
			fmt.Println(b)
		}

		fmt.Println(name)
	}
}
