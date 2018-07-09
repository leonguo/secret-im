package tests

import "testing"
import "../util"

func TestValidPhone(t *testing.T) {
	m := util.IsValidNumber("+8618575682804")
	if !m {
		t.Error("valid 8618575682804 err, got", m)
	}

	hk := util.IsValidNumber("+85229976161")
	if !hk {
		t.Error("valid 85229976161 err, got", hk)
	}

	m1 := util.IsValidNumber("+8522997616a")
	if m1 {
		t.Error("valid 8522997616 err, got", m1)
	}
}
