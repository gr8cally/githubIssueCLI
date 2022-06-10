# gitIssueCLIApp
This project is from excercise 4.11 in [The Go Programming Language book](https://learning.oreilly.com/library/view/the-go-programming/9780134190570).

A tool that lets users
1. create, 
2. read (list all issues assigned to that user), 
3. update and 
4. close Github Issues from the command line.

This project was done 100% in GO with only standard packages used.

Usage
```go run main.go {subcommand} {-flag=value}```
subcommands and their respective flags

| Flags        | Type           | description  |
| ------------- |:-------------:| -----|
| -user     | string | GitHub username |
| -password | string |   [your GitHub personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) |
| -title    | string | Title of the Issue   |
| -body  |  string     | Content of the Issue   |
| -labels  | string      | comma seperated string in list [bug, wontfix, enhancement, help wanted... etc]   |
| -assignees  |  string     |logins for users to assign to the issue passed in as comma seperated strings    |
| -owner  | string      | the account owner of the repository   |
| -repo  |  string     |  the name of the repository  |
| -issueNumber| int | the number that defines the issue|

| subCommands   | flags           | description  | Mandatory  |
| ------------- |:-------------:| -----| -----|
| listIssues     | -user -password | list all issues assigned to that user | -user -password|
| createIssue     | -user -password -owner -repo -title -body -assignees -labels | create a new issue | -user -password -title -owner -repo|
| updateIssue     | -user -password -owner -repo -title -body -assignees -labels -issueNumber | Modify an already existing issue| -user -password -issueNumber -owner -repo |
| closeIssue     |  -user -password -owner -repo -issueNumber | set the state of a pre existing issue to `closed` | -user -password -title -owner -repo -issueNumber |

assuming my GitHub username is `gostar` and i am the owner of `funkyProject` repo.

example: 
```go run main.go createIssue -user=gostar -password={your GitHub personal access token} -title title_from_my_cli_app -owner gostar -repo funkyProject -labels bug,wontfix```

would create a new issue with the values from the flags

