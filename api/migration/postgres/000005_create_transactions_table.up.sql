CREATE TABLE IF NOT EXISTS transactions (
    id_transaction SERIAL,
    product_name VARCHAR(256) NOT NULL,
    customer_name VARCHAR(100) NOT NULL,
    no_transaction VARCHAR(100) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id_transaction)
)