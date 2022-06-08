package github

import "testing"

const username = "gr8cally"
const pass, word = "ghp_1vXdXfbT9j2DXbLM2er", "WXzmyWavrdH2XRBMy"

func TestGetUserIssues(t *testing.T) {
	t.Run("Good Auth details", func(t *testing.T) {
		issues, err := GetUserIssues(username, pass+word)
		if err != nil {
			t.Fatal(err)
		}
		if issues == nil {
			t.Fatal("returned nil issues")
		}
	})

	t.Run("Bad Auth details", func(t *testing.T) {
		issues, err := GetUserIssues(username, pass)
		if err == nil {
			t.Fatal("Got nil error was expecting forbidden error")
		}
		if issues != nil {
			t.Fatal("expected nil issues but got non nil one")
		}
	})
}

func TestCreateIssue(t *testing.T) {
	issue := NewIssue{
		Title:     "this is from a test, please delete",
		Body:      "this is from CLI app, no body",
		Assignees: []string{"gr8cally"},
		Labels:    []string{"bug"},
	}
	t.Run("create new issue with valid body", func(t *testing.T) {
		ok, err := CreateIssue(username, pass+word, issue)
		if !ok {
			t.Fatal("want 201 status didnt get that")
		}
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("create new issue with body missing a required field", func(t *testing.T) {
		ok, err := CreateIssue(username, pass+word, NewIssue{})
		if ok {
			t.Fatal("got 201 which wasn't expected")
		}
		if err == nil {
			t.Fatal("Want a an error 422 here")
		}
	})
}
