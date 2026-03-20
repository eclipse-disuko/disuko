# Backend Quickstart Guide

Welcome to the Disuko backend! This guide will help you set up and run the backend service for local development.

## Prerequisites

- [Go](https://go.dev/dl/) 1.24.x or higher
- [Docker](https://www.docker.com/) (for running dependencies like CouchDB, Valkey, etc.)
- [CouchDB](https://couchdb.apache.org/) and [valkey](https://valkey.io/)

## 1. Install Dependencies

Navigate to the backend directory and install Go modules:

```sh
cd backend
go mod tidy
```

> You can use `go mod tidy` to add missing and remove unused dependencies.

## 2. Configure Environment

- Copy or edit the configuration files in `conf/` (e.g., `config-local.yml`, `config.yml`).
- Set environment variables as needed (see `conf/config.go` for all options).

## 3. Start Required Services


```sh
docker compose -f docker-compose-local.yml up --build
```
or with Podman:
```sh
podman-compose -f docker-compose-local.yml up --build
```

This will start CouchDB, Valkey, and the backend server as defined in the Compose file.

- To stop the services:
  ```sh
  docker compose -f docker-compose-local.yml down
  ```


## 4. TLS

To generate self-signed TLS certificates for local development:

```sh
sh createTLS.sh
```


## 5. Run the Backend Server

To run the Go navigate to Backend folder:

```sh
go run main.go
```

Or, to build and run the binary:

```sh
go build -o app
./app
```



## 6. Project Structure

- `main.go` — Entry point
- `conf/` — Configuration files
- `domain/`, `infra/`, `connector/` — Main application logic
- `resources/` — Static files and templates

---