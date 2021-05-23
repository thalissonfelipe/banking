ALTER TABLE account
ALTER COLUMN id TYPE UUID USING id::UUID;

ALTER TABLE transfer
ALTER COLUMN id TYPE UUID USING id::UUID,
ALTER COLUMN account_origin_id TYPE UUID USING account_origin_id::UUID,
ALTER COLUMN account_destination_id TYPE UUID USING account_destination_id::UUID;
