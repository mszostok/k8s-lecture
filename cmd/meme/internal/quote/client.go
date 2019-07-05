package quote

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

// Doer represents entity making HTTP requests via Do method.
type Doer interface {
	// Do is making call using request and returning response.
	// If err is nil then response body has to be closed or resources will be leaked.
	// To reuse connection make sure that body is drained. If not then underlying connection
	// will be closed upon closing the body.
	// If err is not nil than response is nil and can be ignored.
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	doer    Doer
	baseURL string
}

func NewClient(doer Doer, baseURL string) *Client {
	return &Client{
		doer:    doer,
		baseURL: baseURL,
	}
}

func (c *Client) Get() (string, error) {
	url := fmt.Sprintf("%s/quote", c.baseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "while creating HTTP request")
	}
	resp, err := c.doer.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "while making HTTP call")
	}
	defer func() {
		_ = DrainReader(resp.Body)
		resp.Body.Close()
	}()

	switch {
	case resp.StatusCode == http.StatusOK:
	default:
		return "", fmt.Errorf("Got wrong status code. Expected: [%d], got: [%d] ", http.StatusOK, resp.StatusCode)
	}
	bodyRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "while reading HTTP response body")
	}

	var dto DTO
	if err = json.Unmarshal(bodyRaw, &dto); err != nil {
		return "", errors.Wrap(err, "while decoding HTTP response")
	}

	return dto.Quote, nil
}

type DTO struct {
	Quote string `json:"quote"`
}

// DrainReader reads and discards the remaining part in reader (for example response body data)
// In case of HTTP this ensured that the http connection could be reused for another request if the keepalive http connection behavior is enabled.
func DrainReader(reader io.Reader) error {
	if reader == nil {
		return nil
	}
	_, drainError := io.Copy(ioutil.Discard, io.LimitReader(reader, 4096))
	return drainError
}
