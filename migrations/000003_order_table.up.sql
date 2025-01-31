CREATE TABLE orders(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    order_consignment_id VARCHAR(50) NOT NULL,
    store_id INT NOT NULL,
    merchant_order_id VARCHAR(255),
    recipient_name VARCHAR(255) NOT NULL,
    recipient_phone VARCHAR(15) NOT NULL,
    recipient_address TEXT NOT NULL,
    recipient_city INT NOT NULL,
    recipient_zone INT NOT NULL,
    recipient_area INT NOT NULL,
    delivery_type INT NOT NULL,
    item_type INT NOT NULL,
    special_instruction TEXT,
    item_quantity INT NOT NULL,
    item_weight NUMERIC(10, 2) NOT NULL,
    amount_to_collect NUMERIC(10, 2) NOT NULL,
    item_description TEXT,
    total_fee NUMERIC(10, 2) NOT NULL,
    order_type_id INT NOT NULL,
    cod_fee NUMERIC(10, 2),
    promo_discount NUMERIC(10, 2),
    discount NUMERIC(10, 2),
    delivery_fee NUMERIC(10, 2) NOT NULL,
    order_status VARCHAR(50) NOT NULL,
    archive Boolean DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION
);

CREATE INDEX idx_orders_pending
    ON orders (order_status)
    WHERE order_status = 'Pending' AND archive = FALSE;