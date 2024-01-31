package mscom

import (
	"syscall"
	"unsafe"
)

type IMallocVMT struct {
	IUnknownVMT

	alloc        uintptr
	realloc      uintptr
	free         uintptr
	getSize      uintptr
	didAlloc     uintptr
	heapMinimize uintptr
}

type IMalloc struct{ vt *IMallocVMT }

func (m *IMalloc) IUnknown() *IUnknown {
	return (*IUnknown)(unsafe.Pointer(m))
}

func (m *IMalloc) Alloc(size uintptr) unsafe.Pointer {
	r, _, _ := syscall.SyscallN(m.vt.alloc, uintptr(unsafe.Pointer(m)), size)
	return unsafe.Pointer(unsafe.Add(unsafe.Pointer(nil), r))
}

func (m *IMalloc) Realloc(p unsafe.Pointer, size uintptr) unsafe.Pointer {
	r, _, _ := syscall.SyscallN(m.vt.realloc, uintptr(unsafe.Pointer(m)), uintptr(p), size)
	return unsafe.Add(nil, r)
}

func (m *IMalloc) Free(p unsafe.Pointer) {
	syscall.SyscallN(m.vt.free, uintptr(unsafe.Pointer(m)), uintptr(p))
}

func (m *IMalloc) GetSize(p unsafe.Pointer) uintptr {
	r, _, _ := syscall.SyscallN(m.vt.getSize, uintptr(unsafe.Pointer(m)), uintptr(p))
	return r
}

func (m *IMalloc) DidAlloc(p unsafe.Pointer) int {
	r, _, _ := syscall.SyscallN(m.vt.didAlloc, uintptr(p))
	return int(r)
}

func (m *IMalloc) HeapMinimize() {
	syscall.SyscallN(m.vt.heapMinimize)
}
