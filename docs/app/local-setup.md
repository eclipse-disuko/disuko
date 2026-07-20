# Run it locally

Everything else in these docs is about what DISUKO does and why. This page is the one technical exception: how to get it running on your own machine.

## Option A: full demo stack with Docker (recommended)

This is the fastest way to see DISUKO working end to end — a complete stack (backend, frontend, Keycloak for login, and its databases) started with one command.

**Prerequisites:** [Docker](https://www.docker.com/) and Docker Compose.

From the project root:

```sh
./setup-dev.sh   # Windows: setup-dev.ps1
```

This generates the local TLS certificate the stack needs (choose `mkcert` for a browser-trusted certificate, or a self-signed one).

```sh
docker compose up --build -d
```

Check that everything came up:

```sh
docker compose ps --format "{{.Service}} {{.State}}"
```

Then open `https://localhost:3009/` in your browser and log in with one of the built-in demo accounts:

```
Username: CUSTOMER1   Password: CUSTOMER1
Username: CUSTOMER2   Password: CUSTOMER2
```

If login misbehaves, log out first at `https://localhost:3009/api/v1/oauth/logout` and try again. If the setup wizard asks for an owner or company name you don't have yet, `dummy` works fine for a test run.

**Try uploading an SBOM:** create a project, then go to **Admin** and upload an SBOM schema labeled `common standard` — the official [SPDX 2.3 schema](https://github.com/spdx/spdx-spec/blob/support/2.3/schemas/spdx-schema.json) works. Only after that is a schema in place can a project accept SBOM uploads.

To stop everything:

```sh
docker compose down
```

## Option B: running the code directly (for contributing)

If you want to modify the frontend or backend code itself, run each natively while letting Docker handle just the databases:

```sh
docker compose -f docker-compose-local.yml up --build
```

This starts CouchDB and Valkey, which both the frontend and backend expect to be available.

**Backend** (needs [Go](https://go.dev/dl/) matching the version in `backend/go.mod`):

```sh
cd backend
sh createTLS.sh   # generates server.key / server.crt for local TLS
go build -o /tmp/dps && /tmp/dps
```

**Frontend** (needs the Node.js version in `frontend/.nvmrc`):

```sh
cd frontend
nvm use
npm install
npm run dev:portal-local
```

For the full breakdown of the repository's structure and coding conventions, see [AGENTS.md](https://github.com/eclipse-disuko/disuko/blob/main/AGENTS.md) in the project root.
