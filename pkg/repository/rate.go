package repository

import (
	"database/sql"
	"github.com/huyhvq/eurofxref/pkg/model"
	"time"
)

type RateRepository interface {
	InsertMany([]model.Rate) error
	GetLatestDate() (time.Time, error)
	GetLatestRates() ([]model.Rate, error)
	GetRatesAnalyze() ([]model.RateAnalyze, error)
	GetRatesByDate(date time.Time) ([]model.Rate, error)
}

type rateRepo struct {
	db *sql.DB
}

func NewRate(db *sql.DB) RateRepository {
	return &rateRepo{db: db}
}

func (r *rateRepo) InsertMany(rates []model.Rate) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	q := "INSERT INTO rates(currency, rate, created_at) VALUES (?, ?, ?)"
	stmt, err := tx.Prepare(q)
	for _, rate := range rates {
		if _, err := stmt.Exec(rate.Currency, rate.Rate, rate.Time); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *rateRepo) GetLatestDate() (time.Time, error) {
	var lds time.Time
	if err := r.db.QueryRow("SELECT `created_at` FROM `rates` ORDER BY `created_at` DESC LIMIT 1").Scan(&lds); err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	return lds.UTC(), nil
}

func (r *rateRepo) GetLatestRates() ([]model.Rate, error) {
	d, err := r.GetLatestDate()
	if err != nil {
		return nil, err
	}
	return r.GetRatesByDate(d)
}

func (r *rateRepo) GetRatesByDate(date time.Time) ([]model.Rate, error) {
	rates := make([]model.Rate, 0)
	q := "SELECT `currency`,`rate`,`created_at` from `rates` WHERE `created_at`= ? ORDER BY `currency` ASC"
	results, err := r.db.Query(q, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var rate model.Rate
		if err := results.Scan(&rate.Currency, &rate.Rate, &rate.Time); err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}
	return rates, nil
}

func (r *rateRepo) GetRatesAnalyze() ([]model.RateAnalyze, error) {
	rates := make([]model.RateAnalyze, 0)
	q := "SELECT currency, AVG(rate) as avg_rate, MIN(rate) AS min_rate, MAX(rate) AS max_rate FROM rates GROUP BY currency ORDER BY currency ASC"
	results, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	for results.Next() {
		var rate model.RateAnalyze
		if err := results.Scan(&rate.Currency, &rate.Avg, &rate.Min, &rate.Max); err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}
	return rates, nil
}
