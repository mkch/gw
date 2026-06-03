package win32util

import (
	"unsafe"

	"github.com/mkch/gw/util"
	"github.com/mkch/gw/win32"
)

// WindowProp is a window property with type T. The name of property is *WindowProp[T] itself.
// To create a WindowProp, use [NewWindowProp] function or cast a *win32.WCHAR to *WindowProp[T].
// A windowProp never removes itself. A Set call with nil data is needed to remove it.
type WindowProp[T any] win32.WCHAR

// NewWindowProp creates a new WindowProp with the given name.
func NewWindowProp[T any](name string) *WindowProp[T] {
	var buf []win32.WCHAR
	CString(name, &buf)
	return (*WindowProp[T])(&buf[0])
}

// Set sets a property data of a window.
// If data is nil, the property is removed.
// If data is not nil and the property already exists, the old data will be replaced by the new data.
func (w *WindowProp[T]) Set(hwnd win32.HWND, data *T) error {
	if oldP := win32.RemovePropW(hwnd, (*win32.WCHAR)(w)); oldP != 0 {
		(*util.DataPinner[T])(unsafe.Add(nil, oldP)).Unpin()
	}
	if data == nil {
		return nil
	}
	p := &util.DataPinner[T]{Data: data}
	p.Pin()
	if err := win32.SetPropW(hwnd, (*win32.WCHAR)(w), win32.HANDLE(uintptr(unsafe.Pointer(p)))); err != nil {
		p.Unpin()
		return err
	}
	return nil
}

// Get gets the property data of a window.
// It returns nil if the property does not exist.
func (w *WindowProp[T]) Get(hwnd win32.HWND) *T {
	if p := win32.GetPropW(hwnd, (*win32.WCHAR)(w)); p != 0 {
		return (*util.DataPinner[T])(unsafe.Add(nil, p)).Data
	}
	return nil
}
