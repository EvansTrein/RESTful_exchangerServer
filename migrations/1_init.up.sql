CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100), 
    email VARCHAR(100) UNIQUE, 
    password_hash VARCHAR(255)
);

CREATE TABLE currencies (
    code VARCHAR(5) PRIMARY KEY,
    name VARCHAR(50) NOT NULL    
);

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE, -- If a user is deleted, all his/her accounts will be deleted 
    currency_code VARCHAR(5) REFERENCES currencies(code) ON DELETE RESTRICT, -- you cannot delete a currency if an account is opened for it
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    UNIQUE (user_id, currency_code) -- one account for one currency per user 
);