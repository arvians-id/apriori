CREATE TABLE IF NOT EXISTS categories (
    id_category SERIAL,
    name VARCHAR(20) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id_category)
)