# LU034 - Production Readiness First Guide

> Parent project: **P3 Multi-tenant Workflow Approval System**  
> Roadmap slot: **Month 12 / Week 1**  
> Focus: **Cloud compute, storage, networking, CI/CD, Linux processes and signals**  
> Related lab: **MP25 Deployment Pipeline Lab**  
> Goal: Do LU034 first as a standalone production-readiness slice, then reuse it in P3.

---

# 1. Why Do LU034 First?

Normally LU034 appears near the end of P3 because deployment work is easiest after the app exists.

But you can do LU034 first if you treat it as an infrastructure and operations foundation:

```text
Small runnable API
    v
Docker image
    v
Docker Compose with PostgreSQL
    v
Health and readiness checks
    v
Graceful shutdown
    v
CI quality gates
    v
Cloud mapping and evidence
```

This gives you a working production shell. Later, the full P3 workflow logic can be added inside the same deployable structure.

---

# 2. LU034 Learning Objectives

By the end of LU034, you should be able to explain and demonstrate:

1. How the API maps to cloud compute
2. How PostgreSQL maps to persistent storage
3. Which network paths are public and private
4. How Docker packages the application
5. How CI/CD proves code is safe enough to deploy
6. How `/api/health` differs from `/api/ready`
7. What happens when the API receives `SIGTERM`
8. Why graceful shutdown matters during deployment

---

# 3. LU034 Topic Coverage

The roadmap lists **13 topics** for LU034. This guide covers all of them, but with different depth:

| # | Topic | Role | How this guide covers it | Evidence |
|---:|---|---|---|---|
| 1 | Cloud compute, storage and networking | Core | Map API, PostgreSQL, volume, public endpoint and private network | `docs/deployment/cloud-mapping.md` |
| 2 | CI/CD | Core | Build CI pipeline with format, vet, test, binary build and Docker build gates | `.github/workflows/ci.yml`, `docs/deployment/ci-cd-pipeline.md` |
| 3 | Linux processes and signals | Core | Inspect API process and handle `SIGTERM` graceful shutdown | `docs/evidence/graceful-shutdown.md` |
| 4 | API test | Supporting | Smoke test `/api/health` and `/api/ready`; later extend to business endpoint | `scripts/smoke-test.sh` |
| 5 | Code coverage vs test quality | Supporting | Run coverage, but judge tests by critical behavior and failure cases | CI test output or coverage note |
| 6 | Contract test | Supporting | Verify stable response shape for `/api/health`, `/api/ready` and error responses | API/contract test file or checklist |
| 7 | Database testing | Supporting | Verify readiness against real PostgreSQL in Docker/CI | CI PostgreSQL service and readiness evidence |
| 8 | Deterministic tests | Supporting | Avoid sleeps and external services; use predictable health/readiness checks | Test notes in CI/CD document |
| 9 | Flaky tests | Supporting | Identify timing/order/network risks in tests | Flaky-test notes in CI/CD document |
| 10 | Migration testing | Supporting | Add placeholder gate now; run real migration checks when P3 schema exists | CI/CD pipeline note |
| 11 | Mock, stub and fake | Supporting | Use real PostgreSQL for readiness; use fake dependencies only for unavailable externals | Testing note |
| 12 | Performance test | Supporting | Keep a basic response-time smoke check; defer load testing to later reliability units | Smoke-test note |
| 13 | Regression test | Supporting | Preserve `/api/health`, `/api/ready`, and database-down readiness behavior | Regression checklist |

For LU034-first, the supporting testing topics are intentionally lightweight. The goal is to connect each testing concept to the deployable shell without prematurely building the full P3 workflow.

---

# 4. Recommended Scope

For LU034-first, build only this minimal system:

- Go API
- `GET /api/health`
- `GET /api/ready`
- PostgreSQL connection
- Dockerfile
- Docker Compose
- Graceful shutdown
- GitHub Actions CI
- Basic smoke-test script
- Cloud mapping note
- CI/CD pipeline note
- Graceful shutdown evidence

Do not build these yet unless you already want to continue into full P3:

- Multi-tenancy
- RBAC
- Approval workflow
- Audit trail
- Notifications
- Feature flags
- Rollback automation
- Staging and production deployment

Those belong to later P3 phases or LU035-LU036.

---

# 5. Target Repository Structure

Use this structure for the LU034-first slice:

