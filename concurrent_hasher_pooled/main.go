package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
	"flag"
	"path/filepath"
)

type Result struct{
	path string
	hash [32]byte
	err error
}

func WalkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		errc <- filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.Type().IsRegular() {
				return nil
			}

			select {
			case paths <- path:
			case <-done:
				return fmt.Errorf("walk canceled")
			}
			return nil
		})
	}()

	return paths, errc
}

func hashFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Problem in opening file: ", err)
		return nil, err
	}

	defer file.Close()

	buf := make([]byte, 1024*64)
	hash := sha256.New()
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		hash.Write(buf[:n])
	}
	return hash.Sum(nil), nil
}

// functions will be going cncurrently, so range does not mean it will receive everything, it will receive what is available at the time of execution

func digester(done <-chan struct{}, paths <- chan string, c chan <- Result){
	for path := range paths{
		hashSum, err := hashFile(path)

		var h [32]byte

		if err == nil{
			h = [32]byte(hashSum)
		}

		select{
		case c <- Result{path, h, err}:
		case <-done:
			return
		}
	}
}

// c is a shared channel hence digester is not going to close it, the main function will close it when all the digester functions are done

func main(){
	flag.Parse()

	args := flag.Args()

	root := "."
	if len(args) > 0{
		root = args[0]
	}

	paths, errc := WalkFiles(make(chan struct{}), root)

	var wg sync.WaitGroup

	wg.Add(20)

	c:= make(chan Result)
	done := make(chan struct{})

	for i:=0; i<20; i++{
		go func(){
			digester(done, paths, c)
			wg.Done()
		}()
	}

	go func(){
		wg.Wait()
		close(c)
		close(done)
	}()

	if err := <-errc; err !=nil{
		fmt.Println("Error in walking the files: ", err)
		return
	}

	for r := range c{
		if r.err != nil{
			fmt.Println("Error in hashing the file: ", r.err)
			continue
		}
		fmt.Printf("%x %s\n", r.hash, r.path)
	}
}