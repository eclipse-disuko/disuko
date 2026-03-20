# Project Quickstart Guide

Welcome to the DISUKO project! This guide will help you get both the backend and frontend up and running for local development.

---

## Prerequisites

- [Node.js](https://nodejs.org/) (version 22.15.x as specified in `frontend/.nvmrc`)
- [Go](https://go.dev/dl/) (version 1.24.x or higher as specified in go.mod)
- [Docker](https://www.docker.com/)
- [CouchDB](https://couchdb.apache.org/) and [Valkey](https://valkey.io/)

---

## 1. Start Required Services (CouchDB, Valkey) with Docker Compose

From the project root, run:

```sh
docker compose -f docker-compose-local.yml up --build
```

This will start CouchDB, Valkey, and any other required services as defined in your Compose file.

---

## 2. Backend Setup

Open a new terminal and navigate to the backend folder:

```sh
cd backend
```

To generate self-signed TLS certificates for local development:

```sh
sh createTLS.sh
```

- Install Go dependencies:
  ```sh
  go mod tidy
  ```
- Run the backend server:
  ```sh
  go run main.go
  ```

---

## 3. Frontend Setup

Open another terminal and navigate to the frontend folder:

```sh
cd frontend
```

- **Node.js Version Management:**
  - On macOS/Linux:
    ```sh
    nvm use
    ```
  - On Windows (nvm-windows):
    ```sh
    nvm use 22.15.0
    ```
    (Replace `22.15.0` with the version in your `.nvmrc` if different.)
- Install dependencies:
  ```sh
  npm install
  ```
- Start the development server:
  ```sh
  npm run dev
  ```

---


For more details, see the individual quickstart guides in `docs/app/backend/quickstart.md` and `docs/app/frontend/quickstart.md`.
