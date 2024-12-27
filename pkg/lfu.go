package db 

type Node struct {
	Key        int
	Prev, Next *Node
}

type LinkedList struct {
	Left, Right *Node
	NodeMap     map[int]*Node
}

func NewLinkedList() *LinkedList {
	left := &Node{}
	right := &Node{}
	left.Next = right
	right.Prev = left
	return &LinkedList{
		Left:    left,
		Right:   right,
		NodeMap: make(map[int]*Node),
	}
}

func (ll *LinkedList) Length() int {
	return len(ll.NodeMap)
}

func (ll *LinkedList) PushRight(key int) {
	node := &Node{Key: key, Prev: ll.Right.Prev, Next: ll.Right}
	ll.Right.Prev.Next = node
	ll.Right.Prev = node
	ll.NodeMap[key] = node
}

func (ll *LinkedList) Pop(key int) {
	if node, exists := ll.NodeMap[key]; exists {
		prev, next := node.Prev, node.Next
		prev.Next = next
		next.Prev = prev
		delete(ll.NodeMap, key)
	}
}

func (ll *LinkedList) PopLeft() int {
	if ll.Left.Next == ll.Right {
		return -1 // Empty list
	}
	node := ll.Left.Next
	ll.Pop(node.Key)
	return node.Key
}

type LFUCache struct {
	cap      int
	lfuCnt   int
	countMap map[int]int       // Key -> Frequency
	cache    map[int]int       // Key -> Value
	listMap  map[int]*LinkedList // Frequency -> LinkedList
}

func Constructor(capacity int) LFUCache {
	return LFUCache{
		cap:      capacity,
		lfuCnt:   0,
		countMap: make(map[int]int),
		cache:    make(map[int]int),
		listMap:  make(map[int]*LinkedList),
	}
}

func (lfu *LFUCache) Counter(key int) {
	cnt := lfu.countMap[key]
	lfu.listMap[cnt].Pop(key)
	lfu.countMap[key]++
	cnt++

	if _, exists := lfu.listMap[cnt]; !exists {
		lfu.listMap[cnt] = NewLinkedList()
	}
	lfu.listMap[cnt].PushRight(key)

	if cnt-1 == lfu.lfuCnt && lfu.listMap[cnt-1].Length() == 0 {
		lfu.lfuCnt++
	}
}

func (lfu *LFUCache) Get(key int) int {
	if val, exists := lfu.cache[key]; exists {
		lfu.Counter(key)
		return val
	}
	return -1
}

func (lfu *LFUCache) Put(key int, value int) {
	if lfu.cap == 0 {
		return
	}

	if _, exists := lfu.cache[key]; exists {
		lfu.cache[key] = value
		lfu.Counter(key)
	} else {
		if len(lfu.cache) == lfu.cap {
			lfuKey := lfu.listMap[lfu.lfuCnt].PopLeft()
			delete(lfu.cache, lfuKey)
			delete(lfu.countMap, lfuKey)
		}

		lfu.cache[key] = value
		lfu.countMap[key] = 1
		if _, exists := lfu.listMap[1]; !exists {
			lfu.listMap[1] = NewLinkedList()
		}
		lfu.listMap[1].PushRight(key)
		lfu.lfuCnt = 1
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
