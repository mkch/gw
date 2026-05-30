package win32util_test

import (
	"testing"

	"github.com/mkch/gg"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

func TestKeyName(t *testing.T) {
	cases := []struct {
		name               string
		vk                 win32.VKCode
		doNotCareLeftRight bool
		want               string
		wantErr            bool
	}{
		// Standard keyboard keys.
		{name: "VK_CANCEL", vk: win32.VK_CANCEL, want: "Break"},
		{name: "VK_BACK", vk: win32.VK_BACK, want: "Backspace"},
		{name: "VK_TAB", vk: win32.VK_TAB, want: "Tab"},
		// VK_CLEAR maps to numpad-5 scan code (shared physical key).
		{name: "VK_CLEAR", vk: win32.VK_CLEAR, want: "Num 5"},
		{name: "VK_RETURN", vk: win32.VK_RETURN, want: "Enter"},

		// Modifier keys — generic VK, doNotCareLeftRight=false.
		{name: "VK_SHIFT", vk: win32.VK_SHIFT, want: "Shift"},
		{name: "VK_CONTROL", vk: win32.VK_CONTROL, want: "Ctrl"},
		{name: "VK_MENU", vk: win32.VK_MENU, want: "Alt"},
		// Modifier keys — generic VK, doNotCareLeftRight=true.
		{name: "VK_SHIFT/doNotCare", vk: win32.VK_SHIFT, doNotCareLeftRight: true, want: "Shift"},
		{name: "VK_CONTROL/doNotCare", vk: win32.VK_CONTROL, doNotCareLeftRight: true, want: "Ctrl"},
		{name: "VK_MENU/doNotCare", vk: win32.VK_MENU, doNotCareLeftRight: true, want: "Alt"},

		// Special keys.
		{name: "VK_PAUSE", vk: win32.VK_PAUSE, want: "Pause"},
		{name: "VK_CAPITAL", vk: win32.VK_CAPITAL, want: "Caps Lock"},
		{name: "VK_ESCAPE", vk: win32.VK_ESCAPE, want: "Esc"},
		{name: "VK_SPACE", vk: win32.VK_SPACE, want: "Space"},

		// Navigation / editing keys. MapVirtualKeyW returns non-extended scan codes
		// for these generic VKs, so GetKeyNameTextW returns the numpad key names.
		{name: "VK_PRIOR", vk: win32.VK_PRIOR, want: "Page Up"},
		{name: "VK_NEXT", vk: win32.VK_NEXT, want: "Page Down"},
		{name: "VK_END", vk: win32.VK_END, want: "End"},
		{name: "VK_HOME", vk: win32.VK_HOME, want: "Home"},
		{name: "VK_LEFT", vk: win32.VK_LEFT, want: "Left"},
		{name: "VK_UP", vk: win32.VK_UP, want: "Up"},
		{name: "VK_RIGHT", vk: win32.VK_RIGHT, want: "Right"},
		{name: "VK_DOWN", vk: win32.VK_DOWN, want: "Down"},
		{name: "VK_INSERT", vk: win32.VK_INSERT, want: "Insert"},
		{name: "VK_DELETE", vk: win32.VK_DELETE, want: "Delete"},

		// VK_SNAPSHOT: hardcoded to E0-37 to return "Prnt Scrn" instead of "Sys Req".
		{name: "VK_SNAPSHOT", vk: win32.VK_SNAPSHOT, want: "Prnt Scrn"},

		// Digit keys (VK code == ASCII code of '0'-'9').
		{name: "'0'", vk: '0', want: "0"},
		{name: "'1'", vk: '1', want: "1"},
		{name: "'2'", vk: '2', want: "2"},
		{name: "'3'", vk: '3', want: "3"},
		{name: "'4'", vk: '4', want: "4"},
		{name: "'5'", vk: '5', want: "5"},
		{name: "'6'", vk: '6', want: "6"},
		{name: "'7'", vk: '7', want: "7"},
		{name: "'8'", vk: '8', want: "8"},
		{name: "'9'", vk: '9', want: "9"},

		// Letter keys (VK code == ASCII code of 'A'-'Z').
		{name: "'A'", vk: 'A', want: "A"},
		{name: "'B'", vk: 'B', want: "B"},
		{name: "'C'", vk: 'C', want: "C"},
		{name: "'D'", vk: 'D', want: "D"},
		{name: "'E'", vk: 'E', want: "E"},
		{name: "'F'", vk: 'F', want: "F"},
		{name: "'G'", vk: 'G', want: "G"},
		{name: "'H'", vk: 'H', want: "H"},
		{name: "'I'", vk: 'I', want: "I"},
		{name: "'J'", vk: 'J', want: "J"},
		{name: "'K'", vk: 'K', want: "K"},
		{name: "'L'", vk: 'L', want: "L"},
		{name: "'M'", vk: 'M', want: "M"},
		{name: "'N'", vk: 'N', want: "N"},
		{name: "'O'", vk: 'O', want: "O"},
		{name: "'P'", vk: 'P', want: "P"},
		{name: "'Q'", vk: 'Q', want: "Q"},
		{name: "'R'", vk: 'R', want: "R"},
		{name: "'S'", vk: 'S', want: "S"},
		{name: "'T'", vk: 'T', want: "T"},
		{name: "'U'", vk: 'U', want: "U"},
		{name: "'V'", vk: 'V', want: "V"},
		{name: "'W'", vk: 'W', want: "W"},
		{name: "'X'", vk: 'X', want: "X"},
		{name: "'Y'", vk: 'Y', want: "Y"},
		{name: "'Z'", vk: 'Z', want: "Z"},

		// Windows / application keys.
		// VK_LWIN/VK_RWIN: different scan codes (5B vs 5C), so bit 25 has no effect.
		{name: "VK_LWIN/false", vk: win32.VK_LWIN, want: "Left Windows"},
		{name: "VK_LWIN/true", vk: win32.VK_LWIN, doNotCareLeftRight: true, want: "Left Windows"},
		{name: "VK_RWIN/false", vk: win32.VK_RWIN, want: "Right Windows"},
		{name: "VK_RWIN/true", vk: win32.VK_RWIN, doNotCareLeftRight: true, want: "Right Windows"},
		{name: "VK_APPS", vk: win32.VK_APPS, want: "Application"},

		// Numpad keys.
		{name: "VK_NUMPAD0", vk: win32.VK_NUMPAD0, want: "Num 0"},
		{name: "VK_NUMPAD1", vk: win32.VK_NUMPAD1, want: "Num 1"},
		{name: "VK_NUMPAD2", vk: win32.VK_NUMPAD2, want: "Num 2"},
		{name: "VK_NUMPAD3", vk: win32.VK_NUMPAD3, want: "Num 3"},
		{name: "VK_NUMPAD4", vk: win32.VK_NUMPAD4, want: "Num 4"},
		{name: "VK_NUMPAD5", vk: win32.VK_NUMPAD5, want: "Num 5"},
		{name: "VK_NUMPAD6", vk: win32.VK_NUMPAD6, want: "Num 6"},
		{name: "VK_NUMPAD7", vk: win32.VK_NUMPAD7, want: "Num 7"},
		{name: "VK_NUMPAD8", vk: win32.VK_NUMPAD8, want: "Num 8"},
		{name: "VK_NUMPAD9", vk: win32.VK_NUMPAD9, want: "Num 9"},
		{name: "VK_MULTIPLY", vk: win32.VK_MULTIPLY, want: "Num *"},
		{name: "VK_ADD", vk: win32.VK_ADD, want: "Num +"},
		{name: "VK_SUBTRACT", vk: win32.VK_SUBTRACT, want: "Num -"},
		{name: "VK_DECIMAL", vk: win32.VK_DECIMAL, want: "Num Del"},
		{name: "VK_DIVIDE", vk: win32.VK_DIVIDE, want: "Num /"},

		// Function keys.
		{name: "VK_F1", vk: win32.VK_F1, want: "F1"},
		{name: "VK_F2", vk: win32.VK_F2, want: "F2"},
		{name: "VK_F3", vk: win32.VK_F3, want: "F3"},
		{name: "VK_F4", vk: win32.VK_F4, want: "F4"},
		{name: "VK_F5", vk: win32.VK_F5, want: "F5"},
		{name: "VK_F6", vk: win32.VK_F6, want: "F6"},
		{name: "VK_F7", vk: win32.VK_F7, want: "F7"},
		{name: "VK_F8", vk: win32.VK_F8, want: "F8"},
		{name: "VK_F9", vk: win32.VK_F9, want: "F9"},
		{name: "VK_F10", vk: win32.VK_F10, want: "F10"},
		{name: "VK_F11", vk: win32.VK_F11, want: "F11"},
		{name: "VK_F12", vk: win32.VK_F12, want: "F12"},

		// Lock keys.
		// VK_NUMLOCK: hardcoded to E0-45 to work around MapVirtualKeyW bug.
		{name: "VK_NUMLOCK", vk: win32.VK_NUMLOCK, want: "Num Lock"},
		{name: "VK_SCROLL", vk: win32.VK_SCROLL, want: "Scroll Lock"},

		// Left/right modifier keys.
		// doNotCareLeftRight=true (bit 25) causes GetKeyNameTextW to return the
		// non-sided name for R-side keys that share the same base scan code.
		{name: "VK_LSHIFT/false", vk: win32.VK_LSHIFT, want: "Shift"},
		{name: "VK_LSHIFT/true", vk: win32.VK_LSHIFT, doNotCareLeftRight: true, want: "Shift"},
		{name: "VK_RSHIFT/false", vk: win32.VK_RSHIFT, want: "Right Shift"},
		{name: "VK_RSHIFT/true", vk: win32.VK_RSHIFT, doNotCareLeftRight: true, want: "Shift"},
		{name: "VK_LCONTROL/false", vk: win32.VK_LCONTROL, want: "Ctrl"},
		{name: "VK_LCONTROL/true", vk: win32.VK_LCONTROL, doNotCareLeftRight: true, want: "Ctrl"},
		{name: "VK_RCONTROL/false", vk: win32.VK_RCONTROL, want: "Right Ctrl"},
		{name: "VK_RCONTROL/true", vk: win32.VK_RCONTROL, doNotCareLeftRight: true, want: "Ctrl"},
		{name: "VK_LMENU/false", vk: win32.VK_LMENU, want: "Alt"},
		{name: "VK_LMENU/true", vk: win32.VK_LMENU, doNotCareLeftRight: true, want: "Alt"},
		{name: "VK_RMENU/false", vk: win32.VK_RMENU, want: "Right Alt"},
		{name: "VK_RMENU/true", vk: win32.VK_RMENU, doNotCareLeftRight: true, want: "Alt"},

		// Mouse buttons — no keyboard scan code.
		{name: "VK_LBUTTON", vk: win32.VK_LBUTTON, wantErr: true},
		{name: "VK_RBUTTON", vk: win32.VK_RBUTTON, wantErr: true},
		{name: "VK_MBUTTON", vk: win32.VK_MBUTTON, wantErr: true},
		{name: "VK_XBUTTON1", vk: win32.VK_XBUTTON1, wantErr: true},
		{name: "VK_XBUTTON2", vk: win32.VK_XBUTTON2, wantErr: true},

		// IME keys — no standard scan code on non-CJK layouts.
		{name: "VK_KANA", vk: win32.VK_KANA, wantErr: true},
		{name: "VK_IME_ON", vk: win32.VK_IME_ON, wantErr: true},
		{name: "VK_JUNJA", vk: win32.VK_JUNJA, wantErr: true},
		{name: "VK_FINAL", vk: win32.VK_FINAL, wantErr: true},
		{name: "VK_HANJA", vk: win32.VK_HANJA, wantErr: true},
		{name: "VK_IME_OFF", vk: win32.VK_IME_OFF, wantErr: true},
		{name: "VK_CONVERT", vk: win32.VK_CONVERT, wantErr: true},
		{name: "VK_NONCONVERT", vk: win32.VK_NONCONVERT, wantErr: true},
		{name: "VK_ACCEPT", vk: win32.VK_ACCEPT, wantErr: true},
		{name: "VK_MODECHANGE", vk: win32.VK_MODECHANGE, wantErr: true},

		{name: "VK_SELECT", vk: win32.VK_SELECT, wantErr: true},
		{name: "VK_PRINT", vk: win32.VK_PRINT, wantErr: true},
		{name: "VK_EXECUTE", vk: win32.VK_EXECUTE, wantErr: true},
		{name: "VK_HELP", vk: win32.VK_HELP, wantErr: true},

		{name: "VK_SLEEP", vk: win32.VK_SLEEP, wantErr: true},

		{name: "VK_SEPARATOR", vk: win32.VK_SEPARATOR, wantErr: true},

		// F13-F24 have no scan codes on standard PC keyboards.
		{name: "VK_F13", vk: win32.VK_F13, wantErr: true},
		{name: "VK_F14", vk: win32.VK_F14, wantErr: true},
		{name: "VK_F15", vk: win32.VK_F15, wantErr: true},
		{name: "VK_F16", vk: win32.VK_F16, wantErr: true},
		{name: "VK_F17", vk: win32.VK_F17, wantErr: true},
		{name: "VK_F18", vk: win32.VK_F18, wantErr: true},
		{name: "VK_F19", vk: win32.VK_F19, wantErr: true},
		{name: "VK_F20", vk: win32.VK_F20, wantErr: true},
		{name: "VK_F21", vk: win32.VK_F21, wantErr: true},
		{name: "VK_F22", vk: win32.VK_F22, wantErr: true},
		{name: "VK_F23", vk: win32.VK_F23, wantErr: true},
		{name: "VK_F24", vk: win32.VK_F24, wantErr: true},

		// UI Navigation keys (reserved, no scan codes).
		{name: "VK_NAVIGATION_VIEW", vk: win32.VK_NAVIGATION_VIEW, wantErr: true},
		{name: "VK_NAVIGATION_MENU", vk: win32.VK_NAVIGATION_MENU, wantErr: true},
		{name: "VK_NAVIGATION_UP", vk: win32.VK_NAVIGATION_UP, wantErr: true},
		{name: "VK_NAVIGATION_DOWN", vk: win32.VK_NAVIGATION_DOWN, wantErr: true},
		{name: "VK_NAVIGATION_LEFT", vk: win32.VK_NAVIGATION_LEFT, wantErr: true},
		{name: "VK_NAVIGATION_RIGHT", vk: win32.VK_NAVIGATION_RIGHT, wantErr: true},
		{name: "VK_NAVIGATION_ACCEPT", vk: win32.VK_NAVIGATION_ACCEPT, wantErr: true},
		{name: "VK_NAVIGATION_CANCEL", vk: win32.VK_NAVIGATION_CANCEL, wantErr: true},

		// OEM / Fujitsu keys — no standard scan codes on a PC keyboard.
		{name: "VK_OEM_NEC_EQUAL", vk: win32.VK_OEM_NEC_EQUAL, wantErr: true},
		{name: "VK_OEM_FJ_MASSHOU", vk: win32.VK_OEM_FJ_MASSHOU, wantErr: true},
		{name: "VK_OEM_FJ_TOUROKU", vk: win32.VK_OEM_FJ_TOUROKU, wantErr: true},
		{name: "VK_OEM_FJ_LOYA", vk: win32.VK_OEM_FJ_LOYA, wantErr: true},
		{name: "VK_OEM_FJ_ROYA", vk: win32.VK_OEM_FJ_ROYA, wantErr: true},

		// Browser keys — no standard scan codes.
		{name: "VK_BROWSER_BACK", vk: win32.VK_BROWSER_BACK, wantErr: true},
		{name: "VK_BROWSER_FORWARD", vk: win32.VK_BROWSER_FORWARD, wantErr: true},
		{name: "VK_BROWSER_REFRESH", vk: win32.VK_BROWSER_REFRESH, wantErr: true},
		{name: "VK_BROWSER_STOP", vk: win32.VK_BROWSER_STOP, wantErr: true},
		{name: "VK_BROWSER_SEARCH", vk: win32.VK_BROWSER_SEARCH, wantErr: true},
		{name: "VK_BROWSER_FAVORITES", vk: win32.VK_BROWSER_FAVORITES, wantErr: true},
		// VK_BROWSER_HOME and the following media/launch keys: MapVirtualKeyW returns
		// scan codes that coincide with regular key positions, yielding letter names.
		{name: "VK_BROWSER_HOME", vk: win32.VK_BROWSER_HOME, want: "M"},
		{name: "VK_VOLUME_MUTE", vk: win32.VK_VOLUME_MUTE, want: "D"},
		{name: "VK_VOLUME_DOWN", vk: win32.VK_VOLUME_DOWN, want: "C"},
		{name: "VK_VOLUME_UP", vk: win32.VK_VOLUME_UP, want: "B"},
		{name: "VK_MEDIA_NEXT_TRACK", vk: win32.VK_MEDIA_NEXT_TRACK, want: "P"},
		{name: "VK_MEDIA_PREV_TRACK", vk: win32.VK_MEDIA_PREV_TRACK, want: "Q"},
		{name: "VK_MEDIA_STOP", vk: win32.VK_MEDIA_STOP, want: "J"},
		{name: "VK_MEDIA_PLAY_PAUSE", vk: win32.VK_MEDIA_PLAY_PAUSE, want: "G"},
		{name: "VK_LAUNCH_MAIL", vk: win32.VK_LAUNCH_MAIL, wantErr: true},
		{name: "VK_LAUNCH_MEDIA_SELECT", vk: win32.VK_LAUNCH_MEDIA_SELECT, wantErr: true},
		{name: "VK_LAUNCH_APP1", vk: win32.VK_LAUNCH_APP1, wantErr: true},
		{name: "VK_LAUNCH_APP2", vk: win32.VK_LAUNCH_APP2, want: "F"},

		// OEM punctuation keys (US layout).
		{name: "VK_OEM_1", vk: win32.VK_OEM_1, want: ";"},
		{name: "VK_OEM_PLUS", vk: win32.VK_OEM_PLUS, want: "="},
		{name: "VK_OEM_COMMA", vk: win32.VK_OEM_COMMA, want: ","},
		{name: "VK_OEM_MINUS", vk: win32.VK_OEM_MINUS, want: "-"},
		{name: "VK_OEM_PERIOD", vk: win32.VK_OEM_PERIOD, want: "."},
		{name: "VK_OEM_2", vk: win32.VK_OEM_2, want: "/"},
		{name: "VK_OEM_3", vk: win32.VK_OEM_3, want: "`"},
		{name: "VK_OEM_4", vk: win32.VK_OEM_4, want: "["},
		{name: "VK_OEM_5", vk: win32.VK_OEM_5, want: "\\"},
		{name: "VK_OEM_6", vk: win32.VK_OEM_6, want: "]"},
		{name: "VK_OEM_7", vk: win32.VK_OEM_7, want: "'"},
		{name: "VK_OEM_102", vk: win32.VK_OEM_102, want: "\\"},

		{name: "VK_OEM_8", vk: win32.VK_OEM_8, wantErr: true},

		// Gamepad keys (reserved, no scan codes).
		{name: "VK_GAMEPAD_A", vk: win32.VK_GAMEPAD_A, wantErr: true},
		{name: "VK_GAMEPAD_B", vk: win32.VK_GAMEPAD_B, wantErr: true},
		{name: "VK_GAMEPAD_X", vk: win32.VK_GAMEPAD_X, wantErr: true},
		{name: "VK_GAMEPAD_Y", vk: win32.VK_GAMEPAD_Y, wantErr: true},
		{name: "VK_GAMEPAD_RIGHT_SHOULDER", vk: win32.VK_GAMEPAD_RIGHT_SHOULDER, wantErr: true},
		{name: "VK_GAMEPAD_LEFT_SHOULDER", vk: win32.VK_GAMEPAD_LEFT_SHOULDER, wantErr: true},
		{name: "VK_GAMEPAD_LEFT_TRIGGER", vk: win32.VK_GAMEPAD_LEFT_TRIGGER, wantErr: true},
		{name: "VK_GAMEPAD_RIGHT_TRIGGER", vk: win32.VK_GAMEPAD_RIGHT_TRIGGER, wantErr: true},
		{name: "VK_GAMEPAD_DPAD_UP", vk: win32.VK_GAMEPAD_DPAD_UP, wantErr: true},
		{name: "VK_GAMEPAD_DPAD_DOWN", vk: win32.VK_GAMEPAD_DPAD_DOWN, wantErr: true},
		{name: "VK_GAMEPAD_DPAD_LEFT", vk: win32.VK_GAMEPAD_DPAD_LEFT, wantErr: true},
		{name: "VK_GAMEPAD_DPAD_RIGHT", vk: win32.VK_GAMEPAD_DPAD_RIGHT, wantErr: true},
		{name: "VK_GAMEPAD_MENU", vk: win32.VK_GAMEPAD_MENU, wantErr: true},
		{name: "VK_GAMEPAD_VIEW", vk: win32.VK_GAMEPAD_VIEW, wantErr: true},
		{name: "VK_GAMEPAD_LEFT_THUMBSTICK_BUTTON", vk: win32.VK_GAMEPAD_LEFT_THUMBSTICK_BUTTON, wantErr: true},
		{name: "VK_GAMEPAD_RIGHT_THUMBSTICK_BUTTON", vk: win32.VK_GAMEPAD_RIGHT_THUMBSTICK_BUTTON, wantErr: true},
		{name: "VK_GAMEPAD_LEFT_THUMBSTICK_UP", vk: win32.VK_GAMEPAD_LEFT_THUMBSTICK_UP, wantErr: true},
		{name: "VK_GAMEPAD_LEFT_THUMBSTICK_DOWN", vk: win32.VK_GAMEPAD_LEFT_THUMBSTICK_DOWN, wantErr: true},
		{name: "VK_GAMEPAD_LEFT_THUMBSTICK_RIGHT", vk: win32.VK_GAMEPAD_LEFT_THUMBSTICK_RIGHT, wantErr: true},
		{name: "VK_GAMEPAD_LEFT_THUMBSTICK_LEFT", vk: win32.VK_GAMEPAD_LEFT_THUMBSTICK_LEFT, wantErr: true},
		{name: "VK_GAMEPAD_RIGHT_THUMBSTICK_UP", vk: win32.VK_GAMEPAD_RIGHT_THUMBSTICK_UP, wantErr: true},
		{name: "VK_GAMEPAD_RIGHT_THUMBSTICK_DOWN", vk: win32.VK_GAMEPAD_RIGHT_THUMBSTICK_DOWN, wantErr: true},
		{name: "VK_GAMEPAD_RIGHT_THUMBSTICK_RIGHT", vk: win32.VK_GAMEPAD_RIGHT_THUMBSTICK_RIGHT, wantErr: true},
		{name: "VK_GAMEPAD_RIGHT_THUMBSTICK_LEFT", vk: win32.VK_GAMEPAD_RIGHT_THUMBSTICK_LEFT, wantErr: true},

		// Miscellaneous OEM / vendor keys.
		{name: "VK_OEM_AX", vk: win32.VK_OEM_AX, wantErr: true},
		{name: "VK_ICO_HELP", vk: win32.VK_ICO_HELP, wantErr: true},
		{name: "VK_ICO_00", vk: win32.VK_ICO_00, wantErr: true},
		{name: "VK_PROCESSKEY", vk: win32.VK_PROCESSKEY, wantErr: true},
		{name: "VK_ICO_CLEAR", vk: win32.VK_ICO_CLEAR, wantErr: true},
		{name: "VK_PACKET", vk: win32.VK_PACKET, wantErr: true},
		{name: "VK_OEM_RESET", vk: win32.VK_OEM_RESET, wantErr: true},
		{name: "VK_OEM_JUMP", vk: win32.VK_OEM_JUMP, wantErr: true},
		{name: "VK_OEM_PA1", vk: win32.VK_OEM_PA1, wantErr: true},
		{name: "VK_OEM_PA2", vk: win32.VK_OEM_PA2, wantErr: true},
		{name: "VK_OEM_PA3", vk: win32.VK_OEM_PA3, wantErr: true},
		{name: "VK_OEM_WSCTRL", vk: win32.VK_OEM_WSCTRL, wantErr: true},
		{name: "VK_OEM_CUSEL", vk: win32.VK_OEM_CUSEL, wantErr: true},
		{name: "VK_OEM_ATTN", vk: win32.VK_OEM_ATTN, wantErr: true},
		{name: "VK_OEM_FINISH", vk: win32.VK_OEM_FINISH, wantErr: true},
		{name: "VK_OEM_COPY", vk: win32.VK_OEM_COPY, wantErr: true},
		{name: "VK_OEM_AUTO", vk: win32.VK_OEM_AUTO, wantErr: true},
		{name: "VK_OEM_ENLW", vk: win32.VK_OEM_ENLW, wantErr: true},
		{name: "VK_OEM_BACKTAB", vk: win32.VK_OEM_BACKTAB, wantErr: true},
		{name: "VK_ATTN", vk: win32.VK_ATTN, wantErr: true},
		{name: "VK_CRSEL", vk: win32.VK_CRSEL, wantErr: true},
		{name: "VK_EXSEL", vk: win32.VK_EXSEL, wantErr: true},
		{name: "VK_EREOF", vk: win32.VK_EREOF, wantErr: true},
		{name: "VK_PLAY", vk: win32.VK_PLAY, wantErr: true},
		{name: "VK_ZOOM", vk: win32.VK_ZOOM, wantErr: true},
		{name: "VK_NONAME", vk: win32.VK_NONAME, wantErr: true},
		{name: "VK_PA1", vk: win32.VK_PA1, wantErr: true},
		{name: "VK_OEM_CLEAR", vk: win32.VK_OEM_CLEAR, wantErr: true},
	}

	// Use en-US layout to ensure consistent scan codes for VKs that have different key names
	var buf []win32.WCHAR
	win32util.CString("00000409", &buf) // 0x0409 = en-US
	en_US_HKL := gg.Must(win32.LoadKeyboardLayoutW(&buf[0], win32.KLF_REPLACELANG|win32.KLF_SUBSTITUTE_OK))

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := win32util.KeyName(tc.vk, tc.doNotCareLeftRight, en_US_HKL)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
