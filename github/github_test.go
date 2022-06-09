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
	issue := PostIssue{
		Title:     "this is from a test, please delete",
		Body:      "this is from CLI app, no body",
		Assignees: []string{"gr8cally"},
		Labels:    []string{"bug"},
	}
	t.Run("create new issue with valid body", func(t *testing.T) {
		ok, err := CreateIssue(username, pass+word, "gr8cally", "TAir_Dry", issue)
		if !ok {
			t.Fatal("want 201 status didnt get that")
		}
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("create new issue with body missing a required field", func(t *testing.T) {
		ok, err := CreateIssue(username, pass+word, "gr8cally", "TAir_Dry", PostIssue{})
		if ok {
			t.Fatal("got 201 which wasn't expected")
		}
		if err == nil {
			t.Fatal("Want a an error 422 here")
		}
	})
}

func TestUpdateIssue(t *testing.T) {
	t.Run("update an issue sucessfully", func(t *testing.T) {
		updateIssue := PostIssue{
			Title:       "was new issh, Renamed for test purposes",
			Body:        "body updated from test",
			issueNumber: 2,
		}

		ok, err := UpdateIssue(username, pass+word, "gr8cally", "TAir_Dry", updateIssue)
		if !ok {
			t.Fatal("want 200 status didnt get that")
		}
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("update an issue that does not exist", func(t *testing.T) {
		updateIssue := PostIssue{
			Title:       "was new issh, Renamed for test purposes",
			Body:        "body updated from test",
			issueNumber: 34346790,
		}

		ok, err := UpdateIssue(username, pass+word, "gr8cally", "TAir_Dry", updateIssue)
		if ok {
			t.Fatal("got 200 but wanted a failure")
		}
		if err == nil {
			t.Fatal("want an error not nil error")
		}
	})
}

func TestCloseIssue(t *testing.T) {
	t.Run("close an open issue successfully", func(t *testing.T) {
		ok := CloseIssue(username, pass+word, "gr8cally", "TAir_Dry", 3)
		if !ok {
			t.Fatal("unsuccessful closing of issue")
		}
	})
}
