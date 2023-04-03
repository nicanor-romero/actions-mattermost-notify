# Mattermost Notify action

This action notifies the failures in other Github actions, via Mattermost.

## Inputs

## `access-token`
**Required** The access token for Mattermost, usually saved as a Github secret

## `commit-url`
**Required** The commit URL

## `commit-author-username`
**Required** The commit author username

## `commit-author-email`
**Required** The commit author email

## `commit-message`
**Required** The commit message

## `job-output`
**Required** The output of the job run before this Github Action

## Example usage

```yaml
uses: actions/actions-mattermost-notify@v1
with:
    access-token: ${{ secrets.MATTERMOST_PERSONAL_ACCESS_TOKEN }}
    commit-url: ${{ github.event.head_commit.url }}
    commit-author-username: ${{ github.event.head_commit.author.username }}
    commit-author-email: ${{ github.event.head_commit.author.email }}
    commit-message: ${{ github.event.head_commit.message }}
    job-output: ${{ toJson(needs.test.outputs.steps_output) }}
```