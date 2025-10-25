# Infrastructure

## Architecture Overview

The application consists of three components:

1. **Frontend**: Static HTML/CSS/JS served directly by nginx
2. **Backend API**: Go HTTP server running as a systemd service
3. **Scraper**: Go CLI tool that runs daily via systemd timer

## Frontend

Served directly by nginx from `/srv/techblogs/current/frontend/`.

## Backend API

- **Service**: `techblogs-api.service`
- **Communication**: Unix Domain Socket at `/run/techblogs-api/techblogs-api.sock`
- **Reverse Proxy**: Nginx forwards requests to the socket
- **User**: Runs as `techblogs:www-data`
- **Database**: SQLite at `/srv/techblogs/data/techblogs.db`

## Scraper

- **Service**: `techblogs-scraper.service` (oneshot)
- **Timer**: `techblogs-scraper.timer` (runs daily)
- **Schedule**: 5 minutes after boot, then every 24 hours
- **Function**: Scrapes configured blogs and updates the cache table
- **User**: Runs as `techblogs:www-data`

## Deployment

Managed via GitHub Actions (`.github/workflows/deploy.yml`):

1. **Build**: Compiles Go binaries for Linux amd64 with CGO enabled
2. **Detect changes**: Uses path filters to detect what changed (frontend/backend/migrations/infra)
3. **Migrations**: Runs `golang-migrate` on VPS if migrations changed
4. **Deploy**: Creates timestamped releases in `/srv/techblogs/releases/YYYYMMDDHHMMSS/`
5. **Activate**: Symlinks `/srv/techblogs/current` to new release
6. **Health check**: Verifies backend is responding, rolls back on failure
7. **Cleanup**: Keeps last 10 releases

## Directory Structure on VPS

```
/srv/techblogs/
├── current -> releases/20241025123456/  # Symlink to active release
├── releases/
│   ├── 20241025123456/
│   │   ├── frontend/
│   │   └── backend/
│   └── 20241024091011/
├── migrations/                          # SQL migration files
└── data/
    └── techblogs.db                     # SQLite database
```

## One-time VPS Setup

1. Install `golang-migrate`:
   ```bash
   curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
   sudo mv migrate /usr/local/bin/
   sudo chmod +x /usr/local/bin/migrate
   ```

2. Create `techblogs` user:
   ```bash
   sudo useradd -r -s /bin/bash techblogs
   sudo usermod -a -G www-data techblogs
   ```

3. Create base directories:
   ```bash
   sudo mkdir -p /srv/techblogs/{releases,data,migrations}
   sudo chown -R techblogs:www-data /srv/techblogs
   ```
