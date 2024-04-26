package mock_storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// NewSQLMock возвращает объект sqlmock.
func NewSQLMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, sqlmock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	return sqlxDB, sqlmock
}

// CheckExpectations проверяет невыполненные запросы.
func CheckExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
