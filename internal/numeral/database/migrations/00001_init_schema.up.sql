CREATE TABLE payment_orders
(
    idempotency_key TEXT PRIMARY KEY,
    debtor_iban     TEXT NOT NULL,
    debtor_name     TEXT NOT NULL,
    creditor_iban   TEXT NOT NULL,
    creditor_name   TEXT NOT NULL,
    amount          TEXT NOT NULL,
    status          TEXT NOT NULL,
    created_at      TIME NOT NULL
);