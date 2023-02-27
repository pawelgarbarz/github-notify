package clients

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type GithubHttpClient interface {
	Get(url string) ([]byte, error)
}

type GithubClient struct {
	token string
}

func NewGithubClient(token string) *GithubClient {
	return &GithubClient{token: token}
}

func (c GithubClient) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
