# github-notify

[![Go Report Card](https://goreportcard.com/badge/github.com/pawelgarbarz/github-notify?style=flat-square)](https://goreportcard.com/report/github.com/pawelgarbarz/github-notify)
[![CI](https://github.com/pawelgarbarz/github-notify/actions/workflows/main.yaml/badge.svg)](https://github.com/pawelgarbarz/github-notify/actions/workflows/main.yml)
[![CodeQL](https://github.com/pawelgarbarz/github-notify/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/pawelgarbarz/github-notify/actions/workflows/codeql-analysis.yml)
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
jira-project: JIRA
github-repo: pawelgarbarz/github-notify
```

parameters:
- `token` - GitHub token with access to repository pull requests
- `reviewers` - list of team members who want to be informed about pull requests
  - `github` - user GitHub login
  - `slack` - user Slack login
- `webhook-url` - Slack channel webhook url
- `jira-project` - Jira project code, can be overwritten by parameter  `-p`, `--jira-project`
- `github-repo` - GitHub repository in format brand/name, can be overwritten by parameter  `-u`, `--github-repo`

## Run command
```shell
Usage:
  github-notify pull-request [flags]

Flags:
  -d, --debugLevel int   debug level 0..3
  -h, --help             help for pull-request

Global Flags:
      --config string         config file (default is $HOME/.github-notify.yaml)
  -u, --github-repo string    github repository brand/name
  -p, --jira-project string   jira project code
```