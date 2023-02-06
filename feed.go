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
	Prev string `json:"previous"`
}

type Feed struct {
	Links        Links                        `json:"links"`
	ElementCount int                          `json:"element_count"`
	NeosByDate   map[string][]NearEarthObject `json:"near_earth_objects"`
}

type FeedService struct {
	Client  *Client
	BaseURL string
}

type FeedOptions struct {
	StartDate time.Time
	EndDate   time.Time
	Detailed  bool
}

// NewFeedService creates a new FeedService
func NewFeedService(client *Client) *FeedService {
	return &FeedService{
		Client:  client,
		BaseURL: "https://api.nasa.gov/neo/rest/v1/feed",
	}
}

// NewFeedOptions creates a new FeedOptions with default values
func NewFeedOptions() *FeedOptions {
	return &FeedOptions{
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
		Detailed:  false,
	}
}

// Fetch fetches approaching asteroids for the given date range
func (s *FeedService) Fetch(opts *FeedOptions) (*Feed, error) {
	if opts == nil {
		opts = NewFeedOptions()
	}
	if opts.StartDate.IsZero() {
		opts.StartDate = time.Now()
	}
	if opts.EndDate.IsZero() {
		opts.EndDate = time.Now().AddDate(0, 0, 7)
	}

	req, err := http.NewRequest("GET", s.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("start_date", opts.StartDate.Format("2006-01-02"))
	q.Add("end_date", opts.EndDate.Format("2006-01-02"))
	q.Add("detailed", fmt.Sprintf("%t", opts.Detailed))
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
