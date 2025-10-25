package blogs

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllBlogs() ([]BlogInfo, error) {
	query := `
		SELECT blog_name, blog_href, latest_article_name, latest_article_href, kind
		FROM blog_cache
		ORDER BY updated_at DESC, blog_name ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query blogs: %w", err)
	}
	defer rows.Close()

	var blogs []BlogInfo
	for rows.Next() {
		var blog BlogInfo
		var kind string
		if err := rows.Scan(&blog.BlogName, &blog.BlogHref, &blog.LatestArticleName, &blog.LatestArticleHref, &kind); err != nil {
			return nil, fmt.Errorf("failed to scan blog row: %w", err)
		}
		blog.Kind = Kind(kind)
		blogs = append(blogs, blog)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating blog rows: %w", err)
	}

	return blogs, nil
}

func (r *Repository) GetBlogsByKind(kind Kind) ([]BlogInfo, error) {
	query := `
		SELECT blog_name, blog_href, latest_article_name, latest_article_href, kind
		FROM blog_cache
		WHERE kind = ?
		ORDER BY updated_at DESC, blog_name ASC
	`
	rows, err := r.db.Query(query, string(kind))
	if err != nil {
		return nil, fmt.Errorf("failed to query blogs by kind: %w", err)
	}
	defer rows.Close()

	var blogs []BlogInfo
	for rows.Next() {
		var blog BlogInfo
		var kindStr string
		if err := rows.Scan(&blog.BlogName, &blog.BlogHref, &blog.LatestArticleName, &blog.LatestArticleHref, &kindStr); err != nil {
			return nil, fmt.Errorf("failed to scan blog row: %w", err)
		}
		blog.Kind = Kind(kindStr)
		blogs = append(blogs, blog)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating blog rows: %w", err)
	}

	return blogs, nil
}

func (r *Repository) GetAllBlogConfigs() ([]BlogConfig, error) {
	query := `
		SELECT blog_name, blog_href, kind, article_href_selector, article_name_selector
		FROM blog_configs
		ORDER BY blog_name
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query blog configs: %w", err)
	}
	defer rows.Close()

	var configs []BlogConfig
	for rows.Next() {
		var config BlogConfig
		var kind string
		if err := rows.Scan(&config.BlogName, &config.BlogHref, &kind, &config.ArticleHrefSelector, &config.ArticleNameSelector); err != nil {
			return nil, fmt.Errorf("failed to scan blog config row: %w", err)
		}
		config.Kind = Kind(kind)
		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating blog config rows: %w", err)
	}

	return configs, nil
}

func (r *Repository) GetBlogCache(blogName string) (*BlogInfo, error) {
	query := `
		SELECT blog_name, blog_href, latest_article_name, latest_article_href, kind
		FROM blog_cache
		WHERE blog_name = ?
	`
	var blog BlogInfo
	var kind string
	err := r.db.QueryRow(query, blogName).Scan(&blog.BlogName, &blog.BlogHref, &blog.LatestArticleName, &blog.LatestArticleHref, &kind)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get blog cache: %w", err)
	}
	blog.Kind = Kind(kind)
	return &blog, nil
}

func (r *Repository) UpsertBlogCache(blog BlogInfo) error {
	// Check if data has changed
	existing, err := r.GetBlogCache(blog.BlogName)
	if err != nil {
		return err
	}

	now := time.Now()

	// If exists and data is the same, don't update updated_at
	if existing != nil &&
		existing.BlogHref == blog.BlogHref &&
		existing.LatestArticleName == blog.LatestArticleName &&
		existing.LatestArticleHref == blog.LatestArticleHref &&
		existing.Kind == blog.Kind {
		// No changes, don't update
		return nil
	}

	query := `
		INSERT INTO blog_cache (blog_name, blog_href, latest_article_name, latest_article_href, kind, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(blog_name) DO UPDATE SET
			blog_href = excluded.blog_href,
			latest_article_name = excluded.latest_article_name,
			latest_article_href = excluded.latest_article_href,
			kind = excluded.kind,
			updated_at = excluded.updated_at
	`
	_, err = r.db.Exec(query, blog.BlogName, blog.BlogHref, blog.LatestArticleName, blog.LatestArticleHref, string(blog.Kind), now)
	if err != nil {
		return fmt.Errorf("failed to upsert blog cache: %w", err)
	}
	return nil
}
