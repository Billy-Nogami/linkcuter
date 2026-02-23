-- +goose Up
CREATE TABLE IF NOT EXISTS links (
    code VARCHAR(10) NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT links_code_key PRIMARY KEY (code),
    CONSTRAINT links_url_key UNIQUE (url)
);

-- +goose Down
DROP TABLE IF EXISTS links;
