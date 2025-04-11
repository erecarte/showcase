package test

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/erecarte/showcase/internal/numeral"
	"github.com/erecarte/showcase/internal/numeral/api/payment_orders"
	"github.com/erecarte/showcase/internal/numeral/database"
	"github.com/erecarte/showcase/pkg/api/client"
	"github.com/erecarte/showcase/pkg/api/models"
	"github.com/erecarte/showcase/test/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

//go:embed request_sample.json
var correctJsonSample []byte

//go:embed incorrect_request_sample.json
var incorrectJsonSample []byte

func TestMain(m *testing.M) {
	app := numeral.NewApp(numeral.NewDefaultConfig())
	app.Start()
	bank.NewApp()
	code := m.Run()
	app.Stop()
	os.Exit(code)
}

func TestPaymentOrderIsNotAcceptedIfNotAuthenticated(t *testing.T) {
	initialPO, err := correctPaymentOrder()
	require.NoError(t, err)
	apiClient, err := client.NewNumeralApiClient("http://localhost:8080/v1", "non_existing_username", "non_existing_password")
	require.NoError(t, err)

	createdPO, err := apiClient.PaymentOrders.Create(context.Background(), initialPO)
	require.Error(t, err)
	require.Nil(t, createdPO)
	clientErr := &client.ClientError{}
	require.True(t, errors.As(err, &clientErr))
	assert.Equal(t, http.StatusUnauthorized, clientErr.StatusCode)
}

func TestPaymentOrderIsNotAcceptedIfJsonFormatIsNotCorrect(t *testing.T) {
	initialPO, err := incorrectPaymentOrder()
	require.NoError(t, err)
	apiClient, err := client.NewNumeralApiClient("http://localhost:8080/v1", "RECARTE", "xxxx")
	require.NoError(t, err)

	createdPO, err := apiClient.PaymentOrders.Create(context.Background(), initialPO)
	require.Error(t, err)
	require.Nil(t, createdPO)
	clientErr := &client.ClientError{}
	require.True(t, errors.As(err, &clientErr))
	assert.Equal(t, http.StatusBadRequest, clientErr.StatusCode)
}

func TestPaymentOrderIsSuccessful(t *testing.T) {
	initialPO, err := correctPaymentOrder()
	require.NoError(t, err)
	apiClient, err := client.NewNumeralApiClient("http://localhost:8080/v1", "RECARTE", "xxxx")
	require.NoError(t, err)

	initialPO.IdempotencyUniqueKey = random.String(10)
	createdPO, err := apiClient.PaymentOrders.Create(context.Background(), initialPO)
	require.NoError(t, err)
	require.NotNil(t, createdPO)

	db, err := database.GetConnection("../data/data.sqlite")
	storage := payment_orders.NewSQLLiteStorage(db)
	assert.Eventually(t, func() bool {

		po, err := storage.Get(context.Background(), initialPO.IdempotencyUniqueKey)
		if err != nil {
			return false
		}
		return po.Status == "PROCESSED"
	}, 3*time.Second, 100*time.Millisecond)
}

func correctPaymentOrder() (*models.PaymentOrderApiModel, error) {
	po := &models.PaymentOrderApiModel{}
	err := json.Unmarshal(correctJsonSample, po)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func incorrectPaymentOrder() (*models.PaymentOrderApiModel, error) {
	po := &models.PaymentOrderApiModel{}
	err := json.Unmarshal(incorrectJsonSample, po)
	if err != nil {
		return nil, err
	}
	return po, nil
}
