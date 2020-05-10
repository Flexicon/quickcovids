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

func TestRefreshData(t *testing.T) {
	mService := &mockStatsAPI{}
	a := NewApp(mService)

	mService.On("FetchWorldwideData").Return(sampleWorldStats(), nil)

	a.RefreshData()

	fetchingState := <-a.updateUI
	assert.True(t, fetchingState.IsFetching())

	updatedData := <-a.updateUI
	assert.False(t, updatedData.IsFetching())
	assert.Equal(t, updatedData.GetTitle(), "😷 123,000  ☠️ 2,000  🥳 25,000")
	assert.Equal(t, updatedData.GetSource(), "World")
	assert.Equal(t, updatedData.GetCases(), "150,000")
}

func TestRefreshDataWhileAlreadyFetching(t *testing.T) {
	mService := &mockStatsAPI{}
	a := NewApp(mService)
	a.fetching = true

	mService.On("FetchWorldwideData").Return(&Stats{}, nil)

	a.RefreshData()

	mService.AssertNotCalled(t, "FetchWorldwideData")
}

func TestSelectCountry(t *testing.T) {
	mService := &mockStatsAPI{}
	a := NewApp(mService)

	mService.On("FetchDataForCountry", "Narnia").Return(sampleCountryStats(), nil)

	a.SelectCountry("Narnia")

	fetchingState := <-a.updateUI
	assert.True(t, fetchingState.IsFetching())

	updatedData := <-a.updateUI
	assert.False(t, updatedData.IsFetching())
	assert.Equal(t, updatedData.GetTitle(), "😷 223,000  ☠️ 2,000  🥳 25,000")
	assert.Equal(t, updatedData.GetSource(), "Narnia")
	assert.Equal(t, updatedData.GetCases(), "250,000")
}

func sampleWorldStats() *Stats {
	return &Stats{
		Cases:     150000,
		Deaths:    2000,
		Recovered: 25000,
		Active:    123000,
	}
}

func sampleCountryStats() *Stats {
	return &Stats{
		Cases:     250000,
		Deaths:    2000,
		Recovered: 25000,
		Active:    223000,
		Country:   "Narnia",
	}
}
