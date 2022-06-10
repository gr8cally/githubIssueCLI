package main

import (
	"disposableProject/github"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	listIssuesSubCommand  = "listIssues"
	createIssueSubCommand = "createIssue"
	updateIssueSubCommand = "updateIssue"
	closeIssueSubCommand  = "closeIssue"
	errorParsingFlags     = "There was an error parsing flags"
)

// isFlagPassed loops through the flags parsed and returns true if {name} was set, false otherwise
func isFlagPassed(name string, flagset *flag.FlagSet) bool {
	found := false
	flagset.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	var username, password string
	var title, body, assignees, labels string
	var owner, repo string
	var issueNumber int

	listCmd := flag.NewFlagSet(listIssuesSubCommand, flag.ExitOnError)
	listCmd.StringVar(&username, "user", "", "your username used to login")
	listCmd.StringVar(&password, "password", "", "your password used to login")

	createCmd := flag.NewFlagSet(createIssueSubCommand, flag.ExitOnError)
	createCmd.StringVar(&username, "user", "", "your username used to login")
	createCmd.StringVar(&password, "password", "", "your password used to login")
	createCmd.StringVar(&owner, "owner", "", "owner of the repo")
	createCmd.StringVar(&repo, "repo", "", "repo where new issue will be added")
	createCmd.StringVar(&title, "title", "", "title of new issue")
	createCmd.StringVar(&body, "body", "", "body or description of  new issue")
	createCmd.StringVar(&assignees, "assignees", "", "people the issue will be assigned to, list of comma seperated git usernames")
	createCmd.StringVar(&labels, "labels", "", "labels to associate to this issue,comma seperated string in list [bug, wontfix, enhancement, help wanted... etc]")

	updateCmd := flag.NewFlagSet(updateIssueSubCommand, flag.ExitOnError)
	updateCmd.StringVar(&username, "user", "", "your username used to login")
	updateCmd.StringVar(&password, "password", "", "your password used to login")
	updateCmd.StringVar(&owner, "owner", "", "owner of the repo")
	updateCmd.StringVar(&repo, "repo", "", "repo where new issue will be added")
	updateCmd.IntVar(&issueNumber, "issueNumber", 0, "the number that identifies the issue")
	updateCmd.StringVar(&title, "title", "", "title of new issue")
	updateCmd.StringVar(&body, "body", "", "body or description of  new issue")
	updateCmd.StringVar(&assignees, "assignees", "", "people the issue will be assigned to, list of comma seperated git usernames")
	updateCmd.StringVar(&labels, "labels", "", "labels to associate to this issue,comma seperated string in list [bug, wontfix, enhancement, help wanted... etc]")

	closeCmd := flag.NewFlagSet(closeIssueSubCommand, flag.ExitOnError)
	closeCmd.StringVar(&username, "user", "", "your username used to login")
	closeCmd.StringVar(&password, "password", "", "your password used to login")
	closeCmd.StringVar(&owner, "owner", "", "owner of the repo")
	closeCmd.StringVar(&repo, "repo", "", "repo where new issue will be added")
	closeCmd.IntVar(&issueNumber, "issueNumber", 0, "the number that identifies the issue")

	if len(os.Args) < 2 {
		fmt.Println("expected subcommand [listIssue, createIssue etc]")
		os.Exit(1)
	}

	switch os.Args[1] {

	case listIssuesSubCommand:
		err := listCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(errorParsingFlags)
			return
		}

		results, err := github.GetUserIssues(username, password)
		if err != nil {
			fmt.Println("error encountered while listing issues")
			os.Exit(1)
		}

		for _, v := range *results {
			fmt.Printf("%+v\n", v)
		}

	case createIssueSubCommand:
		err := createCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(errorParsingFlags)
			return
		}

		if !isFlagPassed("owner", createCmd) || !isFlagPassed("repo", createCmd) || !isFlagPassed("title", createCmd) {
			fmt.Println("-owner, -title or -repo not set, they are mandatory")
			os.Exit(1)
		}

		issue := github.PostIssue{
			Title: title,
		}

		if isFlagPassed("body", createCmd) {
			issue.Body = body
		}
		if isFlagPassed("assignees", createCmd) {
			issue.Assignees = strings.Split(assignees, ",")
		}
		if isFlagPassed("labels", createCmd) {
			issue.Labels = strings.Split(labels, ",")
		}

		_, err = github.CreateIssue(username, password, owner, repo, issue)
		if err != nil {
			fmt.Println("Unsuccessful, there was a problem creating issue")
			os.Exit(1)
		}
		fmt.Println("successfully created")

	case updateIssueSubCommand:
		err := updateCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(errorParsingFlags)
			return
		}

		if !isFlagPassed("owner", updateCmd) || !isFlagPassed("repo", updateCmd) || !isFlagPassed("issueNumber", updateCmd) {
			fmt.Println("-owner, -repo or -issueNumber not set, they are mandatory")
			os.Exit(1)
		}

		issue := github.PostIssue{}
		issue.SetIssueNumber(issueNumber)

		if !isFlagPassed("title", createCmd) {
			issue.Title = title
		}
		if isFlagPassed("body", createCmd) {
			issue.Body = body
		}
		if isFlagPassed("assignees", createCmd) {
			issue.Assignees = strings.Split(assignees, ",")
		}
		if isFlagPassed("labels", createCmd) {
			issue.Labels = strings.Split(labels, ",")
		}

		_, err = github.UpdateIssue(username, password, owner, repo, issue)
		if err != nil {
			fmt.Println("unsuccessful, error encountered while updating issue")
			os.Exit(1)
		}
		fmt.Println("successfully updated")

	case closeIssueSubCommand:
		err := closeCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(errorParsingFlags)
			return
		}

		if !isFlagPassed("owner", closeCmd) || !isFlagPassed("repo", closeCmd) || !isFlagPassed("issueNumber", closeCmd) {
			fmt.Println("-owner, -repo or -issueNumber not set, they are mandatoryy")
			os.Exit(1)
		}

		ok := github.CloseIssue(username, password, owner, repo, issueNumber)
		if !ok {
			fmt.Println("unsuccessful, error encountered while setting issue to closed state")
			os.Exit(1)
		}
		fmt.Println("issue successfully closed")
	}
}