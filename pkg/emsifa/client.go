package emsifa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WilayahProvider defines the contract for checking Indonesian regions
type WilayahProvider interface {
	IsValidProvinsi(ctx context.Context, idProvinsi string) (bool, error)
	IsValidKota(ctx context.Context, idProvinsi, idKota string) (bool, error)
}

type emsifaClient struct {
	httpClient *http.Client
	baseURL    string
}

// emsifaRegion represents a generic region item (Province, Regency, etc.)
type emsifaRegion struct {
	ID string `json:"id"`
}

// NewEmsifaClient initializes the external HTTP client targeting Emsifa v2
func NewEmsifaClient() WilayahProvider {
	return &emsifaClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // Prevent hanging requests if Emsifa is down
		},
		baseURL: "https://www.emsifa.com/api-wilayah-indonesia/v2",
	}
}

// fetchAndCheckID is a generic helper to download a JSON array from Emsifa and search for a specific ID
func (c *emsifaClient) fetchAndCheckID(ctx context.Context, endpoint string, targetID string) (bool, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err // Network or timeout error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// If Emsifa returns 404 (e.g., trying to fetch cities for a fake province ID)
		return false, nil
	}

	// Parse the JSON array wrapped in 'data' object
	var response struct {
		Data []emsifaRegion `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to decode json from emsifa: %w", err)
	}

	// Loop through the list to find the matching ID
	for _, region := range response.Data {
		if region.ID == targetID {
			return true, nil // Found it!
		}
	}

	return false, nil // Valid JSON array, but the ID wasn't in it
}

// IsValidProvinsi fetches all provinces and verifies if idProvinsi exists in the array
func (c *emsifaClient) IsValidProvinsi(ctx context.Context, idProvinsi string) (bool, error) {
	return c.fetchAndCheckID(ctx, "provinces.json", idProvinsi)
}

// IsValidKota fetches all regencies for a specific province and verifies if idKota exists in the array
func (c *emsifaClient) IsValidKota(ctx context.Context, idProvinsi, idKota string) (bool, error) {
	endpoint := fmt.Sprintf("regencies/%s.json", idProvinsi)
	return c.fetchAndCheckID(ctx, endpoint, idKota)
}
