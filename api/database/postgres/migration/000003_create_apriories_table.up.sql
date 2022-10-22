CREATE TABLE IF NOT EXISTS apriories (
    id_apriori SERIAL,
    code VARCHAR(10) NOT NULL,
    item VARCHAR(256) NOT NULL,
    discount DECIMAL(6,2) NOT NULL,
    support DECIMAL(6,2) NOT NULL,
    confidence DECIMAL(6,2) NOT NULL,
    range_date VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT,
    image TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (id_apriori)
)