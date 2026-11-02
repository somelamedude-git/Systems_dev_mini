package main

import "fmt"

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
	NilNode := InitNode(10000000000, MAX_LEVEL)
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

	currVal := list.head.value
	currNode := list.head

	var node *Node

	for {
		if currLevel == 1 && currVal > key {
			break
		}

		if currVal == key {
			return node, nil
		}

		nextNode := currNode.next[currLevel-1]

		if nextNode.value <= key {
			currNode = nextNode
		} else {
			currLevel--
		}
	}

	return nil, fmt.Errorf("Key not found in the list")
}
