package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

type Commit struct {
	url            string
	authorUsername string
	authorEmail    string
	commitMessage  string
}

type JobOutput struct {
	Outputs    map[string]string `json:"outputs"`
	Outcome    string            `json:"outcome"`
	Conclusion string            `json:"conclusion"`
}

func (o JobOutput) Failed() bool {
	return o.Outcome == "failure"
}

func main() {
	accessToken := os.Getenv("MATTERMOST_PERSONAL_ACCESS_TOKEN")

	client := model.NewAPIv4Client("https://mattermost.masstack.com")
	client.SetToken(accessToken)

	commit := Commit{
		url:            os.Getenv("COMMIT_URL"),
		authorUsername: os.Getenv("COMMIT_AUTHOR_USERNAME"),
		authorEmail:    os.Getenv("COMMIT_AUTHOR_EMAIL"),
		commitMessage:  os.Getenv("COMMIT_MESSAGE"),
	}

	fmt.Println("commit.url:", commit.url)
	fmt.Println("commit.authorUsername:", commit.authorUsername)
	fmt.Println("commit.authorEmail:", commit.authorEmail)
	fmt.Println("commit.commitMessage:", commit.commitMessage)

	testJobOutputStr := os.Getenv("TEST_JOB_OUTPUT")
	fmt.Println("TEST_JOB_OUTPUT:", testJobOutputStr)
	testJobOutputStr = testJobOutputStr[1 : len(testJobOutputStr)-1]
	fmt.Println("TEST_JOB_OUTPUT:", testJobOutputStr)
	testJobOutputStr = strings.ReplaceAll(testJobOutputStr, "\\n", "")
	fmt.Println("TEST_JOB_OUTPUT:", testJobOutputStr)
	testJobOutputStr = strings.ReplaceAll(testJobOutputStr, "\\", "\"")
	fmt.Println("TEST_JOB_OUTPUT:", testJobOutputStr)

	testJobOutput := make(map[string]JobOutput, 0)
	err := json.Unmarshal([]byte(testJobOutputStr), &testJobOutput)
	if err != nil {
		panic(err)
	}

	// {
	//  "lint": {
	//    "outputs": {},
	//    "outcome": "failure",
	//    "conclusion": "failure"
	//  },
	//  "pytest": {
	//    "outputs": {},
	//    "outcome": "skipped",
	//    "conclusion": "skipped"
	//  }
	//}

	message := fmt.Sprintf(":warning: The commit [%s](%s) by @nicanor.romero (%s - %s) has failed one or more pipeline steps:",
		commit.commitMessage,
		commit.url,
		commit.authorUsername,
		commit.authorEmail,
	)
	for jobKey, jobOutput := range testJobOutput {
		if jobOutput.Failed() {
			message += fmt.Sprintf("\n  * :red_circle: %s", jobKey)
		}
	}

	testChannelId := "audgc68w4pri7eybkt4byg9pze"
	post := &model.Post{
		ChannelId: testChannelId,
		Message:   message,
	}

	_, _ = client.CreatePost(post)
}
