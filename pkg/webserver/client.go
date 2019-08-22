package web

// NewClient returns rest client to interact with embedded REST service
func NewClient() Client {
	return Client{}
}

// Upsert creates new link entity
func (c *Client) Upsert(id string, link LinkItem) error {
	return nil
}

// Delete removes specified link entity
func (c *Client) Delete(id string) error {
	return nil
}
