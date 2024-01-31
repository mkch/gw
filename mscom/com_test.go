package mscom_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/mscom"
	"github.com/mkch/gw/mscom/sys"
)

func ExampleIMalloc() {
	var malloc *mscom.IMalloc
	if r := sys.CoGetMalloc((*unsafe.Pointer)(unsafe.Pointer(&malloc))); r != sys.S_OK {
		panic(r)
	}
	defer malloc.IUnknown().Release()

	p := malloc.Alloc(5)
	if p == nil {
		panic("out of memory")
	}
	defer malloc.Free(p)

	n := malloc.GetSize(p)
	fmt.Println(n)

	// Output:
	// 5
}

func TestBuiltinInterface(t *testing.T) {
	var malloc *mscom.IMalloc
	if r := sys.CoGetMalloc((*unsafe.Pointer)(unsafe.Pointer(&malloc))); r != sys.S_OK {
		t.Fatal(r)
	}
	var unknown *mscom.IUnknown
	err := malloc.IUnknown().QueryInterface(mscom.IID_IUnknown, (*unsafe.Pointer)(unsafe.Pointer(&unknown)))
	if err != nil {
		t.Fatal(err)
	}
	if unsafe.Pointer(unknown) != unsafe.Pointer(malloc) {
		t.Fatal("should be equal")
	}
	(*mscom.IUnknown)(unknown).Release()

	const size = 5
	p := malloc.Alloc(size)
	if p == nil {
		t.Fatal("nil ptr")
	}
	if n := malloc.GetSize(p); n != size {
		t.Fatal(n)
	}
	p = malloc.Realloc(p, size*2)
	if p == nil {
		t.Fatal("nil ptr")
	}
	if n := malloc.GetSize(p); n != size*2 {
		t.Fatal(n)
	}
	malloc.HeapMinimize()
	malloc.Free(p)
	malloc.IUnknown().Release()
}

func TestIUnknownImpl(t *testing.T) {
	var impl *mscom.IUnknown
	err := mscom.CreateIUnknownImpl(&impl)
	if err != nil {
		t.Fatal(err)
	}
	if impl == nil {
		t.Fatal("nil ptr")
	}
	n := impl.AddRef()
	if n != 2 {
		t.Fatal(n)
	}
	var p *mscom.IUnknown
	err = impl.QueryInterface(mscom.IID_IUnknown, (*unsafe.Pointer)(unsafe.Pointer(&p)))
	if err != nil {
		t.Fatal(err)
	}
	if p != impl {
		t.Fatal("should be equal")
	}
	n = p.Release()
	if n != 2 {
		t.Fatal(n)
	}
	impl.Release()
	n = impl.Release()
	if n != 0 {
		t.Fatal(n)
	}
}

func TestIUnknownImplConcurrent(t *testing.T) {
	const N = 99
	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func() {
			time.Sleep(time.Microsecond * time.Duration(rand.Int31n(2000)))
			var impl *mscom.IUnknown
			gg.MustOK(mscom.CreateIUnknownImpl(&impl))
			if impl == nil {
				panic("nil ptr")
			}
			n := impl.AddRef()
			if n != 2 {
				panic(n)
			}
			var p *mscom.IUnknown = &mscom.IUnknown{}
			gg.MustOK(impl.QueryInterface(mscom.IID_IUnknown, (*unsafe.Pointer)(unsafe.Pointer(&p))))
			if p != impl {
				panic(fmt.Sprintf("should be equal %p vs %p", p, impl))
			}
			n = p.Release()
			if n != 2 {
				panic(n)
			}
			impl.Release()
			n = impl.Release()
			if n != 0 {
				panic(n)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
