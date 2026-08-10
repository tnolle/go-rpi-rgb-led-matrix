package animation

import "testing"

func TestDarts180UsesKebabCaseName(t *testing.T) {
	if got := Darts_180.String(); got != "darts-180" {
		t.Fatalf("unexpected animation name: got %q, want %q", got, "darts-180")
	}
	parsed, err := AnimationString("darts-180")
	if err != nil {
		t.Fatal(err)
	}
	if parsed != Darts_180 {
		t.Fatalf("unexpected animation: got %v, want %v", parsed, Darts_180)
	}
}
