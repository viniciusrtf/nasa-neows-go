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

type NeoByDate map[string][]NearEarthObject

func (n NeoByDate) Get(date time.Time) []NearEarthObject {
	return n[date.Format("2006-01-02")]
}

type Feed struct {
	Links            Links     `json:"links"`
	ElementCount     int       `json:"element_count"`
	NearEarthObjects NeoByDate `json:"near_earth_objects"`
}

type FeedService struct {
	Client  *Client
	BaseURL string
}

func NewFeedService(client *Client) *FeedService {
	return &FeedService{
		Client:  client,
		BaseURL: "https://api.nasa.gov/neo/rest/v1/feed",
	}
}

func (s *FeedService) Fetch(startDate time.Time, endDate time.Time) (*Feed, error) {
	req, err := http.NewRequest("GET", s.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("start_date", startDate.Format("2006-01-02"))
	q.Add("end_date", endDate.Format("2006-01-02"))
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
	err = s.ParseResponse(res, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

func (s *FeedService) ParseResponse(res *http.Response, feed *Feed) error {
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
