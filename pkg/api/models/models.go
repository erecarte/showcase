package models

type PaymentOrderApiModel struct {
	DebtorIban           string  `json:"debtor_iban,omitempty"`
	DebtorName           string  `json:"debtor_name,omitempty"`
	CreditorIban         string  `json:"creditor_iban,omitempty"`
	CreditorName         string  `json:"creditor_name,omitempty"`
	Ammount              float64 `json:"ammount,omitempty"`
	IdempotencyUniqueKey string  `json:"idempotency_unique_key,omitempty"`
}
