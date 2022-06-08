package github

import "testing"

func TestGetUserIssues(t *testing.T) {
	t.Run("Good Auth details", func(t *testing.T) {
		var username = "gr8cally"
		var pass, word = "ghp_1vXdXfbT9j2DXbLM2er", "WXzmyWavrdH2XRBMy"
		issues, err := GetUserIssues(username, pass+word)
		if err != nil {
			t.Fatal(err)
		}
		if issues == nil {
			t.Fatal("returned nil issues")
		}
	})

	t.Run("Bad Auth details", func(t *testing.T) {
		var username = "gr8cally"
		var pass = "ghp_1vXdXfbT9j2DXbLM2er"
		issues, err := GetUserIssues(username, pass)
		if err == nil {
			t.Fatal("Got nil error was expecting forbidden error")
		}
		if issues != nil {
			t.Fatal("expected nil issues but got non nil one")
		}
	})
}
