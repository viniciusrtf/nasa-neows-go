package neows

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	c           *http.Client
	middlewares []func(http.RoundTripper) http.RoundTripper
	cache       cache
}

type transport struct {
	http.RoundTripper
	query url.Values
}

const defaultAPIKey = "DEMO_KEY"

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.RawQuery = t.query.Encode() + "&" + req.URL.RawQuery
	return t.RoundTripper.RoundTrip(req)
}

func NewClient(apiKey *string) *Client {
	query := url.Values{}
	if apiKey == nil {
		query.Add("api_key", defaultAPIKey)
	} else {
		query.Add("api_key", *apiKey)
	}

	tr := &transport{http.DefaultTransport, query}
	c := &http.Client{Transport: tr}

	return &Client{c: c}
}

func (c *Client) Use(m func(http.RoundTripper) http.RoundTripper) {
	c.middlewares = append(c.middlewares, m)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	rt := http.DefaultTransport
	for _, m := range c.middlewares {
		rt = m(rt)
	}
	c.c.Transport = rt

	fmt.Println("DOING REQUEST")

	// check if response is in cache (if enabled)
	if c.cache != nil {
		fmt.Println("cache is enabled")
		fmt.Println("looking for: ", req.URL.String())
		cached, found := c.cache.get(req.URL.String())
		fmt.Println("found: ", found)
		fmt.Printf("cached: %s\n", cached)
		if res, ok := c.cache.get(req.URL.String()); ok {
			cachedRes := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(res)),
				Header:     make(http.Header),
			}
			fmt.Println("HIT")
			cachedRes.Header.Set("X-NASA-NeoWs-Go-Cache", "HIT")
			return cachedRes, nil
		}
		fmt.Println("MISS")
	}

	fmt.Println("DOING REQUEST")

	res, err := c.c.Do(req)
	if err != nil {
		return res, err
	}

	// Write the response to the cache (if enabled)
	data := make([]byte, res.ContentLength)
	if c.cache != nil {
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		c.cache.set(req.URL.String(), data)
	}

	// Reinstate the response body
	res.Body = io.NopCloser(bytes.NewBuffer(data))

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
