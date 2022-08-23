# github-notify

This application parse GitHub repository pull requests by title and authors and send notification to Slack channel about pending reviews.

## Configure
Create in Your home directory tile named `.github-notify.yaml` with content
```yaml
token: ghp_github_token
reviewers:
  - github: pawelgarbarz
    slack: garbarz.pawel
webhook-url: https://hooks.slack.com/some-random-code
```

parameters:
- `token` - GitHub token with access to repository pull requests
- `reviewers` - list of team members who want to be informed about pull requests
  - `github` - user GitHub login
  - `slack` - user Slack login
- `webhook-url` - slack channel webhook url

## Usage
run:
```shell
bin/notify get
```

Flags:
```shell
Flags:
  -d, --debugLevel int       debug level 0..3
  -h, --help                 help for get
  -p, --projectCode string   jira project code (default "JIRA")
  -u, --repoUrl string       repository brand/name (default "pawelgarbarz/github-notify")

```