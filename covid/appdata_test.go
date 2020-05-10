package covid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTitle(t *testing.T) {
	var tests = []struct {
		in  appData
		out string
	}{
		{appData{}, "😷 0  ☠️ 0  🥳 0"},
		{appData{fetching: true}, "⏳"},
		{appData{Stats{Active: 12, Deaths: 10, Recovered: 5}, true}, "⏳"},
		{sampleData(), "😷 123,000  ☠️ 2,000  🥳 25,000"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			assert.Equal(t, tt.out, tt.in.GetTitle())
		})
	}
}

func TestGetSource(t *testing.T) {
	var tests = []struct {
		in  appData
		out string
	}{
		{appData{}, "World"},
		{appData{Stats{Country: "Narnia"}, false}, "Narnia"},
		{sampleData(), "Neverland"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			assert.Equal(t, tt.out, tt.in.GetSource())
		})
	}
}

func TestGetCases(t *testing.T) {
	var tests = []struct {
		in  appData
		out string
	}{
		{appData{}, "0"},
		{appData{Stats{Cases: 1}, false}, "1"},
		{appData{Stats{Cases: 100}, false}, "100"},
		{appData{Stats{Cases: 1000}, false}, "1,000"},
		{appData{Stats{Cases: 123999999}, false}, "123,999,999"},
		{sampleData(), "150,000"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			assert.Equal(t, tt.out, tt.in.GetCases())
		})
	}
}

func TestIsFetching(t *testing.T) {
	var tests = []struct {
		in  appData
		out bool
	}{
		{appData{}, false},
		{appData{fetching: false}, false},
		{appData{fetching: true}, true},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			assert.Equal(t, tt.out, tt.in.IsFetching())
		})
	}
}

func sampleData() appData {
	return appData{
		Stats{
			Updated:   1589106525130,
			Cases:     150000,
			Deaths:    2000,
			Recovered: 25000,
			Active:    123000,
			Country:   "Neverland",
		},
		false,
	}
}
