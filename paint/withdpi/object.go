package withdpi

import (
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/util/ref"
	"github.com/mkch/gw/win32"
)

// DPI_INDEPENDENT is a placeholder DPI value for DPI independent objects.
const DPI_INDEPENDENT = ^win32.UINT(0)

// LogStruct contains a win32.LOGxxx struct and the DPI value.
// Type L should be a win32.LOGxxx type(ie. win32.LOGFONTW), and type H is the handle type(ie. win32.HFONT).
type LogStruct[L any, H win32.HGDIOBJ] struct {
	l            L
	dpi          win32.UINT
	changeDPI    func(l *L, oldDPI, newDIP win32.UINT) // change the dimension of l form old dpi to new dip.
	createHandle func(*L) (H, error)                   // ie. win32.CreateFontIndirect
	cache        map[win32.UINT]*ref.Weak[H]           // Key: DPI
}

// NewLogStruct creates a LogStruct from L.
// Parameter changeDPI is a function which changes the dimension of L form old dpi to new dip.
// Parameter createHandle is a function which creates H from L(ie. win32.CreateFontIndirect).
func NewLogStruct[L any, H win32.HGDIOBJ](logFont *L, DPI win32.UINT, changeDPI func(l *L, oldDPI, newDPI win32.UINT), createHandle func(*L) (H, error)) *LogStruct[L, H] {
	return &LogStruct[L, H]{
		*logFont,
		DPI,
		changeDPI,
		createHandle,
		make(map[win32.UINT]*ref.Weak[H]),
	}
}

// Struct returns a copy of the struct in LogStruct.
func (l *LogStruct[L, H]) Struct() *L {
	r := l.l
	return &r
}

// DPI returns the DPI in LogStruct.
func (l *LogStruct[L, H]) DPI() win32.UINT {
	return l.dpi
}

func (l *LogStruct[L, H]) ForDPI(DPI win32.UINT) *L {
	r := l.l
	if l.dpi != DPI_INDEPENDENT && DPI != l.dpi {
		l.changeDPI(&r, l.dpi, DPI)
	}
	return &r
}

// New create an Object[L, H] from LogStruct.
// DPI must be greater than 0, or it panics.
// An object should be released after use.
func New[L any, H win32.HGDIOBJ](l *LogStruct[L, H], DPI win32.UINT) (*Object[L, H], error) {
	if DPI == 0 || l.dpi == 0 {
		panic("invalid DPI")
	}

	obj := Object[L, H]{ref: new(ref.Ref[H]), dpi: DPI, l: l}
	if cached := l.cache[DPI]; cached != nil {
		cached.StrongAssign(obj.ref) // Shared. From cache.
	} else {
		h, err := l.createHandle(l.ForDPI(DPI))
		if err != nil {
			return nil, err
		}
		ref.CreateAssign(obj.ref, h, func(h H) { gg.MustOK(win32.DeleteObject(h)) })
		var w ref.Weak[H]
		w.OnDispose = func() {
			delete(l.cache, DPI)
		}
		obj.ref.WeakAssign(&w)
		l.cache[DPI] = &w
	}
	return &obj, nil
}

// Object holds an H.
// Don't modify the H of the object, because it may be shared by other objects.
type Object[L any, H win32.HGDIOBJ] struct {
	ref *ref.Ref[H]
	dpi win32.UINT
	l   *LogStruct[L, H]
}

// Release releases the resource held by this object.
// Using a released object panics.
func (obj *Object[L, H]) Release() {
	obj.ref.Release()
	*obj = Object[L, H]{}
}

// Handle returns the H held by the object.
// Don't store or modify the returned handle.
func (obj *Object[L, H]) Handle() H {
	return obj.ref.MustData()
}

// Clone make a copy of the object.
// The returned object should be released after use.
func (obj *Object[L, H]) Clone() *Object[L, H] {
	r := Object[L, H]{ref: new(ref.Ref[H]), dpi: obj.dpi, l: obj.l}
	ref.Assign(r.ref, obj.ref)
	return &r
}

// ChangeDPI applies DPI change to the object.
func (obj *Object[L, H]) ChangeDPI(DPI win32.UINT) error {
	if obj.dpi == DPI || obj.dpi == DPI_INDEPENDENT {
		return nil
	}
	l := obj.l // obj.Release() below will clear obj.
	obj.Release()
	n, err := New(l, DPI)
	if err != nil {
		return err
	}
	*obj = *n
	return nil
}

// Debug returns the debug information useful to debug LogStruct itself.
// The return value can be anything and is subject to change at any time.
func (l *LogStruct[L, H]) Debug() unsafe.Pointer {
	return (unsafe.Pointer)(&l.cache)
}
