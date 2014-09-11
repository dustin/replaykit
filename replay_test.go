package replay

import "testing"

func TestNew(t *testing.T) {
	New(1)
	panicked := false
	func() {
		defer func() { panicked = recover() != nil }()
		r := New(0)
		t.Errorf("Failed to panic at 0, got %#v", r)
	}()
	if !panicked {
		t.Errorf("Failed to panick at zero")
	}
	func() {
		defer func() { panicked = recover() != nil }()
		r := New(-1)
		t.Errorf("Failed to panic at -1, got %#v", r)
	}()
	if !panicked {
		t.Errorf("Failed to panick at -1")
	}
}
