package neows

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLookup(t *testing.T) {

	t.Run("Should get a single asteroid by ID", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the last part of the URL
			urlParts := strings.Split(r.URL.Path, "/")
			id := urlParts[len(urlParts)-1]
			if id != "3542519" {
				t.Errorf("Expected ID to be 3542519, got %s", id)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(NewMockAsteroidJson())
		}))

		lookupService := &LookupService{Client: NewClient(defaultAPIKey), BaseURL: ts.URL}
		neo, err := lookupService.Find("3542519")
		if err != nil {
			t.Error(err)
		}

		if neo.ID != "3542519" {
			t.Errorf("Expected ID to be 3542519, got %s", neo.ID)
		}
	})
}
