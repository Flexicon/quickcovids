package covid

import (
	"log"
	"time"
)

// StatsAPI represents the interface for communicating with the covid API
type StatsAPI interface {
	FetchCountries() (CountriesResponse, error)
	FetchWorldwideData() (*Stats, error)
	FetchDataForCountry(country string) (*Stats, error)
}

// App represents the covid stats application state
type App struct {
	covid    StatsAPI
	country  string
	fetching bool
	updateUI chan AppData
}

// NewApp builder
func NewApp(c StatsAPI) *App {
	return &App{
		covid:    c,
		updateUI: make(chan AppData),
	}
}

// BeginDataPolling to keep stats up to date asynchronously
func (a *App) BeginDataPolling() chan AppData {
	t := time.NewTicker(time.Minute * 30)

	go func() {
		a.RefreshData()

		for {
			select {
			case <-t.C:
				a.RefreshData()
			}
		}
	}()

	return a.updateUI
}

// PrepareCountryNames fetches and prepares a list of available countries
func (a *App) PrepareCountryNames() []string {
	log.Println("Fetching countries...")

	countries, err := a.covid.FetchCountries()
	if err != nil {
		log.Fatal(err)
	}

	names := make([]string, len(countries))
	for i, c := range countries {
		names[i] = c.Country
	}

	return names
}

// SelectCountry to fetch data from, triggers a data fetch
func (a *App) SelectCountry(c string) {
	a.country = c
	a.RefreshData()
}

// RefreshData triggers the data to be updated
func (a *App) RefreshData() {
	if a.fetching {
		return
	}
	log.Println("Fetching data...")
	a.fetching = true
	go a.emitUpdate(&appData{fetching: a.fetching})

	data, err := a.fetchStats()
	if err != nil {
		log.Fatalln(err)
	}
	if data == nil {
		log.Fatalln("Empty response")
	}

	log.Printf("Data fetched: %+v\n", data)
	a.fetching = false
	go a.emitUpdate(&appData{
		*data,
		a.fetching,
	})
}

func (a *App) emitUpdate(d *appData) {
	a.updateUI <- d
}

func (a *App) fetchStats() (*Stats, error) {
	if a.country != "" {
		return a.covid.FetchDataForCountry(a.country)
	}
	return a.covid.FetchWorldwideData()
}
