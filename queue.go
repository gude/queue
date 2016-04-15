package queue

import (
	"container/list"
	"sync"
)

type Queue struct {
	// Max Entries is the maximum length of queue
	MaxEntries int
	ll         *list.List                    //队列顺序
	e          map[interface{}]*list.Element //key 和 list ele的map  入队时，如果key已经存在，则更新list顺序使用
	d          map[*list.Element]interface{} //list ele 和 (key,value)的map 出队时，查询list ele对应的key使用，同时用来存储value

	mu sync.RWMutex
}

type Key interface{}

type entry struct {
	key   Key
	value interface{}
}

// New creates a new Queue
func New(maxEntries int) *Queue {
	return &Queue{
		MaxEntries: maxEntries,
		ll:         list.New(),
		e:          make(map[interface{}]*list.Element),
		d:          make(map[*list.Element]interface{}),
	}
}

// Enqueue
func (c *Queue) Enqueue(key Key, value interface{}) (ret bool) {
	ret = true

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.e == nil {
		c.e = make(map[interface{}]*list.Element)
		c.ll = list.New()
		c.d = make(map[*list.Element]interface{})
	}

	//如果已经存在key,覆盖之，并将位置移动到队尾
	if ee, ok := c.e[key]; ok {
		c.ll.MoveToBack(ee)
		c.d[ee] = &entry{key, value}
		return
	}

	//队列已满，不插入
	if c.MaxEntries != 0 && c.Len() >= c.MaxEntries {
		ret = false
		return
	}

	//正常添加
	ele := c.ll.PushBack(key)
	c.e[key] = ele
	c.d[ele] = &entry{key, value}

	//如果超过长度则删除之
	if c.MaxEntries != 0 && c.Len() > c.MaxEntries {
		c.removeOldest()
	}

	return
}

//Dequeue
func (c *Queue) Dequeue() (key Key, value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.e == nil {
		return
	}

	ee := c.ll.Front()
	if d, hit := c.d[ee]; hit {
		value = d.(*entry).value
		key = d.(*entry).key
		c.removeElement(ee)
		return key, value, true
	}

	return
}

//Push Front
func (c *Queue) PushFront(key Key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.e == nil {
		c.e = make(map[interface{}]*list.Element)
		c.d = make(map[*list.Element]interface{})
		c.ll = list.New()

	}

	if _, ok := c.e[key]; ok {
		return
	}
	ele := c.ll.PushFront(key)
	c.e[key] = ele
	c.d[ele] = &entry{key, value}

	if c.MaxEntries != 0 && c.Len() > c.MaxEntries {
		c.removeOldest()
	}
	return
}

//Remove oldest
func (c *Queue) removeOldest() {
	if c.e == nil {
		return
	}
	ele := c.ll.Front()
	if ele != nil {
		c.removeElement(ele)
	}
}

//remove element
func (c *Queue) removeElement(ele *list.Element) {
	c.ll.Remove(ele)
	key := c.d[ele].(*entry).key
	delete(c.e, key)
	delete(c.d, ele)
}

// get len
func (c *Queue) Len() int {
	if c.e == nil {
		return 0
	}
	return c.ll.Len()
}
