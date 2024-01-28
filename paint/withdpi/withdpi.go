package withdpi

import (
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

type Int[T ~int32 | ~uint32] struct {
	Value T
	DPI   win32.UINT
}

func (i *Int[T]) ForDPI(newDPI win32.UINT) T {
	if newDPI == i.DPI {
		return i.Value
	}
	return win32util.DPIConv(i.Value, i.DPI, newDPI)
}

type Value[T any] struct {
	Value    T
	DPI      win32.UINT
	applyDPI func(value *T, oldDPI, newDPI win32.UINT)
}

func NewValue[T any](value *T, DPI win32.UINT, applyDPI func(value *T, oldDPI, newDPI win32.UINT)) *Value[T] {
	return &Value[T]{
		Value:    *value,
		DPI:      DPI,
		applyDPI: applyDPI,
	}
}

func (val Value[T]) ForDPI(newDPI win32.UINT) *T {
	v := val.Value
	if newDPI == val.DPI {
		return &v
	}
	val.applyDPI(&v, val.DPI, newDPI)
	return &v
}
