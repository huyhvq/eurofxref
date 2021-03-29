package ecb

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Service interface {
	FetchRatesAfterDate(date time.Time) ([]Rate, error)
}

type Config struct {
	Endpoint string
}

type ecbService struct {
	cfg    *Config
	client *http.Client
}

var (
	errUnableToConnect = errors.New("unable to connect")
	errCantReadBody    = errors.New("can not read body")
)

func NewService(cfg *Config) Service {
	hc := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &ecbService{
		cfg:    cfg,
		client: hc,
	}
}

func (s ecbService) fetchAllRates() (*HistoryResponse, error) {
	resp, err := s.client.Get(s.cfg.Endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errUnableToConnect
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errCantReadBody
	}
	var h HistoryResponse
	if err := xml.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

func (s ecbService) FetchRatesAfterDate(date time.Time) ([]Rate, error) {
	totalRates, err := s.fetchAllRates()
	if err != nil {
		return nil, err
	}
	rates := make([]Rate, 0, len(totalRates.Cube))
	for _, rs := range totalRates.Cube {
		t, err := time.ParseInLocation("2006-01-02", rs.Time, time.UTC)
		if err != nil {
			return nil, err
		}
		if t.After(date) {
			for _, r := range rs.Cube {
				rates = append(rates, Rate{
					Time:     t,
					Currency: r.Currency,
					Rate:     r.Rate,
				})
			}
		}
	}
	return rates, nil
}