```text
multi-tenant-workflow-approval/
|-- cmd/
|   `-- api/
|       `-- main.go
|-- internal/
|   `-- platform/
|       |-- config.go
|       |-- database.go
|       `-- server.go
|-- docs/
|   |-- deployment/
|   |   |-- cloud-mapping.md
|   |   `-- ci-cd-pipeline.md
|   `-- evidence/
|       `-- graceful-shutdown.md
|-- scripts/
|   `-- smoke-test.sh
|-- .github/
|   `-- workflows/
|       `-- ci.yml
|-- Dockerfile
|-- docker-compose.yml
|-- .env.example
|-- Makefile
|-- go.mod
`-- README.md
```

Keep it small. LU034 is about operational shape, not business features.

---

# 6. Step-by-Step Implementation

## Step 1 - Create a runnable API

Create a Go API with two endpoints:

```text
GET /api/health
GET /api/ready
```

Expected behavior:

| Endpoint | Meaning | Expected result |
|---|---|---|
| `/api/health` | Process is alive | `200` if API process is running |
| `/api/ready` | App can serve real traffic | `200` if database is reachable, `503` if not |

Minimum response:

```json
{
  "status": "ok"
}
```

For readiness failure:

```json
{
  "status": "not_ready",
  "reason": "database unavailable"
}
```

## Step 2 - Add configuration

Read configuration from environment variables:

```text
APP_ENV=local
APP_PORT=8080
DATABASE_URL=postgres://postgres:password@localhost:5432/approval?sslmode=disable
SHUTDOWN_TIMEOUT=10s
```

Rules:

- Do not hardcode database credentials in Go code.
- Commit `.env.example`, not real `.env`.
- Log config names and safe values only.
- Never log `DATABASE_URL` if it contains credentials.

## Step 3 - Add PostgreSQL with Docker Compose

Use PostgreSQL as the first persistent dependency.

```yaml
services:
  db:
    image: postgres:17-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: approval
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d approval"]
      interval: 5s
      timeout: 3s
      retries: 10

volumes:
  postgres_data:
```

Evidence to capture:

```bash
docker compose up -d db
docker compose ps
```

## Step 4 - Dockerize the API

Use a multi-stage Dockerfile.

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:latest

WORKDIR /app

RUN adduser -D appuser
USER appuser

COPY --from=builder /app/api ./api

EXPOSE 8080

CMD ["./api"]
```

Evidence to capture:

```bash
docker build -t approval-api:lu034 .
docker images | grep approval-api
```

## Step 5 - Run API and database together

Add the API service to `docker-compose.yml`.

```yaml
services:
  api:
    build:
      context: .
    ports:
      - "8080:8080"
    environment:
      APP_ENV: local
      APP_PORT: "8080"
      DATABASE_URL: postgres://postgres:password@db:5432/approval?sslmode=disable
      SHUTDOWN_TIMEOUT: 10s
    depends_on:
      db:
        condition: service_healthy
```

Evidence to capture:

```bash
docker compose up -d --build
curl --fail http://localhost:8080/api/health
curl --fail http://localhost:8080/api/ready
```

Then stop the database and prove readiness fails:

```bash
docker compose stop db
curl -i http://localhost:8080/api/ready
```

Expected:

```text
HTTP 503
```

## Step 6 - Add graceful shutdown

The API should handle:

```text
SIGINT
SIGTERM
```

Shutdown order:

```text
Receive SIGTERM
    v
Log shutdown start
    v
Stop accepting new requests
    v
Wait for active requests
    v
Close database pool
    v
Exit before timeout
```

Evidence commands:

```bash
docker compose up -d --build
docker compose top api
docker compose stop -t 10 api
docker compose logs api
```

Save the result in:

```text
docs/evidence/graceful-shutdown.md
```

## Step 7 - Add CI

Create:

```text
.github/workflows/ci.yml
```

Minimum CI pipeline:

```yaml
name: CI

on:
  push:
  pull_request:

jobs:
  test-and-build:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:17-alpine
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: password
          POSTGRES_DB: approval_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd="pg_isready -U postgres -d approval_test"
          --health-interval=5s
          --health-timeout=3s
          --health-retries=10

    env:
      APP_ENV: test
      APP_PORT: 8080
      DATABASE_URL: postgres://postgres:password@localhost:5432/approval_test?sslmode=disable

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.24"

      - name: Download dependencies
        run: go mod download

      - name: Format check
        run: test -z "$(gofmt -l .)"

      - name: Vet
        run: go vet ./...

      - name: Test
        run: go test ./... -race -cover

      - name: Build binary
        run: go build -o bin/api ./cmd/api

      - name: Build Docker image
        run: docker build -t approval-api:${{ github.sha }} .
```

Write what each step proves in:

```text
docs/deployment/ci-cd-pipeline.md
```

Also record how the CI pipeline touches the LU034 supporting testing topics:

| Testing topic | Minimum LU034 implementation |
|---|---|
| API test | Test `/api/health` and `/api/ready` through HTTP |
| Code coverage vs test quality | Save coverage output, then explain which behavior is actually protected |
| Contract test | Assert response status codes and JSON fields |
| Database testing | Run readiness against a real PostgreSQL service |
| Deterministic tests | Avoid `time.Sleep`; use health checks and bounded timeouts |
| Flaky tests | Note possible timing risks from database startup and container networking |
| Migration testing | Add the CI stage placeholder; activate when migrations exist |
| Mock, stub and fake | Prefer real PostgreSQL here; reserve fakes for external email/cloud APIs |
| Performance test | Add a simple timing observation for `/api/health` and `/api/ready` |
| Regression test | Keep tests for health/readiness behavior so deployment changes do not break them |

## Step 8 - Add smoke test script

Create:

```text
scripts/smoke-test.sh
```

