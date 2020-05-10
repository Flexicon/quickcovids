package main

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

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
		go handleItemClicks(app, i)
		go populateCountries(app, i)

		updatesCh := app.BeginDataPolling()
		go listenForUpdates(updatesCh, i)

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
	for {
		select {
		case <-i.refresh.ClickedCh:
			app.RefreshData()
		case <-i.quit.ClickedCh:
			systray.Quit()
		}
	}
}

func populateCountries(a *covid.App, i *menuItems) {
	names := a.PrepareCountryNames()

	worldItem := i.pick.AddSubMenuItem("World", "")
	go listenForCountrySelection(a, "", worldItem)

	i.pick.AddSubMenuItem(strings.Repeat("-", getMaxNameLength(names)+1), "") // Separator

	for _, c := range names {
		countryItem := i.pick.AddSubMenuItem(c, "")
		go listenForCountrySelection(a, c, countryItem)
	}
}

func listenForCountrySelection(a *covid.App, c string, ci *systray.MenuItem) {
	for range ci.ClickedCh {
		a.SelectCountry(c)
	}
}

func getMaxNameLength(names []string) int {
	var max int
	for _, n := range names {
		count := utf8.RuneCountInString(n)
		if count > max {
			max = count
		}
	}

	return max
}

func listenForUpdates(updatesCh chan covid.AppData, i *menuItems) {
	for d := range updatesCh {
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
}
