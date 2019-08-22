package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewClient returns rest client to interact with embedded REST service
func NewClient() Client {
	c := Client{}
	c.baseURL = fmt.Sprintf("http://localhost:%s", port)
	return c
}

func (c *Client) Do(method string, item LinkItem) error {
	var err error
	body, err := json.Marshal(item)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, c.baseURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// NOTE: we return 200 on any internal response for simplicity
	if 200 != resp.StatusCode {
		return fmt.Errorf("%s", respBody)
	}
	return err
}
