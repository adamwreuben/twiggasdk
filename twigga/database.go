package twigga

// This return current database
func (c *Client) CurrentDatabase() string {
	return c.client.Twigga.DefaultDatabase
}
