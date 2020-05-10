package covid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStatsAPI struct {
	mock.Mock
}

func (m *mockStatsAPI) FetchCountries() (CountriesResponse, error) {
	args := m.Called()
	return args.Get(0).(CountriesResponse), args.Error(1)
}

func (m *mockStatsAPI) FetchWorldwideData() (*Stats, error) {
	args := m.Called()
	return args.Get(0).(*Stats), args.Error(1)
}

func (m *mockStatsAPI) FetchDataForCountry(country string) (*Stats, error) {
	args := m.Called(country)
	return args.Get(0).(*Stats), args.Error(1)
}

func TestNewApp(t *testing.T) {
	a := NewApp(&mockStatsAPI{})

	assert.IsType(t, &App{}, a)
}

func TestPrepareCountryNames(t *testing.T) {
	mService := &mockStatsAPI{}
	a := NewApp(mService)

	mService.On("FetchCountries").Return(CountriesResponse{
		{Country: "Narnia"},
		{Country: "Neverland"},
	}, nil)

	names := a.PrepareCountryNames()

	assert.NotEmpty(t, names)
	assert.Len(t, names, 2)

	expected := []string{"Narnia", "Neverland"}
	for i, actual := range names {
		t.Run(expected[i], func(t *testing.T) {
			assert.Equal(t, expected[i], actual)
		})
	}
}
