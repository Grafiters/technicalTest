CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    limit_id BIGINT NOT NULL,
    contract_no BIGINT,
    otr BIGINT,
    admin_fee BIGINT,
    installment BIGINT,
    asset_name VARCHAR(255),
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (limit_id) REFERENCES limits(id) ON DELETE CASCADE
);