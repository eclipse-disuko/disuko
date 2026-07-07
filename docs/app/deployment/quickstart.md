# Deployment Instructions via Docker Compose

Follow these steps to deploy the Disuko backend (and dependencies) using Docker Compose:

---

## 1. Prerequisites

- [Docker] installed on your system
- [Docker Compose] 

---

## 2. Customize Environment (Important)

Before deploying, you may need to customize environment variables and configuration files to match your environment.

- **Environment Variables:**  
  You can define environment variables directly in the `docker-compose.yml` file under each service using the `environment:` key, or by using a `.env` file in your project root.

- **Configuration Files:**  
  Make sure your application’s config files (like `config.yml` or `config-local.yml`) are set up with the correct values for your deployment.

---

## 3. Build and Start the Services

Open a terminal in your project root and run:

```sh
docker compose -f docker-compose.yml up --build
```

- The `--build` flag ensures all images are built before starting.
- Add `-d` to run in detached mode:
  ```sh
  docker compose -f docker-compose.yml up --build -d
  ```

---

## 4. Verify the Deployment

- Check running containers:
  ```sh
  docker compose ps
  ```
- View logs:
  ```sh
  docker compose logs -f
  ```

---

## 5. Stopping the Services

To stop and remove all containers, networks, and volumes created by Compose:

```sh
docker compose -f docker-compose.yml down
```