CREATE TABLE IF NOT EXISTS links (
    code VARCHAR(10) PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE links
    DROP CONSTRAINT IF EXISTS links_code_key,
    ADD CONSTRAINT links_code_key UNIQUE (code);

ALTER TABLE links
    DROP CONSTRAINT IF EXISTS links_url_key,
    ADD CONSTRAINT links_url_key UNIQUE (url);
