package server

import (
	"github.com/huyhvq/eurofxref/pkg/handler"
	"github.com/huyhvq/eurofxref/pkg/model"
	"github.com/huyhvq/eurofxref/pkg/repository"
	"github.com/huyhvq/eurofxref/pkg/service/ecb"
	"net/http"
)

type HttpServer interface {
	Start() error
	Initial(repository.RateRepository) error
}

type httpServer struct {
	handler handler.HttpServerHandler
	ecb     ecb.Service
}

func NewHttpServer(h handler.HttpServerHandler, ecb ecb.Service) HttpServer {
	return &httpServer{
		handler: h,
		ecb:     ecb,
	}
}

func (h *httpServer) Initial(r repository.RateRepository) error {
	t, err := r.GetLatestDate()
	if err != nil {
		return err
	}
	rates, err := h.ecb.FetchRatesAfterDate(t)
	if err != nil {
		return err
	}

	rm := make([]model.Rate, 0, len(rates))

	for _, rate := range rates {
		rm = append(rm, model.Rate{
			Time:     rate.Time.Format("2006-01-02"),
			Currency: rate.Currency,
			Rate:     rate.Rate,
		})
	}
	if err := r.InsertMany(rm); err != nil {
		return err
	}
	return nil
}

func (h *httpServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/rates/latest", h.handler.GetLatestRates)
	mux.HandleFunc("/rates/analyze", h.handler.GetRatesAnalyze)
	mux.HandleFunc("/rates/", h.handler.GetRatesByDate)
	return http.ListenAndServe(":8080", mux)
}
