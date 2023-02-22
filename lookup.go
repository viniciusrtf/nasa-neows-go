package neows

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LookupService struct {
	Client  *Client
	BaseURL string
}

func NewLookupService(client *Client) *LookupService {
	return &LookupService{
		Client:  client,
		BaseURL: "https://api.nasa.gov/neo/rest/v1/neo",
	}
}

func (s *LookupService) Find(id string) (*NearEarthObject, error) {
	req, err := http.NewRequest("GET", s.BaseURL+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check if the response is valid
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid HTTP status code %d while fetching the Feed", res.StatusCode)
	}

	// Parse the response
	neo := &NearEarthObject{}
	err = s.parseResponse(res, neo)
	if err != nil {
		return nil, err
	}
	return neo, nil
}

func (s *LookupService) parseResponse(res *http.Response, neo *NearEarthObject) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, neo)
	if err != nil {
		return err
	}
	return nil
}
