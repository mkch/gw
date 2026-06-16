package listcombo

import (
	"unsafe"

	"github.com/mkch/gw/win32"
	"golang.org/x/exp/constraints"
)

// Adopt generic method if Go added support for that.

// SendMessageRet sends a message to the list box or combo box and returns the result as type RT.
func SendMessageRet[
	RT constraints.Integer, WT constraints.Integer | unsafe.Pointer, LT constraints.Integer | unsafe.Pointer](hwnd win32.HWND, msg win32.UINT, wParam WT, lParam LT) (r RT, err error) {
	if msg == 0 {
		panic("message cannot be 0")
	}
	i, err := win32.SendMessageW(hwnd, msg, win32.WPARAM(wParam), win32.LPARAM(uintptr(lParam)))
	if err != nil {
		return
	}
	return RT(i), nil
}

// SendMessageRetOkay sends a message to the list box or combo box and returns the result as type RT if it is LC_OKAY.
// If the result is not LC_OKAY, it returns [ResultError] with the result as the error code.
func SendMessageRetOkay[
	RT constraints.Integer, WT constraints.Integer | unsafe.Pointer, LT constraints.Integer | unsafe.Pointer](hwnd win32.HWND, msg win32.UINT, wParam WT, lParam LT) (r RT, err error) {
	i, err := SendMessageRet[RT](hwnd, msg, wParam, lParam)
	if err != nil {
		return
	}
	if i != LC_OKAY {
		err = ResultError(i)
		return
	}
	return i, nil
}

// SendMessageOkay sends a message to the list box or combo box and returns nil if the result is LC_OKAY.
// If the result is not LC_OKAY, it returns [ResultError] with the result as the error code.
func SendMessageOkay[
	WT constraints.Integer | unsafe.Pointer, LT constraints.Integer | unsafe.Pointer](hwnd win32.HWND, msg win32.UINT, wParam WT, lParam LT) (err error) {
	i, err := SendMessageRet[win32.LRESULT](hwnd, msg, wParam, lParam)
	if err != nil {
		return
	}
	if i != LC_OKAY {
		err = ResultError(i)
		return
	}
	return
}

// SendMessageRetNoError sends a message to the list box or combo box and returns the result as type RT if it is not LC_ERR.
// If the result is LC_ERR, it returns [ResultError] with the result as the error code.
func SendMessageRetNoError[
	RT constraints.Integer, WT constraints.Integer | unsafe.Pointer, LT constraints.Integer | unsafe.Pointer](hwnd win32.HWND, msg win32.UINT, wParam WT, lParam LT) (r RT, err error) {
	i, err := SendMessageRet[RT](hwnd, msg, wParam, lParam)
	if err != nil {
		return
	}
	if i == RT(*new(LC_ERR)) {
		err = ResultError(i)
		return
	}
	return i, nil
}

// SendMessageNoError sends a message to the list box or combo box and returns nil if the result is not LC_ERR.
// If the result is LC_ERR, it returns [ResultError] with the result as the error code.
func SendMessageNoError[
	WT constraints.Integer | unsafe.Pointer, LT constraints.Integer | unsafe.Pointer](hwnd win32.HWND, msg win32.UINT, wParam WT, lParam LT) (err error) {
	i, err := SendMessageRet[win32.LRESULT](hwnd, msg, wParam, lParam)
	if err != nil {
		return
	}
	if i == LC_ERR {
		err = ResultError(i)
		return
	}
	return
}
