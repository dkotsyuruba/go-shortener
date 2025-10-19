CREATE TABLE IF NOT EXISTS links (
    id VARCHAR(255) PRIMARY KEY,
    original_url VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_links_original_url ON links(original_url);