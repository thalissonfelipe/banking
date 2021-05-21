BEGIN;

CREATE TABLE IF NOT EXISTS "transfer" (
    "id"  VARCHAR(50) PRIMARY KEY,
    "account_origin_id" VARCHAR(50) NOT NULL,
    "account_destination_id" VARCHAR(50) NOT NULL,
    "amount" INT NOT NULL,
    "created_at" DATE
);

COMMIT;