-- Store issued proxy API keys as plaintext (legacy key_hash rows are invalid until keys are recreated).
ALTER TABLE proxy_api_keys RENAME COLUMN key_hash TO api_key;
