package mscom

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/mscom/sys"
)

// IUnknownVMT is the v-table of IUnknown interface.
// The doc(microsoft website) may reorder the method list.
// Refer to C source code to get the correct order.
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

// CreateIUnknownImpl creates an object which implements IUnknown interface.
// This function may also be considered as an example of creating COM objects with [InitIUnknownImpl].
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
		return sys.E_NOINTERFACE
	}, func() {
		Free(mem)
	})

	*ppObject = &mem.IUnknown
	return nil
}

// InitIUnknownImpl initializes an IUnknownVMT with a default implementation.
//
// This function is a helper to implement arbitrary COM object. See [CreateIUnknownImpl] and other examples.
//
// This implementation implements reference counting by atomic.Int32, and calling release when reference
// count reaches 0. It implements QueryInterface method for IID_IUnknown after checking param for E_POINTER
// and then calling queryInterface.
//
// If queryInterface returns S_OK, the reference count is increased by 1 in this func. So don't increase the
// reference count in queryInterface itself.
// The returned MethodCreator is the result of calling [Init](obj) and can be used to add more methods other
// than that of IUnknown.
//
// Pitfall when declaring VMT: The method order of COM interface may be reorder in microsoft website!!
// So refer to the C source code to get the correct order.
func InitIUnknownImpl[T any](obj *T, vt *IUnknownVMT, queryInterface func(sys.REFIID, *unsafe.Pointer) sys.HRESULT, release func()) (mtds *MethodCreator) {
	var refCount atomic.Int32
	refCount.Add(1) // New object has ref count of 1.

	mtds = Init(obj)
	mtds.
		Create(&vt.queryInterface, func(intIID uintptr, intPP uintptr) uintptr {
			iid := sys.REFIID(unsafe.Add(nil, intIID))
			pp := (*unsafe.Pointer)(unsafe.Add(nil, intPP))

			result := sys.E_POINTER
			if pp != nil {
				if *iid == *IID_IUnknown {
					*pp = unsafe.Pointer(obj)
					result = sys.S_OK
				} else {
					result = queryInterface(iid, pp)
				}
			}
			if result == sys.S_OK {
				if *pp == nil {
					panic("return S_OK but pp == nil")
				}
				refCount.Add(1)
			}
			return result.Uintptr()
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
	return
}
