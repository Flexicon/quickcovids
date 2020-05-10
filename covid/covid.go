package covid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://corona.lmao.ninja/v2"

// Service is used to communicate with the corona stats API
type Service struct {
	client *http.Client
}

// NewCovidService builder
func NewCovidService() *Service {
	return &Service{
		client: &http.Client{Timeout: time.Second * 10},
	}
}

// Stats represents the data retrieved from the api
// It should contain aggregated statistical data either worldwide or for a specific country
type Stats struct {
	// Timestamp of last data update in milliseconds
	Updated int `json:"updated"`

	// Count of total cases
	Cases int `json:"cases"`

	// Count of today's new cases
	TodayCases int `json:"todayCases"`

	// Count of today's deaths
	TodayDeaths int `json:"todayDeaths"`

	// Count of total deaths
	Deaths int `json:"deaths"`

	// Count of total recoveries
	Recovered int `json:"recovered"`

	// Count of total active cases
	Active int `json:"active"`

	// Count of total critical cases
	Critical int `json:"critical"`

	// Count of total tests done
	Tests int `json:"tests"`

	// The country for which the dataset is for, only present when fetching data for a specific country
	Country string `json:"country,omitempty"`
}

// Country represents the data for a country in the api
type Country struct {
	Country string `json:"country"`
}

// CountriesResponse represents the list of available countries in the api
type CountriesResponse []*Country

// FetchCountries returns all available countries in the api
func (s *Service) FetchCountries() (CountriesResponse, error) {
	var countries CountriesResponse
	err := s.fetchData("/countries", &countries)

	return countries, err
}

// FetchWorldwideData then parse and return the response
func (s *Service) FetchWorldwideData() (*Stats, error) {
	var data Stats
	err := s.fetchData("/all", &data)

	return &data, err
}

// FetchDataForCountry then parse and return the response
func (s *Service) FetchDataForCountry(country string) (*Stats, error) {
	var data Stats
	err := s.fetchData(fmt.Sprintf("/countries/%s", country), &data)

	return &data, err
}

// fetchData then parse and return the response
func (s *Service) fetchData(endpoint string, data interface{}) error {
	resp, err := s.client.Get(fmt.Sprintf("%s%s", baseURL, endpoint))
	if err != nil {
		return fmt.Errorf("Failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("Failed to decode JSON response: %v", err)
	}

	return nil
}
