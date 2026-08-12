# AGENTS.md

## Repository structure

Monorepo structure with the following packages:

- docs: Vitepress documentation, intended for developers
- frontend: Contains one shared configuration for frontend applications. This includes portal, rlm and cli (3 different apps)
- backend: Go backend
- infra: Kubernetes and Helm Charts

## Project documentation
Per-project details (folder, compilation, coding rules, coding guidelines) live in `docs/agents/`:
- `docs/agents/server.md`: backend
- `docs/agents/frontend.md`: frontend (portal, rlm, cli)

## Contributing

- Write commits message in conventional commit format, e.g. `feat: add new feature` or `fix: resolve bug`
- Write commit messages focused on user impact, not implementation details
- Make sure to have gitleaks configured with `pre-commit install`, do not commit credentials or secrets
- PR should contain short description of the changes and explain why they are needed, maximum 3 bulletpoints.
