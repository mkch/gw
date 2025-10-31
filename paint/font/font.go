package font

import (
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/paint/withdpi"
	"github.com/mkch/gw/win32"
)

// LogFont holds a win32.LOGFONTW and the DPI value.
// See doc of withdpi.LogStruct for details.
type LogFont struct {
	withdpi.LogStruct[win32.LOGFONTW, win32.HFONT]
}

// NewLogFont creates a LogFont.
func NewLogFont(logFont *win32.LOGFONTW, DPI win32.UINT) *LogFont {
	return &LogFont{*withdpi.NewLogStruct[win32.LOGFONTW, win32.HFONT](logFont, DPI,
		func(l *win32.LOGFONTW, oldDPI, newDIP win32.UINT) {
			l.Height = metrics.DPIConv(l.Height, oldDPI, newDIP)
		}, win32.CreateFontIndirectW)}
}

// LOGFONTW returns a copy of the win32.LOGFONTW.
func (l *LogFont) LOGFONTW() *win32.LOGFONTW {
	return l.LogStruct.Struct()
}

// New creates a new Font.
func New(l *LogFont, DPI win32.UINT) (*Font, error) {
	f, err := withdpi.New[win32.LOGFONTW, win32.HFONT](&l.LogStruct, DPI)
	return &Font{*f}, err
}

// Font holds a win32.HFONT.
// See doc of withdpi.Object for details.
type Font struct {
	withdpi.Object[win32.LOGFONTW, win32.HFONT]
}

func (f *Font) HFONT() win32.HFONT {
	return f.Object.Handle()
}

func (f *Font) Clone() *Font {
	return &Font{*f.Object.Clone()}
}
