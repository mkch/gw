/*
Package mscom implements COM(The Microsoft Component Object Model) object creation and invocation.

A COM interface is a struct contains a virtual function table(VMT) pointer. A VMT is a struct of
method pointers.

For creating a COM object, call Init() on the VTM and then create methods use the *MethodCreator returned.

syscall.SyscallN can be used to call methods in VTM. Note be sure to pin go pointers with runtime.Pinner
before passing them to SyscallN.

See type IUnknown and IMalloc for details.
*/
package mscom

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/mkch/gw/mscom/sys"
)

//go:generate go run .scripts/genmethods.go

// Alloc allocates a chunk of memory large enough to hold a T.
// Alloc using the OLE memory allocator(CoTaskMem[Alloc/Free]).
//
// https://learn.microsoft.com/en-us/windows/win32/com/the-ole-memory-allocator
// "Whenever ownership of an allocated chunk of memory is passed through a COM interface or
// between a client and the COM library, you must use this COM allocator to allocate the memory.
// Allocation internal to an object can use any allocation scheme desired, but the COM memory
// allocator is a handy, efficient, and thread-safe allocator."
func Alloc[T any]() *T {
	size := unsafe.Sizeof(*new(T))
	return (*T)(AllocMem(size))
}

// AllocMem allocates a chunk of memory using the OLE memory allocator.
func AllocMem(size uintptr) unsafe.Pointer {
	p := sys.CoTaskMemAlloc(size)
	if p == nil {
		panic("out of memory")
	}
	clear(unsafe.Slice((*byte)(p), size))
	return p
}

// FreeMem frees the memory allocated by Alloc or AllocMem
func FreeMem(p unsafe.Pointer) {
	sys.CoTaskMemFree(p)
}

// Free frees the memory allocated by Alloc or AllocMem.
func Free[T any](p *T) {
	sys.CoTaskMemFree(unsafe.Pointer(p))
}

// methods hold the mapping of method ptr to method function.
// Key: Method ptr. Callback created by windows.NewCallback.
// Value: Function actually executed when ptr is called. Function with zero or more uintptr arguments and a uintptr return value.
type methods struct {
	m map[uintptr]any
}

// exists returns whether ptr is already in the methods.
func (d *methods) Exists(ptr uintptr) bool {
	_, ok := d.m[ptr]
	return ok
}

func newMethods() *methods {
	return &methods{m: make(map[uintptr]any)}
}

// Method returns the function actually executed when ptr is called.
func (d *methods) Method(ptr uintptr) any {
	return d.m[ptr]
}

// setMethod sets the actually function of a method.
func (d *methods) SetMethod(ptr uintptr, f any) {
	d.m[ptr] = f
}

// methodMap is all methods grouped by COM object ptr.
type methodMap struct {
	l sync.RWMutex
	m map[unsafe.Pointer]*methods
}

// Add adds a COM object to the map and returns a newly created methods struct.
func (m *methodMap) Add(obj unsafe.Pointer) (methods *MethodCreator) {
	m.l.Lock()
	defer m.l.Unlock()
	if m.m[obj] != nil {
		panic("duplicated objects")
	}
	mtds := newMethods()
	m.m[obj] = mtds
	return (*MethodCreator)(mtds)
}

// Methods returns all methods of a COM object.
func (m *methodMap) Methods(obj unsafe.Pointer) *methods {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.m[obj]
}

// Remove removes a COM object from the map.
// Should be called after the object is released.
func (m *methodMap) Remove(obj unsafe.Pointer) {
	m.l.Lock()
	defer m.l.Unlock()
	delete(m.m, obj)
}

func newMethodMap() *methodMap {
	return &methodMap{m: make(map[unsafe.Pointer]*methods)}
}

var mtdMap = newMethodMap()

// method is a COM object method.
type method struct {
	nArg int     // count of arguments.
	ptr  uintptr // the callback ptr.
}

func (h *method) Ptr() uintptr {
	return h.ptr
}

// methodCache caches all callback pointers(returned by windows.NewCallback).
// This is necessary because only a limited number of callbacks may be created
// in a single Go process.
// All callbacks are cached by their argument count. All methods  with the same
// prototype use the same callback function. But methods with same prototype in
// one object must have their own callback functions.
// In that function, all methods of the receiver(aka. this pointer) are retrieved
// from methodMap, and then the real code(go function) is found with their callback
// pointer.
type methodCache struct {
	l sync.RWMutex
	m map[int][]method
}

func newMethodCache() *methodCache {
	return &methodCache{m: make(map[int][]method)}
}

func (c *methodCache) Get(nArg int, except func(uintptr) bool) *method {
	c.l.RLock()
	defer c.l.RUnlock()
	for _, h := range c.m[nArg] {
		if !except(h.ptr) {
			return &h
		}
	}
	return nil
}

func (c *methodCache) Add(h method) {
	c.l.Lock()
	defer c.l.Unlock()
	handles := c.m[h.nArg]
	handles = append(handles, h)
	c.m[h.nArg] = handles
}

var mtdCache = newMethodCache()

// MethodCreator can be used to create method of COM object.
type MethodCreator methods

// Create creates a method of an COM object.
// Argument p is the address of a member in v-table.
// Argument f is the function actually executed when the method is called.
// Function f is expected to be a function with zero or more uintptr
// arguments(`this` pointer must be omitted, we have closures in go) and one
// uintptr result.
// The return value is m itself to allow chained calls.
func (m *MethodCreator) Create(p *uintptr, f any) *MethodCreator {
	nArg := reflect.TypeOf(f).NumIn()
	h := mtdCache.Get(nArg, (*methods)(m).Exists)
	if h == nil {
		// Cache miss.
		// Let it panic if too many arguments.
		// Leave prototype check of f to windows.NewCallback.
		h2 := methodFactory[nArg]()
		mtdCache.Add(h2)
		h = &h2
	}
	(*methods)(m).SetMethod(h.ptr, f)
	*p = h.ptr
	return m
}

// Cleanup is expected to be called when an object is released.
func Cleanup[T any](obj *T) {
	mtdMap.Remove(unsafe.Pointer(obj))
}

// Init initialize obj and returns the newly created method set.
func Init[T any](obj *T) *MethodCreator {
	return mtdMap.Add(unsafe.Pointer(obj))
}
