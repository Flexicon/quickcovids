package covid

import (
	"log"
	"time"
)

// App represents the covid stats application state
type App struct {
	covid    *Service
	country  string
	fetching bool

	listeners []chan AppData
}

// NewApp builder
func NewApp(c *Service) *App {
	return &App{
		covid: c,
	}
}

// BeginDataPolling to keep stats up to date asynchronously
func (a *App) BeginDataPolling() {
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

	data, err := a.fetchData()
	if err != nil {
		log.Fatalln(err)
	}
	if data == nil {
		log.Fatalln("Empty response")
	}

	log.Printf("Data fetched: %+v\n", data)
	a.pub(data)
}

// Sub to app state/data changes
func (a *App) Sub(c chan AppData) {
	a.listeners = append(a.listeners, c)
}

// Publish data changes to all subscribers
func (a *App) pub(d *Stats) {
	for _, c := range a.listeners {
		c <- &appData{*d}
	}
}

func (a *App) fetchData() (*Stats, error) {
	if a.country != "" {
		return a.covid.FetchDataForCountry(a.country)
	}
	return a.covid.FetchWorldwideData()
}
