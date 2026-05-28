package font

import (
	"sync"
	"unsafe"

	"github.com/mkch/gw/win32"
)

var defaultFont = sync.OnceValues(func() (*LogFont, error) {
	const DPI = win32.USER_DEFAULT_SCREEN_DPI
	var metrics = win32.NONCLIENTMETRICSW{Size: win32.UINT(unsafe.Sizeof(win32.NONCLIENTMETRICSW{}))}
	if err := win32.SystemParametersInfoForDpi(win32.SPI_GETNONCLIENTMETRICS, win32.UINT(unsafe.Sizeof(metrics)), win32.PVOID(&metrics), 0, DPI); err != nil {
		return nil, err
	}
	return NewLogFont(&metrics.MessageFont, DPI), nil
})

// SysDefault returns the default system font for UI.
func SysDefault() *LogFont {
	lf, err := defaultFont()
	if err != nil {
		panic(err)
	}
	return lf
}
