# Cloud Mapping

## Purpose

This note maps the LU034 production-readiness slice to basic cloud compute, storage, and networking concepts.

The current application is a small Go API with `/api/health` and `/api/ready`. PostgreSQL is the first persistent dependency, and readiness depends on the API being able to reach PostgreSQL.

## Local Runtime

| Component | Local implementation | Current status |
|---|---|---|
| API compute | Go API process from `cmd/api` | Implemented locally |
| API container | Docker container for the Go API | Implemented with `Dockerfile` |
| Database | PostgreSQL container | Implemented in Docker Compose |
| Persistent storage | Docker named volume `postgres_data` | Implemented |
| Private network | Docker Compose service network between API and database | Implemented |
| Public endpoint | `localhost:8080` mapped to the API | Implemented |
| Configuration | Environment variables from `.env` or runtime environment | Implemented |
| Secrets | Local `.env`, excluded from Git | Implemented for local only |
| Logs | Standard output via `slog` and Gin logger | Implemented |

## Cloud Runtime Concept

| Component | Cloud category | Example mapping |
|---|---|---|
| API compute | Container service, VM, or app platform | Run the Go API container |
| Database | Managed PostgreSQL | Store workflow approval data |
| Persistent storage | Managed database storage | Preserve PostgreSQL data across restarts |
| Private network | VPC or private service network | Allow API-to-database traffic without public DB exposure |
| Public endpoint | Load balancer, API gateway, or reverse proxy | Expose HTTPS traffic to clients |
| Secrets | Secret manager or deployment environment secrets | Store `DATABASE_URL` and runtime config |
| Logs | Cloud logging service | Collect API startup, readiness, and shutdown logs |
| CI/CD runner | Hosted or self-hosted build runner | Build, test, and package deployable artifacts |

## Public vs Private Paths

| Path | Visibility | Reason |
|---|---|---|
| Client to API `/api/health` | Public | Used to verify that the API process is alive |
| Client to API `/api/ready` | Public or platform-only | Used to verify that the API can serve traffic |
| API to PostgreSQL | Private | Database credentials and traffic must not be exposed publicly |
| CI runner to container registry | Private/authenticated | Used to publish built images |
| Deployment platform to secret store | Private/authenticated | Used to inject runtime configuration |

## Current Network Shape

Docker Compose runs both the API and PostgreSQL:

```text
client
    |
    | localhost:8080/api/*
    v
API container
    |
    | db:5432 on Docker Compose network
    v
PostgreSQL container
    |
    v
postgres_data Docker volume
```

## Decisions

| Decision | Result |
|---|---|
| Keep the API stateless | Containers can be replaced without losing application data |
| Use PostgreSQL as the source of truth | Persistent state lives in the database, not the API process |
| Use a named Docker volume locally | Database data survives container restarts during local testing |
| Keep `DATABASE_URL` outside source code | Credentials are supplied by environment variables |
| Use `/api/health` for liveness | Confirms the process is running |
| Use `/api/ready` for readiness | Confirms the API can reach PostgreSQL |

## Assumptions

- The API is stateless.
- PostgreSQL is the source of truth.
- Containers can be replaced at any time.
- Database storage must survive restarts and deployments.
- The database should not be exposed publicly in a production cloud environment.
- Local Docker Compose may expose PostgreSQL on `localhost:5432` for developer convenience only.

## Evidence To Add Later

- CI evidence that builds the Docker image.
- CI or script evidence that runs the smoke test automatically.
- Production deployment evidence after this local LU034 slice is reused later in P3.

## Recorded Local Evidence

- `docker build -t approval-api:lu034-check .` completed successfully.
- `docker compose config` validated the Compose file.
- `docker compose up -d --build` started both API and PostgreSQL.
- `docker compose ps` showed PostgreSQL healthy and the API service running.
- `curl --fail http://localhost:8080/api/health` returned HTTP 200.
- `curl --fail http://localhost:8080/api/ready` returned HTTP 200 when PostgreSQL was available.
- `curl -i http://localhost:8080/api/ready` returned HTTP 503 when PostgreSQL was unavailable.
- `docker compose stop -t 10 api` stopped the API container cleanly in about `0.2s`.
