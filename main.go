package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/mattermost/mattermost-server/v5/model"
	"golang.org/x/oauth2"
)

type Commit struct {
	sha            string
	url            string
	authorUsername string
	authorEmail    string
	commitMessage  string
}

type CheckRun struct {
	Name          string
	Conclusion    string
	Url           string
	OutputTitle   string
	OutputText    string
	OutputSummary string
}

func (o CheckRun) Succeeded() bool {
	return o.Conclusion == "success"
}

func (o CheckRun) Failed() bool {
	return o.Conclusion == "failure"
}

func main() {
	fmt.Println("Running actions-mattermost-notify")

	commit := buildCommit()
	checkRun := buildCheckRun()
	message := buildMessage(commit, checkRun)

	sendMessage(message)
}

func sendMessage(message string) {
	testChannelId := "audgc68w4pri7eybkt4byg9pze"
	post := &model.Post{
		ChannelId: testChannelId,
		Message:   message,
	}

	accessToken := os.Getenv("MATTERMOST_ACCESS_TOKEN")
	client := model.NewAPIv4Client("https://mattermost.masstack.com")
	client.SetToken(accessToken)

	_, response := client.CreatePost(post)
	fmt.Println("response:", response)
}

func buildMessage(commit Commit, checkRun CheckRun) (message string) {
	message = fmt.Sprintf(":warning: The commit [%s](%s) by @nicanor.romero (%s - %s) has failed the pipeline step `%s`:",
		commit.commitMessage,
		commit.url,
		commit.authorUsername,
		commit.authorEmail,
		checkRun.Name,
	)
	message += fmt.Sprintf("\n*  [%s](%s): _%s_", checkRun.OutputTitle, checkRun.Url, checkRun.OutputText)
	return
}

func buildCheckRun() (checkRun CheckRun) {
	checkRun = CheckRun{
		Name:          os.Getenv("CHECK_RUN_NAME"),
		Conclusion:    os.Getenv("CHECK_RUN_CONCLUSION"),
		Url:           os.Getenv("CHECK_RUN_URL"),
		OutputTitle:   os.Getenv("CHECK_RUN_OUTPUT_TITLE"),
		OutputText:    os.Getenv("CHECK_RUN_OUTPUT_TEXT"),
		OutputSummary: os.Getenv("CHECK_RUN_OUTPUT_SUMMARY"),
	}
	return
}

func buildCommit() (commit Commit) {

	ctx := context.Background()
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// TODO: Use a regex to get owner and repo name from "https://github.com/masmovil/mm-monorepo"
	repositoryUrl := os.Getenv("REPO_URL")
	repositoryData := strings.Split(strings.Replace(repositoryUrl, "https://github.com/", "", 1), "/")
	repositoryOwner := repositoryData[0]
	repositoryName := repositoryData[1]

	// list all repositories for the authenticated user
	githubCommitData, _, err := client.Repositories.GetCommit(ctx, repositoryOwner, repositoryName, os.Getenv("COMMIT_SHA"), nil)
	if err != nil {
		panic(err)
	}

	commit = Commit{
		sha:            githubCommitData.GetSHA(),
		url:            githubCommitData.GetHTMLURL(),
		authorUsername: githubCommitData.Author.GetLogin(),
		authorEmail:    githubCommitData.Author.GetEmail(),
		commitMessage:  githubCommitData.Commit.GetMessage(),
	}
	return
}
