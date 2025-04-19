CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(30) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY,
    city VARCHAR(30) NOT NULL,
    registration_date TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pvz_registration_date ON pvz(registration_date);

CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY,
    pvz_id UUID REFERENCES pvz(id) ON DELETE CASCADE,
    date_time TIMESTAMP NOT NULL,
    status VARCHAR(30) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_receptions_pvz_id ON receptions(pvz_id);
CREATE INDEX IF NOT EXISTS idx_receptions_date_time ON receptions(date_time);
CREATE INDEX IF NOT EXISTS idx_receptions_pvz_id_date_time ON receptions(pvz_id, date_time);
CREATE INDEX IF NOT EXISTS idx_receptions_open ON receptions(pvz_id, status, date_time DESC);


CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    reception_id UUID REFERENCES receptions(id) ON DELETE CASCADE,
    type VARCHAR(30) NOT NULL,
    date_time TIMESTAMP DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_products_reception_id ON products(reception_id);
CREATE INDEX IF NOT EXISTS idx_products_latest ON products(reception_id, date_time DESC);
