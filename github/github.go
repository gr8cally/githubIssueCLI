// Package github provides a Go API for the GitHub issue tracker.
// See https://developer.github.com/v3/search/#search-issues.
package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login string
	//HTMLURL string `json:"html_url"`
}

type Head struct {
	Limit  int `json:"x-ratelimit-limit"`
	Server string
}

type Label struct {
	Name string
}

type Issu struct {
	Link     string `json:"html_url"`
	Title    string
	Owner    *User `json:"User"`
	Labels   []*Label
	State    string
	Assignee *User
	Body     string
}

func (i Issu) String() string {
	output := "{"
	output += fmt.Sprintf("Link: %v\n", i.Link)
	output += fmt.Sprintf("Title: %v\n", i.Title)
	output += fmt.Sprintf("Owner: %v\n", i.Owner.Login)
	output += "Labels: ["
	for _, v := range i.Labels {
		output += fmt.Sprintf("%v, ", v.Name)
	}
	output = strings.TrimRight(output, ", ")
	output += "]\n"
	output += fmt.Sprintf("State: %v\n", i.State)
	output += fmt.Sprintf("Assignee: %v\n", i.Assignee.Login)
	output += fmt.Sprintf("Body: %v}\n", i.Body)
	return output
}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	//var head Head
	v := resp.Header.Get("x-ratelimit-limit")
	fmt.Println("L: ", v)
	v = resp.Header.Get("server")
	fmt.Println("S: ", v)
	v = resp.Header.Get("x-ratelimit-remaining")
	fmt.Println("R: ", v)

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

func GetUserIssues() (*[]Issu, error) {
	var username = "gr8cally"
	var password = "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"

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.github.com/issues", nil)
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("bad stat code")
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var arr []Issu
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		return nil, err
	}
	return &arr, nil
}
