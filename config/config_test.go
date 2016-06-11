package config

import (
	"os"
	"testing"
)

func TestEnvStrPresent(t *testing.T) {
	os.Setenv("PROJ_FAKE_ENV", "randomness")
	if v := GetEnvStr("PROJ_FAKE_ENV"); v != "randomness" {
		t.Errorf("Expected 'randomness', got %s", v)
	}
	os.Unsetenv("PROJ_FAKE_ENV")
}

func TestEnvStrMissing(t *testing.T) {
	os.Clearenv()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, environment variable is missing")
		}
	}()
	GetEnvStr("PROJ_FAKE_ENV")
}

func TestRequiredIntEnvIsPresent(t *testing.T) {
	os.Setenv("PROJ_FAKE_ENV", "123456")
	if v := GetEnvInt("PROJ_FAKE_ENV"); v != 123456 {
		t.Errorf("Expected integer 123456, got %d", v)
	}
}

func TestRequiredIntEnvIsNotInt(t *testing.T) {
	os.Setenv("PROJ_FAKE_ENV", "124sdf")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic from environemnt variable which is not integer")
		}
	}()
	t.Errorf("Expected panic, got: %d", GetEnvInt("PROJ_FAKE_ENV"))
}

func TestRequiredIntEnvIsMissing(t *testing.T) {
	os.Clearenv()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, env is missing")
		}
	}()
	GetEnvInt("PROJ_FAKE_ENV")
}
