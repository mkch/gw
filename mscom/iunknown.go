package mscom

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/mscom/sys"
)

type IUnknownVMT struct {
	queryInterface MethodPtr
	addRef         MethodPtr
	release        MethodPtr
}

type IUnknown struct{ vt *IUnknownVMT }

func (i *IUnknown) QueryInterface(riid sys.REFIID, pp *unsafe.Pointer) error {
	var pinner runtime.Pinner
	pinner.Pin(pp) // IMPORTANT!
	defer pinner.Unpin()
	r, _ := i.vt.queryInterface.Call(unsafe.Pointer(i), uintptr(unsafe.Pointer(riid)), uintptr(unsafe.Pointer(pp)))
	if sys.HRESULT(r) != sys.S_OK {
		return sys.HResultError(r)
	}
	return nil
}

func (i *IUnknown) AddRef() uint32 {
	r, _ := i.vt.addRef.Call(unsafe.Pointer(i))
	return uint32(r)
}

func (i *IUnknown) Release() uint32 {
	r, _ := i.vt.release.Call(unsafe.Pointer(i))
	return uint32(r)
}

var IID_IUnknown = gg.Must(sys.UuidFromStringW("00000000-0000-0000-C000-000000000046"))

type IUnknownMethods struct {
	QueryInterface func(sys.REFIID, *unsafe.Pointer) int32
	AddRef         func() uint32
	Release        func() uint32
}

// CreateIUnknownImpl creates an object which implements IUnknown interface.
func CreateIUnknownImpl(ppObject **IUnknown) error {
	if ppObject == nil {
		return sys.HResultError(sys.E_POINTER)
	}
	// Alloc the interface and v-table in one block of memory.
	mem := Alloc[struct {
		IUnknown
		IUnknownVMT
	}]()
	if mem == nil {
		panic("out of memory")
	}
	mem.IUnknown.vt = &mem.IUnknownVMT
	InitIUnknownImpl(&mem.IUnknown, &mem.IUnknownVMT, func(id sys.REFIID, p *unsafe.Pointer) sys.HRESULT {
		if *id == *IID_IUnknown {
			mem.AddRef()
			*p = unsafe.Pointer(mem)
			return sys.S_OK
		}
		*p = nil
		return sys.E_NOINTERFACE
	}, func() {
		Free(mem)
	})

	*ppObject = &mem.IUnknown
	return nil
}

// InitIUnknownImpl initializes an IUnknownVMT with a simple implementation.
// The implementation implements reference counting by atomic.Int32, and calling release when reference
// count reaches 0. It implements QueryInterface method by first checking param for E_POINTER and then
// calling queryInterface.
func InitIUnknownImpl[T any](obj *T, vt *IUnknownVMT, queryInterface func(sys.REFIID, *unsafe.Pointer) sys.HRESULT, release func()) {
	var refCount atomic.Int32
	refCount.Add(1) // New object has ref count of 1.

	Init(obj).
		Create(&vt.queryInterface, func(intIID uintptr, intPP uintptr) uintptr {
			iid := sys.REFIID(unsafe.Add(nil, intIID))
			pp := (*unsafe.Pointer)(unsafe.Add(nil, intPP))
			if pp == nil {
				return sys.E_POINTER.Uintptr()
			}
			return uintptr(queryInterface(iid, pp))
		}).
		Create(&vt.addRef, func() uintptr {
			return uintptr(refCount.Add(1))
		}).
		Create(&vt.release, func() uintptr {
			count := refCount.Add(-1)
			if count == 0 {
				// Must be called before Free to avoid race condition.
				// A premature Free call allows another goroutine to reuse
				// the pointer before it is removed from method map.
				Cleanup(obj)
				if release != nil {
					release()
				}
			} else if count < 0 {
				panic("too many releases")
			}
			return uintptr(count)
		})
}
