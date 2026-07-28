# P3 — Multi-tenant Workflow Approval System

> Project type: **Core Full-stack Project**  
> Duration: **Months 10–12**  
> Main focus: **Multi-tenancy, RBAC, Workflow, Audit Trail, Notifications, Testing, Deployment**  
> Final deliverables: **Deployed app, C4 diagrams, ERD, ADRs, Monitoring Plan**  
> Related labs: **MP23 Testing Strategy Lab**, **MP25 Deployment Pipeline Lab**  
> Roadmap learning units covered: **LU028-LU036**. For Month 12, use **LU034-LU036** as the production-readiness checkpoint.
> If you want to do LU034 first, use **LU034_Production_Readiness_Guide.md** as the standalone guide.

---

# 1. Project Objective

Build a production-style workflow approval system that supports multiple tenants.

The project must demonstrate that you can:

1. Design a multi-tenant application
2. Implement authentication and RBAC
3. Implement approval workflow rules
4. Protect tenant data
5. Record audit events
6. Process notifications
7. Test the critical system flow
8. Package the application with Docker
9. Build a CI/CD pipeline
10. Run database migrations
11. Deploy and roll back one version
12. Explain the architecture through C4 diagrams, ERD, ADRs, and a monitoring plan

---

# 2. Recommended Technology

This guide assumes:

```text
Backend          Go
HTTP framework   Gin
Database         PostgreSQL
Database driver  pgx
Migration tool   golang-migrate
Testing          Go testing + testify
Containers       Docker + Docker Compose
CI               GitHub Actions
API style        REST
Authentication   JWT
Documentation    Markdown + Mermaid
```

The structure can be adapted to another stack, but keep the same architectural concepts.

---

# 3. Scope

## 3.1 Required features

The first complete version should contain:

- Tenant management
- User management
- Authentication
- Role-based access control
- Approval request creation
- Approval submission
- Approval or rejection
- Approval history
- Audit trail
- Notification records
- Health and readiness endpoints
- Database migrations
- Automated tests
- Docker deployment
- CI pipeline
- Rollback runbook
- Monitoring plan

## 3.2 Roles

Use four roles:

```text
ADMIN
REQUESTER
APPROVER
AUDITOR
```

## 3.3 Core permissions

| Permission | Admin | Requester | Approver | Auditor |
|---|---:|---:|---:|---:|
| Manage users | Yes | No | No | No |
| Create approval request | Yes | Yes | No | No |
| Submit own request | Yes | Yes | No | No |
| View own request | Yes | Yes | No | No |
| View assigned requests | Yes | No | Yes | No |
| Approve request | Yes | No | Yes | No |
| Reject request | Yes | No | Yes | No |
| View audit trail | Yes | No | No | Yes |

## 3.4 Out of scope for the first version

Do not add these until the required scope is complete:

- Kubernetes
- Microservices
- Multiple databases
- GraphQL
- Complex workflow builder UI
- Real-time WebSocket updates
- External identity provider
- Multi-region deployment
- Advanced analytics
- Mobile application

---

# 4. Critical User Flow

Use one flow as the centre of the project:

```text
Requester logs in
    ↓
Requester creates a draft request
    ↓
Requester submits the request
    ↓
Request becomes PENDING
    ↓
Approver reads the request
    ↓
Approver approves the request
    ↓
Request becomes APPROVED
    ↓
Approval action is stored
    ↓
Audit event is stored
    ↓
Notification is created
```

This flow must be covered by:

- Unit tests
- Integration tests
- API tests
- Contract tests
- E2E tests
- Migration tests
- Regression tests
- Deployment smoke tests

---

# 5. State Model

Use these statuses:

```text
DRAFT
PENDING
APPROVED
REJECTED
CANCELLED
```

Allowed transitions:

```text
DRAFT → PENDING
DRAFT → CANCELLED
PENDING → APPROVED
PENDING → REJECTED
```

Disallowed transitions:

```text
APPROVED → PENDING
REJECTED → APPROVED
CANCELLED → PENDING
PENDING → DRAFT
```

Example domain logic:

```go
func (a *ApprovalRequest) Submit() error {
	if a.Status != StatusDraft {
		return ErrInvalidTransition
	}

	a.Status = StatusPending
	return nil
}

func (a *ApprovalRequest) Approve() error {
	if a.Status != StatusPending {
		return ErrInvalidTransition
	}

	a.Status = StatusApproved
	return nil
}
```

Keep workflow rules in the domain or service layer, not directly in HTTP handlers.

---

# 6. Architecture Overview

## 6.1 System context

```mermaid
flowchart LR
    Requester[Requester]
    Approver[Approver]
    Admin[Admin]
    Auditor[Auditor]
    System[Workflow Approval System]
    Email[Email Provider]

    Requester --> System
    Approver --> System
    Admin --> System
    Auditor --> System
    System --> Email
```

## 6.2 Container diagram

```mermaid
flowchart LR
    Browser[Web Client]
    API[Go REST API]
    DB[(PostgreSQL)]
    Worker[Notification Worker]
    Email[Email Provider]

    Browser --> API
    API --> DB
    API --> Worker
    Worker --> DB
    Worker --> Email
```

For the first version, the notification worker may run inside the same application process.

## 6.3 Component diagram

