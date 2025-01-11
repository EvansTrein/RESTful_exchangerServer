-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100), 
    email VARCHAR(100) UNIQUE, 
    password_hash VARCHAR(255)
);

-- Таблица валют
CREATE TABLE currencies (
    code VARCHAR(5) PRIMARY KEY, 
    name VARCHAR(50) NOT NULL    
);

-- Таблица счетов
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    currency_code VARCHAR(5) REFERENCES currencies(code) ON DELETE RESTRICT,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    UNIQUE (user_id, currency_code)
);

INSERT INTO currencies (code, name)
VALUES
    ('USD', 'US Dollar'),
    ('EUR', 'Euro'),
    ('RUB', 'Russian Ruble'),
    ('CNY', 'Chinese Yuan');


-- Пополнение счета
UPDATE accounts
SET balance = balance + 100.00
WHERE user_id = $
  AND currency_code = 'USD'


-- Получение счета пользователя по ID
SELECT 
    u.id AS user_id,
    u.name AS name,
    u.email AS email,
    a.currency_code AS currency,
    a.balance
FROM users u
JOIN accounts a ON u.id = a.user_id
WHERE u.id = $;