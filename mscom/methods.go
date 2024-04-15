// Code generated by .scripts/genmethods.go DO NOT EDIT.

package mscom

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func newMethod0() (h method) {
	h.nArg = 0
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func() uintptr)
		return f()
	}))
	return
}
func newMethod1() (h method) {
	h.nArg = 1
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer, arg1 uintptr) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func(uintptr) uintptr)
		return f(arg1)
	}))
	return
}
func newMethod2() (h method) {
	h.nArg = 2
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer, arg1 uintptr, arg2 uintptr) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func(uintptr, uintptr) uintptr)
		return f(arg1, arg2)
	}))
	return
}
func newMethod3() (h method) {
	h.nArg = 3
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer, arg1 uintptr, arg2 uintptr, arg3 uintptr) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func(uintptr, uintptr, uintptr) uintptr)
		return f(arg1, arg2, arg3)
	}))
	return
}
func newMethod4() (h method) {
	h.nArg = 4
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer, arg1 uintptr, arg2 uintptr, arg3 uintptr, arg4 uintptr) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func(uintptr, uintptr, uintptr, uintptr) uintptr)
		return f(arg1, arg2, arg3, arg4)
	}))
	return
}
func newMethod5() (h method) {
	h.nArg = 5
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer, arg1 uintptr, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func(uintptr, uintptr, uintptr, uintptr, uintptr) uintptr)
		return f(arg1, arg2, arg3, arg4, arg5)
	}))
	return
}
func newMethod6() (h method) {
	h.nArg = 6
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer, arg1 uintptr, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr, arg6 uintptr) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func(uintptr, uintptr, uintptr, uintptr, uintptr, uintptr) uintptr)
		return f(arg1, arg2, arg3, arg4, arg5, arg6)
	}))
	return
}
func newMethod7() (h method) {
	h.nArg = 7
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer, arg1 uintptr, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr, arg6 uintptr, arg7 uintptr) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func(uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr) uintptr)
		return f(arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	}))
	return
}
func newMethod8() (h method) {
	h.nArg = 8
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer, arg1 uintptr, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr, arg6 uintptr, arg7 uintptr, arg8 uintptr) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func(uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr) uintptr)
		return f(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	}))
	return
}
func newMethod9() (h method) {
	h.nArg = 9
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer, arg1 uintptr, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr, arg6 uintptr, arg7 uintptr, arg8 uintptr, arg9 uintptr) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func(uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr) uintptr)
		return f(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
	}))
	return
}
func newMethod10() (h method) {
	h.nArg = 10
	h.ptr = MethodPtr(windows.NewCallback(func(this unsafe.Pointer, arg1 uintptr, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr, arg6 uintptr, arg7 uintptr, arg8 uintptr, arg9 uintptr, arg10 uintptr) uintptr {
		f := mtdMap.Methods(this).Method(h.ptr).(func(uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr, uintptr) uintptr)
		return f(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10)
	}))
	return
}

var methodFactory = map[int]func() method{
	0:  newMethod0,
	1:  newMethod1,
	2:  newMethod2,
	3:  newMethod3,
	4:  newMethod4,
	5:  newMethod5,
	6:  newMethod6,
	7:  newMethod7,
	8:  newMethod8,
	9:  newMethod9,
	10: newMethod10,
}