CREATE TABLE IF NOT EXISTS comments (
    id_comment SERIAL,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    description TEXT,
    tag VARCHAR(200),
    rating INTEGER NOT NULL,
    created_at TIMESTAMP,
    PRIMARY KEY (id_comment),
    FOREIGN KEY (user_id) REFERENCES users(id_user) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id_product) ON DELETE RESTRICT ON UPDATE CASCADE
)