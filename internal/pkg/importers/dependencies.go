package importers

type debug interface {
	Level() int
}

type githubClient interface {
	Get(url string) ([]byte, error)
}
