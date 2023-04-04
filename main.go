package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	MattermostUrl       = "https://mattermost.masstack.com"
	MattermostChannelId = "audgc68w4pri7eybkt4byg9pze"
)

type Commit struct {
	url            string
	authorUsername string
	authorEmail    string
	commitMessage  string
}

func (c Commit) getCommitMessageTitle() string {
	return strings.Split(c.commitMessage, "\n")[0]
}

type CommitStatus struct {
	Name        string
	Description string
	Conclusion  string
	Url         string
}

func (o CommitStatus) Succeeded() bool {
	return o.Conclusion == "success"
}

func (o CommitStatus) Failed() bool {
	return o.Conclusion == "failure"
}

func main() {
	fmt.Println("Running actions-mattermost-notify")

	mattermostClient := getMattermostClient()
	commit := buildCommit()
	commitStatus := buildCommitStatus()
	message := buildMessage(mattermostClient, commit, commitStatus)

	sendMessage(mattermostClient, message)
}

func getMattermostClient() (client *model.Client4) {
	accessToken := os.Getenv("MATTERMOST_ACCESS_TOKEN")
	client = model.NewAPIv4Client(MattermostUrl)
	client.SetToken(accessToken)

	return client
}

func sendMessage(client *model.Client4, message string) {
	post := &model.Post{
		ChannelId: MattermostChannelId,
		Message:   message,
	}

	_, response := client.CreatePost(post)
	fmt.Println("response:", response)
}

func buildMessage(client *model.Client4, commit Commit, commitStatus CommitStatus) (message string) {
	mattermostUser, resp := client.GetUserByEmail(commit.authorEmail, "")
	if resp.StatusCode != 200 {
		mattermostUser = &model.User{Username: "UNKNOWN"}
	}

	message = fmt.Sprintf(":warning: The commit [_\"%s\"_](%s) by @%s (%s - %s) has failed the pipeline step `%s`:",
		commit.getCommitMessageTitle(),
		commit.url,
		mattermostUser.Username,
		commit.authorUsername,
		commit.authorEmail,
		commitStatus.Name,
	)
	message += fmt.Sprintf("\n*  [%s](%s): _%s_", commitStatus.Name, commitStatus.Url, commitStatus.Description)
	return
}

func buildCommitStatus() (commitStatus CommitStatus) {
	commitStatus = CommitStatus{
		Name:        os.Getenv("STATUS_NAME"),
		Description: os.Getenv("STATUS_DESCRIPTION"),
		Conclusion:  os.Getenv("STATUS_CONCLUSION"),
		Url:         os.Getenv("STATUS_URL"),
	}
	return
}

func buildCommit() (commit Commit) {
	commit = Commit{
		url:            os.Getenv("COMMIT_URL"),
		authorUsername: os.Getenv("COMMIT_AUTHOR_USERNAME"),
		authorEmail:    os.Getenv("COMMIT_AUTHOR_EMAIL"),
		commitMessage:  os.Getenv("COMMIT_MESSAGE"),
	}
	return
}
