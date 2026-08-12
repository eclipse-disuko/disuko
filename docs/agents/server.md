<!--
SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG

SPDX-License-Identifier: Apache-2.0
-->

<!-- TOC -->
* [Server](#server)
  * [About](#about)
    * [Prerequisites](#prerequisites)
    * [Compilation](#compilation)
    * [Start](#start)
  * [Coding rules](#coding-rules)
  * [Coding guidelines](#coding-guidelines)
    * [Adding a new route with its own handler](#adding-a-new-route-with-its-own-handler)
    * [Database entity model & repository pattern](#database-entity-model--repository-pattern)
    * [Adding a database migration step](#adding-a-database-migration-step)
    * [Adding a database seed file](#adding-a-database-seed-file)
<!-- TOC -->

# Server

Go backend of the disuko project.

## About

The backend lives in the `backend` folder and is written in Go. It follows a
4-layer architecture: `domain/` → `infra/repository/` → `infra/service/` → `infra/rest/`.

### Prerequisites

The server reads its config from `./conf/config-local.yml` (optional, gitignored) and
`./conf/config.yml` (relative paths, so the process must be started with `backend` as its
working directory). It also needs a running database (CouchDB) and cache (Valkey/Redis):

```
docker-compose -f docker-compose-local.yml up -d
```

### Compilation

```
cd backend && go build -o /tmp/dps
```

### Start

```
cd backend && /tmp/dps
```

## Coding rules

- Every file must start with the SPDX license header
- Do not write comments
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

## Coding guidelines

### Adding a new route with its own handler

Note the two distinct `rest.go` files involved — do not confuse them:
- `domain/<domain>/rest.go` — the **DTOs** and their `ToDto()`/`ToEntity()` mappers.
- `infra/rest/<name>.go` — the **HTTP handler** (package `rest`) that calls the repository
  and renders DTOs. The filename does not have to be `rest.go` here; any name is fine
  as long as the package is `rest`.

1. **Entity** — the domain entity lives in `domain/<domain>/<domain>.go`, embedding
   `domain.RootEntity` (or `domain.ChildEntity`).
2. **DTO** — add the `<Domain>Dto` struct to `domain/<domain>/rest.go`, embedding
   `domain.BaseDto` where applicable, plus `ToDto()`/`ToEntity()` methods implementing
   `domain.ConvertableEntity[T]` / `domain.ConvertableDto[T]` (see `domain/base_mapper.go`
   and `domain/base_rest.go`). Never expose the entity directly over HTTP.
3. **Repository interface** — add `infra/repository/<domain>/layer.go` with an
   `I<Domain>Repository` interface embedding `base.IBaseRepository[...]` (or
   `IBaseRepositoryWithHardDelete[...]`) plus any domain-specific methods.
4. **Repository implementation** — add `infra/repository/<domain>/repository.go` with the
   private struct and a `New<Domain>Repository(requestSession *logy.RequestSession) I<Domain>Repository`
   constructor, e.g. `department.NewDepartmentRepository`.
5. **Wire the repository** — `server/database.go` is the single place where every repository
   constructor is called: add a field of type `I<Domain>Repository` to `dbRepos` and call
   `New<Domain>Repository(requestSession)` inside `setupDatabase()`. Never construct a
   repository anywhere else.
6. **Handler** — add a file under `infra/rest/` with a `<Name>Handler` struct whose fields are
   repository interfaces (e.g. `Repo departmentRepo.IDepartmentRepository`), and one method per
   route with signature `func (h *<Name>Handler) Method(w http.ResponseWriter, r *http.Request)`.
   Inside: get the session via `logy.GetRequestSession(r)`, check rights via
   `roles.GetAccessAndRolesRightsFromRequest`, call `handler.Repo.Xxx(...)`, map entities to
   DTOs with `domain.ToDtos`/`entity.ToDto()`, and respond with `render.JSON(w, r, dtos)`.
7. **Wire the handler** — `server/handlers.go` is the single place where every handler is
   constructed: add a field to the `handlers` struct and, in `setupHandlers()`, build
   `s.handlers.<name> = rest.<Name>Handler{Repo: s.repos.<domain>}`, injecting the repository
   created in step 5. Never construct a handler anywhere else.
8. **Route** — register the endpoint in `server/routes.go` inside `setupRoutes()` using
   go-chi, referencing the handler method directly, e.g.
   `r.Get("/departments/find/{searchStr}", s.handlers.department.Find)`.

### Database entity model & repository pattern

**Folder placement** (`domain/`)
- Every entity gets its **own subfolder** under `domain/<domain>/`, named after the business
  concept (e.g. `domain/department/`, `domain/customid/`, `domain/checklist/`). Never add a
  new entity as a loose file directly under `domain/`; only the shared base types
  (`base.go`, `base_mapper.go`, `base_rest.go`) live at the package root.
- Inside `domain/<domain>/`, put the entity struct in `<domain>.go` and its DTO(s) +
  `ToDto()`/`ToEntity()` mappers in `rest.go` (see the DTO note above). Split into more files
  in the same folder only if the domain grows large (see `domain/project/` for an example
  with many files).
- `domain/Readme.md` additionally documents (not fully applied everywhere in the existing
  code, but to be followed for new entities): DB-persisted models should be suffixed
  `Entity`, DTOs suffixed `Dto`, root/collection models implement `domain.RootEntity`, and
  child models needing a unique id implement `domain.ChildEntity`.

**Root vs. child entity** (`domain/base.go`)
- `domain.RootEntity` (`ChildEntity` + `Rev`) — for entities that live in their own top-level
  collection and have a dedicated repository. Embed with `` `bson:",inline"` ``, e.g.
  `project.Project`, `sbomlist.SbomList`.
- `domain.ChildEntity` (`Key`/`Created`/`Updated`) — for entities nested/embedded inside a
  parent `RootEntity`'s document (no own collection, no revision, no repository), e.g.
  `project.Token`, `project.ProjectVersion`, `pdocument.PDocument`.
- Rule of thumb: if it needs its own `I<Domain>Repository`, it's a `RootEntity`; if it only
  ever appears as a field inside another entity, it's a `ChildEntity`.

**Soft delete vs. hard delete** (`infra/repository/base/base_repository_softdelete.go`,
`base_repository_harddelete.go`)
- **Soft delete is the default.** Embed `domain.SoftDelete` (`Deleted bool`, implements
  `domain.ISoftDelete`) next to `domain.RootEntity` in the entity struct, e.g.
  `project.Project` embeds both `domain.RootEntity` and `domain.SoftDelete`. The repository
  interface embeds `base.IBaseRepositoryWithSoftDelete[*Entity]` and the constructor uses
  `base.CreateRepositoryWithSoftDelete(...)`. `Delete()` then only sets `Deleted = true`;
  the row stays in the DB and is excluded from normal `Find*`/`Count*` calls. Use the
  `*WithDeleted` variants (`FindByKeyWithDeleted`, `FindAllWithDeleted`, `CountAllWithDeleted`,
  `ExistByKeyWithDeleted`) to include soft-deleted rows, and `DeleteHard()` to purge one for
  real.
- **Hard delete** is opt-in for entities that don't need an undo/history (e.g. lookup data
  like `department.Department`, which loads everything into memory) or for **migration
  steps**. In that case the repository interface embeds
  `base.IBaseRepositoryWithHardDelete[*Entity]` and the constructor uses
  `base.CreateRepositoryWithHardDelete(...)`; `Delete()` physically removes the document.
- **Migration on an otherwise soft-delete collection**: add a second, migration-only
  interface embedding `base.IBaseRepositoryWithHardDelete[*Entity]` for the same entity, e.g.
  `approvallist.IApprovalListRepositoryMigration`, so migration code can hard-delete rows
  from a collection whose regular repository (`IApprovalListRepository`) stays soft-delete.

**Creating a new `I<Domain>Repository`**
1. Pick `RootEntity`/`ChildEntity` and soft-/hard-delete per the rules above.
2. `infra/repository/<domain>/layer.go`: declare `I<Domain>Repository` embedding
   `base.IBaseRepositoryWithSoftDelete[*Entity]` (default) or
   `base.IBaseRepositoryWithHardDelete[*Entity]`, plus domain-specific methods.
3. `infra/repository/<domain>/repository.go`: private struct embedding
   `base.BaseRepositoryWithSoftDelete[*Entity]` (or `BaseRepositoryWithHardDelete[*Entity]`),
   and `New<Domain>Repository(requestSession *logy.RequestSession) I<Domain>Repository` calling
   `base.CreateRepositoryWithSoftDelete`/`CreateRepositoryWithHardDelete` with the collection
   name, an `entityCreator func() *Entity`, an optional `preDelete` hook, optimize-unset
   attributes, and indexes.
4. Wire it into `dbRepos`/`setupDatabase()` in `server/database.go` as described above — this
   is still the single place repository constructors are called.

### Adding a database migration step

Migrations live in `infra/service/startup/startup.go` on `StartUpHandler` and run once at
server startup via `MigrateDatabase()`.

1. **Register the step** — add one entry to the `steps` slice inside `MigrateDatabase()`:
   ```go
   {Name: "MIGRATE_REMOVE_ORPHANED_SBOM_FILES", Do: startUpHandler.migrateRemoveOrphanedSbomFiles},
   ```
   `Name` must be a **unique key** across all entries. It is hashed into the `Migration`
   entity's `Key` (`domain.migration.New`, `infra/repository/migration`) and is also the
   attribute `MigrateDatabase()` queries on to decide whether the step already ran
   (`startUpHandler.MigrationRepository.Query(...)` filtered by `Name`). Reusing an existing
   name will make your step silently skipped; a colliding hash would corrupt the dedupe
   record — always pick a new, descriptive, upper-snake-case name.
   Only after `step.Do(requestSession)` **returns** does `MigrateDatabase()` call
   `MigrationRepository.Save(...)` to mark the step done — so a step that never returns
   (e.g. an unrecovered panic) will re-run on every restart, and later steps in the slice
   are never reached.
2. **Dependencies** — if the migration needs a repository/service not yet on `StartUpHandler`,
   add a field there and wire it in `server/handlers.go` where `startup.StartUpHandler{...}`
   is constructed (same place all its other repositories are injected).
3. **Method signature and logging** — every migration method has this exact shape and must
   log a `START` and an `END` line with `logy.Infof`, using the method name as prefix:
   ```go
   func (startUpHandler *StartUpHandler) migrateRemoveOrphanedSbomFiles(requestSession *logy.RequestSession) {
       logy.Infof(requestSession, "migrateRemoveOrphanedSbomFiles - START")

       // ... migration logic ...

       logy.Infof(requestSession, "migrateRemoveOrphanedSbomFiles - END")
   }
   ```
4. **Make it crash-proof** — the `END` log line must always be reached, and one failing item
   must not abort the whole step. Wrap the body in `exception.TryCatchAndLog` (or
   `exception.TryCatch` with a custom catch) so a panic is recovered and logged instead of
   propagating:
   ```go
   func (startUpHandler *StartUpHandler) migrateXxx(requestSession *logy.RequestSession) {
       logy.Infof(requestSession, "migrateXxx - START")

       exception.TryCatchAndLog(requestSession, func() {
           // top-level guard for the whole step
           items := startUpHandler.SomeRepository.FindAll(requestSession, false)
           for _, item := range items {
               exception.TryCatch(func() {
                   // per-item logic; one bad item must not stop the loop
               }, func(e exception.Exception) {
                   exception.LogException(requestSession, e)
               })
           }
       })

       logy.Infof(requestSession, "migrateXxx - END")
   }
   ```
   See `migrateSyncProjectAndSbomRetentionFlags` and `migrateRemoveOrphanedSbomFiles` in
   `infra/service/startup/startup.go` for full reference implementations of this pattern,
   including counting processed/succeeded/failed items in the logs.
5. **Soft-delete vs. hard-delete repositories when reading "all" data** — know which kind of
   repository you're querying (see the entity model section above):
   - **Soft-delete repositories** (`IBaseRepositoryWithSoftDelete`) hide entities with
     `Deleted = true` from `FindAll`/`FindByKey`/etc. If the migration must also consider
     soft-deleted rows (e.g. to avoid deleting S3 files that still belong to a soft-deleted
     SBOM list), use the `*WithDeleted` variants, e.g.
     `startUpHandler.SbomListRepository.FindAllWithDeleted(requestSession, false)`.
     Using plain `FindAll` here would silently ignore soft-deleted rows and could delete data
     that is still supposed to be recoverable.
   - **Hard-delete repositories** (`IBaseRepositoryWithHardDelete`) have no `Deleted` flag at
     all — deleted rows are physically gone, so plain `FindAll` already returns everything
     that exists.
   - Never call `Delete()`/mutate an entity in a migration without first checking which
     delete semantics its repository has — soft-deleting where a hard delete (or vice versa)
     was intended silently changes behavior for the rest of the application.

### Adding a database seed file

Seeding runs once at startup via `(dbRepos).seedDb()` in `server/seeding.go`, called from
`setupDatabase()` in `server/database.go`. There are two independent mechanisms.

**Regular entity seeds (`.jsonl`)**
1. **Register the collection** — add an entry to `entityCreatorMap` in `server/seeding.go`:
   ```go
   "myCollection": func() interface{} {
       return &mydomain.MyEntity{}
   },
   ```
   Without an entry here, a matching file is silently skipped (no error).
2. **Add the file** — `conf/dbseeds/defaultdb/<collectionName>.jsonl` (and, if the data should
   also apply to the `VanillaDisuko` config, `conf/dbseeds/disuko/<collectionName>.jsonl`).
   Format is **JSONL** (one JSON object per line, no surrounding array), with field names
   matching the entity's `json` tags. Every object needs a **non-empty** `Key` on its embedded
   `RootEntity` — `insertIfNotExists` reads it via reflection
   (`FieldByName("RootEntity").FieldByName("Key")`), so this only works for entities embedding
   `domain.RootEntity` directly.
3. **No further wiring needed** — `seedDb()` scans every `*.jsonl` file in the seed folder on
   startup and calls `processSeedFile` automatically; there is no central file registry beyond
   `entityCreatorMap`.
4. **Idempotency** — by default `insertIfNotExists` queries by key and skips insertion if a
   row already exists, so it's safe to leave seed files in the repo and re-run on every
   startup. A handful of collections (`labels`, `spdxSchemas`, `jobs`, `licenses`) instead use
   `legacyInsertIfNotExists` with custom lookup logic (e.g. by name/type instead of key) — add
   a new `else if collName == "..."` branch there (and to the dispatch condition in
   `processSeedFile`) if your entity needs a non-key existence check.

**i18n seeds (`.json`, separate mechanism)**
- File under `conf/dbseeds/i18n/<name>.<localeCode>.json` (e.g. `strings.en.json`), format
  `{"KEY": "value", ...}`.
- Loaded by `seedI18n()`, guarded by `db.i18nLocale.GetLocaleCount(...) > 0` — i.e. it only
  runs while the `i18n` collection is **completely empty**; unlike the `.jsonl` seeds this is
  a one-time bootstrap, not idempotent per key. Adding entries to an existing locale later on
  should go through the `/dashboard/admin/i18n` admin UI (see the frontend i18n docs) or the
  `/api/v1/i18n` endpoints instead of new seed files.
