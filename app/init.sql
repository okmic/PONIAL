-- init.sql
-- Здесь можно добавить начальные данные или создать дополнительные таблицы
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Пример создания дополнительной таблицы
-- CREATE TABLE IF NOT EXISTS users (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     email VARCHAR(255) UNIQUE NOT NULL,
--     created_at TIMESTAMP DEFAULT NOW()
-- );