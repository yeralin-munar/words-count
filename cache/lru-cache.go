package cache

type LRUCache struct {
	cap            int
	m              map[string]*Node
	head, tail     *Node
	removeCallback func(string, int)
}

type Node struct {
	Key        string
	Val        int
	Next, Prev *Node
}

/*
	removeCallback needs to hook removes from cache
*/
func NewLRUCache(limit int, removeCallback func(string, int)) *LRUCache {
	m := make(map[string]*Node)
	head := Node{}
	tail := Node{}
	head.Next = &tail
	tail.Prev = &head

	return &LRUCache{
		cap:            limit,
		m:              m,
		head:           &head,
		tail:           &tail,
		removeCallback: removeCallback,
	}
}

func (l *LRUCache) add(node *Node) {
	l.m[node.Key] = node

	headNext := l.head.Next
	headNext.Prev = node

	node.Next = headNext
	node.Prev = l.head

	l.head.Next = node
}

/*
	node1			<->			node2		<->		node3
	next=node2				next=node3				next=...
	prev=...				prev=node1				prev=node2

	> remove
	next=node3										prev=node1

*/
func (l *LRUCache) remove(node *Node) {
	delete(l.m, node.Key)
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

func (l *LRUCache) Clear() {
	for _, node := range l.m {
		l.removeCallback(node.Key, node.Val)
		l.remove(node)
	}
}

func (l *LRUCache) Store(key string, val int) {
	if node, ok := l.m[key]; ok {
		l.remove(node)
	}

	if l.cap == len(l.m) {
		l.removeCallback(l.tail.Prev.Key, l.tail.Prev.Val)
		l.remove(l.tail.Prev)
	}

	node := &Node{
		Key: key,
		Val: val,
	}
	l.add(node)
}

func (l *LRUCache) Get(key string) int {
	node, ok := l.m[key]
	if !ok {
		return 0
	}

	l.remove(node)
	l.add(node)

	return node.Val
}
