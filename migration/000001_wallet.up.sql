CREATE TABLE wallets (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance DECIMAL(20, 4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_uuid UUID REFERENCES wallets(uuid),
    transaction_type VARCHAR(10) NOT NULL,
    amount DECIMAL(20, 4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_uuid) REFERENCES wallets(uuid)
);