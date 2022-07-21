CREATE TABLE card (
    id UUID PRIMARY KEY,
    card_number VARCHAR(80) UNIQUE,
    balance BIGINT CHECK (balance > 0)
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL DEFAULT '',
    password VARCHAR(10) NOT NULL,
    email VARCHAR(80) NOT NULL UNIQUE,
    card_id uuid REFERENCES card (id),
    is_admin BOOLEAN DEFAULT FALSE
);

CREATE TABLE product (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE,
    quantity INTEGER,
    price INTEGER,
    original_price INTEGER
);

CREATE TABLE sales_history (
    customer_id UUID REFERENCES users(id),
    product_id UUID REFERENCES product (id),
    sold_quantity INTEGER,
    profit INTEGER,
    time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- DROP TABLE sales_history;
-- DROP TABLE product;
-- DROP TABLE users;
-- DROP TABLE card;

SELECT * FROM users;
-- SELECT * FROM card;
-- SELECT * FROM product;