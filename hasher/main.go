package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"io/fs"
	"path/filepath"
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

	root := "."
	if len(args) > 0{
		root = args[0]
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error{
		if err != nil{
			fmt.Fprintln(os.Stderr, err)
			return nil
		}
		if !d.Type().IsRegular(){
			return nil
		}

		sum, er := hashFile(path)
		if er != nil{
			fmt.Println(er)
			return nil
		}

		fmt.Printf("%x  %s\n", sum, path)

		return nil
	})


	if err != nil{
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
