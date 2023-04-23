package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)

const (
	GitHubOrganization = "masmovil"
	MattermostUrl      = "https://mattermost.masstack.com"
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

type WebhookMessage struct {
	Text    string `json:"text"`
	Channel string `json:"channel,omitempty"`
}

func (o CommitStatus) Succeeded() bool {
	return o.Conclusion == "success"
}

func (o CommitStatus) Failed() bool {
	return o.Conclusion == "failure"
}

// GithubUserSSO is used to unmarshall GitHub API response
type GithubUserSSO struct {
	Data struct {
		Organization struct {
			SAMLIdentityProvider struct {
				ExternalIdentities struct {
					Edges []struct {
						Node struct {
							SamlIdentity struct {
								NameId string `json:"nameId"`
							} `json:"samlIdentity"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"externalIdentities"`
			} `json:"samlIdentityProvider"`
		} `json:"organization"`
	} `json:"data"`
}

func main() {
	fmt.Println("Running actions-mattermost-notify")

	mattermostClient := getMattermostClient()
	commit := buildCommit()
	commitStatus := buildCommitStatus()
	message := buildMessage(mattermostClient, commit, commitStatus)

	sendMessage(message)
	return
}

func getMattermostClient() (client *model.Client4) {
	accessToken := os.Getenv("MATTERMOST_ACCESS_TOKEN")
	client = model.NewAPIv4Client(MattermostUrl)
	client.SetToken(accessToken)

	return client
}

func sendMessage(message string) {
	webHookUrl := os.Getenv("MATTERMOST_INCOMING_WEBHOOK_URL")
	webHookMesasge := WebhookMessage{
		Text:    message,
		Channel: os.Getenv("MATTERMOST_CHANNEL_NAME"),
	}
	body, err := json.Marshal(webHookMesasge)
	if err != nil {
		fmt.Println("got error marshaling the request body:", err)
		return
	}

	response, err := http.Post(webHookUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Println("got error posting to webhook:", err)
		return
	}
	fmt.Println(response)
	return
}

func buildMessage(client *model.Client4, commit Commit, commitStatus CommitStatus) (message string) {
	mattermostUser, _, _ := client.GetUserByEmail(commit.authorEmail, "")

	userMention := buildUserMention(mattermostUser, commit.authorUsername)

	message = fmt.Sprintf(":warning: The commit [_\"%s\"_](%s) by %s has failed the pipeline step `%s`:",
		commit.getCommitMessageTitle(),
		commit.url,
		userMention,
		commitStatus.Name,
	)
	message += fmt.Sprintf("\n*  [%s](%s): _%s_", commitStatus.Name, commitStatus.Url, commitStatus.Description)
	return
}

func buildUserMention(mattermostUser *model.User, githubAuthorUsername string) (mention string) {
	githubAuthorUrl := "https://github.com/" + githubAuthorUsername
	if mattermostUser != nil {
		mention += fmt.Sprintf("@%s ([%s](%s))", mattermostUser.Username, githubAuthorUsername, githubAuthorUrl)
	} else {
		mention += fmt.Sprintf("[%s](%s)", githubAuthorUsername, githubAuthorUrl)
	}
	return mention
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

	authorEmail, err := getAuthorEmailFromGithubSSO(commit.authorUsername)
	if err != nil {
		// If we are unable to get email from GitHub SSO, we will use the one specified in the commit metadata
		fmt.Println("got error getting email from github SSO:", err)
		return
	}
	// Replace the email from the commit with the one from GitHub SSO
	commit.authorEmail = authorEmail

	return
}

func getAuthorEmailFromGithubSSO(authorUsername string) (authorEmail string, err error) {
	// Get email from organization SSO, using GitHub username as key
	queryBody := fmt.Sprintf("{\"query\": \"query {organization(login: \\\"%s\\\"){samlIdentityProvider{externalIdentities(first: 1, login: \\\"%s\\\") {edges {node {samlIdentity {nameId}}}}}}}\"}", GitHubOrganization, authorUsername)
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer([]byte(queryBody)))
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("got error while doing request to github API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("got error reading github API response body:", err)
		return
	}

	var githubAuthorSSO GithubUserSSO
	err = json.Unmarshal(body, &githubAuthorSSO)
	if err != nil {
		fmt.Println("got error unmarshalling github API response body:", err)
		return
	}

	if len(githubAuthorSSO.Data.Organization.SAMLIdentityProvider.ExternalIdentities.Edges) == 0 {
		err = errors.New("no external identity edges")
		fmt.Println("got zero external identity edges from github api response:", err)
		return
	}

	authorEmail = githubAuthorSSO.Data.Organization.SAMLIdentityProvider.ExternalIdentities.Edges[0].Node.SamlIdentity.NameId
	return
}
