package win32

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var lzComdlg32 = windows.NewLazySystemDLL("Comdlg32.dll")

type CHOOSE_FONT_FLAG DWORD

const (
	CF_SCREENFONTS          CHOOSE_FONT_FLAG = 0x00000001
	CF_PRINTERFONTS         CHOOSE_FONT_FLAG = 0x00000002
	CF_BOTH                 CHOOSE_FONT_FLAG = (CF_SCREENFONTS | CF_PRINTERFONTS)
	CF_SHOWHELP             CHOOSE_FONT_FLAG = 0x00000004
	CF_ENABLEHOOK           CHOOSE_FONT_FLAG = 0x00000008
	CF_ENABLETEMPLATE       CHOOSE_FONT_FLAG = 0x00000010
	CF_ENABLETEMPLATEHANDLE CHOOSE_FONT_FLAG = 0x00000020
	CF_INITTOLOGFONTSTRUCT  CHOOSE_FONT_FLAG = 0x00000040
	CF_USESTYLE             CHOOSE_FONT_FLAG = 0x00000080
	CF_EFFECTS              CHOOSE_FONT_FLAG = 0x00000100
	CF_APPLY                CHOOSE_FONT_FLAG = 0x00000200
	CF_ANSIONLY             CHOOSE_FONT_FLAG = 0x00000400
	CF_SCRIPTSONLY          CHOOSE_FONT_FLAG = CF_ANSIONLY
	CF_NOVECTORFONTS        CHOOSE_FONT_FLAG = 0x00000800
	CF_NOOEMFONTS           CHOOSE_FONT_FLAG = CF_NOVECTORFONTS
	CF_NOSIMULATIONS        CHOOSE_FONT_FLAG = 0x00001000
	CF_LIMITSIZE            CHOOSE_FONT_FLAG = 0x00002000
	CF_FIXEDPITCHONLY       CHOOSE_FONT_FLAG = 0x00004000
	CF_WYSIWYG              CHOOSE_FONT_FLAG = 0x00008000 // must also have CF_SCREENFONTS & CF_PRINTERFONTS
	CF_FORCEFONTEXIST       CHOOSE_FONT_FLAG = 0x00010000
	CF_SCALABLEONLY         CHOOSE_FONT_FLAG = 0x00020000
	CF_TTONLY               CHOOSE_FONT_FLAG = 0x00040000
	CF_NOFACESEL            CHOOSE_FONT_FLAG = 0x00080000
	CF_NOSTYLESEL           CHOOSE_FONT_FLAG = 0x00100000
	CF_NOSIZESEL            CHOOSE_FONT_FLAG = 0x00200000
	CF_SELECTSCRIPT         CHOOSE_FONT_FLAG = 0x00400000
	CF_NOSCRIPTSEL          CHOOSE_FONT_FLAG = 0x00800000
	CF_NOVERTFONTS          CHOOSE_FONT_FLAG = 0x01000000
	CF_INACTIVEFONTS        CHOOSE_FONT_FLAG = 0x02000000
)

type CHOOSE_FONT_TYPE WORD

const (
	SIMULATED_FONTTYPE CHOOSE_FONT_TYPE = 0x8000
	PRINTER_FONTTYPE   CHOOSE_FONT_TYPE = 0x4000
	SCREEN_FONTTYPE    CHOOSE_FONT_TYPE = 0x2000
	BOLD_FONTTYPE      CHOOSE_FONT_TYPE = 0x0100
	ITALIC_FONTTYPE    CHOOSE_FONT_TYPE = 0x0200
	REGULAR_FONTTYPE   CHOOSE_FONT_TYPE = 0x0400
)

type CHOOSEFONTW struct {
	StructSize   DWORD
	Owner        HWND
	DC           HDC
	LogFont      *LOGFONTW
	PointSize    INT
	Flags        CHOOSE_FONT_FLAG
	Color        COLORREF
	CustomData   LPARAM
	Hook         uintptr
	TemplateName *WCHAR
	Instance     HINSTANCE
	Style        *WCHAR
	FontType     CHOOSE_FONT_TYPE
	_            WORD
	SizeMin      INT
	SizeMax      INT
}

var lzCommDlgExtendedError = lzComdlg32.NewProc("CommDlgExtendedError")

type ComDlgExtError DWORD

func (err ComDlgExtError) Error() string {
	return fmt.Sprintf("CommDlgExtendedError: 0x%X", DWORD(err))
}

func CommDlgExtendedError() error {
	if r, _, _ := lzCommDlgExtendedError.Call(); r == 0 {
		return nil
	} else {
		return ComDlgExtError(r)
	}
}

var lzChooseFontW = lzComdlg32.NewProc("ChooseFontW")

func ChooseFontW(font *CHOOSEFONTW) (bool, error) {
	if r, _, _ := lzChooseFontW.Call(uintptr(unsafe.Pointer(font))); r == 0 {
		if err := CommDlgExtendedError(); err != nil {
			return false, err
		}
		return false, nil
	} else {
		return true, nil
	}
}
