# Techblogs Backend

Go backend for the techblogs application, featuring a REST API and web scraper.

## Components

- **API** (`cmd/api`) - HTTP server serving blog data
- **Scraper** (`cmd/scraper`) - Web scraper to fetch latest blog articles

## Database

Uses SQLite with migrations managed by [golang-migrate](https://github.com/golang-migrate/migrate).

### Running Migrations

```bash
# Install golang-migrate CLI
brew install golang-migrate  # macOS
# or
go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path ./migrations -database "sqlite3://./data/techblogs.db" up

# Rollback
migrate -path ./migrations -database "sqlite3://./data/techblogs.db" down 1
```

## Development

```bash
# Start API
go run cmd/api/main.go

# Run scraper manually
go run cmd/scraper/main.go
```

## Deployment

### Cron Job Setup

After deploying, set up a cron job on the server to run the scraper every 12 hours:

1. SSH into your server:
   ```bash
   ssh user@your-server.com
   ```

2. Edit crontab for the `techblogs` user:
   ```bash
   sudo -u techblogs crontab -e
   ```

3. Add the following line to run the scraper every 12 hours (at midnight and noon):
   ```cron
   0 0,12 * * * /srv/techblogs/current/backend/techblogs-scraper >> /var/log/techblogs-scraper.log 2>&1
   ```

4. Verify the cron job is installed:
   ```bash
   sudo -u techblogs crontab -l
   ```

### Environment Variables

- `DB_PATH` - Path to SQLite database (default: `./data/techblogs.db`)
- `LISTEN_ADDR` - Server listen address (default: `127.0.0.1:5011`)

### Log Monitoring

View scraper logs:
```bash
sudo tail -f /var/log/techblogs-scraper.log
```

View API logs:
```bash
sudo journalctl -u techblogs-api -f
```

## Adding New Blogs

1. Add blog configuration to `migrations/002_seed_blogs.up.sql`:
   ```sql
   INSERT INTO blog_configs (blog_name, blog_href, kind, article_selector) VALUES
   ('Blog Name', 'https://example.com/blog', 'organization', '.article-link a');
   ```

2. Find the correct CSS selector:
   - Inspect the blog's HTML
   - Find the first article link (`<a>` tag)
   - Use browser DevTools to copy the selector
   - Simplify the selector (remove unnecessary classes)

3. Drop and re-run migrations locally to test:
   ```bash
   rm ./data/techblogs.db
   migrate -path ./migrations -database "sqlite3://./data/techblogs.db" up
   go run cmd/scraper/main.go
   ```

4. Deploy and the scraper will pick up the new blog automatically.
