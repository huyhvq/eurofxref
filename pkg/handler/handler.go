package handler

import (
	"encoding/json"
	"errors"
	"github.com/huyhvq/eurofxref/pkg/model"
	"github.com/huyhvq/eurofxref/pkg/repository"
	"net/http"
	"time"
)

type HttpServerHandler interface {
	GetLatestRates(w http.ResponseWriter, r *http.Request)
	GetRatesByDate(w http.ResponseWriter, r *http.Request)
	GetRatesAnalyze(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	rateRepo repository.RateRepository
}

type ExchangeRate struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

type ExchangeRateAnalyze struct {
	Base         string                 `json:"base"`
	RatesAnalyze map[string]RateAnalyze `json:"rates_analyze"`
}

type RateAnalyze struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
	Avg float64 `json:"avg"`
}

var (
	errInvalidMethod  = errors.New("invalid method in request")
	errInvalidRequest = errors.New("invalid request")
)

func NewHandler(r repository.RateRepository) HttpServerHandler {
	return &handler{rateRepo: r}
}

func (h *handler) GetLatestRates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorRespond(w, http.StatusMethodNotAllowed, errInvalidMethod.Error())
		return
	}
	rates, err := h.rateRepo.GetLatestRates()
	if err != nil {
		errorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonRespond(w, http.StatusOK, exchangeRateTransform(rates))
	return
}

func (h *handler) GetRatesByDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorRespond(w, http.StatusMethodNotAllowed, errInvalidMethod.Error())
		return
	}
	date := r.URL.Path[len("/rates/"):]
	t, err := time.ParseInLocation("2006-01-02", date, time.UTC)
	if err != nil {
		errorRespond(w, http.StatusNotFound, errInvalidRequest.Error())
		return
	}
	rates, err := h.rateRepo.GetRatesByDate(t)
	if err != nil {
		errorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonRespond(w, http.StatusOK, exchangeRateTransform(rates))
	return
}

func (h *handler) GetRatesAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorRespond(w, http.StatusMethodNotAllowed, errInvalidMethod.Error())
		return
	}
	rates, err := h.rateRepo.GetRatesAnalyze()
	if err != nil {
		errorRespond(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonRespond(w, http.StatusOK, exchangeRateAnalyzeTransform(rates))
	return
}

func exchangeRateTransform(rates []model.Rate) *ExchangeRate {
	rs := make(map[string]float64, len(rates))
	for _, rate := range rates {
		rs[rate.Currency] = rate.Rate
	}
	return &ExchangeRate{
		Base:  "EUR",
		Rates: rs,
	}
}

func exchangeRateAnalyzeTransform(rates []model.RateAnalyze) *ExchangeRateAnalyze {
	r := make(map[string]RateAnalyze, len(rates))
	for _, rate := range rates {
		r[rate.Currency] = RateAnalyze{
			Min: rate.Min,
			Max: rate.Max,
			Avg: rate.Avg,
		}
	}
	return &ExchangeRateAnalyze{
		Base:         "EUR",
		RatesAnalyze: r,
	}
}

func jsonRespond(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func errorRespond(w http.ResponseWriter, code int, message string) {
	jsonRespond(w, code, map[string]string{"error": message})
}
