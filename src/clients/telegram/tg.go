package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	data, err := c.doRequest("getUpdates", q)
	if err != nil {
		return nil, err
	}
	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("TG request error", err)
	}

	req.URL.RawQuery = query.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TG request error", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("TG request error", err)
	}
	return body, nil
}

func (c *Client) SendMessage(chat_id int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chat_id))
	q.Add("text", text)

	_, err := c.doRequest("sendMessage", q)
	if err != nil {
		return fmt.Errorf("TG request error", err)
	}
	return err

}
