# github-notify
![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/pawelgarbarz/github-notify?style=flat-square)](https://goreportcard.com/report/github.com/pawelgarbarz/github-notify)
[![CI](https://github.com/pawelgarbarz/github-notify/actions/workflows/main.yaml/badge.svg)](https://github.com/pawelgarbarz/github-notify/actions/workflows/main.yaml)
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
cache-enabled: true
cache-ttl: 30
```

parameters:
- `token` - GitHub token with access to repository pull requests
- `reviewers` - list of team members who want to be informed about pull requests
  - `github` - user GitHub login
  - `slack` - user Slack login
- `webhook-url` - Slack channel webhook url
- `jira-project` - Jira project code, can be overwritten by parameter  `-p`, `--jira-project`
- `github-repo` - GitHub repository in format brand/name, can be overwritten by parameter  `-u`, `--github-repo`
- `cache-enabled` - Enable/Disable local persistent cache file to store lat send timestamp. Default `true`
- `cache-ttl` - Time to leave data in cache, set by int number of seconds. Default `14400` which is `4h`

## Run command
```shell
  github-notify [command]

Available Commands:
  cache-clear-all      clear ALL local cache
  cache-clear-outdated clear outdated local cache
  completion           Generate the autocompletion script for the specified shell
  help                 Help about any command
  pull-request         get pull requests from github

Flags:
      --config string         config file (default is $HOME/.github-notify.yaml)
  -d, --debugLevel int        debug level 0..3
  -u, --github-repo string    github repository brand/name
  -h, --help                  help for github-notify
  -p, --jira-project string   jira project code
  -t, --toggle                Help message for toggle
```

## Troubleshooting

### Got Slack message only at first program run
Most probably this is caused because anti-spam cache feature. By default, next message will be sent after `4h` You can change `cache-ttl` value to lewer that time. There is also possibility to clear whole cache by calling `github-notify cache clear all`

### local cache file is very big
Because local cache is build on top of `SQLite` and this database do not support TTL cache mechanism, I recommend to run `github-notify cache clear outdated` command from time to time.
If program execution is automated for example by `crone`, I recommend to add e.g. daily clean of outdated cache.

## Example of `crone` configuration

Run: At every 5th minute past every hour from 7 through 17 on every day-of-week from Monday through Friday.
```shell
*/5 7-17 * * 1-5 /home/pawel.garbarz/github-notify pull-request
0 6,18 * * 1-5 /home/pawel.garbarz/github-notify cache clear outdated
```
