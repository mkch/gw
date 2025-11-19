package objectmap

import (
	"math/rand"
)

type Handle uintptr

const MaxHandle = ^Handle(0)

type ObjectMap[V any] struct {
	minHandle Handle
	maxHandle Handle
	m         map[Handle]V
}

func New[V any](minHandle Handle, maxHandle Handle) *ObjectMap[V] {
	if maxHandle < minHandle {
		panic("invalid parameters")
	}
	return &ObjectMap[V]{maxHandle: maxHandle, m: make(map[Handle]V)}
}

func (m *ObjectMap[V]) Add(value V) Handle {
	var count = m.maxHandle - m.minHandle
	for i := Handle(0); i < count; i++ {
		h := Handle(rand.Uint64()%uint64(count) + uint64(m.minHandle))
		if _, ok := m.m[h]; !ok {
			m.m[h] = value
			return h
		}
	}
	panic("too many objects")
}

func (m *ObjectMap[V]) Value(h Handle) (value V, ok bool) {
	value, ok = m.m[h]
	return
}

func (m *ObjectMap[V]) Remove(h Handle) {
	delete(m.m, h)
}

// Len returns the number of objects in the map.
func (m *ObjectMap[V]) Len() int {
	return len(m.m)
}
