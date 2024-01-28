/*
Package ref implements reference counted resource management.

A Ref[T] is a reference counted T. Struct Ref must not be copied(go vet detects this).

Creation of Ref:

	CreateAssign(&result, ...)

Assignment of two Refs:

	Assign(&refLeft, &refRight)
*/
package ref

import "slices"

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type ref[T any] struct {
	dispose func(T)
	count   uint
	data    *T         // nil if disposed
	w       []*Weak[T] // all weak references
}

// Ref is a reference counted T.
// Zero Ref is an empty reference.
type Ref[T any] struct {
	noCopy  noCopy
	*ref[T] // nil if empty
}

// AddRef increases the reference count by 1.
// Panics if r is empty.
func (r *Ref[T]) AddRef() {
	r.count++
}

// Release decreases the reference count by 1.
// The reference is disposed and set to empty if it's reference count becomes zero.
// Panics if r is empty.
func (r *Ref[T]) Release() {
	r.count--
	if r.count == 0 {
		r.dispose(*r.data)

		for _, w := range r.w {
			if w.OnDispose != nil {
				w.OnDispose()
			}
			w.ref = nil
		}
		r.w = nil
		r.dispose = nil
		r.data = nil
	}
}

// Empty returns whether r is empty.
func (r *Ref[T]) Empty() bool {
	return r.ref == nil || r.data == nil
}

// Data returns the source of r.
// Zero T and false is returned if r is empty.
func (r *Ref[T]) Data() (value T, ok bool) {
	if r.ref == nil || r.data == nil {
		return
	}
	return *r.data, true
}

// WeakAssign assigns a weak reference of r to left.
// Panics if r is empty.
func (r *Ref[T]) WeakAssign(left *Weak[T]) {
	if r.Empty() {
		panic("nil ref")
	}
	if !left.Empty() {
		left.ref.w = slices.DeleteFunc(left.ref.w, func(w *Weak[T]) bool { return left == w })
	}

	left.ref = r
	r.w = append(r.w, left)
}

// MustData returns the source of r.
// Panic if r is empty.
func (r *Ref[T]) MustData() T {
	if r.ref == nil {
		panic("nil ref")
	}
	return *r.data
}

// Assign assigns right to left.
// After successful assignment, left references the same data as right, and the original left is released.
func Assign[T any](left *Ref[T], right *Ref[T]) {
	if left.ref == right.ref {
		return // NOP: self assignment.
	}
	if left.ref != nil {
		left.Release()
	}
	if right.ref != nil {
		right.AddRef()
	}
	left.ref = right.ref
}

// CreateAssign creates a ref-counted data and assigns[using the semantics of Assign()] it to left.
// Disposer is a function which is called when reference count becomes zero.
// Panics if disposer is nil.
func CreateAssign[T any](left *Ref[T], data T, disposer func(T)) {
	if disposer == nil {
		panic("nil disposer")
	}
	if !left.Empty() {
		left.Release()
	}
	dataCopy := data
	left.ref = &ref[T]{dispose: disposer, count: 1, data: &dataCopy}
}

// Weak is a wake reference.
// A wake reference does nothing to the reference count of the source.
type Weak[T any] struct {
	noCopy    noCopy
	ref       *Ref[T]
	OnDispose func() // OnDispose will be called when the source reference is disposed.
}

// StrongAssign assigns the source of w to left.
// If w is empty, left will be empty too after this method is called.
func (w *Weak[T]) StrongAssign(left *Ref[T]) {
	if w.Empty() {
		Assign(left, &Ref[T]{})
	}
	Assign(left, w.ref)
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

// WeakAssign creates a new weak reference of R and assigns[using the semantics of WeakAssign()] it
// to left, which R is the source of right.
func WeakAssign[T any](left *Weak[T], right *Weak[T]) {
	if left.ref == right.ref {
		return // NOP: self assignment.
	}
	right.ref.WeakAssign(left)
}
