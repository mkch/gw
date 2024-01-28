package win32

type INT_PTR int32  // A signed integer type for pointer precision.
type LONG_PTR int32 // A signed long type for pointer precision.

// SetWindowLongPtrW is not available on 386.
var lzGetWindowLongPtrW = lzUser32.NewProc("GetWindowLongW")

// SetWindowLongPtrW is not available on 386.
var lzSetWindowLongPtrW = lzUser32.NewProc("SetWindowLongW")
