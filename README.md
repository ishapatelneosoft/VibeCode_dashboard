# DASHBOARD Project — Go API + Next.js Kanban

Full-stack app with:

- **Backend**: Go + Gin + PostgreSQL (session-cookie auth, Kanban board, tasks, worklogs, time reports, Swagger)
- **Frontend**: Next.js (App Router) Kanban UI (Azure/Jira-style) with drag & drop + reports

## What you get

- **Auth**: register/login/logout/session validation/current user
- **Board**: columns + tasks, create task, move task (drag & drop), assign users, history
- **Time tracking**: log work per task and view a time report (per-task + grand total)

## Repository structure

```
.
├── README.md
├── AI/                         # AI-related documentation and features
├── backend/                    # Backend Go application
│   ├── cmd/server/             # Backend entrypoint (Go)
│   ├── internal/               # Backend code (controller/service/repository/domain)
│   ├── migrations/             # SQL migrations (PostgreSQL)
│   ├── docs/                   # Swagger artifacts (committed)
│   ├── tests/                  # Go tests
│   ├── k8s/                    # Kubernetes manifests (Kustomize)
│   ├── docker/                 # Docker-related files
│   ├── Dockerfile              # Backend Docker image
│   ├── docker-compose.yml      # Local Postgres + backend
│   ├── go.mod
│   └── server
├── frontend/                   # Next.js frontend
└── plans/                      # Project plans and architecture
```

## Quickstart (recommended)

Start backend + DB with Docker, run frontend locally.

### 1) Backend + Postgres (Docker)

From backend directory:

```bash
cd backend
docker-compose up --build
```

Backend URLs:

- `http://localhost:8080` (API landing page)
- `http://localhost:8080/health`
- `http://localhost:8080/swagger/index.html`

### 2) Frontend (local)

```bash
cd frontend
source "$HOME/.nvm/nvm.sh"
nvm use 25
npm install
npm run dev
```

Frontend URL:

- `http://localhost:3000`

## Configuration

### Backend env vars (Go)

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | PostgreSQL user |
| `DB_PASSWORD` | `p@ssw0rd` | PostgreSQL password |
| `DB_NAME` | `auth_db` | Database name |
| `DB_SSL_MODE` | `disable` | SSL mode |
| `SERVER_PORT` | `8080` | HTTP server port |

### Frontend env vars (Next.js)

Create/update `frontend/.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### (check URLs)

- `http://localhost:8080`
- `http://localhost:8080/health`
- `http://localhost:8080/swagger/index.html`

### (DB persistence)

1. Create some data (via UI or API).
2. Stop containers:

```bash
docker-compose down
```

3. Start again:

```bash
docker-compose up
```

Notes:
- Postgres data is stored in Docker volume **`pgdata`**.
- `migrations/*.sql` run automatically **only on first DB init** (when `pgdata` is empty).

### KPI 38 (smoke test via UI)

1. Open `http://localhost:3000`
2. Register → Login
3. Create a task
4. Drag the task to another column (move)
5. Open the task details and log time
6. Click **View reports** and verify totals

## Kubernetes (Kustomize)

Manifests are in `k8s/`:

- `k8s/base`: Deployment/Service/ConfigMap/Secret (probes use `/health`)
- `k8s/overlays/dev`: 1 replica
- `k8s/overlays/prod`: 2 replicas

Apply:

```bash
kubectl apply -k k8s/overlays/dev
```

Important:
- Assumes an **external PostgreSQL** (managed DB or separate in-cluster DB chart).
- Update `k8s/base/configmap.yaml` and `k8s/base/secret.yaml` with real DB settings.

## Development

### Backend tests

```bash
go test ./tests/...
go test -cover ./...
```

### Frontend build

```bash
cd frontend
source "$HOME/.nvm/nvm.sh"
nvm use 25
npm run build
```

## Notes

More detailed feature notes and flows live under `AI/`.