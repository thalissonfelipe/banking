BEGIN;

CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY,
    account_origin_id UUID REFERENCES accounts,
    account_destination_id UUID REFERENCES accounts,
    amount INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
