package payment_orders

import (
	"context"
	"fmt"
	"log/slog"
)

type Service struct {
	storage              PaymentOrderStorage
	paymentFileGenerator *BankFileSender
}

func NewService(storage PaymentOrderStorage, generator *BankFileSender) *Service {
	return &Service{storage: storage, paymentFileGenerator: generator}
}

func (s *Service) CreatePaymentOrder(ctx context.Context, po *PaymentOrder) (*PaymentOrder, error) {
	slog.Info("creating payment order", "id", po.IdempotencyUniqueKey, "status", po.Status)
	err := s.storage.Store(ctx, po)
	if err != nil {
		return nil, fmt.Errorf("failed to store payment order: %w, %w", err, ErrRecordAlreadyExists)
	}
	err = s.paymentFileGenerator.SendFileToBank(po)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment file: %w", err)
	}
	return po, nil
}

func (s *Service) GetPaymentOrder(ctx context.Context, id string) (*PaymentOrder, error) {
	po, err := s.storage.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment order: %w", err)
	}
	return po, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status string) error {
	slog.Info("updating payment order", "id", id, "status", status)
	return s.storage.UpdateStatus(ctx, id, status)
}
