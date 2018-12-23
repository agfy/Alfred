CREATE TABLE IF NOT EXISTS goods (
    id SERIAL PRIMARY KEY,
    name TEXT, 
    class TEXT, 
    shop TEXT, 
    volume TEXT, 
    price INTEGER,
    foodType TEXT
);