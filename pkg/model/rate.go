package model

type Rate struct {
	Time     string
	Currency string
	Rate     float64
}

type RateAnalyze struct {
	Currency string
	Min      float64
	Max      float64
	Avg      float64
}
