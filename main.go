package main

import (
	"log"

	"github.com/getlantern/systray"
)

func main() {
	log.Println("Setting up...")

	app := NewApp()
	systray.Run(onReady(app), onExit())
}

func onReady(app *App) func() {
	return func() {
		systray.SetTitle("⏳")
		systray.SetTooltip("Quick Covid Stats")

		mCurrent := systray.AddMenuItem("Current stats: World", "Where the current dataset comes from")
		mCurrent.Disable()

		mTotalCases := systray.AddMenuItem("Cases: -", "Total cases for the current dataset")
		mTotalCases.Disable()

		systray.AddSeparator()

		mPick := systray.AddMenuItem("Pick a country", "Select country to fetch data from")
		mRefresh := systray.AddMenuItem("Refresh", "Fetch fresh data")

		systray.AddSeparator()

		mQuit := systray.AddMenuItem("Quit", "")

		app.CurrentCountryItem = mCurrent
		app.TotalCasesItem = mTotalCases
		app.PickACountryItem = mPick
		app.RefreshItem = mRefresh

		go func() {
			for {
				select {
				case <-mRefresh.ClickedCh:
					app.RefreshData()
				case <-mQuit.ClickedCh:
					app.Quit()
				}
			}
		}()

		app.BeginDataPolling()
		app.PopulateCountries()
		log.Println("Ready and set up!")
	}
}

func onExit() func() {
	return func() {
		log.Println("Exiting...")
	}
}
