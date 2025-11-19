package brush

import (
	"errors"

	"github.com/mkch/gw/util/ref"
	"github.com/mkch/gw/win32"
)

type Brush struct {
	ref *ref.Ref[win32.HBRUSH]
}

func NewStock(stockType win32.StockObjectType) (*Brush, error) {
	h := win32.GetStockObject[win32.HBRUSH](stockType)
	if h == 0 {
		return nil, errors.New("failed to get stock brush")
	}
	return &Brush{
		ref: ref.New(h, func(h win32.HBRUSH) { /* no-op for stock objects */ }),
	}, nil
}

func New(logBrush *win32.LOGBRUSH) (*Brush, error) {
	h, err := win32.CreateBrushIndirect(logBrush)
	if err != nil {
		return nil, err
	}
	return &Brush{
		ref: ref.New(h, func(h win32.HBRUSH) { win32.DeleteObject(h) }),
	}, nil
}

func (b *Brush) HBRUSH() win32.HBRUSH {
	return b.ref.MustData()
}

func (b *Brush) Clone() *Brush {
	return &Brush{
		ref: b.ref.AddRef(),
	}
}

func (b *Brush) Release() {
	b.ref.Release()
	b.ref = nil
}
