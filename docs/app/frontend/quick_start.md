# Frontend Quickstart Guide

Welcome to the DISUKO frontend! This guide will help you get the frontend up and running for local development.

---

## Prerequisites

- [Node.js](https://nodejs.org/) (version 22.15.x as specified in [`frontend/package.json`](../../frontend/package.json))
- (Optional) [yarn](https://yarnpkg.com/) or [pnpm](https://pnpm.io/) if you prefer

---

## 1. Getting Started

**Open a terminal and navigate to the frontend folder:**

```sh
cd frontend
```

---

## 2. Node.js Version Management with nvm

This project uses [nvm (Node Version Manager)](https://github.com/nvm-sh/nvm) to ensure you are running the correct Node.js version as specified in the `.nvmrc` file.

- **On macOS/Linux:**
  ```sh
  nvm use
  ```
- **On Windows (nvm-windows):**
  ```sh
  nvm use 22.15.0
  ```
  (Replace `22.15.0` with the version in your `.nvmrc` if different.)

If you do not have the required version installed, nvm will prompt you to install it.

---

## 3. Install Dependencies

Install project dependencies:

```sh
npm install
```

---

## 4. Start the Development Server

Start the portal app using Vite:

```sh
npm run dev
```

## 5. Project Structure

- `apps/portal/` — App wrapper and configuration for the DISUKO
- `libs/portal/` — Core content and components of the DISUKO application

---

## 6. Useful Commands

- **Linting:**
  ```sh
  npm run lint
  ```
- **Build for production:**
  ```sh
  npm run build
  ```
- **Preview production build:**
  ```sh
  npm run preview
  ```

---

## 7. Features

- Vue 3 + Vuetify 3 UI
- Pinia for state management
- Vue Router for navigation
- Automated component importing
- TypeScript support

---