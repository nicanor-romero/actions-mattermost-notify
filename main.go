package main

import (
	"fmt"
	"os"

	"github.com/mattermost/mattermost-server/v5/model"
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
	accessToken := os.Getenv("MATTERMOST_ACCESS_TOKEN")

	client := model.NewAPIv4Client("https://mattermost.masstack.com")
	client.SetToken(accessToken)

	commit := Commit{
		sha: os.Getenv("COMMIT_SHA"),
		url: os.Getenv("REPO_URL") + "/commit/" + os.Getenv("COMMIT_SHA"),

		// TODO: Get the rest via Github API
		authorUsername: "missing_username",
		authorEmail:    "missing_email",
		commitMessage:  "missing_message",
	}

	checkRun := CheckRun{
		Name:          os.Getenv("CHECK_RUN_NAME"),
		Conclusion:    os.Getenv("CHECK_RUN_CONCLUSION"),
		Url:           os.Getenv("CHECK_RUN_URL"),
		OutputTitle:   os.Getenv("CHECK_RUN_OUTPUT_TITLE"),
		OutputText:    os.Getenv("CHECK_RUN_OUTPUT_TEXT"),
		OutputSummary: os.Getenv("CHECK_RUN_OUTPUT_SUMMARY"),
	}

	message := fmt.Sprintf(":warning: The commit [%s](%s) by @nicanor.romero (%s - %s) has failed the pipeline step %s",
		commit.commitMessage,
		commit.url,
		commit.authorUsername,
		commit.authorEmail,
		checkRun.Name,
	)

	if checkRun.Failed() {
		message += "\n  :red_circle: "
	} else {
		// TODO: Delete this, only debugging
		message += "\n  :large_green_circle: "
	}
	message += fmt.Sprintf("%s\n    _%s_", checkRun.OutputTitle, checkRun.OutputText)

	testChannelId := "audgc68w4pri7eybkt4byg9pze"
	post := &model.Post{
		ChannelId: testChannelId,
		Message:   message,
	}

	_, _ = client.CreatePost(post)
}
