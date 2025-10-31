package pen

import (
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/paint/withdpi"
	"github.com/mkch/gw/win32"
)

// LogPen holds a win32.EXTLOGPEN and the DPI value.
// See doc of withdpi.LogStruct for details.
type LogPen struct {
	withdpi.LogStruct[win32.EXTLOGPEN, win32.HPEN]
}

// NewLogFont creates a LogPen from win32.EXTLOGPEN.
func NewExtLogPen(l *win32.EXTLOGPEN, DPI win32.UINT) *LogPen {
	// https://learn.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-extcreatepen
	// "If the dwPenStyle parameter is PS_GEOMETRIC, the width is given in logical units. If dwPenStyle is PS_COSMETIC, the width must be set to 1."
	// "If elpWidth specifies geometric lines, the lengths are in logical units. Otherwise, the lines are cosmetic and lengths are in device units."
	if (l.PenStyle&win32.PS_NULL != 0) || // NULL pen
		(l.PenStyle&win32.PS_GEOMETRIC == 0) { // Not GEOMETRIC
		DPI = withdpi.DPI_INDEPENDENT
	}
	return &LogPen{*withdpi.NewLogStruct[win32.EXTLOGPEN, win32.HPEN](l, DPI, func(l *win32.EXTLOGPEN, oldDPI, newDPI win32.UINT) {
		if oldDPI != withdpi.DPI_INDEPENDENT {
			l.Width = metrics.DPIConv(l.Width, oldDPI, newDPI)
			for i, entry := range l.StyleEntry {
				l.StyleEntry[i] = metrics.DPIConv(entry, oldDPI, newDPI)
			}
		}
	}, func(l *win32.EXTLOGPEN) (win32.HPEN, error) {
		return win32.ExtCreatePen(l.PenStyle, l.Width,
			&win32.LOGBRUSH{
				Style: l.BrushStyle,
				Color: l.Color,
				Hatch: l.Hatch,
			}, l.StyleEntry)
	})}
}

// NewLogPen creates a LogPen from win32.LOGPEN.
func NewLogPen(logPen *win32.LOGPEN, DPI win32.UINT) *LogPen {
	extLogPen := win32.EXTLOGPEN{
		PenStyle: win32.PS_GEOMETRIC,
		Width:    win32.DWORD(logPen.Width),
		Color:    logPen.Color,
	}
	if (logPen.Style == win32.PS_DASH ||
		logPen.Style == win32.PS_DOT ||
		logPen.Style == win32.PS_DASHDOT ||
		logPen.Style == win32.PS_DASHDOTDOT) &&
		logPen.Width > 1 {
		extLogPen.PenStyle |= win32.PS_SOLID
	} else {
		extLogPen.PenStyle |= logPen.Style
	}
	return NewExtLogPen(&extLogPen, DPI)
}

// EXTLOGPEN returns a copy of the win32.EXTLOGPEN.
func (l *LogPen) EXTLOGPEN() *win32.EXTLOGPEN {
	return l.LogStruct.Struct()
}

// New creates a new Pen.
// See doc of withdpi.Object for details.
func New(l *LogPen, DPI win32.UINT) (*Pen, error) {
	obj, err := withdpi.New(&l.LogStruct, DPI)
	if err != nil {
		return nil, err
	}
	return &Pen{*obj}, nil
}

// NewCosmetic creates a DPI independent cosmetic pen.
func NewCosmetic(style win32.PEN_STYLE, color win32.COLORREF) (*Pen, error) {
	return New(NewExtLogPen(&win32.EXTLOGPEN{
		PenStyle: win32.PS_COSMETIC | style,
		Color:    color,
		Width:    1,
	}, withdpi.DPI_INDEPENDENT), withdpi.DPI_INDEPENDENT)
}

type Pen struct {
	withdpi.Object[win32.EXTLOGPEN, win32.HPEN]
}

func (pen *Pen) HPEN() win32.HPEN {
	return pen.Handle()
}

func (pen *Pen) Clone() *Pen {
	return &Pen{*pen.Object.Clone()}
}
