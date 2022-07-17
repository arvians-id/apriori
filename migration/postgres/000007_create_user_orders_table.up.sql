CREATE TABLE IF NOT EXISTS user_orders (
    id_order SERIAL,
    payload_id INTEGER NOT NULL,
    code VARCHAR(256),
    name VARCHAR(256),
    price INTEGER,
    image VARCHAR(256),
    quantity INTEGER,
    total_price_item INTEGER,
    PRIMARY KEY (id_order),
    FOREIGN KEY (payload_id) REFERENCES payloads(id_payload) ON DELETE CASCADE ON UPDATE CASCADE
)