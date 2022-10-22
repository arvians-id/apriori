CREATE TABLE IF NOT EXISTS products (
    id_product SERIAL,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price INTEGER NOT NULL,
    category VARCHAR(100) NOT NULL,
    is_empty BOOLEAN NOT NULL DEFAULT FALSE,
    mass INTEGER NOT NULL,
    image TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id_product)
)