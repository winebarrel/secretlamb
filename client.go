package secretlamb

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

type client struct {
	url *url.URL
}

func newClient(path string) (*client, error) {
	port := os.Getenv("PARAMETERS_SECRETS_EXTENSION_HTTP_PORT")

	if port == "" {
		port = "2773"
	}

	url, err := url.Parse("http://localhost:" + port + path)

	if err != nil {
		return nil, err
	}

	client := &client{
		url: url,
	}

	return client, nil
}

func (client *client) get(query *url.Values) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, client.url.String(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Aws-Parameters-Secrets-Token", os.Getenv("AWS_SESSION_TOKEN"))
	req.URL.RawQuery = query.Encode()
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		text := res.Status

		if len(body) > 0 {
			text += "\n" + string(body)
		}

		return nil, errors.New(text)
	}

	return body, nil
}
