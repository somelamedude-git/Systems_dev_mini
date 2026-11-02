package main

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
