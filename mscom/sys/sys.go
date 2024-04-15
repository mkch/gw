package sys

import (
	"unicode/utf16"
	"unsafe"

	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/sysutil"
	"github.com/mkch/gw/win32/win32util"
	"golang.org/x/sys/windows"
)

//go:generate stringer -output sys_string.go -type=HRESULT,RPC_STATUS

type RPC_STATUS win32.LONG

type RpcStatusError RPC_STATUS

func (err RpcStatusError) Error() string {
	return RPC_STATUS(err).String()
}

const (
	RPC_S_OK                  RPC_STATUS = 0
	RPC_S_INVALID_STRING_UUID RPC_STATUS = 1705
)

var lzRpcrt4 = windows.NewLazySystemDLL("Rpcrt4.dll")

var lzUuidFromStringW = lzRpcrt4.NewProc("UuidFromStringW")

func UuidFromStringW(uuid string) (*UUID, error) {
	str := append(utf16.Encode(([]rune)(uuid)), 0)
	var ret UUID
	if r, _, _ := lzUuidFromStringW.Call(uintptr(unsafe.Pointer(&str[0])), uintptr(unsafe.Pointer(&ret))); r != uintptr(RPC_S_OK) {
		return nil, RpcStatusError(r)
	}
	return &ret, nil
}

var lzUuidToStringW = lzRpcrt4.NewProc("UuidToStringW")

func UuidToStringW(uuid *UUID) (string, error) {
	var str *win32.WCHAR
	if r, _, _ := lzUuidToStringW.Call(uintptr(unsafe.Pointer(uuid)), uintptr(unsafe.Pointer(&str))); r != uintptr(RPC_S_OK) {
		return "", RpcStatusError(r)
	}
	defer RpcStringFreeW(&str)

	strlen := 0
	for p := str; *p != 0; p = (*win32.WCHAR)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(uint16(0)))) {
		strlen++
	}
	return win32util.GoString(str, strlen), nil
}

var lzRpcStringFreeW = lzRpcrt4.NewProc("RpcStringFreeW")

func RpcStringFreeW(str **win32.WCHAR) error {
	if r, _, _ := lzRpcStringFreeW.Call(uintptr(unsafe.Pointer(str))); r != 0 {
		return RpcStatusError(r)
	}
	return nil
}

type HRESULT win32.HRESULT

func (h HRESULT) Uintptr() uintptr {
	return uintptr(h)
}

type HResultError HRESULT

func (err HResultError) Error() string {
	return HRESULT(err).String()
}

const (
	S_OK           HRESULT = 0
	S_FALSE        HRESULT = 1
	E_NOINTERFACE  HRESULT = -(^0x80004002 & 0x7FFFFFFF) - 1 //0x80004002
	E_POINTER      HRESULT = -(^0x80004003 & 0x7FFFFFFF) - 1 //0x80004003
	E_UNEXPECTED   HRESULT = -(^0x8000FFFF & 0x7FFFFFFF) - 1 //0x8000FFFF
	E_NOTIMPL      HRESULT = -(^0x80004001 & 0x7FFFFFFF) - 1 //0x80004001
	E_OUTOFMEMORY  HRESULT = -(^0x8007000E & 0x7FFFFFFF) - 1 //0x8007000E
	E_INVALIDARG   HRESULT = -(^0x80070057 & 0x7FFFFFFF) - 1 //0x80070057
	E_HANDLE       HRESULT = -(^0x80070006 & 0x7FFFFFFF) - 1 //0x80070006
	E_ABORT        HRESULT = -(^0x80004004 & 0x7FFFFFFF) - 1 //0x80004004
	E_FAIL         HRESULT = -(^0x80004005 & 0x7FFFFFFF) - 1 //0x80004005
	E_ACCESSDENIED HRESULT = -(^0x80070005 & 0x7FFFFFFF) - 1 //0x80070005
)

type UUID struct {
	_ win32.ULONG
	_ win32.USHORT
	_ win32.USHORT
	_ [8]win32.UCHAR
}

type GUID = UUID
type REFIID = *GUID

var lzOle32 = windows.NewLazySystemDLL("Ole32.dll")

var lzCoInitialize = lzOle32.NewProc("CoInitialize")

func CoInitialize() uintptr {
	return sysutil.As[uintptr](lzCoInitialize.Call(0))
}

var lzCoGetMalloc = lzOle32.NewProc("CoGetMalloc")

func CoGetMalloc(ppMalloc *unsafe.Pointer) HRESULT {
	return sysutil.As[HRESULT](lzCoGetMalloc.Call(1, uintptr(unsafe.Pointer(ppMalloc))))
}

var lzCoTaskMemAlloc = lzOle32.NewProc("CoTaskMemAlloc")

func CoTaskMemAlloc(size uintptr) unsafe.Pointer {
	r, _, _ := lzCoTaskMemAlloc.Call(size)
	return unsafe.Add(unsafe.Pointer(nil), r)
}

var lzCoTaskMemFree = lzOle32.NewProc("CoTaskMemFree")

func CoTaskMemFree(p unsafe.Pointer) {
	lzCoTaskMemFree.Call(uintptr(p))
}

var lzCoTaskMemRealloc = lzOle32.NewProc("CoTaskMemRealloc")

func CoTaskMemRealloc(p unsafe.Pointer, size uintptr) unsafe.Pointer {
	r, _, _ := lzCoTaskMemRealloc.Call(uintptr(p), size)
	return unsafe.Add(unsafe.Pointer(nil), r)
}
