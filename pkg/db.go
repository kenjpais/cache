package db

import (
	"fmt"
	"sync"
)

// Node represents a doubly linked list node.
type Node struct {
	Key, Value string
	Prev, Next *Node
}

// EvictionPolicy is the strategy interface for different eviction policies.
type EvictionPolicy interface {
	Insert(node *Node)
	Remove(node *Node)
	Evict(cache *Cache)
}

// Cache represents the in-memory key-value store with an eviction policy.
type Cache struct {
	Capacity int
	Data     map[string]*Node
	Policy   EvictionPolicy
	Mutex    sync.Mutex
}

// NewCache initializes the cache with the given capacity and eviction policy.
func NewCache(capacity int, policy EvictionPolicy) *Cache {
	return &Cache{
		Capacity: capacity,
		Data:     make(map[string]*Node),
		Policy:   policy,
	}
}

// Get retrieves a value from the cache and marks it as recently used.
func (cache *Cache) Get(key string) (string, bool) {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	node, exists := cache.Data[key]
	if !exists {
		return "", false
	}

	// Move the node to the most recently used position.
	cache.Policy.Remove(node)
	cache.Policy.Insert(node)

	return node.Value, true
}

// AddEntry adds or updates a key-value pair in the cache.
func (cache *Cache) AddEntry(key, value string) {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	if node, exists := cache.Data[key]; exists {
		node.Value = value
		cache.Policy.Remove(node)
		cache.Policy.Insert(node)
		fmt.Println("Updated entry:", "Key:", key, "Value:", value)
		return
	}

	newNode := &Node{Key: key, Value: value}
	cache.Data[key] = newNode
	cache.Policy.Insert(newNode)

	// Evict the least recently used entry if the cache exceeds capacity.
	if len(cache.Data) > cache.Capacity {
		cache.Policy.Evict(cache)
	}

	fmt.Println("Added entry:", "Key:", key, "Value:", value)
}

// Delete removes a key-value pair from the cache.
func (cache *Cache) Delete(key string) {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	node, exists := cache.Data[key]
	if !exists {
		fmt.Println("Key does not exist:", key)
		return
	}

	cache.Policy.Remove(node)
	delete(cache.Data, key)
	fmt.Println("Deleted entry:", "Key:", key)
}

// Clear removes all entries from the cache.
func (cache *Cache) Clear() {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	cache.Data = make(map[string]*Node)
	cache.Policy = NewLRU() // Reset the eviction policy
	fmt.Println("Cleared all entries in the cache")
}

// ChangePolicy dynamically changes the eviction policy.
func (cache *Cache) ChangePolicy(policy EvictionPolicy) {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	// Apply the new policy while keeping existing data.
	cache.Policy = policy
	fmt.Println("Changed eviction policy")
}
