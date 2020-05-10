package covid

import (
	"golang.org/x/text/message"
)

// AppData that is passed to consumers/subscribers for interpreting
type AppData interface {
	GetTitle() string
	GetSource() string
	GetCases() string
}

// appData that implements AppData and is passed to consumers/subscribers for interpreting
type appData struct {
	Stats
}

func (d *appData) GetTitle() string {
	return d.printer().Sprintf("😷 %d  ☠️ %d  🥳 %d", d.Active, d.Deaths, d.Recovered)
}

func (d *appData) GetSource() string {
	if d.Country != "" {
		return d.Country
	}
	return "World"
}

func (d *appData) GetCases() string {
	return d.printer().Sprint(d.Cases)
}

func (d *appData) printer() *message.Printer {
	p := message.NewPrinter(message.MatchLanguage("en"))

	return p
}
