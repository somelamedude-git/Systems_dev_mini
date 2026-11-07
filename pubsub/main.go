package main

import (
	"sync"
)

type Topic struct {
	name   string
	mu     sync.Mutex
	subs   []chan string
	closed bool
}
