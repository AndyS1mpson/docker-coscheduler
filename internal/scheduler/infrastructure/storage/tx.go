package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Tx оборачивает в транзакцию выполнение метода
func (s *Storage) Tx(ctx context.Context, f func(context.Context, *sqlx.Tx) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction failed: '%w'", err)
	}

	defer tx.Rollback() //nolint:errcheck

	if err = f(ctx, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: '%w'", err)
	}

	return nil
}
