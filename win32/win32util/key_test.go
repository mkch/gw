package win32util

import "testing"

func TestKeyMessageLParamRepeatAndScanCode(t *testing.T) {
	var l KeyMessageLParam

	l.SetRepeatCount(0x1234)
	l.SetScanCode(0x56)

	if got := l.RepeatCount(); got != 0x1234 {
		t.Fatalf("RepeatCount() = %d, want %d", got, 0x1234)
	}
	if got := l.ScanCode(); got != 0x56 {
		t.Fatalf("ScanCode() = 0x%02X, want 0x56", got)
	}
}

func TestKeyMessageLParamFlags(t *testing.T) {
	var l KeyMessageLParam

	if l.Extended() || l.DialogMode() || l.MenuMode() || l.AltDown() || l.PreviousDown() || l.KeyUp() {
		t.Fatal("all flags should be false by default")
	}

	l.SetExtended(true)
	l.SetDialogMode(true)
	l.SetMenuMode(true)
	l.SetAltDown(true)
	l.SetPreviousDown(true)
	l.SetKeyUp(true)

	if !l.Extended() || !l.DialogMode() || !l.MenuMode() || !l.AltDown() || !l.PreviousDown() || !l.KeyUp() {
		t.Fatal("all flags should be true after setting to true")
	}

	l.SetExtended(false)
	l.SetDialogMode(false)
	l.SetMenuMode(false)
	l.SetAltDown(false)
	l.SetPreviousDown(false)
	l.SetKeyUp(false)

	if l.Extended() || l.DialogMode() || l.MenuMode() || l.AltDown() || l.PreviousDown() || l.KeyUp() {
		t.Fatal("all flags should be false after setting to false")
	}
}

func TestKeyMessageLParamFieldIsolation(t *testing.T) {
	var l KeyMessageLParam

	l.SetRepeatCount(0x00FF)
	l.SetScanCode(0x44)
	l.SetExtended(true)
	l.SetDialogMode(true)
	l.SetMenuMode(false)
	l.SetAltDown(true)
	l.SetPreviousDown(false)
	l.SetKeyUp(true)

	if got := l.RepeatCount(); got != 0x00FF {
		t.Fatalf("RepeatCount() changed unexpectedly, got %d", got)
	}
	if got := l.ScanCode(); got != 0x44 {
		t.Fatalf("ScanCode() changed unexpectedly, got 0x%02X", got)
	}
	if !l.Extended() || !l.DialogMode() || l.MenuMode() || !l.AltDown() || l.PreviousDown() || !l.KeyUp() {
		t.Fatal("flag values are not as expected")
	}
}

func TestKeyMessageLParamString(t *testing.T) {
	var l KeyMessageLParam
	l.SetRepeatCount(4660)
	l.SetScanCode(0x56)
	l.SetExtended(true)
	l.SetDialogMode(false)
	l.SetMenuMode(true)
	l.SetAltDown(false)
	l.SetPreviousDown(true)
	l.SetKeyUp(true)

	const want = "RepeatCount: 4660, ScanCode: 0x56, Extended: true, DialogMode: false, MenuMode: true, AltDown: false, PreviousDown: true, KeyUp: true"
	if got := l.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
