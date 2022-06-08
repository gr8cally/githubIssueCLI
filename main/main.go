package main

import (
	"disposableProject/github"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	//auth()
	//result, err := github.SearchIssues(os.Args[1:])
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Printf("%d issues:\n", result.TotalCount)
	//
	//var monthOld, yearOld, overYear []*github.Issue
	//for _, item := range result.Items {
	//	now := time.Now().UTC()
	//	if DifferenceInDays(now, item) <= 30 {
	//		monthOld = append(monthOld, item)
	//	} else if DifferenceInDays(now, item) <= 365 {
	//		yearOld = append(yearOld, item)
	//	} else {
	//		overYear = append(overYear, item)
	//	}
	//}
	//printGroup(monthOld, "less than a month")
	//printGroup(yearOld, "less than a Year")
	//printGroup(overYear, "Over a Year")
	var username = "gr8cally"
	var pass, word = "ghp_1vXdXfbT9j2DXbLM2er", "WXzmyWavrdH2XRBMy"
	//results, err := github.GetUserIssues(username, pass+word)
	//if err != nil {
	//	fmt.Println("Er ret")
	//	os.Exit(1)
	//}
	//for _, v := range *results {
	//	fmt.Printf("%+v\n", v)
	//}

	issue := github.NewIssue{
		Title:     "issue 1 lege",
		Body:      "this is from CLI app, no body",
		Assignees: []string{"gr8cally"},
		Labels:    []string{"bug"},
	}

	_, err := github.CreateIssue(username, pass+word, issue)
	if err != nil {
		fmt.Println("Er ret")
		os.Exit(1)
	}
	fmt.Println("sucessfully created")

}

func printGroup(arr []*github.Issue, str string) {
	fmt.Println("\n", str)
	for _, item := range arr {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
		fmt.Println(item.CreatedAt)
	}
}

func DifferenceInDays(now time.Time, item *github.Issue) float64 {
	return now.Sub(item.CreatedAt).Hours() / 24
}

func auth() {
	var username = "gr8cally"
	var password = "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"
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode, " c")
	resp.Body.Close()
	v := resp.Header.Get("x-ratelimit-limit")
	fmt.Println("L: ", v)
	v = resp.Header.Get("server")
	fmt.Println("S: ", v)
	v = resp.Header.Get("x-ratelimit-remaining")
	fmt.Println("R: ", v)
}
