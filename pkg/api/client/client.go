package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/erecarte/showcase/pkg/api/models"
	"io"
	"net/http"
	"net/url"
)

const (
	headerContentType = "Content-Type"
	headerAccept      = "Accept"
)

type ClientError struct {
	models.ApiError
	StatusCode int
}

type NumeralApiClientOpt func(*NumeralApiClient)
type NumeralApiClient struct {
	PaymentOrders *PaymentOrderApi
	baseURL       *url.URL
	httpClient    *http.Client
	username      string
	password      string
}

func WithHTTPClient(httpClient *http.Client) NumeralApiClientOpt {
	return func(n *NumeralApiClient) {
		n.httpClient = httpClient
	}
}

func NewNumeralApiClient(baseURL string, username, password string, opts ...NumeralApiClientOpt) (*NumeralApiClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("create new numeral client: %w", err)
	}
	client := &NumeralApiClient{
		baseURL:    u,
		httpClient: &http.Client{},
		username:   username,
		password:   password,
	}
	paymentOrderApi := NewPaymentOrderApi(client)
	client.PaymentOrders = paymentOrderApi
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

func (c NumeralApiClient) post(ctx context.Context, path string, payload any, v any) error {
	joinPath := c.baseURL.JoinPath(path)
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, joinPath.String(), bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("create new request: %w", err)
	}
	req.Header.Set(headerContentType, "application/json")
	return c.sendRequest(req, v)
}

func (c NumeralApiClient) get(ctx context.Context, path string, id string, v any) error {
	joinPath := c.baseURL.JoinPath(path, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, joinPath.String(), nil)
	if err != nil {
		return fmt.Errorf("create new request: %w", err)
	}
	err = c.sendRequest(req, v)
	if err != nil {
		return err
	}
	return nil
}

func (c NumeralApiClient) sendRequest(req *http.Request, result any) error {
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set(headerAccept, "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		v := &ClientError{
			StatusCode: resp.StatusCode,
		}
		err := json.Unmarshal(responseBody, v)
		if err != nil {
			return fmt.Errorf("unmarshal error response body: %w", err)
		}
		return v
	}
	err = json.Unmarshal(responseBody, result)
	if err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}
	return nil
}
