package ecb

import (
	"time"
)

type Rate struct {
	Time     time.Time
	Currency string
	Rate     float64
}

type HistoryResponse struct {
	Cube []struct {
		Time string                `xml:"time,attr"`
		Cube []HistoryCubeResponse `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

type HistoryCubeResponse struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}
