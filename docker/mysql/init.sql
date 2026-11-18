CREATE TABLE IF NOT EXISTS clients (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL,
    balance DOUBLE NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    CONSTRAINT fk_accounts_client FOREIGN KEY (client_id) REFERENCES clients(id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    account_id_from VARCHAR(255) NOT NULL,
    account_id_to VARCHAR(255) NOT NULL,
    amount DOUBLE NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT fk_transactions_account_from FOREIGN KEY (account_id_from) REFERENCES accounts(id),
    CONSTRAINT fk_transactions_account_to FOREIGN KEY (account_id_to) REFERENCES accounts(id)
);
