package main

import (
	"fmt"
	"log"

	"github.com/flexicon/quickcovids/covid"
	"github.com/getlantern/systray"
)

type menuItems struct {
	current *systray.MenuItem
	total   *systray.MenuItem
	pick    *systray.MenuItem
	refresh *systray.MenuItem
	quit    *systray.MenuItem
}

func main() {
	log.Println("Setting up...")

	c := covid.NewService()
	app := covid.NewApp(c)

	systray.Run(onReady(app), onExit)
}

func onReady(app *covid.App) func() {
	return func() {
		systray.SetTitle("⏳")
		systray.SetTooltip("Quick Covid Stats")

		i := setupMenuItems()
		handleItemClicks(app, i)
		populateCountries(app, i)

		app.Sub(listenForUpdates(i))
		app.BeginDataPolling()

		log.Println("Ready and set up!")
	}
}

func onExit() {
	log.Println("Exiting...")
}

func setupMenuItems() *menuItems {
	current := systray.AddMenuItem("Current stats: World", "Where the current dataset comes from")
	current.Disable()

	total := systray.AddMenuItem("Cases: -", "Total cases for the current dataset")
	total.Disable()

	systray.AddSeparator()

	pick := systray.AddMenuItem("Pick a country", "Select country to fetch data from")
	refresh := systray.AddMenuItem("Refresh", "Fetch fresh data")

	systray.AddSeparator()

	quit := systray.AddMenuItem("Quit", "")

	return &menuItems{
		current: current,
		total:   total,
		pick:    pick,
		refresh: refresh,
		quit:    quit,
	}
}

func handleItemClicks(app *covid.App, i *menuItems) {
	go func() {
		for {
			select {
			case <-i.refresh.ClickedCh:
				app.RefreshData()
			case <-i.quit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func populateCountries(a *covid.App, i *menuItems) {
	go func() {
		names := a.PrepareCountryNames()

		for _, c := range names {
			countryItem := i.pick.AddSubMenuItem(c, "")

			go func(country string) {
				for range countryItem.ClickedCh {
					a.SelectCountry(country)
				}
			}(c)
		}
	}()
}

func listenForUpdates(i *menuItems) chan covid.AppData {
	updateUI := make(chan covid.AppData)

	go func() {
		for d := range updateUI {
			systray.SetTitle(d.GetTitle())
			i.current.SetTitle(fmt.Sprintf("Current stats: %s", d.GetSource()))
			i.total.SetTitle(fmt.Sprintf("Cases: %s", d.GetCases()))

			if d.IsFetching() {
				i.pick.Disable()
				i.refresh.Disable()
			} else {
				i.pick.Enable()
				i.refresh.Enable()
			}
		}
	}()

	return updateUI
}
