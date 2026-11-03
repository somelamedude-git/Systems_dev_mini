package main

import (
	"fmt"
	"math/rand"
)

type Node struct {
	value int
	next  []*Node
}

func InitNode(value int, MAX_LEVEL int) *Node {
	node := &Node{
		value: value,
		next:  make([]*Node, MAX_LEVEL),
	}

	return node
}

type SkipList struct {
	currLevel int
	maxLevel  int
	head      *Node
}

func InitSkipList(MAX_LEVEL int) *SkipList {
	head := InitNode(-1, MAX_LEVEL)
	NilNode := InitNode(100000, MAX_LEVEL)
	for i := 0; i < MAX_LEVEL; i++ {
		head.next[i] = NilNode
	}

	currLevel := 1

	sl := &SkipList{
		currLevel: currLevel,
		maxLevel:  MAX_LEVEL,
		head:      head,
	}
	return sl
}

func SearchKey(list *SkipList, key int) (keyPointer *Node, err error) {
	currLevel := list.currLevel
	currNode := list.head

	for currLevel > 0 {
		nextNode := currNode.next[currLevel-1]

		if nextNode.value < key {
			currNode = nextNode
		} else {
			currLevel--
		}
	}

	candidateNode := currNode.next[0]
	if candidateNode.value == key {
		return candidateNode, nil
	}

	return nil, fmt.Errorf("Key not found in the list")
}

func randomLevel(maxLevel int) int {
	p := 0.5
	newLevel := 1
	for rand.Float64() < p {
		newLevel++
		if newLevel >= maxLevel {
			break
		}
	}
	return newLevel
}

func InsertKey(list *SkipList, key int) {
	currLevel := list.currLevel
	update := make([]*Node, list.maxLevel)

	currNode := list.head

	for currLevel > 0 {
		nextNode := currNode.next[currLevel-1]

		if nextNode.value < key {
			currNode = nextNode
		} else {
			update[currLevel-1] = currNode
			currLevel--
		}
	}

	currLevel = list.currLevel
	newLevel := randomLevel(list.maxLevel)

	newNode := InitNode(key, newLevel)

	for i := 0; i < min(currLevel, newLevel); i++ {
		nextNode := update[i].next[i]
		newNode.next[i] = nextNode
		update[i].next[i] = newNode
	}

	for currLevel < newLevel {
		currLevel++

		nextNode := list.head.next[currLevel-1]
		newNode.next[currLevel-1] = nextNode
		list.head.next[currLevel-1] = newNode
	}

	list.currLevel = currLevel
}
