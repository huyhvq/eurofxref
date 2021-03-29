package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/huyhvq/eurofxref/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var rates = []model.Rate{
	{
		Time:     "2021-03-22",
		Currency: "USD",
		Rate:     1.345,
	},
	{
		Time:     "2021-03-23",
		Currency: "VND",
		Rate:     27212.22,
	},
	{
		Time:     "2021-03-24",
		Currency: "SGD",
		Rate:     1.59,
	}, {
		Time:     "2021-03-25",
		Currency: "JPY",
		Rate:     128.86,
	},
}

func TestNewRate(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	r := NewRate(db)
	assert.NotNil(t, r)
}

func TestRateRepo_GetLatestDate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error when opening a stub database connection")
	defer db.Close()

	et, err := time.ParseInLocation("2006-01-02", "2021-03-25", time.UTC)
	assert.Nil(t, err, "Error when parse time")
	columns := []string{"created_at"}
	mock.ExpectQuery("SELECT `created_at` FROM `rates` ORDER BY `created_at` DESC LIMIT 1").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(et))
	r := NewRate(db)
	rt, err := r.GetLatestDate()
	assert.Nil(t, err)
	assert.Equal(t, et, rt)
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestRateRepo_GetLatestDate_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error when opening a stub database connection")
	defer db.Close()
	emptyErr := errors.New("empty")
	mock.ExpectQuery("SELECT `created_at` FROM `rates` ORDER BY `created_at` DESC LIMIT 1").
		WillReturnError(emptyErr)
	r := NewRate(db)
	rt, err := r.GetLatestDate()
	assert.Equal(t, emptyErr, err)
	assert.Equal(t, time.Time{}, rt)
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
	t.Run("Case Error No Rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT `created_at` FROM `rates` ORDER BY `created_at` DESC LIMIT 1").
			WillReturnError(sql.ErrNoRows)
		rt, err = r.GetLatestDate()
		assert.Nil(t, err)
		assert.Equal(t, time.Time{}, rt)
	})
}

