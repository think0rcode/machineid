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
