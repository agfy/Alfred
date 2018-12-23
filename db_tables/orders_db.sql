CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    telegram_id INTEGER,
    goods_id INTEGER, 
    amount INTEGER, 
    create_time TIMESTAMP WITHOUT TIME ZONE
);