package objectmap_test

import (
	"testing"

	"github.com/mkch/gw/internal/objectmap"
)

func panics(f func()) (result bool) {
	defer func() {
		result = recover() != nil
	}()
	f()
	return
}

func TestObjectMap_NewInvalidParametersPanics(t *testing.T) {
	if !panics(func() { objectmap.New[int](10, 9) }) {
		t.Fatal("should panic when maxHandle < minHandle")
	}
}

func TestObjectMap_AddExhaustSmallRange(t *testing.T) {
	const min = objectmap.Handle(100)
	const max = objectmap.Handle(104)

	count := int(max - min)
	om := objectmap.New[int](min, max)
	handles := make([]objectmap.Handle, 0, count)
	for i := 0; i < count; i++ {
		h := om.Add(i)
		handles = append(handles, h)
	}

	for _, h := range handles {
		if h < min || h >= max {
			t.Fatalf("handle out of range: got=%d, want in [%d, %d)", h, min, max)
		}
	}
	if om.Len() != count {
		t.Fatalf("wrong map length after exhaustion: got=%d want=%d", om.Len(), count)
	}
	if !panics(func() { om.Add(999) }) {
		t.Fatal("should panic when map is full")
	}
}

func TestObjectMap_ValueRemoveLen(t *testing.T) {
	const min = objectmap.Handle(0xFF)
	const max = objectmap.Handle(0x1FF)

	om := objectmap.New[string](min, max)
	if om.Len() != 0 {
		t.Fatalf("initial length should be 0, got %d", om.Len())
	}

	h := om.Add("v")
	if h < min || h >= max {
		t.Fatalf("handle out of range: got=%d, want in [%d, %d)", h, min, max)
	}
	if om.Len() != 1 {
		t.Fatalf("wrong length after Add: got=%d want=1", om.Len())
	}

	if v, ok := om.Value(h); !ok {
		t.Fatal("Value should return ok=true for existing handle")
	} else if v != "v" {
		t.Fatalf("wrong value: got=%q want=%q", v, "v")
	}

	om.Remove(h)
	if om.Len() != 0 {
		t.Fatalf("wrong length after Remove: got=%d want=0", om.Len())
	}
	if v, ok := om.Value(h); ok {
		t.Fatalf("Value should return ok=false after Remove, got value=%q", v)
	}

	om.Remove(h)
	if om.Len() != 0 {
		t.Fatalf("length should remain 0 after removing non-existing handle, got=%d", om.Len())
	}
}

func TestObjectMap_AddRespectsNonZeroMinHandle(t *testing.T) {
	const min = objectmap.Handle(0xFF)
	const max = objectmap.Handle(0x104)

	om := objectmap.New[int](min, max)
	h := om.Add(1)
	if h < min || h >= max {
		t.Fatalf("handle out of range for non-zero minHandle: got=%d, want in [%d, %d)", h, min, max)
	}
}
