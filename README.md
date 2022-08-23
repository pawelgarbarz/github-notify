# github-notify

[![Go Report Card](https://goreportcard.com/badge/github.com/pawelgarbarz/github-notify?style=flat-square)](https://goreportcard.com/report/github.com/pawelgarbarz/github-notify)
[![Test](https://github.com/pawelgarbarz/github-notify/actions/workflows/main.yaml/badge.svg)](https://github.com/pawelgarbarz/github-notify/actions/workflows/main.yml)
[![Release](https://img.shields.io/github/release/pawelgarbarz/github-notify.svg?style=flat-square)](https://github.com/pawelgarbarz/github-notify/releases/latest)

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
bin/github-notify get
```

Flags:
```shell
Flags:
  -d, --debugLevel int       debug level 0..3
  -h, --help                 help for get
  -p, --projectCode string   jira project code (default "JIRA")
  -u, --repoUrl string       repository brand/name (default "pawelgarbarz/github-notify")

```