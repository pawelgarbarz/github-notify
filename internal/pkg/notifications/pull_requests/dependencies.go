package pull_requests

type debug interface {
	Level() int
}

type config interface {
	SlackLoginByGithub(ghUsername string) (string, error)
}

type notificationClient interface {
	Send(text string) error
}
