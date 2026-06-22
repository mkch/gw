package msghandler_test

import (
	"testing"

	"github.com/mkch/gw/internal/msghandler"
	"github.com/mkch/gw/win32"
)

func TestChain_NumHandlers(t *testing.T) {
	var chain msghandler.Chain
	if n := chain.NumHandlers(); n != 0 {
		t.Errorf("expected 0 handlers, got %d", n)
	}
	key1 := chain.AddHandler(func(arg *msghandler.Arg, callPrev func(*msghandler.Arg) win32.LRESULT) win32.LRESULT {
		return 0
	})
	if n := chain.NumHandlers(); n != 1 {
		t.Errorf("expected 1 handler, got %d", n)
	}
	chain.RemoveHandler(key1)
	if n := chain.NumHandlers(); n != 0 {
		t.Errorf("expected 0 handlers, got %d", n)
	}

	// Duplicate removal should not decrease the count below 0.
	chain.RemoveHandler(key1)
	if n := chain.NumHandlers(); n != 0 {
		t.Errorf("expected 0 handlers, got %d", n)
	}
}
