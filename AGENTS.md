# AGENTS.md

## Repository structure
Monorepo structure with the following packages:
- docs: Vitepress documentation, intended for developers
- frontend: Contains one shared configuration for frontend applications. This includes portal, rlm and cli (3 different apps)
- backend: Go backend
- infra: Kubernetes and Helm Charts

## Commands
- `cd frontend && npm install && npm run dev:portal-local` to start the frontend
- `cd backend && go build -o /tmp/dps && /tmp/dps` to start backend

## General coding guidelines
- Do not write comments
- Every file must start with the SPDX license header (same format as backend)

## Frontend coding guidelines
- Use vuetify components where possible
- Avoid using v-col and v-row, use flexbox or grid instead, see @shared/layouts
- Use tailwindcss, do not write inline styles
- Services live in `libs/portal/services/`; implement as classes calling `getApi()` at module level; export a default singleton instance
- Pinia stores in `stores/`; name them `use*Store`; use setup-function style; wrap state in `reactive({})`, expose via `toRefs()`; use `storeToRefs()` when destructuring in components
- Use `useSnackbar()` (`info()` / `error()`) for user-facing feedback
- Models live in `libs/portal/model/`; Use `interface` for types.
- Use `@shared/types/table` (`DataTableHeader`, `DataTableItem`) for all data table column definitions
- All UI strings through `vue-i18n`; translation keys are SCREAMING_SNAKE_CASE; portal keys in `libs/portal/i18n/locales/{en,de}.json`, shared keys in `libs/shared/i18n/locales/{en,de}.json`
- Use event bus (`@disclosure-portal/utils/eventbus`) for cross-component communication when props/emits/store can't handle it
- Use `dayjs` with shared constants (`DATE_FORMAT`, `DATE_FORMAT_SHORT`, etc.) from `@shared/utils/constant` for dates
- Use layout components from `@shared/layouts` (`Stack`, `TableLayout`, `DialogLayout`, `ReactiveDialogLayout`)
- Code shared across apps in `libs/shared/`; portal-specific in `libs/portal/`; `apps/portal/` is the entry point (router, plugins only)

## Backend coding guidelines
- Every file must start with the SPDX license header
- 4-layer architecture: `domain/` → `infra/repository/` → `infra/service/` → `infra/rest/`; never skip or reverse layers
- Each domain area under `infra/repository/` must have a `layer.go` (interface) and a separate implementation file
- `I` prefix for interfaces (e.g., `IProjectRepository`); `Struct` suffix for private implementations (e.g., `projectRepositoryStruct`)
- Pass `*logy.RequestSession` as the first parameter to all repository and service methods
- Use `exception.ThrowException*` in HTTP handlers; never return raw Go errors from handlers
- All error codes and i18n message keys as constants in `helper/message/messages.go`
- Separate entities from REST DTOs; use `ToDto()`/`ToEntity()` via `ConvertableEntity`/`ConvertableDto` interfaces in `domain/base_mapper.go`
- Use `logy` package for logging; never use `fmt.Print*` or the standard `log` package
- Use `go-chi/chi/v5` for routing, `rest.Val` for validation, `go-chi/render` for JSON responses
- Use `New*` constructors; never instantiate entities with raw struct literals
- Use `domain.MapTo`, `domain.ToDtos`, `domain.ToEntities` helpers for slice conversions; never write manual loops
- Obtain database via `base.NewDatabase()` using the `IDatabase` interface; never reference specific drivers directly
- Config via `jinzhu/configor`; access through `conf.Config`

## Contributing
- Write commits message in conventional commit format, e.g. `feat: add new feature` or `fix: resolve bug`
- Write commit messages focused on user impact, not implementation details
- Make sure to have gitleaks configured with `pre-commit install`, do not commit credentials or secrets
- PR should contain short description of the changes and explain why they are needed, maximum 3 bulletpoints.