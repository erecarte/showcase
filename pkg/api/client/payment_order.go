package client

import (
	"context"
	"github.com/erecarte/showcase/pkg/api/models"
)

type PaymentOrderApi struct {
	client *NumeralApiClient
}

func NewPaymentOrderApi(client *NumeralApiClient) *PaymentOrderApi {
	return &PaymentOrderApi{client: client}
}

func (a *PaymentOrderApi) Create(ctx context.Context, paymentOrder *models.PaymentOrderApiModel) (*models.PaymentOrderApiModel, error) {
	result := &models.PaymentOrderApiModel{}
	err := a.client.post(ctx, "/payment_orders", paymentOrder, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *PaymentOrderApi) Retrieve(ctx context.Context, id string) (*models.PaymentOrderApiModel, error) {
	result := &models.PaymentOrderApiModel{}
	err := a.client.get(ctx, "/payment_orders", id, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
