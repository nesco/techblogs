CREATE TABLE IF NOT EXISTS blog_configs (
    id INTEGER PRIMARY KEY,
    blog_name TEXT NOT NULL UNIQUE,
    blog_href TEXT NOT NULL,
    kind TEXT NOT NULL CHECK (kind IN ('organization', 'individual')),
    article_href_selector TEXT NOT NULL,
    article_name_selector TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS blog_cache (
    id INTEGER PRIMARY KEY,
    blog_name TEXT NOT NULL UNIQUE,
    blog_href TEXT NOT NULL,
    latest_article_name TEXT,
    latest_article_href TEXT,
    kind TEXT NOT NULL CHECK (kind IN ('organization', 'individual')),
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (blog_name) REFERENCES blog_configs (
        blog_name
    ) ON DELETE CASCADE
);
