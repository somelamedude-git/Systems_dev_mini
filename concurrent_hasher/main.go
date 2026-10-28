package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type result struct {
	path string
	hash [32]byte
	err  error
}

func hashFile(path string) ([]byte, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf := make([]byte, 1024*64)
	hash := sha256.New()
	for {
		n, err := file.Read(buf)
		hash.Write(buf[:n])

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	return hash.Sum(nil), nil
}

func hashHelperConc(root string, done <-chan struct{}) (<-chan result, <-chan error) {
	c := make(chan result)
	errc := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup

		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, er error) error {
			if er != nil {
				fmt.Println(er)
				return nil
			}

			if !d.Type().IsRegular() {
				return nil
			}

			wg.Add(1)

			go func() {
				defer wg.Done()

				hashSum, erro := hashFile(path)

				var h [32]byte

				if erro == nil {
					copy(h[:], hashSum)
				}

				select {
				case c <- result{path, h, erro}:
				case <-done:
				}
			}()

			select {
			case <-done:
				return fmt.Errorf("walk canceled")
			default:
				return nil
			}
		})

		go func() {
			wg.Wait()
			close(c)
		}()

		errc <- err
	}()

	return c, errc
}

func actualHash(root string) (map[string][32]byte, error) {
	done := make(chan struct{})
	defer close(done)

	c, errc := hashHelperConc(root, done)

	mp := make(map[string][32]byte)
	for r := range c {
		if r.err != nil {
			fmt.Println(r.err)
			continue
		}

		mp[r.path] = r.hash
	}

	if err := <-errc; err != nil {
		return nil, err
	}

	return mp, nil
}

func main() {
	flag.Parse()

	args := flag.Args()
	root := "."

	if len(args) > 0 {
		root = args[0]
	}

	map_, err := actualHash(root)

	if err != nil {
		fmt.Println(err)
		return
	}

	for key, val := range map_ {
		fmt.Printf("%x  %s\n", val, key)
	}
}
