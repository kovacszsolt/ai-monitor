ALTER TABLE providers ADD COLUMN credential TEXT NOT NULL DEFAULT '';
ALTER TABLE providers DROP COLUMN credential_ref;
