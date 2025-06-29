package twigga

func (c *Client) CurrentDatabase() string {
	return c.client.Twigga.DefaultDatabase
}
