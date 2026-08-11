# Unit test plan — `libs/portal/views`

Status: planning document for coordinating multiple agents. Not itself a test file (excluded from vitest by name, but move/delete once work is done if you'd rather not keep it in the tree).

## 1. Scope

Frontend has ~0 test coverage. This plan covers unit tests for `libs/portal/views/**/*.vue` only.

**Only views that are identical downstream are in scope** — a downstream-specific view would make the test (and any future edits to keep it green) diverge from what disukorlm actually ships, so it's not worth the maintenance burden here.

### In scope — identical in both repos (33 views)

```
Dashboard.vue
Home.vue
tasks/Tasks.vue
licenses/Licenses.vue
projects/Projects.vue
admin/AdminProfile.vue
admin/ReviewTemplates.vue
admin/MailTemplates.vue
admin/Users.vue
admin/UpcomingDeletions.vue
admin/AdminProjects.vue
admin/FeatureFlags.vue
admin/AdminClassifications.vue
admin/Newsbox.vue
admin/CustomIds.vue
announcements/Announcements.vue
user/UserProfile.vue
policies/PolicyRules.vue
analytics/Stats.vue
admin/tools/SampleData.vue
admin/tools/Mail.vue
admin/tools/PolicyRuleClassificationMatrix.vue
admin/tools/S3DataIntegrity.vue
admin/tools/TermsOfUseManagement.vue
admin/tools/SystemStatisticTable.vue
admin/tools/Jobs.vue
admin/tools/NotificationBarManagement.vue
admin/tools/AccessRights.vue
admin/schema/Schemas.vue
admin/schema/SchemaMain.vue
admin/i18n/I18n.vue
admin/i18n/LocaleDetails.vue
admin/checklist/ChecklistMain.vue
admin/checklist/Checklist.vue
```

### Out of scope — diverges downstream, skip

```
licenses/LicenseMain.vue
projects/ProjectsDetail.vue
projects/ProjectsVersions.vue
admin/AdminTools.vue
admin/Labels.vue
admin/AdminDashboard.vue
policies/PolicyRulesDetail.vue
analytics/AnalyticMain.vue
admin/tools/ExportImportTools.vue
admin/tools/Statistics.vue
```

Re-run the diff before starting work in case either repo has moved on:

```bash
cd /Users/MINHCNG/Projects/github/disuko/frontend/libs/portal/views
for f in $(find . -type f -name "*.vue"); do
  d="/Users/MINHCNG/Projects/ghe/disukorlm/frontend/libs/portal/views/$f"
  [ -f "$d" ] && diff -q "$f" "$d" || echo "MISSING_DOWNSTREAM: $f"
done
```

## 2. Existing conventions to reuse (don't reinvent)

- Runner: vitest (`npm run test` / `test:run` / `test:coverage`), jsdom env, config at `frontend/vitest.config.mts`.
- Global setup `frontend/vitest.setup.ts` already: mocks `vue-router` (`useRoute`/`useRouter` return static stubs), installs a real `vue-i18n` instance with the `en` locale loaded, and stubs `$t` to return the key. `ResizeObserver` is mocked globally.
- `libs/portal/test-utils/vuetify-stubs.ts` exports `vuetifyStubs` — a shared map of lightweight Vuetify component stubs (`v-btn`, `v-card`, `v-data-table`, `v-dialog`, etc.). Spread this into `global.stubs` instead of hand-rolling new stubs.
- Pinia: use `createTestingPinia({createSpy: vi.fn, stubActions: false})` from `@pinia/testing` as a `global.plugins` entry when a component touches a real Pinia store (see `libs/portal/components/new-wizard/tests/WizardStepOwner.test.ts`). `stubActions: false` so store actions still run — most view-level tests only need to spy on `setCurrentBreadcrumbs` calls, not stub everything.
- Service/composable mocking: `vi.mock('@disclosure-portal/services/...', () => ({default: {...}}))`, combined with `vi.hoisted()` for the mock fns so they're reset in `beforeEach` (see `NewOrEditLicenseDialog.test.ts`).
- Test file location: co-locate under a `tests/` subfolder next to what's tested (matches `components/dialog/tests/`, `components/new-wizard/tests/`). For views, mirror that per-directory, e.g. `views/tests/Dashboard.test.ts`, `views/admin/tools/tests/Jobs.test.ts`.

## 3. Step 0 — shared view test helper (must land first, single agent, blocks nothing else logically but avoid duplicate work)

Almost every view follows the same shape: it calls `useBreadcrumbsStore()` and sets breadcrumbs (via `onMounted`, `onBeforeMount`, or synchronously in `<script setup>`), then renders one big child component (a `Grid*`, `Table*`, or feature component) that should be **stubbed**, not exercised — the child's own behavior is out of scope for a view test.

One agent should add `libs/portal/test-utils/view-test-utils.ts` before batch work starts, to avoid 33 slightly-different copies of the same boilerplate:

```ts
// mountView(Component, { childStubs, piniaOptions, mountOptions }) that:
// - installs createTestingPinia (createSpy: vi.fn, stubActions: false)
// - spreads vuetifyStubs
// - merges in view-specific child component stubs (e.g. { GridProjects: true })
// - returns { wrapper, pinia } so callers can grab useBreadcrumbsStore(pinia) and assert on it
```

Because `useBreadcrumbsStore` is a real Pinia store (not mocked), the pattern for asserting breadcrumbs is:

```ts
const pinia = createTestingPinia({stubActions: false});
const wrapper = mount(AdminProjects, {global: {plugins: [pinia], stubs: {...vuetifyStubs, AdminGridProjects: true}}});
const breadcrumbs = useBreadcrumbsStore(pinia);
expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([...]);
```

Note `useRoute`/`useRouter` are already globally mocked in `vitest.setup.ts`, so `dashboardCrumbs`/`isInAdminArea` will resolve against `path: '/'` unless a test overrides the mock's return value for a specific case.

This step is small (~1-2 hours) and everything else depends on its output existing (or agents proceed without it and accept minor duplication — call this out if Step 0 isn't done yet when a batch starts).

## 4. Batches (parallelizable once Step 0 lands)

Views are tiered by size/complexity of their own logic (not their rendered children). Each batch is a separate agent; batches don't touch shared files except adding new `*.test.ts` files, so they're safe to run concurrently. Within a batch, each view's test is an independent file — no ordering needed.

### Batch A — thin wrapper views (breadcrumbs + one child render only)

Pattern: assert breadcrumbs are set correctly, assert the expected child component renders. No service/store mocking beyond breadcrumbs.

```
Dashboard.vue is NOT here (see Batch C) — everything else "under 100 lines" from wc -l:
admin/AdminProjects.vue        (15 lines)
projects/Projects.vue          (18 lines)
admin/UpcomingDeletions.vue    (21 lines)
admin/CustomIds.vue            (23 lines)
admin/checklist/Checklist.vue  (23 lines)
admin/MailTemplates.vue        (24 lines)
admin/Users.vue                (27 lines)
tasks/Tasks.vue                (34 lines)
admin/tools/Mail.vue           (49 lines)
admin/AdminProfile.vue         (59 lines)
user/UserProfile.vue           (64 lines)
admin/tools/NotificationBarManagement.vue (82 lines)
Home.vue                       (91 lines, but see note below)
admin/tools/S3DataIntegrity.vue        (97 lines)
admin/tools/TermsOfUseManagement.vue   (101 lines)
```

Note: `Home.vue` pulls in `ProfileService`, `useAppStore`, `useNewsboxStore`, `useUserStore`, `useWizardStore` — more store/service surface than the rest of this batch. Still small, but budget a bit more time for it and mock those stores/services per the Step 0 pattern.

### Batch B — medium views (some local state/computed, still single responsibility)

```
admin/tools/SampleData.vue                 (147 lines)
admin/tools/PolicyRuleClassificationMatrix.vue (157 lines)
admin/FeatureFlags.vue                     (166 lines)
admin/i18n/I18n.vue                        (168 lines)
announcements/Announcements.vue            (186 lines)
admin/checklist/ChecklistMain.vue          (187 lines) — uses useChecklistsStore, useLabelStore
analytics/Stats.vue                        (198 lines)
admin/tools/SystemStatisticTable.vue       (199 lines)
```

These need `vi.mock` for at least one service (e.g. `admin/tools/SampleData.vue`, `SystemStatisticTable.vue`) or a real store dependency (`ChecklistMain.vue`). Assert: initial render doesn't throw, key data-fetch call fires on mount, and any user action (button click, row select) triggers the expected service call — mirror the assertion style in `NewOrEditLicenseDialog.test.ts`.

### Batch C — larger views (real business logic worth covering beyond smoke)

```
admin/schema/SchemaMain.vue     (249 lines) — AdminService, useUserStore, useRoute/useRouter
admin/schema/Schemas.vue        (283 lines) — AdminService, useUserStore, useRouter
admin/ReviewTemplates.vue       (290 lines) — adminService, useFormSubmission, useSnackbar
admin/AdminClassifications.vue  (339 lines) — useView composable, confirmation dialog
admin/tools/Jobs.vue            (343 lines) — AdminService, confirmation dialog, snackbar
Dashboard.vue                   (388 lines) — useAppStore, useUserStore, useThemeStore, useLanguageStore, eventBus, logout util
policies/PolicyRules.vue        (418 lines)
admin/Newsbox.vue               (431 lines)
admin/tools/AccessRights.vue    (448 lines)
admin/i18n/LocaleDetails.vue    (720 lines) — i18nService, DialogLayout, form state
licenses/Licenses.vue           (897 lines) — heaviest: license service, filtersets service, admin service, snackbar, confirmation dialog, table state
```

For this batch, don't aim for full line coverage — pick the 3-5 behaviors that matter per view (data loads on mount and populates the table/grid props; a primary action like create/delete/confirm calls the right service and shows a snackbar; an error response surfaces via `useSnackbar().error`; any guarded/rights-based conditional rendering branches both ways). Treat `Licenses.vue` and `LocaleDetails.vue` as their own sub-task each — don't bundle either in with three other views in the same agent turn.

## 5. Per-view test checklist (apply to every view regardless of batch)

1. `it('sets the expected breadcrumbs on mount')` — assert `setCurrentBreadcrumbs` call args (skip only if the view truly has no breadcrumbs call).
2. `it('renders without throwing')` — baseline smoke test with all children/services stubbed.
3. For views with a data fetch on mount: assert the service/store fetch method is called once, and that a stubbed successful response reaches the child component's props (don't assert on the child's internal rendering — stub it and check the prop it received).
4. For views with a user-triggered action (create/delete/confirm/export): assert the service call fires with expected args, and that success/error both drive the right snackbar/dialog outcome (mock `useSnackbar`).
5. Skip: styling/CSS assertions, deep-diving into stubbed child components' internals, testing Vuetify's own behavior.

## 6. Definition of done

- New `*.test.ts` files only — no edits to the views themselves (unless a genuine bug is uncovered; flag it separately rather than silently changing behavior).
- `npm run test:run` passes locally for the new files.
- `npm run test:coverage` — spot check the touched views actually gained coverage (no false-positive smoke tests that stub away everything and assert nothing meaningful).
- Existing lint/format rules pass (`SPDX` header comment at top of each new file, matching the pattern in every existing test file).

## 7. Suggested agent split (5 agents, in this order)

1. **Agent 0 (sequential, blocks 2-4 ideally, but not strictly)**: Step 0 helper + Batch A's first 2-3 views as a reference implementation others can copy.
2. **Agent 1 (parallel)**: rest of Batch A.
3. **Agent 2 (parallel)**: Batch B.
4. **Agent 3 (parallel)**: Batch C minus `Licenses.vue` and `LocaleDetails.vue`.
5. **Agent 4 (parallel, or sequential after 3 if context/time constrained)**: `Licenses.vue` and `LocaleDetails.vue` only — these two alone are comparable in size to all of Batch A combined.

Agents 1-4 can start as soon as Agent 0's helper file exists; if Agent 0 is slow, any of them can proceed without it (inline the boilerplate) and it can be refactored into the shared helper afterward without behavior changes.
