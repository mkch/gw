package win32

import "unsafe"

// Source: https://learn.microsoft.com/en-us/windows/win32/winprog/windows-data-types

type BYTE uint8  // A byte (8 bits).
type WORD uint16 // A 16-bit unsigned integer.
type ATOM WORD
type INT int32 // A 32-bit signed integer.
type BOOL INT
type BOOLEAN BYTE
type CHAR int8    // An 8-bit Windows (ANSI) character.
type UCHAR uint8  // An unsigned CHAR.
type WCHAR int16  // A 16-bit Unicode character.
type DWORD uint32 // A 32-bit unsigned integer.
type COLORREF DWORD
type DWORDLONG uint64  // A 64-bit unsigned integer.
type DWORD_PTR uintptr // An unsigned long type for pointer precision.
type ULONG_PTR uintptr // An unsigned LONG_PTR.
type FLOAT float32     // A floating-point variable.
type HANDLE uintptr
type HWND HANDLE
type LONG int32     // A 32-bit signed integer.
type LONGLONG int64 // A 64-bit signed integer.
type UINT_PTR uintptr
type LPARAM LONG_PTR
type WPARAM UINT_PTR
type LRESULT LONG_PTR
type SIZE_T ULONG_PTR
type SSIZE_T LONG_PTR
type SHORT int16 // A 16-bit integer.
type UINT uint32
type ULONG uint32     // An unsigned LONG.
type ULONGLONG uint64 // A 64-bit unsigned integer.
type USHORT uint16
type HRESULT LONG

type PVOID unsafe.Pointer

type HMENU HANDLE
type HPEN HANDLE
type HBITMAP HANDLE
type HFONT HANDLE
type HBRUSH HANDLE

type HGDIOBJ interface {
	HMENU | HPEN | HBITMAP | HFONT | HBRUSH
}

type HDC HANDLE // DeleteDC
type HINSTANCE HANDLE
type HACCEL HANDLE
type HMODULE HANDLE
type HGLOBAL HANDLE
type HRSRC HANDLE
type HICON HANDLE   // DestroyIcon
type HCURSOR HANDLE // DestroyCursor

func RGB(r, g, b byte) COLORREF {
	return COLORREF(r) | (COLORREF(g) << 8) | (COLORREF(b) << 16)
}

func R(color COLORREF) byte {
	return byte(color & 0xFF)
}

func G(color COLORREF) byte {
	return byte((color >> 8) & 0xFF)
}

func B(color COLORREF) byte {
	return byte((color >> 16) & 0xFF)
}
