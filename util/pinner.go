package util

import (
	"runtime"
	"sync"
)

// keepAliveRefs is a map that keeps references to runtime.Pinner objects
// to prevent them from being garbage collected while they are still in use.
var keepAliveRefs sync.Map

func keepAlive(p *runtime.Pinner) {
	keepAliveRefs.Store(p, struct{}{})
}

func stopKeepAlive(p *runtime.Pinner) {
	keepAliveRefs.Delete(p)
}

// DataPinner is a wrapper of runtime.Pinner that pins the data of type *T.
type DataPinner[T any] struct {
	runtime.Pinner
	Data *T
}

func (w *DataPinner[T]) Pin() {
	// Pin w is enough because w.Data is not accessed in C code.
	w.Pinner.Pin(w)
	keepAlive(&w.Pinner)
}

func (w *DataPinner[T]) Unpin() {
	w.Pinner.Unpin()
	stopKeepAlive(&w.Pinner)
}
