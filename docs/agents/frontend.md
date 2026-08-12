<!--
SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG

SPDX-License-Identifier: Apache-2.0
-->
<!-- TOC -->
* [Frontend](#frontend)
  * [About](#about)
    * [Compilation](#compilation)
    * [Start](#start)
  * [Coding rules](#coding-rules)
  * [Coding guidelines](#coding-guidelines)
    * [Creating a typical admin list page (title + add button + search + grid)](#creating-a-typical-admin-list-page-title--add-button--search--grid)
    * [Adding a new tile to the admin dashboard (#/dashboard/admin)](#adding-a-new-tile-to-the-admin-dashboard-dashboardadmin)
    * [Extending the breadcrumbs](#extending-the-breadcrumbs)
    * [Adding / extending i18n translations](#adding--extending-i18n-translations)
<!-- TOC -->

# Frontend

Shared Node.js configuration hosting the frontend applications `portal`, `rlm` and `cli`.

## About

The frontend lives in the `frontend` folder and is written in TypeScript/Vue using
Node.js tooling. Code shared across apps is in `libs/shared/`, app-specific code
in `libs/<app>/`, and each `apps/<app>/` folder is the entry point (router, plugins only).

### Compilation

```
cd frontend && npm install
```

### Start

```
npm run dev:portal-local
```

## Coding rules

- Every file must start with the SPDX license header
- Do not write comments
- Use vuetify components where possible
- Avoid using v-col and v-row, use flexbox or grid instead, see @shared/layouts
- Use tailwindcss, do not write inline styles
- Services live in `libs/portal/services/`; implement as classes calling `getApi()` at module level; export a default singleton instance
- Pinia stores in `stores/`; name them `use*Store`; use setup-function style; wrap state in `reactive({})`, expose via `toRefs()`; use `storeToRefs()` when destructuring in components
- Use `useSnackbar()` (`info()` / `error()`) for user-facing feedback
- Models live in `libs/portal/model/`; use `interface` for types
- Use `@shared/types/table` (`DataTableHeader`, `DataTableItem`) for all data table column definitions
- All UI strings through `vue-i18n`; translation keys are SCREAMING_SNAKE_CASE; portal keys in `libs/portal/i18n/locales/{en,de}.json`, shared keys in `libs/shared/i18n/locales/{en,de}.json`
- Use event bus (`@disclosure-portal/utils/eventbus`) for cross-component communication when props/emits/store can't handle it
- Use `dayjs` with shared constants (`DATE_FORMAT`, `DATE_FORMAT_SHORT`, etc.) from `@shared/utils/constant` for dates
- Use layout components from `@shared/layouts` (`Stack`, `TableLayout`, `DialogLayout`, `ReactiveDialogLayout`)
- Code shared across apps in `libs/shared/`; app-specific code in `libs/<app>/`; `apps/<app>/` is the entry point (router, plugins only)

## Coding guidelines

### Creating a typical admin list page (title + add button + search + grid)

Reference implementation: `libs/portal/views/admin/Newsbox.vue`, routed at
`/dashboard/admin/newsbox` (same pattern as `/dashboard/admin/users`).

1. **View file** — place the page under `libs/<app>/views/<area>/<Name>.vue` (e.g.
   `libs/portal/views/admin/Newsbox.vue`). Every file starts with the SPDX header, then a
   `<script lang="ts" setup>` block, then `<template>`.
2. **Route** — register it in `apps/<app>/src/router/index.ts`: add a lazy import
   (`const Name = () => import('@disclosure-portal/views/admin/Name.vue');`) and a child
   route entry with `path`, `name`, `component`, and `meta.title` (`en`/`de`) under the
   existing `admin` parent route.
3. **Store + service** — fetch/mutate data through a Pinia store (`use<Name>Store`) that
   calls a service class (`<name>.service.ts`, `getApi()` singleton); never call the API
   directly from the view (see the general coding rules above for the exact store/service
   shape).
4. **Page layout** — wrap everything in `@shared/layouts` `TableLayout`. It has two named
   slots:
   - `#buttons`: page title (`<h1 class="text-h5">`), an optional `DCActionButton` to open the
     create dialog, then `<v-spacer />`, then `DSearchField v-model="search"` right-aligned.
   - `#table`: the `v-data-table` (or similar grid).
   `TableLayout` measures the `#buttons` slot and window size and sets the `#table` wrapper's
   height so the grid always fills the remaining viewport height exactly (pixel-perfect,
   no page-level scroll) — you get this for free, do not set a fixed height yourself.
   ```html
   <TableLayout data-testid="newsbox">
     <template #buttons>
       <h1 class="text-h5">{{ t('NEWSBOX') }}</h1>
       <DCActionButton large icon="mdi-plus" :text="t('BTN_ADD')" class="mx-2" @click="openCreateDialog" />
       <v-spacer></v-spacer>
       <DSearchField v-model="search" />
     </template>
     <template #table>
       <v-data-table v-model:search="search" :headers="headers" :items="filteredItems" item-value="_key" />
     </template>
   </TableLayout>
   ```
5. **Column headers with a filter icon** — define columns as `DataTableHeader[]`
   (`@shared/types/table`) via a `computed()`. For a column that needs a filter dropdown
   instead of plain text, override its header slot with `GridFilterHeader` (keeps
   sort-icon/alignment behavior) wrapping a `GridHeaderFilterIcon` (the actual filter
   dropdown, bound with `v-model` to a `ref<string[]>`):
   ```html
   <template #[`header.expiry`]="{column, getSortIcon, toggleSort}">
     <GridFilterHeader :column="column" :toggleSort="toggleSort" :getSortIcon="getSortIcon">
       <template #filter>
         <GridHeaderFilterIcon
           v-model="selectedFilterStatus"
           :column="column"
           :label="t('STATUS')"
           :allItems="[{text: t('ACTIVE'), value: 'false'}, {text: t('EXPIRED'), value: 'true'}]" />
       </template>
     </GridFilterHeader>
   </template>
   ```
   Apply the selected filter values in a `computed()` (e.g. `filteredItems`) that filters the
   store's raw items before passing them to `:items` on the `v-data-table`.
6. **Actions in a cell** — add an `actions` column to `headers` (`value: 'actions'`, not
   sortable), then override its item slot with `TableActionButtons`, passing a list of
   `{icon, hint, event}` buttons and listening to the emitted events:
   ```html
   <template #[`item.actions`]="{item}">
     <TableActionButtons
       variant="slider"
       :buttons="getActionButtons(item)"
       @edit="openEditDialog(item)"
       @delete="showConfirmDelete(item)" />
   </template>
   ```
   Use `variant="slider"` together with `useTableActionSlider()` when the actions column
   width needs to grow/shrink to reveal extra buttons; use `width: sliderWidth.value` on the
   `actions` header definition in that case.
7. **Opening a dialog from a row/button action** — keep a `showDialog`/`isEditMode`/
   `editingItem` ref set, put a `<v-dialog v-model="showDialog">` next to (not inside)
   `TableLayout`, and render a `DialogLayout` inside it with a `config` (`title`,
   `primaryButton`, `secondaryButton`, `loading`, `showIdle`) and the form as default slot:
   ```html
   <v-dialog v-model="showDialog" content-class="large" scrollable width="600">
     <DialogLayout :config="dialogConfig" @primary-action="submit" @secondary-action="closeDialog" @close="closeDialog">
       <v-form ref="formRef" @submit.prevent="submit">
         <Stack>
           <TextField v-model="form.title" :label="t('TITLE')" required />
         </Stack>
       </v-form>
     </DialogLayout>
   </v-dialog>
   ```
   For destructive actions (delete), use the shared `ConfirmationDialog` component/
   `ConfirmationType` instead of a custom dialog.

### Adding a new tile to the admin dashboard (`#/dashboard/admin`)

Reference implementation: `libs/portal/views/admin/AdminDashboard.vue`. The page renders one
`ITile` (`libs/portal/model/ITile.ts`) per admin module, gated by a right/role check.

1. **Have a target page first** — the tile just links to an existing route (e.g.
   `/dashboard/admin/newsbox`); create that page/route first if it doesn't exist yet (see the
   "Creating a typical admin list page" section above).
2. **Decide the access check** — in `AdminDashboard.vue`'s `tiles` computed, every tile is
   wrapped in an `if (RightsUtils.xxx())` from `@shared/user/utils/RightsUtils`:
   - Prefer a **role check** if one already fits: `RightsUtils.isApplicationAdmin()`,
     `isDomainAdmin()`, `isFOSSOffice()`, `isProjectAnalyst()`, etc. (`RightsUtils.hasRole`
     against the `Group` enum in `@shared/user/models/Rights`).
   - Otherwise add a **CRUD-based check**: add a `allow<Feature>: CRUDRights` field to the
     `Rights` class and a `has<Feature>Access()` method there (e.g. `hasLabelAccess()`), then
     mirror that method as a static wrapper in `RightsUtils` (e.g.
     `RightsUtils.hasLabelAccess = () => RightsUtils.rights().allow<Feature> && (...)`). The
     `Rights` object itself is populated from the backend profile response
     (`@shared/user/services/profile.service.ts`) — a new `allow<Feature>` right must also be
     added/returned there on the backend side for the check to ever be `true`.
3. **Push the tile** — inside the `tiles` computed, guarded by the check from step 2:
   ```ts
   if (RightsUtils.isApplicationAdmin()) {
     res.push({
       color: 'primary',
       cnt: -1, // or counts.value?.xxxCount ?? -1, see step 4
       visible: true,
       title: 'FEATURE_TITLE_KEY',
       url: '/dashboard/admin/feature',
       icon: 'mdi-icon-name',
       expandGroup: false,
       expand: false,
     });
   }
   ```
   `title` is an i18n key (SCREAMING_SNAKE_CASE, add it to
   `libs/portal/i18n/locales/{en,de}.json`), `icon` is any `mdi-*` Material Design Icon, and
   `cnt: -1` means "no count badge" (rendered as an invisible placeholder character instead of
   a number).
4. **Optional live count badge** — to show a count instead of `-1`: add a field to
   `DashboardCounts` (`libs/shared/types/DashboardCounts.ts`), have the backend's
   `GET /api/v1/counts/dashboard` handler populate it (see `s.handlers.count` /
   `GetDashboardCountsHandler` in `server/routes.go` on the backend), then reference it as
   `counts.value?.<field> || -1` in the tile — `counts` is already fetched in `onMounted()` via
   `adminService.getDashboardCounts()`.

### Extending the breadcrumbs

Breadcrumbs are driven by `useBreadcrumbsStore` (`libs/shared/stores/breadcrumbs.store.ts`,
Pinia) and rendered by `DBreadcrumb.vue`, which reads `currentBreadcrumbs` and auto-disables
(makes non-clickable) the last entry — you never set `disabled` on the last item yourself.

1. **Set the trail from the page** — every view that shows in the breadcrumb bar defines an
   `initBreadcrumbs()` function and calls `breadcrumbs.setCurrentBreadcrumbs([...])` with an
   array of `{title, href}` (title = i18n key, prefixed `BC_`), called once in `onMounted()`
   for static pages. Call `initBreadcrumbs()` again whenever data the trail depends on
   changes (e.g. a project name loaded async) — see `initBreadcrumbs`/`watch(currentProject, ...)`
   in `ProjectsDetail.vue`.
2. **Reuse the shared prefix segments** — don't hand-write the leading `Dashboard`/`Admin`
   crumbs; the store already exposes them:
   - `dashboard` — `{title: t('BC_Dashboard'), href: '/dashboard/home'}`
   - `adminDashboard` — `{title: t('BC_ADMIN'), href: '/dashboard/admin'}`
   - `dashboardCrumbs` — a computed that resolves to `[dashboard, adminDashboard]` when the
     current route is under `.../admin/...`, otherwise just `[dashboard]`. Use this for any
     new admin page instead of hard-coding both segments:
     ```ts
     const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();
     const initBreadcrumbs = () => breadcrumbs.setCurrentBreadcrumbs(dashboardCrumbs);
     ```
   - `projectsCrumb` — the same idea for the `Projects` section root, spread together with
     `dashboardCrumbs` plus the page-specific tail crumb(s), e.g.
     `breadcrumbs.setCurrentBreadcrumbs([...dashboardCrumbs, projectsCrumb, currentProjectCrumb])`.
3. **Adding a new reusable section-root crumb** — if you're introducing a whole new area (not
   just one page) that multiple pages will link back to, add a new named constant next to
   `dashboard`/`adminDashboard`/`projectsCrumb` in `breadcrumbs.store.ts` and return it from
   the store, so every page under that section imports and spreads it instead of duplicating
   the `{title, href}` literal.
4. **Page-specific tail crumb(s)** — for a single page (or a page whose last crumb depends on
   loaded data, e.g. a detail page's entity name), just append plain `{title, href}` objects
   after the shared prefix — no store changes needed for these, e.g.
   `{title: t('BC_NEWSBOX'), href: '/dashboard/admin/newsbox'}` in `Newsbox.vue`.

### Adding / extending i18n translations

There was no prior documentation for this — the setup has three layers: static shared keys,
static app-specific keys, and dynamic keys served by the backend at runtime.

1. **Static shared keys** (`libs/shared/i18n/locales/{en,de}.json`) — use these for strings
   used by more than one app (`portal`, `rlm`, `cli`), e.g. common buttons/errors. Add the
   same key to both `en.json` and `de.json`.
2. **Static app-specific keys** (`libs/<app>/i18n/locales/{en,de}.json`, e.g.
   `libs/portal/i18n/locales/en.json`) — use these for strings only relevant to one app.
   `libs/<app>/i18n/index.ts` builds the `vue-i18n` instance by spreading shared messages
   first, then the app's own messages on top (`{...sharedEN, ...en}`), so an app-specific key
   **overrides** a shared key with the same name. This `i18n` instance is registered once via
   `app.use(i18n)` in `apps/<app>/src/plugins/index.ts`. Compile-time optimization is provided
   by the `VueI18nPlugin` (`unplugin-vue-i18n`) configured in `apps/<app>/vite.config.mts`.
3. **Key naming & usage** — keys are SCREAMING_SNAKE_CASE (breadcrumb keys additionally get a
   `BC_` prefix). Never hardcode UI strings in a template; call `const {t} = useI18n();` in
   `<script setup>` and use `t('YOUR_KEY')`.
4. **Dynamic keys from the backend** — on top of the two static layers, translations can also
   be managed at runtime through the backend's `i18n` collection (admin UI at
   `/dashboard/admin/i18n`, `I18n.vue`/`LocaleDetails.vue`, backed by `i18n.service.ts` calling
   `GET/PUT/DELETE /api/v1/i18n/...`). On every app start, `App.vue`'s `onMounted` calls
   `loadLocales()` → `i18nService.getLocales()` to list all backend-known locales, then
   `loadI18nLocale(code)` for each one, which fetches `GET /api/v1/i18n/{code}` and merges the
   result into the running `vue-i18n` instance:
   ```ts
   const existing = i18n.global.getLocaleMessage(code);
   i18n.global.setLocaleMessage(code, {...res.data.entries, ...existing});
   ```
   Note the spread order: **static keys (`existing`) win** over backend keys on a collision —
   the backend collection is effectively a runtime-extendable **fallback/extension** layer for
   keys that don't already exist in the static JSON bundles, not an override mechanism. Use it
   for content that should be editable by admins without a frontend redeploy (e.g. adding a
   brand-new locale, or filling in missing translations) rather than for changing existing
   static keys.
5. **Adding a new locale entirely at runtime** — do this through the `/dashboard/admin/i18n`
   admin UI (`i18nService.upsertTranslation`/`importLocale`) instead of touching the static
   JSON files; it only requires a new locale row in the backend's `i18n` collection, no
   frontend deploy.
