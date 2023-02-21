package neows

import (
	"bytes"
	"io"
	"net/http"
)

type Client struct {
	c      *http.Client
	cache  cache
	apiKey string

	Feed *FeedService
}
 
// defaultAPIKey is the default API key to use if none is provided. Its usage 
// results in restrictions on the number of requests that can be made per hour.
// Obtain your own API key at https://api.nasa.gov/index.html#apply-for-an-api-key
const defaultAPIKey = "DEMO_KEY"

// NewClient creates a new Client
func NewClient(apiKey string) *Client {
	c := &Client{c: &http.Client{}, apiKey: apiKey}
	c.Feed = NewFeedService(c, NewDefaultFeedOptions())
	return c
}

// Do executes the given HTTP request, apply default middlewares and (optionally) cache the response
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Add the api_key query parameter to the request
	q := req.URL.Query()
	if q.Get("api_key") == "" {
		if c.apiKey == "" {
			q.Add("api_key", defaultAPIKey)
		} else {
			q.Add("api_key", c.apiKey)
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

// EnableCache enables the in-memory cache
func (c *Client) EnableCache() {
	c.cache = newInMemoryCache()
}

// DisableCache disables the in-memory cache
func (c *Client) DisableCache() {
	c.cache = nil
}

// DelCachedEntry deletes the given entry from the cache
func (c *Client) DelCachedEntry(fullUrl string) {
	c.cache.del(fullUrl)
}

// FlushCache flushes the cache
func (c *Client) FlushCache() {
	c.cache.flush()
}
