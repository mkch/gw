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
	return win32util.GoString(str, strlen+1), nil
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
	unused1 win32.ULONG
	unused2 win32.USHORT
	unused3 win32.USHORT
	unused4 [8]win32.UCHAR
	// Don't make these fields blanks(_).
	// Blank fields are not considered when comparing equality.
}

type GUID = UUID
type REFIID = *GUID
type REFCLSID = *GUID

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

var lzCoCreateInstance = lzOle32.NewProc("CoCreateInstance")

type CLSCTX win32.DWORD

const (
	CLSCTX_INPROC_SERVER                  CLSCTX = 1
	CLSTX_INPROC_HANDLER                  CLSCTX = 0x2
	CLSTX_LOCAL_SERVER                    CLSCTX = 0x4
	CLSTX_INPROC_SERVER16                 CLSCTX = 0x8
	CLSTX_REMOTE_SERVER                   CLSCTX = 0x10
	CLSTX_INPROC_HANDLER16                CLSCTX = 0x20
	CLSTX_RESERVED1                       CLSCTX = 0x40
	CLSTX_RESERVED2                       CLSCTX = 0x80
	CLSTX_RESERVED3                       CLSCTX = 0x100
	CLSTX_RESERVED4                       CLSCTX = 0x200
	CLSTX_NO_CODE_DOWNLOAD                CLSCTX = 0x400
	CLSTX_RESERVED5                       CLSCTX = 0x800
	CLSTX_NO_CUSTOM_MARSHAL               CLSCTX = 0x1000
	CLSTX_ENABLE_CODE_DOWNLOAD            CLSCTX = 0x2000
	CLSTX_NO_FAILURE_LOG                  CLSCTX = 0x4000
	CLSTX_DISABLE_AAA                     CLSCTX = 0x8000
	CLSTX_ENABLE_AAA                      CLSCTX = 0x10000
	CLSTX_FROM_DEFAULT_CONTEXT            CLSCTX = 0x20000
	CLSTX_ACTIVATE_X86_SERVER             CLSCTX = 0x40000
	CLSCTX_ACTIVATE_32_BIT_SERVER         CLSCTX = CLSTX_ACTIVATE_X86_SERVER
	CLSTX_ACTIVATE_64_BIT_SERVER          CLSCTX = 0x80000
	CLSTX_ENABLE_CLOAKING                 CLSCTX = 0x100000
	CLSTX_APPCONTAINER                    CLSCTX = 0x400000
	CLSTX_ACTIVATE_AAA_AS_IU              CLSCTX = 0x800000
	CLSTX_RESERVED6                       CLSCTX = 0x1000000
	CLSTX_ACTIVATE_ARM32_SERVER           CLSCTX = 0x2000000
	CLSCTX_ALLOW_LOWER_TRUST_REGISTRATION CLSCTX = 0x4000000
	CLSTX_PS_DLL                          CLSCTX = 0x80000000
)

func CoCreateInstance(clsid *UUID, outer *unsafe.Pointer /*IUnknown*/, ctx CLSCTX, riid REFIID, ppv *unsafe.Pointer) HRESULT {
	return sysutil.As[HRESULT](lzCoCreateInstance.Call(uintptr(unsafe.Pointer(clsid)), uintptr(unsafe.Pointer(outer)), uintptr(ctx), uintptr(unsafe.Pointer(riid)), uintptr(unsafe.Pointer(ppv))))
}
