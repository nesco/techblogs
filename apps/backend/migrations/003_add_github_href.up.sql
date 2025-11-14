-- Add github_href column to blog_configs and blog_cache tables
ALTER TABLE blog_configs ADD COLUMN github_href TEXT NOT NULL DEFAULT '';
ALTER TABLE blog_cache ADD COLUMN github_href TEXT NOT NULL DEFAULT '';
