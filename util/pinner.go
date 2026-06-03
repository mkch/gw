package util

import "runtime"

// DataPinner is a wrapper of runtime.Pinner that pins the data of type *T.
type DataPinner[T any] struct {
	runtime.Pinner
	Data *T
}

func (w *DataPinner[T]) Pin() {
	// Pine w is enough because w.Data is not accessed in C code.
	w.Pinner.Pin(w)
}

func (w *DataPinner[T]) Unpin() {
	w.Pinner.Unpin()
}
