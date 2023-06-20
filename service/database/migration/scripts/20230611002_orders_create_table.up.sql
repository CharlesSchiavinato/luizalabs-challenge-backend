CREATE TABLE orders (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "buy_date" date NOT NULL,
    "total" real NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
	        REFERENCES users(id)
);
