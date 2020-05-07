package main

import (
	"fmt"
	"log"
	"time"

	"github.com/flexicon/quickcovids/covid"
	"github.com/getlantern/systray"
)

// App represents the main application state
type App struct {
	covid    *covid.Service
	country  string
	fetching bool
	data     *covid.Stats

	CurrentCountryItem *systray.MenuItem
	PickACountryItem   *systray.MenuItem
	RefreshItem        *systray.MenuItem
}

// NewApp builder
func NewApp() *App {
	return &App{
		covid:   covid.NewCovidService(),
		country: "",
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

// PopulateCountries asynchronously prepares a list of available countries and adds them as options to select from
func (a *App) PopulateCountries() {
	log.Println("Fetching countries...")

	go func() {
		countries, err := a.covid.FetchCountries()
		if err != nil {
			log.Fatal(err)
		}

		for _, c := range countries {
			countryItem := a.PickACountryItem.AddSubMenuItem(c.Country, "")
			go a.listenForCountrySelection(c.Country, countryItem)
		}
	}()
}

// RefreshData triggers the data to be updated
func (a *App) RefreshData() {
	if a.fetching {
		return
	}
	log.Println("Fetching data...")
	a.toggleFetchingState()

	data, err := a.fetchData()
	if err != nil {
		log.Fatalln(err)
	}
	if data == nil {
		log.Fatalln("Empty response")
	}

	log.Printf("Data fetched: %+v\n", data)
	a.data = data
	a.updateUI()
	a.toggleFetchingState()
}

// Quit does any necessary App cleanups and quits the systray
func (a *App) Quit() {
	systray.Quit()
}

func (a *App) fetchData() (*covid.Stats, error) {
	if a.country != "" {
		return a.covid.FetchDataForCountry(a.country)
	}
	return a.covid.FetchWorldwideData()
}

func (a *App) toggleFetchingState() {
	a.fetching = !a.fetching

	if a.RefreshItem != nil {
		if a.fetching {
			a.RefreshItem.Disable()
		} else {
			a.RefreshItem.Enable()
		}
	}
}

func (a *App) listenForCountrySelection(country string, item *systray.MenuItem) {
	for range item.ClickedCh {
		log.Printf("Country: %s selected\n", country)

		a.country = country
		a.RefreshData()
	}
}

func (a *App) updateUI() {
	systray.SetTitle(fmt.Sprintf("😷 %d  ☠️ %d  %d", a.data.Active, a.data.Deaths, a.data.Recovered))

	if a.country != "" {
		a.CurrentCountryItem.SetTitle(fmt.Sprintf("Current stats: %s", a.country))
	} else {
		a.CurrentCountryItem.SetTitle("Current stats: World")
	}
}
