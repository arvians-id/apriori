CREATE TABLE IF NOT EXISTS comments (
    id_comment SERIAL,
    user_order_id INTEGER NOT NULL,
    product_code VARCHAR(10) NOT NULL,
    description TEXT,
    tag VARCHAR(200),
    rating INTEGER NOT NULL,
    created_at TIMESTAMP,
    PRIMARY KEY (id_comment),
    FOREIGN KEY (user_order_id) REFERENCES user_orders(id_order) ON DELETE CASCADE ON UPDATE CASCADE
)