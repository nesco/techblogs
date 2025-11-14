-- Remove github_href column from blog_configs and blog_cache tables
ALTER TABLE blog_cache DROP COLUMN github_href;
ALTER TABLE blog_configs DROP COLUMN github_href;
