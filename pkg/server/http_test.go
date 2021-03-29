package server

import (
	"errors"
	"github.com/huyhvq/eurofxref/pkg/model"
	"github.com/huyhvq/eurofxref/pkg/service/ecb"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

type mockSrv struct {
}

var (
	fetchRatesAfterDateErr = errors.New("fetch rates after date error")
	getLatestRatesErr      = errors.New("get latest rate error")
	insertManyErr          = errors.New("insert many")
	mt, _                  = time.ParseInLocation("2006-01-02", "2021-03-04", time.UTC)
	ft, _                  = time.ParseInLocation("2006-01-02", "2021-03-03", time.UTC)
)

func (m mockSrv) FetchRatesAfterDate(date time.Time) ([]ecb.Rate, error) {
	if date == ft {
		return []ecb.Rate{{
			Time:     date,
			Currency: "USD",
			Rate:     1,
		}}, nil
	}
	if date == mt {
		return nil, fetchRatesAfterDateErr
	}
	return []ecb.Rate{{
		Time:     date,
		Currency: "USD",
		Rate:     1,
	}, {
		Time:     date,
		Currency: "USD",
		Rate:     1,
	}}, nil
}

type mockHandler struct {
}

func (m mockHandler) GetLatestRates(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (m mockHandler) GetRatesByDate(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (m mockHandler) GetRatesAnalyze(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

type mockRepo struct {
	Date time.Time
}

func (m mockRepo) InsertMany(rates []model.Rate) error {
	if len(rates) == 2 {
		return nil
	}
	return insertManyErr
}

func (m mockRepo) GetLatestDate() (time.Time, error) {
	t := time.Time{}
	if m.Date == t {
		return m.Date, getLatestRatesErr
	}
	return m.Date, nil
}

func (m mockRepo) GetLatestRates() ([]model.Rate, error) {
	panic("implement me")
}

func (m mockRepo) GetRatesAnalyze() ([]model.RateAnalyze, error) {
	panic("implement me")
}

func (m mockRepo) GetRatesByDate(date time.Time) ([]model.Rate, error) {
	panic("implement me")
}

func TestNewHttpServer(t *testing.T) {
	h := NewHttpServer(mockHandler{}, mockSrv{})
	assert.NotNil(t, h)
}

func TestHttpServer_Initial(t *testing.T) {
	t.Run("Initial successful", func(t *testing.T) {
		mtf, _ := time.ParseInLocation("2006-01-02", "2021-03-05", time.UTC)
		err := NewHttpServer(mockHandler{}, mockSrv{}).Initial(mockRepo{Date: mtf})
		assert.Nil(t, err)
	})
	t.Run("Initial failed on GetLatestDate", func(t *testing.T) {
		err := NewHttpServer(mockHandler{}, mockSrv{}).Initial(mockRepo{})
		assert.NotNil(t, err)
		assert.Equal(t, getLatestRatesErr, err)
	})
	t.Run("Initial failed on FetchRatesAfterDate", func(t *testing.T) {
		err := NewHttpServer(mockHandler{}, mockSrv{}).Initial(mockRepo{Date: mt})
		assert.NotNil(t, err)
		assert.Equal(t, fetchRatesAfterDateErr, err)
	})
	t.Run("Initial failed on InsertMany", func(t *testing.T) {
		err := NewHttpServer(mockHandler{}, mockSrv{}).Initial(mockRepo{Date: ft})
		assert.NotNil(t, err)
		assert.Equal(t, insertManyErr, err)
	})
}
