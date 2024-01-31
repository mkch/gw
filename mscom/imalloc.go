package mscom

import (
	"unsafe"
)

type IMallocVMT struct {
	IUnknownVMT

	alloc        MethodPtr
	realloc      MethodPtr
	free         MethodPtr
	getSize      MethodPtr
	didAlloc     MethodPtr
	heapMinimize MethodPtr
}

type IMalloc struct{ vt *IMallocVMT }

func (m *IMalloc) IUnknown() *IUnknown {
	return (*IUnknown)(unsafe.Pointer(m))
}

func (m *IMalloc) Alloc(size uintptr) unsafe.Pointer {
	r, _ := m.vt.alloc.Call(unsafe.Pointer(m), size)
	return unsafe.Pointer(unsafe.Add(unsafe.Pointer(nil), r))
}

func (m *IMalloc) Realloc(p unsafe.Pointer, size uintptr) unsafe.Pointer {
	r, _ := m.vt.realloc.Call(unsafe.Pointer(m), uintptr(p), size)
	return unsafe.Add(nil, r)
}

func (m *IMalloc) Free(p unsafe.Pointer) {
	m.vt.free.Call(unsafe.Pointer(m), uintptr(p))
}

func (m *IMalloc) GetSize(p unsafe.Pointer) uintptr {
	r, _ := m.vt.getSize.Call(unsafe.Pointer(m), uintptr(p))
	return r
}

func (m *IMalloc) DidAlloc(p unsafe.Pointer) int {
	r, _ := m.vt.didAlloc.Call(unsafe.Pointer(m), uintptr(p))
	return int(r)
}

func (m *IMalloc) HeapMinimize() {
	m.vt.heapMinimize.Call(unsafe.Pointer(m))
}
