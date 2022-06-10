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
	"time"
)

const (
	BaseUrl    = "https://api.github.com"
	issuesPath = "/issues"
	reposPath  = "/repos"
)

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

type GetIssue struct {
	Link     string `json:"html_url"`
	Title    string
	Owner    *User `json:"User"`
	Labels   []*Label
	State    string
	Assignee *User
	Body     string
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

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(BaseUrl + "?q=" + q)
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

func buildUrl(paths ...string) *url.URL {
	currentUrl, err := url.Parse(BaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	currentUrl.Path = path.Join(paths...)
	return currentUrl
}

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

func (issue PostIssue) SetIssueNumber(i int) {
	issue.issueNumber = i
}

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
