ALTER TABLE providers ADD COLUMN credential_ref TEXT NOT NULL DEFAULT '';
ALTER TABLE providers DROP COLUMN credential;
