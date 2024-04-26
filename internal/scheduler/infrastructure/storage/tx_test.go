package storage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	. "github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/storage"
	mock "github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/infrastructure/storage/mocks"
)

// tagID 7852 проверка ошибки начала транзакции
func TestTxBeginError(t *testing.T) {
	db, sqlMock := mock.NewSQLMock(t)
	defer db.Close()

	sqlMock.ExpectBegin().WillReturnError(errors.New(`expected error`))

	s := New(db)
	err := s.Tx(context.Background(), func(ctx context.Context, tx *sqlx.Tx) error {
		return nil
	})

	assert.Error(t, err)
	mock.CheckExpectations(t, sqlMock)
}

// tagID 7852 проверка ошибки коллбека
func TestTxCallbackError(t *testing.T) {
	expErr := errors.New(`expected error`)

	db, sqlMock := mock.NewSQLMock(t)
	defer db.Close()

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	s := New(db)
	err := s.Tx(context.Background(), func(ctx context.Context, tx *sqlx.Tx) error {
		return expErr
	})

	assert.EqualError(t, err, expErr.Error())
	mock.CheckExpectations(t, sqlMock)
}

// tagID 7852 проверка ошибки коммита
func TestTxCommitError(t *testing.T) {

	db, sqlMock := mock.NewSQLMock(t)
	defer db.Close()

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit().WillReturnError(errors.New(`expected error`))

	s := New(db)
	err := s.Tx(context.Background(), func(ctx context.Context, tx *sqlx.Tx) error {
		return nil
	})

	assert.Error(t, err)
	mock.CheckExpectations(t, sqlMock)
}

// tagID 7852 проверка успешной транзакции
func TestTxCommitSuccess(t *testing.T) {

	db, sqlMock := mock.NewSQLMock(t)
	defer db.Close()

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	s := New(db)
	err := s.Tx(context.Background(), func(ctx context.Context, tx *sqlx.Tx) error {
		return nil
	})

	assert.NoError(t, err)
	mock.CheckExpectations(t, sqlMock)
}
