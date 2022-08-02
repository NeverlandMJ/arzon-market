CREATE TABLE users (
    id UUID PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL DEFAULT '',
    password VARCHAR(255) NOT NULL,
    email VARCHAR(80) NOT NULL UNIQUE,
    is_admin BOOLEAN DEFAULT FALSE
);

CREATE TABLE card (
    id UUID PRIMARY KEY,
    card_number VARCHAR(80) UNIQUE,
    balance BIGINT CHECK (balance > 0),
    owner UUID REFERENCES users (id)
);


CREATE TABLE product (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE,
    description TEXT,
    quantity INTEGER, 
    price INTEGER,
    original_price INTEGER,
    img VARCHAR(255)
);

CREATE TABLE sales_history (
    customer_id UUID REFERENCES users(id),
    product_id UUID REFERENCES product (id),
    sold_quantity INTEGER,
    profit INTEGER,
    time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

