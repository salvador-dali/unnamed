package testHelpers

import "testing"

func TestRandomString(t *testing.T) {
	length := 50

	// bigger, not on the edge
	for i := 0; i < 300; i++ {
		s := RandomString(length, 1, 0)
		if len(s) <= length {
			t.Errorf("String should be bigger %v, got %v. IMPORTANT function is random, investigate, do not blindly rerun", length, len(s))
		}
	}

	// bigger, on the edge
	for i := 0; i < 20; i++ {
		s := RandomString(length, 1, 1)
		if len(s) != length+1 {
			t.Errorf("String should be equal to %v, got %v. IMPORTANT function is random, investigate, do not blindly rerun", length+1, len(s))
		}
	}

	// smaller, not on the edge
	for i := 0; i < 300; i++ {
		s := RandomString(length, 0, 0)
		if len(s) > length {
			t.Errorf("String should be smaller or equal to %v, got %v. IMPORTANT function is random, investigate, do not blindly rerun", length, len(s))
		}
	}

	// smaller, on the edge
	for i := 0; i < 20; i++ {
		s := RandomString(length, 0, 1)
		if len(s) != length {
			t.Errorf("String should be equal to %v, got %v. IMPORTANT function is random, investigate, do not blindly rerun", length, len(s))
		}
	}
}