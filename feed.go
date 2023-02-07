package neows

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Links struct {
	Self string `json:"self"`
	Next string `json:"next"`
	Prev string `json:"prev"`
}

type Feed struct {
	Links        Links                        `json:"links"`
	ElementCount int                          `json:"element_count"`
	NeosByDate   map[string][]NearEarthObject `json:"near_earth_objects"`
}

type FeedService struct {
	Client  *Client
	BaseURL string
	opts    *FeedOptions
}

type FeedOptions struct {
	Detailed bool
}

// NewFeedService creates a new FeedService
func NewFeedService(client *Client, opts *FeedOptions) *FeedService {
	return &FeedService{
		Client:  client,
		BaseURL: "https://api.nasa.gov/neo/rest/v1/feed",
		opts:    opts,
	}
}

// NewDefaultFeedOptions creates a new FeedOptions with default values
func NewDefaultFeedOptions() *FeedOptions {
	return &FeedOptions{
		Detailed: false,
	}
}

// Fetch fetches approaching asteroids for the given date range
func (s *FeedService) Fetch(start time.Time, end time.Time) (*Feed, error) {
	if start.IsZero() {
		start = time.Now()
	}
	if end.IsZero() {
		end = start.AddDate(0, 0, 7)
	}

	req, err := http.NewRequest("GET", s.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("start_date", start.Format("2006-01-02"))
	q.Add("end_date", end.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	return s.doRequest(req)
}

// Today fetches approaching asteroids for today
func (s *FeedService) Today() (*Feed, error) {
	return s.Fetch(time.Now(), time.Now())
}

// Next fetches the next page of the given Feed
func (s *FeedService) Next(current *Feed) (*Feed, error) {
	// Last page was reached
	if current.Links.Next == "" {
		return nil, nil
	}

	req, err := http.NewRequest("GET", current.Links.Next, nil)
	if err != nil {
		return nil, err
	}

	return s.doRequest(req)
}

// Prev fetches the previous page of the given Feed
func (s *FeedService) Prev(current *Feed) (*Feed, error) {
	// This is the first page
	if current.Links.Prev == "" {
		return nil, nil
	}

	req, err := http.NewRequest("GET", current.Links.Prev, nil)
	if err != nil {
		return nil, err
	}

	return s.doRequest(req)
}

func (s *FeedService) doRequest(req *http.Request) (*Feed, error) {
	// Add optional params
	q := req.URL.Query()
	if s.opts != nil {
		q.Add("detailed", fmt.Sprintf("%t", s.opts.Detailed))
	}
	req.URL.RawQuery = q.Encode()

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check if the response is valid
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid HTTP status code %d while fetching the Feed", res.StatusCode)
	}

	// Parse the response
	var feed Feed
	err = s.parseResponse(res, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

func (s *FeedService) parseResponse(res *http.Response, feed *Feed) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, feed)
	if err != nil {
		return err
	}
	return nil
}
