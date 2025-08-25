package machineid

import "testing"

func TestIDStable(t *testing.T) {
	id1, err := ID()
	if err != nil {
		t.Fatalf("ID returned error: %v", err)
	}
	id2, err := ID()
	if err != nil {
		t.Fatalf("ID returned error: %v", err)
	}
	if id1 != id2 {
		t.Fatalf("expected stable ID, got %q and %q", id1, id2)
	}
	if id1 == "" {
		t.Fatal("ID is empty")
	}
}

func TestRawIDStable(t *testing.T) {
	b1, i1, _, _ := RawID()
	b2, i2, _, _ := RawID()
	if b1 != b2 || i1 != i2 {
		t.Fatalf("expected stable raw IDs, got %q/%q and %q/%q", b1, i1, b2, i2)
	}
}
