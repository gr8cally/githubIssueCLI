// Package github provides a Go API for the GitHub issue tracker.
// See https://developer.github.com/v3/search/#search-issues.
package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

const (
	BaseUrl    = "https://api.github.com"
	issuesPath = "/issues"
	reposPath  = "/repos"
)

type User struct {
	Login string
}

type Label struct {
	Name string
}

type GetIssue struct {
	Link     string `json:"html_url"`
	Title    string
	Owner    *User `json:"User"`
	Labels   []*Label
	State    string
	Assignee *User
	Body     string
}

type IssueState string

const (
	open   IssueState = "open"
	closed            = "closed"
)

type PostIssue struct {
	Title       string     `json:"title,omitempty"`
	Body        string     `json:"body,omitempty"`
	Assignees   []string   `json:"assignees,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
	State       IssueState `json:"state,omitempty"`
	issueNumber int
}

func (issue *PostIssue) SetIssueNumber(i int) {
	issue.issueNumber = i
}

func (i GetIssue) String() string {
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

// GetUserIssues queries GitHub and returns a list of issues assigned to the authenticated user
func GetUserIssues(username, password string) (*[]GetIssue, error) {
	currentUrl := buildUrl(issuesPath)
	resp := getResponse(username, password, currentUrl, "GET", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("bad stat code")
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var arr []GetIssue
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		return nil, err
	}
	return &arr, nil
}

// buildUrl takes in a path parameters as strings in the order which they should appear with or without leading `/`
// and returns a GitHub api url with those paths
func buildUrl(paths ...string) *url.URL {
	currentUrl, err := url.Parse(BaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	currentUrl.Path = path.Join(paths...)
	return currentUrl
}

// getResponse takes in username+password for authentication, and sends a {method} request to the url and returns the response
// the body parameter must be of type json
func getResponse(username string, password string, url *url.URL, method string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		log.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

// CreateIssue takes in username+password for authentication, and
// creates an issue in the specified {owner}/{repo} using the marshalled issue passed in
// it returns true if successfully created, false otherwise with the error encountered
func CreateIssue(username, password, owner, repo string, issue PostIssue) (bool, error) {
	jsonIssue, err := json.Marshal(issue)
	if err != nil {
		log.Fatal(err)
	}

	currentUrl := buildUrl(reposPath, owner, repo, issuesPath)
	resp := getResponse(username, password, currentUrl, "POST", bytes.NewBuffer(jsonIssue))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return false, errors.New("issue not created")
	}
	return true, nil
}

// UpdateIssue takes in username+password for authentication,
// and updates an issue in {owner}/{repo}, the issue number will be defined within the issue argument
// it updates the issue with the values passed in and returns true, nil if successful, false and error otherwise
func UpdateIssue(username, password, owner, repo string, issue PostIssue) (bool, error) {
	jsonIssue, err := json.Marshal(issue)
	if err != nil {
		log.Fatal(err)
	}

	currentUrl := buildUrl(reposPath, owner, repo, issuesPath, strconv.Itoa(issue.issueNumber))
	resp := getResponse(username, password, currentUrl, "PATCH", bytes.NewBuffer(jsonIssue))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("issue not updated")
	}
	return true, nil
}

// CloseIssue takes in username+password for authentication,
// it sets the issue with {issueNumber} in {owner}/{repo} to closed, if successful returns true, false otherwise
func CloseIssue(username, password, owner, repo string, issueNumber int) bool {
	issue := PostIssue{
		State:       closed,
		issueNumber: issueNumber,
	}
	updatedIssue, err := UpdateIssue(username, password, owner, repo, issue)
	if err != nil {
		return false
	}
	return updatedIssue
}
