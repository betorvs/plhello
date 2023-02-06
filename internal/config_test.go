package internal

import "testing"

func TestGetEnv(t *testing.T) {
	t.Log("Given a empty variable")
	{
		empty := GetEnv("EMPTYENVVAR", "defaultvalue")
		if empty != "defaultvalue" {
			t.Fatalf("\tGetEnv test should returned defaultvalue")
		}
	}
	t.Log("Given a environment variable")
	{
		t.Setenv("ENVVAR", "notdefaultvalue")
		notempty := GetEnv("ENVVAR", "defaultvalue")
		if notempty != "notdefaultvalue" {
			t.Fatalf("\tGetEnv test should returned notdefaultvalue")
		}
	}

}