func TestRateRepo_GetRatesByDate(t *testing.T) {
	columns := []string{"currency", "rate", "created_at"}
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error when opening a stub database connection")
	defer db.Close()
	et, err := time.ParseInLocation("2006-01-02", "2021-03-25", time.UTC)
	assert.Nil(t, err, "Error when parse time")
	mock.ExpectQuery("SELECT (.+) from `rates` (.+) ORDER BY `currency` ASC").WithArgs("2021-03-25").
		WillReturnRows(sqlmock.NewRows(columns).AddRow("USD", 1.345, "2021-03-25"))
	r := NewRate(db)
	rs, err := r.GetRatesByDate(et)
	assert.Nil(t, err)
	assert.NotNil(t, rs)
	assert.Equal(t, 1, len(rs))
	assert.Equal(t, model.Rate{
		Time:     "2021-03-25",
		Currency: "USD",
		Rate:     1.345,
	}, rs[0])
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestRateRepo_GetRatesByDate_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error when opening a stub database connection")
	defer db.Close()
	et, err := time.ParseInLocation("2006-01-02", "2021-03-25", time.UTC)
	assert.Nil(t, err, "Error when parse time")
	expectedErr := errors.New("expected error")
	mock.ExpectQuery("SELECT (.+) from `rates` (.+) ORDER BY `currency` ASC").WithArgs("2021-03-25").
		WillReturnError(expectedErr)
	r := NewRate(db)
	rs, err := r.GetRatesByDate(et)
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, rs)
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")

	columns := []string{"currency", "rate", "created_at"}
	mock.ExpectQuery("SELECT (.+) from `rates` (.+) ORDER BY `currency` ASC").WithArgs("2021-03-25").
		WillReturnRows(sqlmock.NewRows(columns).AddRow("USD", "ahihi", "2021-03-25"))
	rs, err = r.GetRatesByDate(et)
	assert.NotNil(t, err)
	assert.Nil(t, rs)
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestRateRepo_InsertMany(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error when opening a stub database connection")
	defer db.Close()

	mock.ExpectBegin()
	ep := mock.ExpectPrepare("INSERT INTO rates\\(currency, rate, created_at\\) VALUES \\(\\?, \\?, \\?\\)")
	for _, rate := range rates {
		ep.ExpectExec().WithArgs(rate.Currency, rate.Rate, rate.Time).WillReturnResult(sqlmock.NewResult(1, 1))
	}
	mock.ExpectCommit()
	r := NewRate(db)
	err = r.InsertMany(rates)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestRateRepo_InsertMany_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error when opening a stub database connection")
	defer db.Close()

	mock.ExpectBegin()
	epf := mock.ExpectPrepare("INSERT INTO rates\\(currency, rate, created_at\\) VALUES \\(\\?, \\?, \\?\\)")
	expectedErr := errors.New("expected error")
	epf.ExpectExec().WillReturnError(expectedErr)
	mock.ExpectRollback()
	r := NewRate(db)
	err = r.InsertMany(rates)
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)

	mock.ExpectBegin().WillReturnError(expectedErr)
	err = r.InsertMany(rates)
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestRateRepo_GetLatestRates(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error when opening a stub database connection")
	defer db.Close()

	et, err := time.ParseInLocation("2006-01-02", "2021-03-25", time.UTC)
	assert.Nil(t, err, "Error when parse time")
	mock.ExpectQuery("SELECT `created_at` FROM `rates` ORDER BY `created_at` DESC LIMIT 1").
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(et))

	mock.ExpectQuery("SELECT (.+) from `rates` (.+) ORDER BY `currency` ASC").WithArgs("2021-03-25").
		WillReturnRows(sqlmock.NewRows([]string{"currency", "rate", "created_at"}).AddRow("USD", 1.345, "2021-03-25"))
	rs, err := NewRate(db).GetLatestRates()
	assert.Nil(t, err)
	assert.NotNil(t, rs)
	assert.Equal(t, 1, len(rs))
	assert.Equal(t, model.Rate{
		Time:     "2021-03-25",
		Currency: "USD",
		Rate:     1.345,
	}, rs[0])
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestRateRepo_GetLatestRates_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error when opening a stub database connection")
	defer db.Close()
	expectedErr := errors.New("expected error")
	mock.ExpectQuery("SELECT `created_at` FROM `rates` ORDER BY `created_at` DESC LIMIT 1").
		WillReturnError(expectedErr)

	rs, err := NewRate(db).GetLatestRates()
	assert.Nil(t, rs)
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestRateRepo_GetRatesAnalyze(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error when opening a stub database connection")
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM rates GROUP BY currency ORDER BY currency ASC").
		WillReturnRows(sqlmock.NewRows([]string{"currency", "avg_rate", "min_rate", "max_rate"}).AddRow("USD", 1.456, 1.345, 1.567))
	r := NewRate(db)
	rs, err := r.GetRatesAnalyze()
	assert.Nil(t, err)
	assert.NotNil(t, rs)
	assert.Equal(t, 1, len(rs))
	assert.Equal(t, []model.RateAnalyze{{
		Currency: "USD",
		Min:      1.345,
		Max:      1.567,
		Avg:      1.456,
	}}, rs)
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestRateRepo_GetRatesAnalyze_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error when opening a stub database connection")
	defer db.Close()
	expectedErr := errors.New("expected error")

	mock.ExpectQuery("SELECT (.+) FROM rates GROUP BY currency ORDER BY currency ASC").
		WillReturnError(expectedErr)
	r := NewRate(db)
	rs, err := r.GetRatesAnalyze()
	assert.Nil(t, rs)
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)

	mock.ExpectQuery("SELECT (.+) FROM rates GROUP BY currency ORDER BY currency ASC").
		WillReturnRows(sqlmock.NewRows([]string{"currency", "avg_rate", "min_rate", "max_rate"}).AddRow("USD", 1.456, 1.345, "error"))
	rs, err = r.GetRatesAnalyze()
	assert.Nil(t, rs)
	assert.NotNil(t, err)
	assert.NotEqual(t, expectedErr, err)
	assert.Nil(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}
