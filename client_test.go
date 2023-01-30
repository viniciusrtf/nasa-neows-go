package neows

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {

	t.Run("Should fallback to defaultAPIKey in every request if no apiKey specified in NewClient", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the api_key query parameter is present
			apiKey := r.URL.Query().Get("api_key")
			if apiKey != defaultAPIKey {
				t.Errorf("Want api_key = %s, but got api_key = %s", defaultAPIKey, apiKey)
			}
		}))
		defer ts.Close()

		client := NewClient(nil)

		// Make a request to the test server with a different query parameter
		url := ts.URL + "?param=value"
		_, err := client.c.Get(url)
		if err != nil {
			t.Error(err)
		}

		// Make another request to the test server without any query parameters
		url = ts.URL
		_, err = client.c.Get(url)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Should append the specified apiKey to every request", func(t *testing.T) {
		testAPIKey := "abracadabra"

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the api_key query parameter is present
			apiKey := r.URL.Query().Get("api_key")
			if apiKey != testAPIKey {
				t.Errorf("Want api_key = %s, but got api_key = %s", testAPIKey, apiKey)
			}
		}))
		defer ts.Close()

		client := NewClient(&testAPIKey)

		// Make a request to the test server with a different query parameter
		url := ts.URL + "?param=value"
		_, err := client.c.Get(url)
		if err != nil {
			t.Error(err)
		}

		// Make another request to the test server without any query parameters
		url = ts.URL
		_, err = client.c.Get(url)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Should retrieve data from cache if cache is enabled", func(t *testing.T) {
		testAPIKey := "abracadabra"
		cachedData := []byte("cached data")

		// Create a mock server to return cached data
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(cachedData)
		}))
		defer ts.Close()

		// Enable cache
		client := NewClient(&testAPIKey)
		client.EnableCache()

		// Make a request to the mock server to cache the data
		url := ts.URL
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Error(err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		resBody, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(resBody, cachedData) {
			t.Errorf("Want %s, but got %s", cachedData, resBody)
		}

		// Make another request to the mock server to get the cached data
		res, err = client.Do(req)
		if err != nil {
			t.Error(err)
		}
		resBody, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(resBody, cachedData) {
			t.Errorf("Want %s, but got %s", cachedData, resBody)
		}
		if v := res.Header.Get("X-NASA-NeoWs-Go-Cache"); v != "HIT" {
			t.Errorf("Want X-NASA-NeoWs-Go-Cache = HIT, but got %s", v)
		}
	})
}
