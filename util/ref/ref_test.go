package ref_test

import (
	"testing"

	"github.com/mkch/gw/util/ref"
)

func NoCopy() {
	var l, r ref.Ref[int]
	// go vet should report the following lines.
	l = r // go vet: assignment copies lock value to l
	r = l // go vet: assignment copies lock value to r
}

func TestEmpty(t *testing.T) {
	var empty ref.Ref[int]
	if !empty.Empty() {
		t.Fatal("should be empty")
	}
	if !panics(func() { empty.AddRef() }) {
		t.Fatal("should panic")
	}
	if !panics(func() { empty.Release() }) {
		t.Fatal("should panic")
	}
	if data, ok := empty.Data(); ok {
		t.Fatal("should be false")
	} else if data != 0 {
		t.Fatalf("wrong data: %v", data)
	}
}

func panics(f func()) (result bool) {
	defer func() {
		result = recover() != nil
	}()
	f()
	return
}

func TestRelease(t *testing.T) {
	var r ref.Ref[string]
	var disposed = false
	ref.CreateAssign(&r, "str", func(s string) { disposed = true })
	if r.Empty() {
		t.Fatal("should not be empty")
	}
	if data, ok := r.Data(); !ok {
		t.Fatal("should not be empty")
	} else if data != "str" {
		t.Fatalf("wrong data: %v", data)
	}
	r.Release()
	if !disposed {
		t.Fatal("should be disposed")
	}
	if !r.Empty() {
		t.Fatal("should be empty")
	}
}

func TestAssign(t *testing.T) {
	var empty ref.Ref[int]
	var l ref.Ref[int]
	var lDisposed = false
	ref.CreateAssign(&l, 1, func(i int) { lDisposed = true })
	var r ref.Ref[int]
	var rDisposed = false
	ref.CreateAssign(&r, 2, func(i int) { rDisposed = true })

	ref.Assign(&l, &r)
	if !lDisposed {
		t.Fatal("should be disposed")
	}
	if rDisposed {
		t.Fatal("should not be disposed")
	}
	if data, ok := r.Data(); !ok {
		t.Fatal("should not be empty")
	} else if data != 2 {
		t.Fatalf("wrong data: %v", data)
	}

	if data, ok := l.Data(); !ok {
		t.Fatal("should not be empty")
	} else if data != 2 {
		t.Fatalf("wrong data: %v", data)
	}

	ref.Assign(&l, &empty)
	if rDisposed {
		t.Fatal("should not be disposed")
	}
	r.Release()
	if !rDisposed {
		t.Fatal("should be disposed")
	}
}

func TestWeak(t *testing.T) {
	var r ref.Ref[int]
	rDisposed := false
	ref.CreateAssign(&r, 1, func(i int) { rDisposed = true })
	var wr ref.Weak[int]
	r.WeakAssign(&wr)
	if data, ok := wr.Data(); !ok {
		t.Fatal("should not be empty")
	} else if data != 1 {
		t.Fatalf("wrong data: %v", data)
	}
	wrNotified := false
	wr.OnDispose = func() { wrNotified = true }

	var r2 ref.Ref[int]
	wr.StrongAssign(&r2)
	if data, ok := r2.Data(); !ok {
		t.Fatal("should not be empty")
	} else if data != 1 {
		t.Fatalf("wrong data: %v", data)
	}

	r2.Release()
	r.Release()

	if !r.Empty() {
		t.Fatal("should be empty")
	}
	if !rDisposed {
		t.Fatal("should be disposed")
	}

	if !r2.Empty() {
		t.Fatal("should be empty")
	}

	if !wrNotified {
		t.Fatal("should be notified")
	}
	if !wr.Empty() {
		t.Fatal("should be empty")
	}
}
