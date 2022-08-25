package clients

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

type SlackHttpClient interface {
	Send(text string) error
}

type SlackClient struct {
	webhookUrl string
}

func NewSlackClient(webhookUrl string) *SlackClient {
	return &SlackClient{webhookUrl: webhookUrl}
}

func (c *SlackClient) Send(text string) error {
	var jsonStr = []byte(fmt.Sprintf("{\"text\": \"%s\"}", text))

	req, err := http.NewRequest("POST", c.webhookUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	return nil
}