```mermaid
flowchart TB
    Handler[HTTP Handlers]
    Auth[Authentication Middleware]
    RBAC[Authorization Policy]
    ApprovalService[Approval Service]
    AuditService[Audit Service]
    NotificationService[Notification Service]
    ApprovalRepo[Approval Repository]
    AuditRepo[Audit Repository]
    NotificationRepo[Notification Repository]
    DB[(PostgreSQL)]

    Handler --> Auth
    Auth --> RBAC
    RBAC --> ApprovalService
    ApprovalService --> ApprovalRepo
    ApprovalService --> AuditService
    ApprovalService --> NotificationService
    AuditService --> AuditRepo
    NotificationService --> NotificationRepo
    ApprovalRepo --> DB
    AuditRepo --> DB
    NotificationRepo --> DB
```

---

# 7. Repository Structure

```text
multi-tenant-workflow-approval/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── jwt.go
│   │   ├── middleware.go
│   │   └── password.go
│   ├── authorization/
│   │   ├── permissions.go
│   │   └── policy.go
│   ├── tenant/
│   ├── user/
│   ├── approval/
│   │   ├── domain.go
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── postgres_repository.go
│   ├── audit/
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── postgres_repository.go
│   ├── notification/
│   │   ├── service.go
│   │   ├── sender.go
│   │   ├── repository.go
│   │   └── worker.go
│   ├── platform/
│   │   ├── database.go
│   │   ├── config.go
│   │   ├── logger.go
│   │   └── server.go
│   └── testutil/
├── migrations/
├── tests/
│   ├── integration/
│   ├── api/
│   ├── contract/
│   ├── e2e/
│   ├── migration/
│   └── performance/
├── docs/
│   ├── architecture/
│   ├── adr/
│   ├── testing/
│   ├── deployment/
│   └── monitoring/
├── scripts/
│   ├── migrate-up.sh
│   ├── migrate-down.sh
│   ├── seed.sh
│   ├── smoke-test.sh
│   ├── deploy.sh
│   └── rollback.sh
├── deploy/
├── .github/
│   └── workflows/
│       └── ci.yml
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

Do not create all packages on the first day. Add each package only when the corresponding feature is implemented.

---

# 8. Database Design

## 8.1 Main tables

```text
tenants
users
roles
permissions
user_roles
role_permissions
approval_requests
approval_actions
audit_logs
notifications
```

## 8.2 ERD

```mermaid
erDiagram
    TENANTS ||--o{ USERS : contains
    TENANTS ||--o{ APPROVAL_REQUESTS : owns
    TENANTS ||--o{ AUDIT_LOGS : owns
    TENANTS ||--o{ NOTIFICATIONS : owns

    USERS ||--o{ USER_ROLES : has
    ROLES ||--o{ USER_ROLES : assigned
    ROLES ||--o{ ROLE_PERMISSIONS : grants
    PERMISSIONS ||--o{ ROLE_PERMISSIONS : contains

    USERS ||--o{ APPROVAL_REQUESTS : creates
    APPROVAL_REQUESTS ||--o{ APPROVAL_ACTIONS : records
    USERS ||--o{ APPROVAL_ACTIONS : performs

    USERS ||--o{ AUDIT_LOGS : acts
    APPROVAL_REQUESTS ||--o{ NOTIFICATIONS : triggers
```

## 8.3 Tenant rule

Every tenant-owned business table must include:

```sql
tenant_id UUID NOT NULL
```

Every tenant-owned query must filter by `tenant_id`.

Correct:

```sql
SELECT *
FROM approval_requests
WHERE tenant_id = $1
  AND id = $2;
```

Incorrect:

```sql
SELECT *
FROM approval_requests
WHERE id = $1;
```

---

# 9. Implementation Order

Follow this order. Do not start deployment before the critical application flow works.

---

# Phase 0 — Prepare the Project

## Goal

Create the repository and a runnable empty API.

## Tasks

- [ ] Create Git repository
- [ ] Initialise Go module
- [ ] Add Gin
- [ ] Add pgx
- [ ] Add configuration loader
- [ ] Add structured logger
- [ ] Add `/health`
- [ ] Add `/ready`
- [ ] Add `.env.example`
- [ ] Add initial README

## Commands

```bash
mkdir multi-tenant-workflow-approval
cd multi-tenant-workflow-approval

git init

go mod init github.com/<username>/multi-tenant-workflow-approval

go get github.com/gin-gonic/gin
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/google/uuid
go get github.com/stretchr/testify
```

## Evidence

```text
GET /health → 200
GET /ready → 503 before database exists
```

## Definition of Done

- API starts
- Configuration is loaded from environment variables
- Logs show application start
- Health endpoint works

---

# Phase 1 — Database and Migrations

## Goal

Create a repeatable PostgreSQL schema.

## Tasks

- [ ] Add PostgreSQL through Docker Compose
- [ ] Add migration tool
- [ ] Create tenant migration
- [ ] Create user migration
- [ ] Create RBAC migrations
- [ ] Create approval request migration
- [ ] Create approval action migration
- [ ] Create audit log migration
- [ ] Create notification migration
- [ ] Add seed data
- [ ] Test migration up
- [ ] Test migration down

## Migration order

```text
001_create_tenants
002_create_users
003_create_roles
004_create_permissions
005_create_user_roles
006_create_role_permissions
007_create_approval_requests
008_create_approval_actions
009_create_audit_logs
010_create_notifications
```

## Definition of Done

- Fresh database can reach the latest schema
- One migration can be rolled back
- Seed creates one tenant, admin, requester, approver, and auditor
- `/ready` returns `200`

---

# Phase 2 — Authentication

## Goal

Allow users to log in and receive an authenticated identity.

## Tasks

- [ ] Store password hash
- [ ] Create login endpoint
- [ ] Verify password
- [ ] Generate JWT
- [ ] Add authentication middleware
- [ ] Put user ID and tenant ID in request context
- [ ] Reject invalid or expired tokens
- [ ] Add authentication tests

## Endpoint

```http
POST /api/auth/login
```

Request:

```json
{
  "email": "requester@tenant-a.test",
  "password": "password"
}
```

Response:

```json
{
  "access_token": "<jwt>",
  "token_type": "Bearer"
}
```

## JWT claims

```json
{
  "sub": "user-id",
  "tenant_id": "tenant-id",
  "roles": ["REQUESTER"],
  "exp": 1780000000
}
```

## Security rules

- Do not store plain-text passwords
- Do not trust tenant ID from request body
- Tenant identity must come from the authenticated token or server-side context
- Do not log passwords or JWT tokens

## Definition of Done

- Valid user can log in
- Invalid password is rejected
- Protected endpoint requires JWT
- Tenant ID is available in request context

---

# Phase 3 — RBAC

## Goal

Protect actions with centralised permission checks.

## Tasks

- [ ] Define permission constants
- [ ] Define role-permission mapping
- [ ] Create authorisation policy
- [ ] Add middleware or service-level checks
- [ ] Test every role
- [ ] Test denied access
- [ ] Record denied access in logs

## Permission constants

```go
const (
	PermissionUserManage       = "user.manage"
	PermissionApprovalCreate   = "approval.create"
	PermissionApprovalSubmit   = "approval.submit"
	PermissionApprovalReadOwn  = "approval.read_own"
	PermissionApprovalReadAll  = "approval.read_all"
	PermissionApprovalApprove  = "approval.approve"
	PermissionApprovalReject   = "approval.reject"
	PermissionAuditRead        = "audit.read"
)
```

## Recommended approach

Use:

```go
authorizer.HasPermission(user, PermissionApprovalApprove)
```

Avoid role checks spread across handlers:

```go
if user.Role == "ADMIN" || user.Role == "APPROVER" {
	// avoid this pattern
}
```

## Definition of Done

- Requester cannot approve
- Approver cannot manage users
- Auditor cannot modify approval requests
- Admin can manage users
- Denied actions return `403`

---

# Phase 4 — Approval Request Creation

## Goal

Allow a requester to create and read a draft request.

## Tasks

- [ ] Create approval domain model
- [ ] Create repository interface
- [ ] Create PostgreSQL repository
- [ ] Create service
- [ ] Create handler
- [ ] Add create endpoint
- [ ] Add get endpoint
- [ ] Add list endpoint
- [ ] Filter every query by tenant
- [ ] Add tests

## Endpoints

```http
POST /api/approvals
GET /api/approvals/:id
GET /api/approvals
```

## Create request

```json
{
  "title": "Purchase laptop",
  "description": "Laptop for the engineering team"
}
```

## Response

```json
{
  "id": "approval-id",
  "title": "Purchase laptop",
  "description": "Laptop for the engineering team",
  "status": "DRAFT"
}
```

## Definition of Done

- Requester creates a draft
- Draft belongs to authenticated tenant and user
- Tenant A cannot read Tenant B's request
- List endpoint only returns current tenant data

---

# Phase 5 — Workflow Transitions

## Goal

Implement submit, approve, reject, and cancel rules.

## Tasks

- [ ] Add domain transition methods
- [ ] Submit draft
- [ ] Approve pending request
- [ ] Reject pending request
- [ ] Cancel draft
- [ ] Reject invalid transitions
- [ ] Use database transaction
- [ ] Lock row during approval decision
- [ ] Add tests for every transition

## Endpoints

```http
POST /api/approvals/:id/submit
POST /api/approvals/:id/approve
POST /api/approvals/:id/reject
POST /api/approvals/:id/cancel
```

## Transaction flow

```text
BEGIN
  ↓
SELECT approval FOR UPDATE
  ↓
Verify tenant
  ↓
Verify permission
  ↓
Verify current state
  ↓
Update approval status
  ↓
Insert approval action
  ↓
Insert audit event
  ↓
Insert notification
  ↓
COMMIT
```

On failure:

```text
ROLLBACK
```

## Definition of Done

- Draft can be submitted
- Pending request can be approved or rejected
- Invalid transition returns `409`
- Action is recorded
- Partial data is never committed

---

# Phase 6 — Audit Trail

## Goal

Record important actions in an append-only log.

## Events to record

```text
USER_CREATED
ROLE_ASSIGNED
LOGIN_FAILED
APPROVAL_CREATED
APPROVAL_SUBMITTED
APPROVAL_APPROVED
APPROVAL_REJECTED
APPROVAL_CANCELLED
PERMISSION_DENIED
```

## Audit data

```json
{
  "tenant_id": "tenant-id",
  "actor_id": "user-id",
  "action": "APPROVAL_APPROVED",
  "resource_type": "approval_request",
  "resource_id": "approval-id",
  "request_id": "request-id",
  "metadata": {
    "previous_status": "PENDING",
    "new_status": "APPROVED"
  }
}
```

## Tasks

- [ ] Create audit table
- [ ] Create audit repository
- [ ] Create audit service
- [ ] Write events in the same transaction when required
- [ ] Add audit list endpoint
- [ ] Protect endpoint with auditor permission
- [ ] Add tests

## Endpoint

```http
GET /api/audit-logs
```

## Definition of Done

- Approval events create audit records
- Tenant A cannot read Tenant B's audit logs
- Audit log cannot be modified through the API
- Auditor can read audit records

---

# Phase 7 — Notifications

## Goal

Create reliable notification records from workflow events.

## First implementation

Do not connect a real email provider immediately.

Use:

```text
Approval event
    ↓
Create notification row
    ↓
Worker reads pending notification
    ↓
LogEmailSender sends simulated email
    ↓
Mark notification SENT
```

## Notification statuses

```text
PENDING
PROCESSING
SENT
FAILED
```

## Tasks

- [ ] Create notification table
- [ ] Create sender interface
- [ ] Create log sender
- [ ] Create notification worker
- [ ] Add retry count
- [ ] Add next retry time
- [ ] Add maximum retry rule
- [ ] Add worker tests
- [ ] Add deterministic retry tests

## Sender interface

```go
type EmailSender interface {
	Send(ctx context.Context, message EmailMessage) error
}
```

## Retry example

```text
Attempt 1 → immediate
Attempt 2 → after 1 minute
Attempt 3 → after 5 minutes
Final failure → FAILED
```

Use an injected clock so tests do not depend on real time.

## Definition of Done

- Workflow creates notification
- Worker processes pending notification
- Failed send is retried
- Retry exhaustion produces `FAILED`
- Tests do not use real email

---

# Phase 8 — Error Response and Request ID

## Goal

Use a stable error format across the API.

## Error response

```json
{
  "error": {
    "code": "APPROVAL_NOT_FOUND",
    "message": "approval request not found",
    "request_id": "req-123"
  }
}
```

## Tasks

- [ ] Add request ID middleware
- [ ] Add central error model
- [ ] Map domain errors to HTTP status
- [ ] Avoid leaking database errors
- [ ] Log internal error with request ID
- [ ] Add contract tests

## Recommended status codes

| Case | Status |
|---|---:|
| Validation failed | 400 |
| Invalid or missing token | 401 |
| Permission denied | 403 |
| Resource not found | 404 |
| Invalid state transition | 409 |
| Unexpected error | 500 |
| Database unavailable | 503 |

## Definition of Done

- All errors use one format
- Every response has a request ID
- Internal details are not returned to clients
- Error contract is tested

---

# Phase 9 — MP23 Testing Strategy Lab

## Goal

Create a test pyramid and critical-path test suite.

## Test pyramid

```text
               E2E
          API / Contract
      Integration / Database
              Unit
```

## Unit tests

Test:

- State transitions
- Permission policy
- Retry calculation
- Validation
- Error mapping

## Integration tests

Test:

- PostgreSQL repositories
- Transactions
- Tenant isolation
- Audit persistence
- Notification persistence

## API tests

Test:

- Login
- Create approval
- Submit approval
- Approve approval
- Reject invalid role
- Tenant isolation
- Error format

## Contract tests

Test:

- Required response fields
- Data types
- HTTP status
- Error structure
- Backward-compatible API behaviour

## E2E test

One complete test:

```text
Login requester
    ↓
Create request
    ↓
Submit request
    ↓
Login approver
    ↓
Approve request
    ↓
Verify status
    ↓
Verify action
    ↓
Verify audit log
    ↓
Verify notification
```

## Migration tests

Test:

- Empty database
- Upgrade from previous version
- Rollback
- Existing data remains usable

## Regression test

Create a permanent test for:

```text
Tenant A must not approve Tenant B's request
```

## Flaky-test review

Check:

- No `time.Sleep`
- No test order dependency
- No shared database state
- No real external API
- Fixed clock is used
- Unique test data is generated

## Coverage

Run:

```bash
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Do not treat 100% coverage as the goal.

## MP23 Deliverables

- [ ] `docs/testing/testing-strategy.md`
- [ ] Test pyramid
- [ ] Critical-path test list
- [ ] Unit tests
- [ ] Integration tests
- [ ] API tests
- [ ] Contract tests
- [ ] E2E test
- [ ] Migration tests
- [ ] Retry tests
- [ ] Regression test
- [ ] Coverage report
- [ ] Flaky-test analysis

---

# Phase 10 — Dockerisation

## Goal

Package the application and database environment.

## Dockerfile

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

## Docker Compose

```yaml
services:
  api:
    build:
      context: .
    ports:
      - "8080:8080"
    environment:
      APP_PORT: "8080"
      DATABASE_URL: postgres://postgres:password@db:5432/approval?sslmode=disable
      JWT_SECRET: development-only-secret
      SHUTDOWN_TIMEOUT: 10s
    depends_on:
      db:
        condition: service_healthy

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

## Cloud mapping

| Local component | Cloud category |
|---|---|
| API container | Compute |
| PostgreSQL | Database storage |
| Named volume | Persistent storage |
| Docker network | Private networking |
| Port mapping | Public endpoint |
| Docker image | Deployment artifact |
| Environment variables | Runtime configuration |

## Definition of Done

- Application builds as image
- API and database run with one command
- Data survives container restart
- API reports readiness failure when database stops

---

# Phase 11 — Linux Processes and Graceful Shutdown

## Goal

Stop the application safely during deployment.

## Signals

```text
SIGINT  → Ctrl+C
SIGTERM → graceful termination request
SIGKILL → immediate forced termination
```

## Tasks

- [ ] Inspect process with `ps`
- [ ] Inspect container process
- [ ] Handle `SIGINT`
- [ ] Handle `SIGTERM`
- [ ] Stop accepting new requests
- [ ] Wait for active requests
- [ ] Close database pool
- [ ] Exit within timeout
- [ ] Test shutdown

## Test commands

```bash
docker compose top api
docker compose stop -t 10 api
```

## Definition of Done

- SIGTERM is logged
- Active request completes
- Database pool closes
- Container exits cleanly

---

# Phase 12 — MP25 Deployment Pipeline Lab

## Goal

Build, test, package, migrate, deploy, verify, and roll back.

## Roadmap alignment for Month 12

This phase must cover the Month 12 roadmap item:

```text
Deployment and production readiness:
Docker, CI/CD, migrations, rollback, staging/prod, feature flags
```

It also maps directly to:

| Learning unit | Main topic | Required evidence |
|---|---|---|
| LU034 | Cloud compute, storage, networking; CI/CD; Linux processes and signals | Cloud mapping note, CI pipeline, graceful shutdown evidence |
| LU035 | IAM and least privilege; rollback; TCP connection lifecycle | Permission matrix, rollback runbook, timeout/reuse note |
| LU036 | Docker image layers, multi-stage builds; Linux filesystem permissions; TLS/HTTPS basics | Production Dockerfile, non-root container, TLS termination note |

For LU034 specifically, do not stop at reading. Produce these three artifacts:

```text
docs/deployment/cloud-mapping.md
docs/deployment/ci-cd-pipeline.md
docs/evidence/graceful-shutdown.md
```

LU034 is complete only when you can explain:

1. Which cloud service category runs the API container
2. Which storage is persistent and which storage is disposable
3. Which network paths are public and private
4. What happens when the process receives `SIGTERM`
5. What each CI/CD stage proves before deployment

## Pipeline

```text
Push code
   ↓
Format check
   ↓
Static analysis
   ↓
Unit tests
   ↓
Integration tests
   ↓
API and contract tests
   ↓
Build Go binary
   ↓
Build Docker image
   ↓
Run migration check
   ↓
Deploy
   ↓
Health check
   ↓
Readiness check
   ↓
Smoke test
   ↓
Rollback if verification fails
```

## CI/CD quality gates

Use explicit gates. A deployment candidate should not be promoted unless all required gates pass.

| Gate | Command or check | Blocks deploy? |
|---|---|---:|
| Format | `test -z "$(gofmt -l .)"` | Yes |
| Static analysis | `go vet ./...` | Yes |
| Unit tests | `go test ./...` | Yes |
| Race-sensitive tests | `go test ./... -race` | Yes |
| Integration tests | Run with PostgreSQL service | Yes |
| Migration up | Apply migrations to empty test database | Yes |
| Migration rollback | Roll back at least one version | Yes |
| Docker build | Build image tagged with commit SHA | Yes |
| Smoke test | `/health`, `/ready`, and one business operation | Yes |
| Evidence capture | Save command output or screenshots in `docs/evidence/` | No, but required for portfolio |

## GitHub Actions

Create:

```text
.github/workflows/ci.yml
```

Minimum steps:

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

      - name: Build image
        run: docker build -t approval-api:${{ github.sha }} .
```

## Environment strategy

Create separate environment files and document the difference between local, staging, and production.

```text
.env.example
.env.local.example
.env.staging.example
.env.production.example
```

| Setting | Local | Staging | Production |
|---|---|---|---|
| `APP_ENV` | `local` | `staging` | `production` |
| `DATABASE_URL` | Local Compose database | Managed or staging database | Production managed database |
| `JWT_SECRET` | Local development secret | Secret manager or CI secret | Secret manager only |
| `LOG_LEVEL` | `debug` | `info` | `info` or `warn` |
| `FEATURE_APPROVAL_COMMENT_REQUIRED` | Either value for testing | Mirrors next release behaviour | Controlled release value |

Rules:

- Never commit real secrets.
- Production and staging must not share a database.
- Staging should run the same Docker image tag that will be promoted.
- Configuration should come from environment variables, not hardcoded constants.

## IAM and deployment identity checklist

For LU035, create:

```text
docs/deployment/iam-permission-matrix.md
```

Minimum identities:

| Identity | Allowed actions | Must not allow |
|---|---|---|
| Application runtime | Connect to database, write logs, read runtime config | Deploy new versions, change database schema outside migration path |
| CI test job | Start test database, run tests, build image | Access production secrets |
| CI deploy job | Pull approved image, run migration, deploy selected version | Read user data directly unless required by deploy platform |
| Human maintainer | Trigger deployment, trigger rollback, view logs | Use shared permanent credentials |

The goal is least privilege, not a perfect cloud-specific policy.

## Migration safety checklist

Before deployment:

- [ ] Migration files are immutable once merged
- [ ] Migration names are ordered and descriptive
- [ ] `migrate up` works on an empty database
- [ ] `migrate up` works on seeded previous-version data
- [ ] One-step rollback works for the latest migration
- [ ] Backward-incompatible changes have a two-step plan
- [ ] Migration failure leaves the application in a known state

Two-step migration example:

```text
Release A:
  Add nullable column
  Backfill data
  Deploy code that writes both old and new shape

Release B:
  Make column required after data is complete
  Remove old read path
```

Avoid this in one release:

```text
Drop old column
Deploy new code
Hope nothing still reads the old column
```

## Smoke tests

```bash
curl --fail http://localhost:8080/health
curl --fail http://localhost:8080/ready
```

Then test one business operation.

Minimum business smoke test:

```text
1. Login as requester
2. Create approval request
3. Submit approval request
4. Login as approver
5. Approve approval request
6. Verify final status is APPROVED
7. Verify audit log exists
```

Save the smoke-test script as:

```text
scripts/smoke-test.sh
```

## Rollback

Use versioned images:

```text
approval-api:v1
approval-api:v2
```

Rollback flow:

```text
v2 deployed
    ↓
Smoke test fails
    ↓
Stop v2
    ↓
Deploy v1
    ↓
Run health check
    ↓
Run readiness check
    ↓
Verify data
```

Create:

```text
docs/deployment/rollback-runbook.md
scripts/rollback.sh
```

Rollback runbook must include:

- Current version
- Target previous version
- Reason for rollback
- Migration compatibility status
- Exact commands
- Health check result
- Readiness check result
- Business smoke-test result
- Data verification result
- Follow-up issue or incident note

Important rule:

```text
Application rollback is easy only when the database schema remains compatible with both versions.
```

## Deployment strategy study

Implement only one local deployment strategy.

Explain these three:

- Rolling deployment
- Blue-green deployment
- Canary deployment

For local practice, blue-green is easiest:

```text
Blue  → API v1 on port 8081
Green → API v2 on port 8082
Proxy or manual switch selects active version
```

## Feature flag

Add one simple feature flag:

```env
FEATURE_APPROVAL_COMMENT_REQUIRED=true
```

Use it to enable or disable a validation rule.

Acceptance criteria:

- [ ] When the flag is `false`, approver can approve without comment
- [ ] When the flag is `true`, approver must provide a comment
- [ ] Both behaviours are covered by tests
- [ ] The active flag value is logged at startup without exposing secrets
- [ ] The final report explains why this flag is safe to toggle

Do not use a feature flag to bypass tenant isolation, authentication, permission checks, audit logging, or migration safety.

## Network and TLS notes

For LU035-LU036, add a short note:

```text
docs/deployment/network-and-tls.md
```

Include:

- Public path: client to API or reverse proxy
- Private path: API to PostgreSQL
- Timeout settings for HTTP server, database connection, and outbound email sender
- Connection reuse expectations
- Where TLS terminates in a real deployment
- Why local Docker Compose may use plain HTTP while production must use HTTPS

## Linux filesystem and container permissions

Verify:

- [ ] Container runs as non-root user
- [ ] Application binary is executable by the app user
- [ ] Runtime does not require writing to the application directory
- [ ] Logs go to stdout/stderr
- [ ] Uploaded files are out of scope unless P3.3 is explicitly implemented

Evidence commands:

```bash
docker compose exec api id
docker compose exec api ps
docker compose exec api ls -la /app
```

## MP25 Deliverables

- [ ] Dockerfile
- [ ] Docker Compose
- [ ] Environment configuration
- [ ] Staging and production configuration notes
- [ ] Cloud compute/storage/network mapping
- [ ] IAM least-privilege matrix
- [ ] Migration scripts
- [ ] Migration safety checklist
- [ ] CI pipeline
- [ ] CI/CD quality gates
- [ ] Docker image build
- [ ] Deployment script
- [ ] Smoke-test script
- [ ] Rollback script
- [ ] Rollback runbook
- [ ] Deployment strategy explanation
- [ ] One feature flag
- [ ] Network, TCP timeout, and TLS note
- [ ] Linux process, signal, and filesystem-permission evidence
- [ ] Deployment evidence

---

# Phase 13 — Monitoring Plan

## Goal

Define how the system will be observed.

## Logs

Every request log should include:

```text
timestamp
level
request_id
tenant_id
user_id
method
path
status_code
duration_ms
error_code
```

Do not log:

- Passwords
- JWT tokens
- Sensitive request bodies
- Database passwords

## Metrics

Technical metrics:

```text
HTTP request count
HTTP error rate
Request latency
Database connection usage
Readiness failures
Worker failures
```

Business metrics:

```text
Approvals created
Approvals submitted
Approvals approved
Approvals rejected
Average approval duration
Notification success rate
```

## Alerts

Define alerts for:

```text
5xx error rate above threshold
Readiness failure
Database unavailable
High p95 latency
Notification retries exhausted
Worker stopped
Migration failed
```

## Monitoring deliverable

Create:

```text
docs/monitoring/monitoring-plan.md
```

Include:

- Logs
- Metrics
- Dashboards
- Alerts
- Runbook links

---

# Phase 14 — Architecture Documentation

## Required documents

```text
docs/architecture/
├── c4-context.md
├── c4-container.md
├── c4-component.md
└── erd.md
```

Each document should explain:

1. What the diagram shows
2. Important boundaries
3. Security assumptions
4. Runtime dependencies
5. Known limitations

---

# Phase 15 — ADRs

Create only decisions that affect the architecture.

Recommended ADRs:

```text
ADR 001 — Shared PostgreSQL database with tenant_id
ADR 002 — Centralised permission-based RBAC
ADR 003 — Transactional workflow, action, and audit updates
ADR 004 — Notification outbox or persisted notification queue
ADR 005 — Docker Compose for local deployment
```

## ADR template

```markdown
# ADR XXX: Decision title

## Status

Accepted

## Context

What problem must be solved?

## Decision

What was selected?

## Reasons

Why was it selected?

## Alternatives considered

What other options were considered?

## Consequences

What improves and what becomes harder?
```

---

# 10-Week Implementation Plan

Use this plan across Months 10–12.

## Week 1 — Scope and architecture

- [ ] Confirm roles
- [ ] Confirm permissions
- [ ] Define critical flow
- [ ] Draft C4 context
- [ ] Draft ERD
- [ ] Create repository
- [ ] Add health endpoint

## Week 2 — Database foundation

- [ ] Add PostgreSQL
- [ ] Create migrations
- [ ] Create seed data
- [ ] Add readiness endpoint
- [ ] Test migration up and down

## Week 3 — Authentication

- [ ] Password hashing
- [ ] Login
- [ ] JWT
- [ ] Authentication middleware
- [ ] Authentication tests

## Week 4 — RBAC and tenant isolation

- [ ] Roles
- [ ] Permissions
- [ ] Authorisation policy
- [ ] Tenant context
- [ ] Tenant isolation tests
- [ ] ADR 001 and ADR 002

## Week 5 — Approval workflow

- [ ] Create draft
- [ ] Read and list
- [ ] Submit
- [ ] Approve
- [ ] Reject
- [ ] Invalid transition tests

## Week 6 — Audit and notifications

- [ ] Approval actions
- [ ] Audit logs
- [ ] Notification records
- [ ] Worker
- [ ] Retry logic
- [ ] ADR 003 and ADR 004

## Week 7 — MP23 testing strategy

- [ ] Test pyramid
- [ ] Unit tests
- [ ] Integration tests
- [ ] API tests
- [ ] Contract tests
- [ ] E2E test
- [ ] Migration tests
- [ ] Regression test

## Week 8 — Docker and runtime operations

- [ ] Dockerfile
- [ ] Docker Compose
- [ ] Environment configuration
- [ ] Persistent volume
- [ ] Graceful shutdown
- [ ] Runtime dependency documentation

## Week 9 — MP25 deployment pipeline

- [ ] GitHub Actions
- [ ] Docker build in CI
- [ ] Migration step
- [ ] Deployment script
- [ ] Smoke test
- [ ] Rollback test
- [ ] Rollback runbook

## Week 10 — Monitoring and final evidence

- [ ] Monitoring plan
- [ ] Complete C4 diagrams
- [ ] Complete ERD
- [ ] Complete ADRs
- [ ] Performance baseline
- [ ] Final README
- [ ] Project demo
- [ ] Final report

---

# Recommended Order for Each Feature

For every feature, use this implementation loop:

```text
1. Define behaviour
2. Write domain rule
3. Define repository interface
4. Add database migration if required
5. Implement PostgreSQL repository
6. Implement service
7. Implement handler
8. Add unit tests
9. Add integration tests
10. Add API test
11. Add logs and metrics
12. Update documentation
```

Do not start from the HTTP handler and place all business logic there.

---

# Git Branch and Commit Strategy

Use small branches:

```text
feature/auth-login
feature/rbac-policy
feature/approval-create
feature/approval-submit
feature/approval-approve
feature/audit-log
feature/notification-worker
feature/docker-deployment
feature/ci-pipeline
```

Commit examples:

```text
feat(auth): add JWT login flow
feat(approval): implement submit transition
test(tenant): prevent cross-tenant approval access
feat(audit): record approval decision events
build(docker): add multi-stage API image
ci: run tests and build Docker image
docs(adr): record tenant isolation decision
```

---

# Security Checklist

- [ ] Passwords are hashed
- [ ] JWT has expiration
- [ ] JWT secret comes from environment
- [ ] Tenant ID comes from authenticated context
- [ ] Every tenant query filters by tenant ID
- [ ] Permission checks are centralised
- [ ] Sensitive fields are not logged
- [ ] Database errors are not exposed
- [ ] Input is validated
- [ ] SQL queries use parameters
- [ ] Invalid transitions are rejected
- [ ] Audit logs record important actions
- [ ] CORS is configured deliberately
- [ ] Development secrets are not committed

---

# Final Definition of Done

## Product

- [ ] Tenant can be created
- [ ] Users can authenticate
- [ ] Roles and permissions work
- [ ] Requester can create a request
- [ ] Requester can submit a request
- [ ] Approver can approve or reject
- [ ] Invalid transitions are blocked
- [ ] Tenant isolation works
- [ ] Approval actions are stored
- [ ] Audit trail is stored
- [ ] Notification is created
- [ ] Notification retry works

## Testing

- [ ] Unit tests
- [ ] Integration tests
- [ ] API tests
- [ ] Contract tests
- [ ] E2E critical-path test
- [ ] Migration tests
- [ ] Async worker tests
- [ ] Retry tests
- [ ] Regression test
- [ ] Coverage report
- [ ] Flaky-test review

## Deployment

- [ ] Dockerfile
- [ ] Docker Compose
- [ ] Environment config
- [ ] Local, staging, and production config notes
- [ ] Cloud compute, storage, and network mapping
- [ ] IAM least-privilege matrix
- [ ] Migrations
- [ ] Migration safety checklist
- [ ] Persistent storage
- [ ] Health endpoint
- [ ] Readiness endpoint
- [ ] Graceful shutdown
- [ ] CI pipeline
- [ ] CI/CD quality gates
- [ ] Docker image build
- [ ] Smoke test
- [ ] Rollback test
- [ ] Rollback runbook
- [ ] Feature flag example
- [ ] Network, timeout, and TLS note
- [ ] Container non-root and filesystem permission evidence

## Documentation

- [ ] C4 context diagram
- [ ] C4 container diagram
- [ ] C4 component diagram
- [ ] ERD
- [ ] 3–5 ADRs
- [ ] Testing strategy
- [ ] Deployment pipeline document
- [ ] Rollback runbook
- [ ] Cloud mapping note
- [ ] IAM permission matrix
- [ ] Network and TLS note
- [ ] Monitoring plan
- [ ] README
- [ ] Final project report

---

# Evidence Checklist

Keep visible evidence in `docs/evidence/` or README.

- [ ] Login response
- [ ] Create approval response
- [ ] Submit approval response
- [ ] Approve response
- [ ] Cross-tenant access rejection
- [ ] Audit record
- [ ] Notification record
- [ ] Migration output
- [ ] Test output
- [ ] Coverage report
- [ ] E2E result
- [ ] Docker build
- [ ] Docker Compose status
- [ ] Health response
- [ ] Readiness response
- [ ] Database failure readiness response
- [ ] Graceful shutdown log
- [ ] CI result
- [ ] Deployment result
- [ ] Rollback result
- [ ] Performance baseline
- [ ] Architecture diagrams
- [ ] ADR links

---

# Final Demo Script

Use this flow to demonstrate the finished project:

```text
1. Start the system with Docker Compose
2. Show health and readiness
3. Login as requester
4. Create an approval request
5. Submit the request
6. Login as an unauthorised user and show 403
7. Login as approver
8. Approve the request
9. Show approval action
10. Show audit log
11. Show notification record
12. Attempt cross-tenant access and show rejection
13. Run automated tests
14. Show CI pipeline
15. Show graceful shutdown
16. Deploy a new version
17. Run smoke test
18. Roll back to the previous version
19. Show C4 diagrams, ERD, ADRs, and monitoring plan
```

---

# Final Project Report Template

```markdown
# P3 Final Project Report

## Project summary

Describe the problem and target users.

## Scope

List implemented and excluded features.

## Architecture

Link C4 diagrams and explain major components.

## Data model

Link ERD and explain tenant isolation.

## Authentication and RBAC

Explain identity, roles, permissions, and policy checks.

## Workflow

Explain states and allowed transitions.

## Audit trail

Explain recorded events and append-only behaviour.

## Notifications

Explain worker, retries, and failure handling.

## Testing strategy

Summarise MP23 test pyramid and critical-path coverage.

## Deployment pipeline

Summarise MP25 build, migration, deployment, verification, and rollback flow.

## Monitoring plan

Describe logs, metrics, and alerts.

## Architectural decisions

Link ADRs.

## Evidence

Link test, pipeline, deployment, and rollback evidence.

## Limitations

List known limitations.

## Next improvements

List the next three architectural improvements.
```

---

# Start Here

Begin with only these tasks:

```text
1. Create the repository
2. Add GET /health
3. Start PostgreSQL with Docker Compose
4. Add GET /ready
5. Create the tenants migration
6. Create the users migration
7. Create seed users for four roles
8. Draft the ERD
9. Draft the critical approval flow
10. Commit the first working version
```

Do not start notification retries, CI/CD, or deployment strategies before authentication, tenant isolation, and the approval flow are working.

---

# LU034 Start Here

Use this section when you are specifically studying **LU034**.

LU034 assumes the core P3 application already has:

- Authentication
- RBAC
- Tenant isolation
- Approval create, submit, and approve flow
- Audit logs
- Notification records
- Basic tests
- `/health` and `/ready`

If those are not done yet, finish the critical flow first. LU034 is about proving the application can be packaged and operated.

## LU034 study objective

By the end of LU034, you should be able to say:

```text
I can map this app to cloud compute, storage and networking.
I can explain what CI/CD proves before deployment.
I can inspect the running process and shut it down safely.
```

## LU034 implementation checklist

- [ ] Build the API Docker image
- [ ] Run API and PostgreSQL with Docker Compose
- [ ] Confirm `/health` returns `200`
- [ ] Confirm `/ready` returns `200` when PostgreSQL is available
- [ ] Confirm `/ready` fails when PostgreSQL is unavailable
- [ ] Add graceful shutdown handling for `SIGTERM`
- [ ] Run `docker compose stop -t 10 api`
- [ ] Save shutdown log evidence
- [ ] Create GitHub Actions CI
- [ ] Run format, vet, tests, migration check, and Docker build in CI
- [ ] Write `docs/deployment/cloud-mapping.md`
- [ ] Write `docs/deployment/ci-cd-pipeline.md`
- [ ] Write `docs/evidence/graceful-shutdown.md`

## LU034 evidence checklist

Save command output or screenshots for:

- [ ] Docker image build
- [ ] Docker Compose startup
- [ ] Health check
- [ ] Readiness check
- [ ] Database-down readiness failure
- [ ] Process inspection with `ps`
- [ ] Graceful shutdown log
- [ ] CI pipeline run
- [ ] Docker image tag using commit SHA

## LU034 self-test questions

Answer these in your own words:

1. What is the compute unit for this application in local Docker and in a cloud deployment?
2. Which data survives a container restart?
3. Which network path should be private?
4. Why should CI run migration checks before deployment?
5. What is the difference between `SIGTERM` and `SIGKILL`?
6. Why should the app stop accepting new requests before closing the database pool?
7. What would make `/ready` fail while `/health` still passes?
