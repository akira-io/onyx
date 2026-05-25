package appearance

import "testing"

func TestIsDark_ReturnsWithoutPanic(t *testing.T) {
	first := IsDark()
	second := IsDark()
	if first != second {
		t.Fatalf("IsDark not stable within a test run: %v then %v", first, second)
	}
}
