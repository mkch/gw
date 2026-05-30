package win32

type VKCode byte

const (
	VK_LBUTTON VKCode = 0x01
	VK_RBUTTON VKCode = 0x02
	VK_CANCEL  VKCode = 0x03
	VK_MBUTTON VKCode = 0x04 // NOT contiguous with L & RBUTTON

	VK_XBUTTON1 VKCode = 0x05 // NOT contiguous with L & RBUTTON
	VK_XBUTTON2 VKCode = 0x06 // NOT contiguous with L & RBUTTON

	// 0x07 : reserved

	VK_BACK VKCode = 0x08
	VK_TAB  VKCode = 0x09

	// 0x0A - 0x0B : reserved

	VK_CLEAR  VKCode = 0x0C
	VK_RETURN VKCode = 0x0D

	// 0x0E - 0x0F : unassigned

	VK_SHIFT   VKCode = 0x10
	VK_CONTROL VKCode = 0x11
	VK_MENU    VKCode = 0x12
	VK_PAUSE   VKCode = 0x13
	VK_CAPITAL VKCode = 0x14

	VK_KANA    VKCode = 0x15
	VK_HANGEUL VKCode = 0x15 // old name - should be here for compatibility
	VK_HANGUL  VKCode = 0x15
	VK_IME_ON  VKCode = 0x16
	VK_JUNJA   VKCode = 0x17
	VK_FINAL   VKCode = 0x18
	VK_HANJA   VKCode = 0x19
	VK_KANJI   VKCode = 0x19
	VK_IME_OFF VKCode = 0x1A

	VK_ESCAPE VKCode = 0x1B

	VK_CONVERT    VKCode = 0x1C
	VK_NONCONVERT VKCode = 0x1D
	VK_ACCEPT     VKCode = 0x1E
	VK_MODECHANGE VKCode = 0x1F

	VK_SPACE    VKCode = 0x20
	VK_PRIOR    VKCode = 0x21
	VK_NEXT     VKCode = 0x22
	VK_END      VKCode = 0x23
	VK_HOME     VKCode = 0x24
	VK_LEFT     VKCode = 0x25
	VK_UP       VKCode = 0x26
	VK_RIGHT    VKCode = 0x27
	VK_DOWN     VKCode = 0x28
	VK_SELECT   VKCode = 0x29
	VK_PRINT    VKCode = 0x2A
	VK_EXECUTE  VKCode = 0x2B
	VK_SNAPSHOT VKCode = 0x2C
	VK_INSERT   VKCode = 0x2D
	VK_DELETE   VKCode = 0x2E
	VK_HELP     VKCode = 0x2F

	// VK_0 - VK_9 are the same as ASCII '0' - '9' (0x30 - 0x39)
	// 0x3A - 0x40 : unassigned
	// VK_A - VK_Z are the same as ASCII 'A' - 'Z' (0x41 - 0x5A)

	VK_LWIN VKCode = 0x5B
	VK_RWIN VKCode = 0x5C
	VK_APPS VKCode = 0x5D

	// 0x5E : reserved

	VK_SLEEP VKCode = 0x5F

	VK_NUMPAD0   VKCode = 0x60
	VK_NUMPAD1   VKCode = 0x61
	VK_NUMPAD2   VKCode = 0x62
	VK_NUMPAD3   VKCode = 0x63
	VK_NUMPAD4   VKCode = 0x64
	VK_NUMPAD5   VKCode = 0x65
	VK_NUMPAD6   VKCode = 0x66
	VK_NUMPAD7   VKCode = 0x67
	VK_NUMPAD8   VKCode = 0x68
	VK_NUMPAD9   VKCode = 0x69
	VK_MULTIPLY  VKCode = 0x6A
	VK_ADD       VKCode = 0x6B
	VK_SEPARATOR VKCode = 0x6C
	VK_SUBTRACT  VKCode = 0x6D
	VK_DECIMAL   VKCode = 0x6E
	VK_DIVIDE    VKCode = 0x6F
	VK_F1        VKCode = 0x70
	VK_F2        VKCode = 0x71
	VK_F3        VKCode = 0x72
	VK_F4        VKCode = 0x73
	VK_F5        VKCode = 0x74
	VK_F6        VKCode = 0x75
	VK_F7        VKCode = 0x76
	VK_F8        VKCode = 0x77
	VK_F9        VKCode = 0x78
	VK_F10       VKCode = 0x79
	VK_F11       VKCode = 0x7A
	VK_F12       VKCode = 0x7B
	VK_F13       VKCode = 0x7C
	VK_F14       VKCode = 0x7D
	VK_F15       VKCode = 0x7E
	VK_F16       VKCode = 0x7F
	VK_F17       VKCode = 0x80
	VK_F18       VKCode = 0x81
	VK_F19       VKCode = 0x82
	VK_F20       VKCode = 0x83
	VK_F21       VKCode = 0x84
	VK_F22       VKCode = 0x85
	VK_F23       VKCode = 0x86
	VK_F24       VKCode = 0x87

	// 0x88 - 0x8F : UI navigation

	VK_NAVIGATION_VIEW   VKCode = 0x88 // reserved
	VK_NAVIGATION_MENU   VKCode = 0x89 // reserved
	VK_NAVIGATION_UP     VKCode = 0x8A // reserved
	VK_NAVIGATION_DOWN   VKCode = 0x8B // reserved
	VK_NAVIGATION_LEFT   VKCode = 0x8C // reserved
	VK_NAVIGATION_RIGHT  VKCode = 0x8D // reserved
	VK_NAVIGATION_ACCEPT VKCode = 0x8E // reserved
	VK_NAVIGATION_CANCEL VKCode = 0x8F // reserved

	VK_NUMLOCK VKCode = 0x90
	VK_SCROLL  VKCode = 0x91

	// NEC PC-9800 kbd definitions
	VK_OEM_NEC_EQUAL VKCode = 0x92 // '=' key on numpad

	// Fujitsu/OASYS kbd definitions
	VK_OEM_FJ_JISHO   VKCode = 0x92 // 'Dictionary' key
	VK_OEM_FJ_MASSHOU VKCode = 0x93 // 'Unregister word' key
	VK_OEM_FJ_TOUROKU VKCode = 0x94 // 'Register word' key
	VK_OEM_FJ_LOYA    VKCode = 0x95 // 'Left OYAYUBI' key
	VK_OEM_FJ_ROYA    VKCode = 0x96 // 'Right OYAYUBI' key

	// 0x97 - 0x9F : unassigned

	// VK_L* & VK_R* - left and right Alt, Ctrl and Shift virtual keys.
	// Used only as parameters to GetAsyncKeyState() and GetKeyState().
	// No other API or message will distinguish left and right keys in this way.
	VK_LSHIFT   VKCode = 0xA0
	VK_RSHIFT   VKCode = 0xA1
	VK_LCONTROL VKCode = 0xA2
	VK_RCONTROL VKCode = 0xA3
	VK_LMENU    VKCode = 0xA4
	VK_RMENU    VKCode = 0xA5

	VK_BROWSER_BACK      VKCode = 0xA6
	VK_BROWSER_FORWARD   VKCode = 0xA7
	VK_BROWSER_REFRESH   VKCode = 0xA8
	VK_BROWSER_STOP      VKCode = 0xA9
	VK_BROWSER_SEARCH    VKCode = 0xAA
	VK_BROWSER_FAVORITES VKCode = 0xAB
	VK_BROWSER_HOME      VKCode = 0xAC

	VK_VOLUME_MUTE         VKCode = 0xAD
	VK_VOLUME_DOWN         VKCode = 0xAE
	VK_VOLUME_UP           VKCode = 0xAF
	VK_MEDIA_NEXT_TRACK    VKCode = 0xB0
	VK_MEDIA_PREV_TRACK    VKCode = 0xB1
	VK_MEDIA_STOP          VKCode = 0xB2
	VK_MEDIA_PLAY_PAUSE    VKCode = 0xB3
	VK_LAUNCH_MAIL         VKCode = 0xB4
	VK_LAUNCH_MEDIA_SELECT VKCode = 0xB5
	VK_LAUNCH_APP1         VKCode = 0xB6
	VK_LAUNCH_APP2         VKCode = 0xB7

	// 0xB8 - 0xB9 : reserved

	VK_OEM_1      VKCode = 0xBA // ';:' for US
	VK_OEM_PLUS   VKCode = 0xBB // '+' any country
	VK_OEM_COMMA  VKCode = 0xBC // ',' any country
	VK_OEM_MINUS  VKCode = 0xBD // '-' any country
	VK_OEM_PERIOD VKCode = 0xBE // '.' any country
	VK_OEM_2      VKCode = 0xBF // '/?' for US
	VK_OEM_3      VKCode = 0xC0 // '`~' for US

	// 0xC1 - 0xC2 : reserved

	// 0xC3 - 0xDA : Gamepad input

	VK_GAMEPAD_A                       VKCode = 0xC3 // reserved
	VK_GAMEPAD_B                       VKCode = 0xC4 // reserved
	VK_GAMEPAD_X                       VKCode = 0xC5 // reserved
	VK_GAMEPAD_Y                       VKCode = 0xC6 // reserved
	VK_GAMEPAD_RIGHT_SHOULDER          VKCode = 0xC7 // reserved
	VK_GAMEPAD_LEFT_SHOULDER           VKCode = 0xC8 // reserved
	VK_GAMEPAD_LEFT_TRIGGER            VKCode = 0xC9 // reserved
	VK_GAMEPAD_RIGHT_TRIGGER           VKCode = 0xCA // reserved
	VK_GAMEPAD_DPAD_UP                 VKCode = 0xCB // reserved
	VK_GAMEPAD_DPAD_DOWN               VKCode = 0xCC // reserved
	VK_GAMEPAD_DPAD_LEFT               VKCode = 0xCD // reserved
	VK_GAMEPAD_DPAD_RIGHT              VKCode = 0xCE // reserved
	VK_GAMEPAD_MENU                    VKCode = 0xCF // reserved
	VK_GAMEPAD_VIEW                    VKCode = 0xD0 // reserved
	VK_GAMEPAD_LEFT_THUMBSTICK_BUTTON  VKCode = 0xD1 // reserved
	VK_GAMEPAD_RIGHT_THUMBSTICK_BUTTON VKCode = 0xD2 // reserved
	VK_GAMEPAD_LEFT_THUMBSTICK_UP      VKCode = 0xD3 // reserved
	VK_GAMEPAD_LEFT_THUMBSTICK_DOWN    VKCode = 0xD4 // reserved
	VK_GAMEPAD_LEFT_THUMBSTICK_RIGHT   VKCode = 0xD5 // reserved
	VK_GAMEPAD_LEFT_THUMBSTICK_LEFT    VKCode = 0xD6 // reserved
	VK_GAMEPAD_RIGHT_THUMBSTICK_UP     VKCode = 0xD7 // reserved
	VK_GAMEPAD_RIGHT_THUMBSTICK_DOWN   VKCode = 0xD8 // reserved
	VK_GAMEPAD_RIGHT_THUMBSTICK_RIGHT  VKCode = 0xD9 // reserved
	VK_GAMEPAD_RIGHT_THUMBSTICK_LEFT   VKCode = 0xDA // reserved

	VK_OEM_4 VKCode = 0xDB //  '[{' for US
	VK_OEM_5 VKCode = 0xDC //  '\|' for US
	VK_OEM_6 VKCode = 0xDD //  ']}' for US
	VK_OEM_7 VKCode = 0xDE //  ''"' for US
	VK_OEM_8 VKCode = 0xDF

	// 0xE0 : reserved

	// Various extended or enhanced keyboards
	VK_OEM_AX   VKCode = 0xE1 //  'AX' key on Japanese AX kbd
	VK_OEM_102  VKCode = 0xE2 //  "<>" or "\|" on RT 102-key kbd.
	VK_ICO_HELP VKCode = 0xE3 //  Help key on ICO
	VK_ICO_00   VKCode = 0xE4 //  00 key on ICO

	VK_PROCESSKEY VKCode = 0xE5

	VK_ICO_CLEAR VKCode = 0xE6

	VK_PACKET VKCode = 0xE7

	// 0xE8 : unassigned

	// Nokia/Ericsson definitions
	VK_OEM_RESET   VKCode = 0xE9
	VK_OEM_JUMP    VKCode = 0xEA
	VK_OEM_PA1     VKCode = 0xEB
	VK_OEM_PA2     VKCode = 0xEC
	VK_OEM_PA3     VKCode = 0xED
	VK_OEM_WSCTRL  VKCode = 0xEE
	VK_OEM_CUSEL   VKCode = 0xEF
	VK_OEM_ATTN    VKCode = 0xF0
	VK_OEM_FINISH  VKCode = 0xF1
	VK_OEM_COPY    VKCode = 0xF2
	VK_OEM_AUTO    VKCode = 0xF3
	VK_OEM_ENLW    VKCode = 0xF4
	VK_OEM_BACKTAB VKCode = 0xF5

	VK_ATTN      VKCode = 0xF6
	VK_CRSEL     VKCode = 0xF7
	VK_EXSEL     VKCode = 0xF8
	VK_EREOF     VKCode = 0xF9
	VK_PLAY      VKCode = 0xFA
	VK_ZOOM      VKCode = 0xFB
	VK_NONAME    VKCode = 0xFC
	VK_PA1       VKCode = 0xFD
	VK_OEM_CLEAR VKCode = 0xFE
)
