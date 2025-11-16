package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nesco/techblogs/backend/internal/blogs"
	"github.com/nesco/techblogs/backend/internal/database"
)

func main() {
	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/techblogs.db"
	}

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	repo := blogs.NewRepository(db)

	// Get all blog configurations
	configs, err := repo.GetAllBlogConfigs()
	if err != nil {
		log.Fatalf("Failed to get blog configs: %v", err)
	}

	log.Printf("Starting scraper for %d blogs...\n", len(configs))

	// Scrape each blog
	for _, config := range configs {
		log.Printf("Scraping %s (%s)...\n", config.BlogName, config.BlogHref)

		articleName, articleHref, err := scrapeBlog(config)
		if err != nil {
			log.Printf("Error scraping %s: %v\n", config.BlogName, err)
			continue
		}

		// Update cache
		blogInfo := blogs.BlogInfo{
			BlogName:          config.BlogName,
			BlogHref:          config.BlogHref,
			LatestArticleName: articleName,
			LatestArticleHref: articleHref,
			Kind:              config.Kind,
			GitHubHref:        config.GitHubHref,
		}

		if err := repo.UpsertBlogCache(blogInfo); err != nil {
			log.Printf("Error updating cache for %s: %v\n", config.BlogName, err)
			continue
		}

		log.Printf("Successfully scraped %s: %s\n", config.BlogName, articleName)
	}

	log.Println("Scraping complete!")
}

func scrapeBlog(config blogs.BlogConfig) (articleName string, articleHref string, err error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if config.ArticleHrefSelector == "" {
		return "", "", nil
	}

	req, err := http.NewRequest("GET", config.BlogHref, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set custom user agent
	req.Header.Set("User-Agent", "TechBlogs-Scraper/1.0 (+https://github.com/nesco/techblogs)")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch blog: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Find the article href using CSS selector
	hrefSelection := doc.Find(config.ArticleHrefSelector).First()
	if hrefSelection.Length() == 0 {
		return "", "", fmt.Errorf("no articles found with href selector: %s", config.ArticleHrefSelector)
	}

	// Extract article href from the link
	articleHref, exists := hrefSelection.Attr("href")
	if !exists {
		return "", "", fmt.Errorf("article href not found")
	}

	// Make href absolute if it's relative
	if !strings.HasPrefix(articleHref, "https://") && !strings.HasPrefix(articleHref, "http://") {
		articleHref = normalizeURL(config.BlogHref, articleHref)
	}

	// Find the article name using CSS selector
	nameSelection := doc.Find(config.ArticleNameSelector).First()
	if nameSelection.Length() == 0 {
		return "", "", fmt.Errorf("no articles found with name selector: %s", config.ArticleNameSelector)
	}

	// Extract article name by recursively finding the first leaf element with text
	articleName = nameSelection.Text() // extractTextFromFirstLeaf(nameSelection)

	// If no text found, try aria-label attribute
	if articleName == "" {
		articleName, _ = nameSelection.Attr("aria-label")
	}

	// Clean up excessive whitespace and newlines
	articleName = strings.Join(strings.Fields(articleName), " ")

	if articleName == "" {
		return "", "", fmt.Errorf("article name is empty")
	}

	return articleName, articleHref, nil
}

func normalizeURL(baseURL, path string) string {
	var scheme, host string
	if strings.HasPrefix(baseURL, "https://") {
		scheme = "https://"
		host = strings.TrimPrefix(baseURL, "https://")
	} else if strings.HasPrefix(baseURL, "http://") {
		scheme = "http://"
		host = strings.TrimPrefix(baseURL, "http://")
	} else {
		scheme = "https://"
		host = baseURL
	}

	// Remove trailing slash and path from host
	host = strings.TrimSuffix(strings.Split(host, "/")[0], "/")

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return scheme + host + path
}
