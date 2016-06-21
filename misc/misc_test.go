package misc

import (
	"testing"
)

func TestIsIdValid(t *testing.T) {
	for _, v := range []int{1, 5, 2, 9, 193, 17} {
		if !IsIdValid(v) {
			t.Errorf("Id %v should be valid", v)
		}
	}

	for _, v := range []int{-9, -3, -1, 0, -93145} {
		if IsIdValid(v) {
			t.Errorf("Id %v should be invalid", v)
		}
	}
}

func TestIsPasswordValid(t *testing.T) {
	for _, v := range []string{"password", "a123fsdf3", "  sadf3fs", "13ds45sdfdfadf"} {
		if !IsPasswordValid(v) {
			t.Errorf("Id %v should be valid", v)
		}
	}

	for _, v := range []string{"1234567", "", "345", "uestheI", "35", "g"} {
		if IsPasswordValid(v) {
			t.Errorf("Password %v should be valid", v)
		}
	}
}

func TestValidateString(t *testing.T) {
	table := []struct {
		input  string
		n      int
		output string
		ok     bool
	}{
		{"something", 9, "something", true},
		{"  something", 9, "something", true},
		{"  something   ", 9, "something", true},
		{"  something", 8, "", false},
		{"12  something ds", 16, "12  something ds", true},
		{"         ", 1, "", false},
		{"		sdfsd", 5, "sdfsd", true},
	}
	for _, v := range table {
		s, ok := ValidateString(v.input, v.n)
		if s != v.output || ok != v.ok {
			t.Errorf("Expected %v, %v. Got %v, %v", v.output, v.ok, s, ok)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	table := []struct {
		input  string
		output string
		ok     bool
	}{
		{"some@email.com", "some@email.com", true},
		{" sOMe@EMaiL.com  ", "some@email.com", true},
		{" SOME@EMAIL.cOm  ", "some@email.com", true},
		{" SOME@EMAIL123.cOm  ", "some@email123.com", true},
		{" SOMEEMAIL123.cOm  ", "", false},
		{"SOMEEMAIL123.cOm", "", false},
	}

	for _, v := range table {
		s, ok := ValidateEmail(v.input)
		if s != v.output || ok != v.ok {
			t.Errorf("Expected %v, %v. Got %v, %v", v.output, v.ok, s, ok)
		}
	}
}
