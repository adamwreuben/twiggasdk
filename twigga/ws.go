package twigga

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// im
func (c *Client) ListenToDocumentChanges(db, table, id string) (*websocket.Conn, error) {
	endpoint := fmt.Sprintf("%s/document/%s/%s/%s/changes", c.wSBaseURL, db, table, id)
	return c.openWS(endpoint)
}

func (c *Client) ListenToCollectionChanges(db, table string) (*websocket.Conn, error) {
	endpoint := fmt.Sprintf("%s/document/%s/%s/changes", c.wSBaseURL, db, table)
	return c.openWS(endpoint)
}

func (c *Client) openWS(rawurl string) (*websocket.Conn, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	if c.client.Token != "" {
		header.Set("BONGO-TOKEN", c.client.Token)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	return conn, err
}
