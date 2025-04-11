package payment_orders

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DbStorage struct {
	db *sql.DB
}

func NewSQLLiteStorage(db *sql.DB) *DbStorage {
	return &DbStorage{db: db}
}

func (s *DbStorage) Store(ctx context.Context, po *PaymentOrder) error {
	result, err := s.db.ExecContext(
		ctx,
		`INSERT INTO payment_orders (
                            idempotency_key,
                            debtor_iban, 
                            debtor_name,
                            creditor_iban, 
                            creditor_name,
                            amount,
                            status,
                            created_at
                            ) VALUES (?,?,?,?,?,?,?,?);`,
		po.IdempotencyUniqueKey,
		po.DebtorIban,
		po.DebtorName,
		po.CreditorIban,
		po.CreditorName,
		po.Amount,
		po.Status,
		po.CreatedAt,
	)
	if err != nil {
		fmt.Println("error inserting payment_order", err)
		return err
	}
	if numRows, err := result.RowsAffected(); numRows == 0 || err != nil {
		return sql.ErrNoRows
	}
	return nil
}

func (s *DbStorage) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE payment_orders SET status = ? WHERE idempotency_key = ?", status, id)
	if err != nil {
		return fmt.Errorf("error updating payment order: %w", err)
	}
	return nil
}

func (s *DbStorage) Get(ctx context.Context, id string) (*PaymentOrder, error) {
	var po PaymentOrder
	row := s.db.QueryRowContext(
		ctx,
		`SELECT * FROM payment_orders WHERE idempotency_key=?;`, id,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var createdAt string
	err := row.Scan(&po.IdempotencyUniqueKey, &po.DebtorIban, &po.DebtorName, &po.CreditorIban, &po.CreditorName, &po.Amount, &po.Status, &createdAt)
	if err != nil {
		return nil, err
	}
	layout := "2006-01-02T15:04:05-0700"
	date, err := time.Parse(layout, createdAt)
	po.CreatedAt = &date
	return &po, nil
}
