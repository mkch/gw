/*
Package ref implements reference counted resource management.

A Ref[T] is a reference counted T. Struct Ref must not be copied(go vet detects this).

Creation of Ref:

	CreateAssign(&result, ...)

Assignment of two Refs:

	Assign(&refLeft, &refRight)
*/
package ref

import (
	"errors"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// Ref is a reference counted T.
type Ref[T any] struct {
	noCopy  noCopy
	dispose func(T)
	count   uint
	data    *T         // nil if disposed
	w       []*Weak[T] // all weak references
}

var errEmptyRef = errors.New("empty Ref")
var errNilDisposer = errors.New("nil disposer")

// AddRef increases the reference count by 1 and returns r itself.
// Panics if r is empty.
func (r *Ref[T]) AddRef() *Ref[T] {
	if r.count == 0 {
		panic(errEmptyRef)
	}
	r.count++
	return r
}

// Release decreases the reference count by 1.
// The reference is disposed and set to empty if it's reference count becomes zero.
// Panics if r is empty.
func (r *Ref[T]) Release() {
	if r.count == 0 {
		panic(errEmptyRef)
	}
	r.count--
	if r.count == 0 {
		r.dispose(*r.data)

		for _, w := range r.w {
			if w.onDispose != nil {
				w.onDispose()
			}
			w.ref = nil
		}
		*r = Ref[T]{}
	}
}

// Empty returns whether r is empty.
func (r *Ref[T]) Empty() bool {
	return r.count == 0
}

// Data returns the source of r.
// Zero T and false is returned if r is empty.
func (r *Ref[T]) Data() (value T, ok bool) {
	if r.data == nil {
		return
	}
	return *r.data, true
}

// Weak creates a weak reference of r.
// Panics if r is empty.
// Parameter onDispose is function which will be called when the source reference is disposed
func (r *Ref[T]) Weak(onDispose func()) *Weak[T] {
	if r.Empty() {
		panic(errEmptyRef)
	}

	w := &Weak[T]{ref: r, onDispose: onDispose}
	r.w = append(r.w, w)
	return w
}

// MustData returns the source of r.
// Panic if r is empty.
func (r *Ref[T]) MustData() T {
	if r.data == nil {
		panic(errEmptyRef)
	}
	return *r.data
}

// Assign assigns right to left.
// After successful assignment, left references the same data as right, and the original left is released.
func Assign[T any](left **Ref[T], right *Ref[T]) {
	if *left == right {
		return // NOP: self assignment.
	}
	if *left != nil {
		(*left).Release()
	}
	if right != nil {
		right.AddRef()
	}
	*left = right
}

// New creates a ref-counted data.
// Parameter disposer is a function which is called when reference count becomes zero.
// Panics if disposer is nil.
func New[T any](data T, disposer func(T)) *Ref[T] {
	if disposer == nil {
		panic(errNilDisposer)
	}
	dataCopy := data
	return &Ref[T]{dispose: disposer, count: 1, data: &dataCopy}
}

// Weak is a wake reference.
// A wake reference does nothing to the reference count of the source.
type Weak[T any] struct {
	noCopy    noCopy
	ref       *Ref[T]
	onDispose func() // OnDispose will be called when the source reference is disposed.
}

// Strong returns the source of w and adds its reference count.
// If w is empty, returns nil.
func (w *Weak[T]) Strong() *Ref[T] {
	if w.Empty() {
		return nil
	}
	w.ref.AddRef()
	return w.ref
}

// Empty returns whether w references nothing or the source reference is empty.
func (w *Weak[T]) Empty() bool {
	if w.ref == nil {
		return true
	}
	return w.ref.Empty()
}

// Data returns the source of w.
// Zero T and false is returned if w is empty.
func (w *Weak[T]) Data() (value T, ok bool) {
	if w.ref == nil {
		return
	}
	return w.ref.Data()
}

// MustData returns the source of the source reference of w.
// Panics if w is empty.
func (w *Weak[T]) MustData() T {
	if w.ref == nil {
		panic("nil ref")
	}
	return w.ref.MustData()
}
