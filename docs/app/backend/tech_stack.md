# Backend Tech Stack

This document outlines the main technologies and tools used in the Disuko backend.

## Core Language & Framework

- **Go (Golang)**  
  The backend is written in Go, a statically typed, compiled language known for its simplicity and performance.

## Database

- **CouchDB**  
  Used as the primary document-oriented database for storing application data.

## Caching & Messaging

- **Valkey**  
  Used for caching, session management, and lightweight message brokering.

## Containerization & Orchestration

- **Docker**  
  Used for containerizing backend services and dependencies.
- **Kubernetes**  
  Used for orchestrating containers in development and production environments.

## Configuration & Environment

- **YAML**  
  Configuration files are managed in YAML format (see the `conf/` directory).
- **Environment Variables**  
  Used for sensitive and environment-specific settings.

## DevOps & Tooling

- **Helm**  
  For managing Kubernetes deployments.
- **Shell Scripts**  
  Used for automation and environment setup.



For more details, see the [Quickstart Guide](./quick_start.md)