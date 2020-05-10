package covid

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestNewService(t *testing.T) {
	s := NewService()

	assert.IsType(t, &Service{}, s)
}

func TestFetchCountries(t *testing.T) {
	mHTTP := new(mockHTTPClient)
	mRes := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(sampleCountriesResponse()))),
	}

	mHTTP.On("Get", "https://corona.lmao.ninja/v2/countries").Return(mRes, nil)

	s := &Service{client: mHTTP}

	res, err := s.FetchCountries()

	assert.Nil(t, err)
	assert.IsType(t, CountriesResponse{}, res)
	assert.Len(t, res, 2)

	assert.Equal(t, "Narnia", res[0].Country)
	assert.Equal(t, "Neverland", res[1].Country)
}

func TestFetchWorldwideData(t *testing.T) {
	mHTTP := new(mockHTTPClient)
	mRes := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(sampleWorldwideData()))),
	}

	mHTTP.On("Get", "https://corona.lmao.ninja/v2/all").Return(mRes, nil)

	s := &Service{client: mHTTP}

	res, err := s.FetchWorldwideData()

	assert.Nil(t, err)
	assert.IsType(t, &Stats{}, res)
	assert.Equal(t, "", res.Country)
	assert.Equal(t, 25821, res.Cases)
}

func TestFetchDataForCountry(t *testing.T) {
	mHTTP := new(mockHTTPClient)
	mRes := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(sampleCountryData()))),
	}

	mHTTP.On("Get", "https://corona.lmao.ninja/v2/countries/Narnia").Return(mRes, nil)

	s := &Service{client: mHTTP}

	res, err := s.FetchDataForCountry("Narnia")

	assert.Nil(t, err)
	assert.IsType(t, &Stats{}, res)
	assert.Equal(t, "Narnia", res.Country)
	assert.Equal(t, 15821, res.Cases)
}

func TestFetchDataHTTPError(t *testing.T) {
	mHTTP := new(mockHTTPClient)
	mRes := &http.Response{}
	mErr := errors.New("Failed to make request")

	mHTTP.On("Get", "https://corona.lmao.ninja/v2/some-endpoint").Return(mRes, mErr)

	s := &Service{client: mHTTP}

	var res struct{}
	err := s.fetchData("/some-endpoint", &res)

	assert.Empty(t, res)
	assert.Equal(t, fmt.Errorf("Failed to fetch data: %v", mErr), err)
}

func TestFetchDataJSONError(t *testing.T) {
	mHTTP := new(mockHTTPClient)
	mRes := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("{"))),
	}

	mHTTP.On("Get", "https://corona.lmao.ninja/v2/some-endpoint").Return(mRes, nil)

	s := &Service{client: mHTTP}

	var res struct{}
	err := s.fetchData("/some-endpoint", &res)

	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to decode JSON response")
}

func sampleCountriesResponse() string {
	return `[
		{
			"updated": 1589109722180,
			"country": "Narnia",
			"cases": 15821,
			"todayCases": 170,
			"deaths": 791,
			"todayDeaths": 6,
			"recovered": 5698,
			"active": 9332,
			"critical": 160,
			"casesPerOneMillion": 418,
			"deathsPerOneMillion": 21,
			"tests": 460686,
			"testsPerOneMillion": 12172,
			"continent": "Europe"
		},
		{
			"updated": 1589109722232,
			"country": "Neverland",
			"cases": 5558,
			"todayCases": 0,
			"deaths": 494,
			"todayDeaths": 0,
			"recovered": 2546,
			"active": 2518,
			"critical": 22,
			"casesPerOneMillion": 127,
			"deathsPerOneMillion": 11,
			"tests": 6500,
			"testsPerOneMillion": 148,
			"continent": "North America"
		}
	]`
}

func sampleCountryData() string {
	return `{
		"updated": 1589109722180,
		"country": "Narnia",
		"cases": 15821,
		"todayCases": 170,
		"deaths": 791,
		"todayDeaths": 6,
		"recovered": 5698,
		"active": 9332,
		"critical": 160,
		"casesPerOneMillion": 418,
		"deathsPerOneMillion": 21,
		"tests": 460686,
		"testsPerOneMillion": 12172,
		"continent": "Europe"
	}`
}

func sampleWorldwideData() string {
	return `{
		"updated": 1589109722180,
		"cases": 25821,
		"todayCases": 170,
		"deaths": 791,
		"todayDeaths": 6,
		"recovered": 5698,
		"active": 9332,
		"critical": 160,
		"casesPerOneMillion": 418,
		"deathsPerOneMillion": 21,
		"tests": 460686,
		"testsPerOneMillion": 12172
	}`
}
