package db

import (
	"fmt"
)

// LRU represents the LRU eviction policy using a doubly linked list.
type LRU struct {
	Left  *Node // Least recently used node
	Right *Node // Most recently used node
}

func NewLRU() *LRU {
	left := &Node{}
	right := &Node{}
	left.Next, right.Prev = right, left
	return &LRU{Left: left, Right: right}
}

// Insert adds a node to the right end of the LRU list.
func (lru *LRU) Insert(node *Node) {
	prev := lru.Right.Prev
	prev.Next, node.Prev = node, prev
	node.Next, lru.Right.Prev = lru.Right, node
}

// Remove deletes a node from the LRU list.
func (lru *LRU) Remove(node *Node) {
	prev, next := node.Prev, node.Next
	prev.Next, next.Prev = next, prev
}

// Evict removes the least recently used (LRU) node from the cache.
func (lru *LRU) Evict(cache *Cache) {
	// Evict the leftmost node (Least Recently Used)
	lruNode := lru.Left.Next
	if lruNode == lru.Right {
		return // No nodes to evict
	}
	lru.Remove(lruNode)
	delete(cache.Data, lruNode.Key)
	fmt.Println("Evicted entry:", "Key:", lruNode.Key, "Value:", lruNode.Value)
}
