# Graceful Shutdown Evidence

## Date

2026-07-29

## Commands

```bash
docker compose up -d --build
docker compose ps
curl --fail http://localhost:8080/api/health
curl --fail http://localhost:8080/api/ready
docker compose stop db
curl -i http://localhost:8080/api/ready
docker compose start db
curl --fail http://localhost:8080/api/ready
docker compose stop -t 10 api
docker compose logs api
```

## Expected Result

- API starts in Docker Compose.
- `/api/health` returns HTTP 200.
- `/api/ready` returns HTTP 200 when PostgreSQL is reachable.
- `/api/ready` returns HTTP 503 when PostgreSQL is unavailable.
- API receives `SIGTERM` when Docker Compose stops the service.
- API stops accepting new requests.
- API closes the database pool.
- API exits before the 10 second shutdown timeout.

## Actual Result

Docker Compose stopped the API container cleanly:

```text
[+] Stopping 1/1
 ✔ Container multi-tenant-workflow-approval-api-1  Stopped  0.2s
```

The API started with environment configuration from Docker Compose:

```text
level=INFO msg="configuration loaded" app_env=local app_port=8080 shutdown_timeout=10s
level=INFO msg="server starting" address=:8080
```

The liveness endpoint returned HTTP 200:

```text
[GIN] 2026/07/28 - 17:22:24 | 200 | GET "/api/health"
```

When PostgreSQL was unavailable, readiness returned HTTP 503:

```text
level=ERROR msg="database readiness check failed" error="failed to connect to `user=postgres database=approval`: hostname resolving error: lookup db on 127.0.0.11:53: no such host"
[GIN] 2026/07/28 - 17:22:53 | 503 | GET "/api/ready"
```

After PostgreSQL was restored, readiness returned HTTP 200:

```text
[GIN] 2026/07/28 - 17:23:03 | 200 | GET "/api/ready"
[GIN] 2026/07/28 - 17:23:23 | 200 | GET "/api/ready"
```

When Docker Compose stopped the API service, the API handled the shutdown signal and closed the database pool:

```text
level=INFO msg="shutdown signal received"
level=INFO msg="closing database pool"
level=INFO msg="server shutdown completed"
```

## Notes

- The shutdown timeout was configured with `SHUTDOWN_TIMEOUT=10s`.
- The observed API container stop time was about `0.2s`, which is under the configured timeout.
- Docker Compose injects runtime environment variables into the container, so the warning about `.env` not being found inside the container is expected.
- Gin is currently running in debug mode. Set `GIN_MODE=release` in Docker Compose later for a more production-like runtime.
