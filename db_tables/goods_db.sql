CREATE TABLE IF NOT EXISTS goods (
    id SERIAL PRIMARY KEY,
    name TEXT, 
    class TEXT, 
    shop TEXT, 
    volume INTEGER, 
    price INTEGER,
    foodType TEXT
);