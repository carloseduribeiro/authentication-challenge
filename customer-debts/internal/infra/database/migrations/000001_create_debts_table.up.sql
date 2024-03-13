CREATE SCHEMA IF NOT EXISTS customer;

CREATE TABLE IF NOT EXISTS customer.debts
(
    id         UUID PRIMARY KEY,
    document   VARCHAR(11)    NOT NULL,
    dueDate    DATE           NOT NULL,
    amount     DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP      NOT NULL
);
