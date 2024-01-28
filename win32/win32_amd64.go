package win32

type INT_PTR int64  // A signed integer type for pointer precision.
type LONG_PTR int64 // A signed long type for pointer precision.

var lzGetWindowLongPtrW = lzUser32.NewProc("GetWindowLongPtrW")
var lzSetWindowLongPtrW = lzUser32.NewProc("SetWindowLongPtrW")
