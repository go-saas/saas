package saas

import (
	"container/list"
	"fmt"
	"io"
	"sync"
)

// Cache is used a LRU (Least recently used) cache replacement policy. adapted from https://github.com/Code-Hex/go-generics-cache/blob/main/policy/lru/lru.go
//
// Discards the least recently used items first. This algorithm requires
// keeping track of what was used when, which is expensive if one wants
// to make sure the algorithm always discards the least recently used item.
type Cache[K comparable, V io.Closer] struct {
	cap   int
	list  *list.List
	items map[K]*list.Element
	mu    sync.Mutex
}

type entry[K comparable, V any] struct {
	key K
	val V
}

// Option is an option for LRU cache.
type Option func(*options)

type options struct {
	capacity int
}

func newOptions() *options {
	return &options{
		capacity: 128,
	}
}

// WithCapacity is an option to set cache capacity.
func WithCapacity(cap int) Option {
	return func(o *options) {
		o.capacity = cap
	}
}

// NewCache creates a new thread safe LRU cache whose capacity is the default size (128).
func NewCache[K comparable, V io.Closer](opts ...Option) *Cache[K, V] {
	o := newOptions()
	for _, optFunc := range opts {
		optFunc(o)
	}
	return &Cache[K, V]{
		cap:   o.capacity,
		list:  list.New(),
		items: make(map[K]*list.Element, o.capacity),
	}
}

// Get looks up a key's value from the cache.
func (c *Cache[K, V]) Get(key K) (zero V, _ bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.get(key)
}

func (c *Cache[K, V]) get(key K) (zero V, _ bool) {
	e, ok := c.items[key]
	if !ok {
		return
	}
	// updates cache order
	c.list.MoveToFront(e)
	return e.Value.(*entry[K, V]).val, true
}

func (c *Cache[K, V]) Set(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, val)
}

// Set sets a value to the cache with key. replacing any existing value.
func (c *Cache[K, V]) set(key K, val V) {

	if e, ok := c.items[key]; ok {
		// updates cache order
		c.list.MoveToFront(e)
		entry := e.Value.(*entry[K, V])
		entry.val = val
		return
	}

	newEntry := &entry[K, V]{
		key: key,
		val: val,
	}
	e := c.list.PushFront(newEntry)
	c.items[key] = e

	if c.list.Len() > c.cap {
		c.deleteOldest()
	}
}

// GetOrSet combine Get and Set
func (c *Cache[K, V]) GetOrSet(key K, factory func() (V, error)) (zero V, set bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.get(key); ok {
		return v, false, nil
	}
	//use factory
	v, err := factory()
	if err != nil {
		return zero, false, err
	}
	c.set(key, v)
	return v, true, nil
}

// Keys returns the keys of the cache. the order is from oldest to newest.
func (c *Cache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.keys()
}

func (c *Cache[K, V]) keys() []K {
	keys := make([]K, 0, len(c.items))
	for ent := c.list.Back(); ent != nil; ent = ent.Prev() {
		entry := ent.Value.(*entry[K, V])
		keys = append(keys, entry.key)
	}
	return keys
}

// Len returns the number of items in the cache.
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.list.Len()
}

// Delete deletes the item with provided key from the cache.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deleteKey(key)
}

func (c *Cache[K, V]) deleteKey(key K) error {
	if e, ok := c.items[key]; ok {
		return c.delete(e)
	}
	return nil
}

// Flush delete all items
func (c *Cache[K, V]) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error
	for _, k := range c.keys() {
		nerr := c.deleteKey(k)
		if nerr != nil {
			if err == nil {
				err = nerr
			} else {
				err = fmt.Errorf("%w; ", err)
			}
		}
	}
	return err
}

func (c *Cache[K, V]) deleteOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()
	e := c.list.Back()
	c.delete(e)
}

func (c *Cache[K, V]) delete(e *list.Element) error {
	c.list.Remove(e)
	entry := e.Value.(*entry[K, V])
	delete(c.items, entry.key)

	return entry.val.Close()
}
