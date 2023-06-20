CREATE TABLE orders_product (
    "id" bigserial PRIMARY KEY,
    "order_id" bigint NOT NULL,
    "product_id" bigint NOT NULL,
    "product_value" real NOT NULL,
    CONSTRAINT fk_order
        FOREIGN KEY(order_id) 
	        REFERENCES orders(id)
);

CREATE INDEX "idx_order_id_product_id" ON orders_product (order_id, product_id);