Minimum checks:

```bash
#!/usr/bin/env sh
set -eu

BASE_URL="${BASE_URL:-http://localhost:8080}"

curl --fail "$BASE_URL/api/health"
curl --fail "$BASE_URL/api/ready"
```

Run:

```bash
chmod +x scripts/smoke-test.sh
./scripts/smoke-test.sh
```

---

# 7. Cloud Mapping Note

Create:

```text
docs/deployment/cloud-mapping.md
```

Use this template:

```markdown
# Cloud Mapping

## Local Runtime

| Component | Local implementation |
|---|---|
| API compute | Docker container |
| Database | PostgreSQL container |
| Persistent storage | Docker named volume |
| Private network | Docker Compose network |
| Public endpoint | localhost port mapping |

## Cloud Runtime Concept

| Component | Cloud category |
|---|---|
| API compute | Container service, VM, or app platform |
| Database | Managed PostgreSQL |
| Persistent storage | Managed database storage |
| Private network | VPC or private service network |
| Public endpoint | Load balancer, API gateway, or reverse proxy |
| Secrets | Secret manager or deployment environment secrets |
| Logs | Cloud logging service |

## Public vs Private Paths

- Public: browser/client to API endpoint
- Private: API to PostgreSQL
- Private: deployment runner to deployment platform

## Assumptions

- The API is stateless.
- PostgreSQL is the source of truth.
- Containers can be replaced at any time.
- Database storage must survive restarts and deployments.
```

---

# 8. CI/CD Pipeline Note

Create:

```text
docs/deployment/ci-cd-pipeline.md
```

Use this template:

```markdown
# CI/CD Pipeline

## Pipeline Flow

1. Checkout code
2. Setup Go
3. Download dependencies
4. Check formatting
5. Run static analysis
6. Run tests
7. Build binary
8. Build Docker image
9. Run smoke test before deployment

## Quality Gates

| Gate | What it catches |
|---|---|
| Format check | Inconsistent Go formatting |
| Vet | Suspicious Go code |
| Unit tests | Broken local behavior |
| Race tests | Unsafe concurrent code |
| Build binary | Compilation or packaging issue |
| Docker build | Broken container artifact |
| Smoke test | App cannot start or serve basic traffic |

## LU034 Testing Topic Notes

| Topic | Decision |
|---|---|
| API test | Test HTTP behavior, not only Go functions |
| Code coverage vs test quality | Coverage is evidence, but critical behavior matters more |
| Contract test | Response fields and status codes should stay stable |
| Database testing | Readiness must use a real PostgreSQL connection |
| Deterministic tests | Tests should not depend on real time or test order |
| Flaky tests | Container startup and network timing need bounded retries |
| Migration testing | Enable real migration checks after the first schema migration exists |
| Mock, stub and fake | Use fakes only for dependencies that are out of scope for LU034 |
| Performance test | Record a simple baseline before optimizing |
| Regression test | Keep permanent tests for health/readiness behavior |

## Deployment Rule

Only deploy an image that passed CI and is tagged with the commit SHA.
```

---

# 9. Graceful Shutdown Evidence Template

Create:

```text
docs/evidence/graceful-shutdown.md
```

Use this template:

````markdown
# Graceful Shutdown Evidence

## Date

YYYY-MM-DD

## Commands

```bash
docker compose up -d --build
docker compose top api
docker compose stop -t 10 api
docker compose logs api
```

## Expected Result

- API receives `SIGTERM`
- API logs shutdown start
- API stops accepting new requests
- API closes database pool
- API exits before 10 seconds

## Actual Result

Paste summarized command output here.

## Notes

Any issue found and how it was fixed.
````

---

# 10. Definition of Done

LU034 is done when all of these are true:

- [ ] API runs locally
- [ ] `/api/health` works
- [ ] `/api/ready` checks PostgreSQL
- [ ] Docker image builds
- [ ] Docker Compose starts API and PostgreSQL
- [ ] Data is stored in a persistent named volume
- [ ] Readiness fails when PostgreSQL is unavailable
- [ ] API handles `SIGTERM`
- [ ] Container exits cleanly within timeout
- [ ] CI runs format, vet, tests, binary build, and Docker build
- [ ] `scripts/smoke-test.sh` works
- [ ] `docs/deployment/cloud-mapping.md` exists
- [ ] `docs/deployment/ci-cd-pipeline.md` exists
- [ ] `docs/evidence/graceful-shutdown.md` exists
- [ ] All 13 LU034 topics are mapped to evidence or a deliberate lightweight note
- [ ] Supporting testing topics are connected to the deployable shell

---

# 11. What to Bring Back Into P3 Later

When you continue full P3, keep this LU034 foundation and add:

- Tenant migrations
- User and auth migrations
- RBAC tables
- Approval workflow endpoints
- Audit logs
- Notification worker
- Full critical-path smoke test
- Rollback runbook
- Feature flag validation
- Monitoring plan

The key idea:

```text
LU034 creates the deployable shell.
P3 fills that shell with the multi-tenant workflow product.
```
