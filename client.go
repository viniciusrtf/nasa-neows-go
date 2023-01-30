package neows

import (
	"bytes"
	"io"
	"net/http"
)

type Client struct {
	c      *http.Client
	cache  cache
	apiKey *string

	Feed *FeedService
}

const defaultAPIKey = "DEMO_KEY"

func NewClient(apiKey *string) *Client {
	c := &Client{c: &http.Client{}, apiKey: apiKey}
	c.Feed = NewFeedService(c)
	return c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Add the api_key query parameter to the request
	q := req.URL.Query()
	if q.Get("api_key") == "" {
		if c.apiKey == nil {
			q.Add("api_key", defaultAPIKey)
		} else {
			q.Add("api_key", *c.apiKey)
		}
		req.URL.RawQuery = q.Encode()
	}

	// check if response is in cache (if enabled)
	if c.cache != nil {
		if res, ok := c.cache.get(req.URL.String()); ok {
			cachedRes := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(res)),
				Header:     make(http.Header),
			}
			cachedRes.Header.Set("X-NASA-NeoWs-Go-Cache", "HIT")
			return cachedRes, nil
		}
	}

	res, err := c.c.Do(req)
	if err != nil {
		return res, err
	}

	// Write the response to the cache (if enabled)
	if c.cache != nil {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		c.cache.set(req.URL.String(), data)

		// Reinstate the response body
		res.Body = io.NopCloser(bytes.NewBuffer(data))
	}

	return res, nil
}

func (c *Client) EnableCache() {
	c.cache = newInMemoryCache()
}

func (c *Client) DisableCache() {
	c.cache = nil
}

func (c *Client) DelCachedEntry(fullUrl string) {
	c.cache.del(fullUrl)
}

func (c *Client) FlushCache() {
	c.cache.flush()
}
