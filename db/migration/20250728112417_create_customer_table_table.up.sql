CREATE TABLE IF NOT EXISTS customers (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255),
    nik INT,
    full_name VARCHAR(255),
    legal_name VARCHAR(255),
    birth_place VARCHAR(255),
    birth_date DATE,
    salary BIGINT,
    ktp_image_url TEXT,
    selfie_image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email ON customers(email);
CREATE INDEX idx_nik ON customers(nik);