package events

import "testing"

func TestMouseClickOptDefaults(t *testing.T) {
	var m MouseClickOpt

	if m.Control() || m.LButton() || m.MButton() || m.RButton() || m.Shift() || m.XButton1() || m.XButton2() {
		t.Fatal("all mouse options should be false by default")
	}
}

func TestMouseClickOptSetters(t *testing.T) {
	var m MouseClickOpt

	m.SetControl(true)
	m.SetLButton(true)
	m.SetMButton(true)
	m.SetRButton(true)
	m.SetShift(true)
	m.SetXButton1(true)
	m.SetXButton2(true)

	if !m.Control() || !m.LButton() || !m.MButton() || !m.RButton() || !m.Shift() || !m.XButton1() || !m.XButton2() {
		t.Fatal("all mouse options should be true after setting to true")
	}

	m.SetControl(false)
	m.SetLButton(false)
	m.SetMButton(false)
	m.SetRButton(false)
	m.SetShift(false)
	m.SetXButton1(false)
	m.SetXButton2(false)

	if m.Control() || m.LButton() || m.MButton() || m.RButton() || m.Shift() || m.XButton1() || m.XButton2() {
		t.Fatal("all mouse options should be false after setting to false")
	}
}

func TestMouseClickOptFieldIsolation(t *testing.T) {
	var m MouseClickOpt

	m.SetControl(true)
	m.SetShift(true)
	m.SetXButton2(true)

	if !m.Control() {
		t.Fatal("Control() should be true")
	}
	if m.LButton() || m.MButton() || m.RButton() || m.XButton1() {
		t.Fatal("unrelated button flags should remain false")
	}
	if !m.Shift() {
		t.Fatal("Shift() should be true")
	}
	if !m.XButton2() {
		t.Fatal("XButton2() should be true")
	}

	m.SetShift(false)
	if !m.Control() || !m.XButton2() {
		t.Fatal("clearing Shift should not affect other flags")
	}
	if m.Shift() {
		t.Fatal("Shift() should be false")
	}
}
