# Component Map — `ProjectsVersions.vue`

> **Scope**: Only portal-specific (`@disclosure-portal`) and notable third-party components are listed.
> Vuetify (`v-*`) and `@shared` components (e.g. `DCActionButton`, `DIconButton`, `TableLayout`, `DialogLayout`, `MenuItem`, `Tooltip`, `TableActionButtons`, etc.) are intentionally omitted.

> **Legend**
>
> - 🏪 **Store-based** — all data read from Pinia stores; dialog can open without the parent passing anything
> - 🔗 **Parent-dependent** — needs a specific selected item (e.g. a clicked row) not held in a shared store; must receive data from the parent
> - ⚠️ **Partial** — create-mode is store-based, but edit/detail-mode is parent-dependent
> - ✅ **Done** — refactored in the latest commit

---

## Root: `ProjectsVersions.vue`

### Header / toolbar components

| Component     | Description                                                                  |
| ------------- | ---------------------------------------------------------------------------- |
| `ProjectMenu` | "More ▾" dropdown menu; slots used by parent to inject version-level actions |

Menu items slotted in from `ProjectsVersions.vue`:

| Slot action      | What it triggers                    |
| ---------------- | ----------------------------------- |
| "Overall Review" | Opens `OverallReviewDialog` at root |
| "Delete Version" | Opens `ConfirmationDialog` at root  |

### Dialogs at root level

| Dialog                | Trigger                 | Dependency                                                                                                   | Can use shared state?                                                                               |
| --------------------- | ----------------------- | ------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------- |
| `VersionDialogForm`   | Edit version button     | `version` is still passed via `open(config)`; `projectId` now read from `projectStore` directly              | ⚠️ Partial — `projectID` was removed from config, but edit mode still depends on a passed `version` |
| `ConfirmationDialog`  | "Delete Version" action | Config uses `versionDetails._key` / `.name` from store                                                       | 🏪 Yes                                                                                              |
| `OverallReviewDialog` | "Overall Review" action | Reads project key, version key, SBOM history, selected SBOM, approvable key — all from store inside `open()` | ✅ 🏪 Yes — `open()` now takes zero arguments                                                       |

---

## `ProjectMenu`

> Source: `components/projects/ProjectMenu.vue`

### Dialogs inside `ProjectMenu`

| Dialog                   | Dependency                                                | Can use shared state? |
| ------------------------ | --------------------------------------------------------- | --------------------- |
| `RequestFOSSDD`          | Reads project / version / SBOM data from store internally | 🏪 Yes                |
| `RequestReview`          | Reads project / version data from store internally        | 🏪 Yes                |
| `RequestApproval`        | Reads project / version data from store internally        | 🏪 Yes                |
| `AddChildrenErrorDialog` | Reads current project from store                          | 🏪 Yes                |
| `ConfirmationDialog`     | Config set reactively on `currentProject.deptMissing`     | 🏪 Yes                |

---

## Tab Components

### `TabOverview`

> Source: `components/projects/projectsVersions/charts/TabOverview.vue`

| Component     | Description                                                      |
| ------------- | ---------------------------------------------------------------- |
| `ChartHeader` | Chart title + navigation link (local sub-component in `charts/`) |
| `Bar`         | Bar chart (`vue-chartjs`)                                        |
| `Doughnut`    | Doughnut chart (`vue-chartjs`)                                   |

**No dialogs.** Chart clicks navigate to another tab via router.

**Data source**: ✅ 🏪 Fully store-based — stats now fetched via `sbomStore.fetchSBOMStats()` and `sbomStore.fetchGeneralVersionStats()`; results cached in `sbomStore.sbomStats` / `sbomStore.generalStats`. Direct `versionService` calls removed.

---

### `TabComponentList`

> Source: `components/projects/projectsVersions/TabComponentList.vue`

| Component            | Description                                        |
| -------------------- | -------------------------------------------------- |
| `PolicyDecisionCell` | Renders the policy decision icon in each table row |

#### Dialogs inside `TabComponentList`

| Dialog                      | Trigger                                                 | Dependency                                                                                       | Can use shared state?                             |
| --------------------------- | ------------------------------------------------------- | ------------------------------------------------------------------------------------------------ | ------------------------------------------------- |
| `ComponentDetailsDialog`    | Row click                                               | Needs fetched `ComponentDetails` plus row-derived policy/license context assembled by the parent | 🔗 Parent-dependent                               |
| `LicenseRuleDialog`         | Click license-rule icon in `showLicenseDecision` column | Needs `ComponentInfoSlim` + license ID from that row                                             | 🔗 Parent-dependent                               |
| `PolicyDecisionDialog`      | `PolicyDecisionCell` emits `open-policy-decision`       | Needs `ComponentInfoSlim` + `PolicyRuleStatus` from that row                                     | 🔗 Parent-dependent                               |
| `BulkPolicyDecisionsDialog` | "Bulk Policy Decision" button                           | Needs list of `DialogBulkPolicyDecisionEntry` computed from the filtered table                   | ⚠️ Partial — list could be moved to a store slice |

---

### `TabSBOMHistory`

> Source: `components/projects/projectsVersions/TabSBOMHistory.vue`

Thin wrapper — renders a single child:

```
TabSBOMHistory
  └── GridSBOM  [channelView=true]
```

#### `GridSBOM`

> Source: `components/grids/GridSBOM.vue`

No portal-specific UI components outside of dialogs.

#### Dialogs inside `GridSBOM`

| Dialog                       | Trigger                      | Dependency                                         | Can use shared state? |
| ---------------------------- | ---------------------------- | -------------------------------------------------- | --------------------- |
| `ReviewRemarkDialog`         | "Add Remark" row action      | Needs selected `SpdxFile`                          | 🔗 Parent-dependent   |
| `ConfirmationDialog`         | "Delete" row action          | Needs selected `SpdxFile`                          | 🔗 Parent-dependent   |
| `SbomValidationErrorsDialog` | SBOMs with validation errors | Needs validation error data from the selected SBOM | 🔗 Parent-dependent   |

> **Note**: `GridSBOM.reloadSboms()` now delegates to `sbomStore.fetchSBOMHistory()` (✅ done). `SpdxTagDialog` (portal-specific duplicate) was deleted; `DSpdxTagDialog` from `@shared` is used instead (✅ done).

---

### `TabSBOMCompare`

> Source: `components/projects/projectsVersions/TabSBOMCompare.vue`

No portal-specific UI components outside of dialogs.

#### Dialogs inside `TabSBOMCompare`

| Dialog                   | Trigger                           | Dependency                                   | Can use shared state? |
| ------------------------ | --------------------------------- | -------------------------------------------- | --------------------- |
| `ComponentCompareDialog` | Row click on a component diff row | Needs `ComponentMultiDiff` data for that row | 🔗 Parent-dependent   |

---

### `TabSBOMQualityMain`

> Source: `components/projects/projectsVersions/TabSBOMQualityMain.vue`

Container with four inner sub-tabs loaded dynamically via `defineAsyncComponent`.

| Sub-tab ID       | Component                                | Description               |
| ---------------- | ---------------------------------------- | ------------------------- |
| `scanRemarks`    | `TabScanRemarks`                         | Scan remarks table        |
| `licenseRemarks` | `TabLicenseRemarks`                      | License remarks table     |
| `reviewRemarks`  | `TabReviewRemarks` → `GridReviewRemarks` | Review remarks management |
| `generalRemarks` | `TabGeneralRemarks`                      | Static i18n text only     |

#### `TabScanRemarks`

> Source: `sbom-quality/TabScanRemarks.vue`

No portal-specific components or dialogs. Row clicks navigate to the Component tab via router.
**Data source**: 🏪 Store-based.

---

#### `TabLicenseRemarks`

> Source: `sbom-quality/TabLicenseRemarks.vue`

No portal-specific components or dialogs. Expandable rows with obligation details.
**Data source**: 🏪 Store-based.

---

#### `TabReviewRemarks` → `GridReviewRemarks`

> Source: `sbom-quality/TabReviewRemarks.vue` (thin wrapper)
> Grid: `components/grids/GridReviewRemarks.vue`

No portal-specific UI components outside of dialogs.

#### Dialogs inside `GridReviewRemarks`

| Dialog                             | Trigger                                                                             | Dependency                                                     | Can use shared state?                                         |
| ---------------------------------- | ----------------------------------------------------------------------------------- | -------------------------------------------------------------- | ------------------------------------------------------------- |
| `ReviewRemarkDialog`               | "Edit" row action                                                                   | Needs selected `ReviewRemark` for pre-fill                     | 🔗 Parent-dependent                                           |
| `ConfirmationDialog` (close)       | "Close" row action                                                                  | Needs selected remark                                          | 🔗 Parent-dependent                                           |
| `ConfirmationDialog` (cancel)      | "Cancel" row action                                                                 | Needs selected remark                                          | 🔗 Parent-dependent                                           |
| `ConfirmationDialog` (reopen)      | "Reopen" row action                                                                 | Needs selected remark                                          | 🔗 Parent-dependent                                           |
| `ConfirmationDialog` (bulk close)  | "Close" bulk button (selected rows)                                                 | Needs selected remark list                                     | 🔗 Parent-dependent                                           |
| `ConfirmationDialog` (bulk cancel) | "Cancel" bulk button (selected rows)                                                | Needs selected remark list                                     | 🔗 Parent-dependent                                           |
| `ReviewRemarksDetailsDialog`       | "View" row action                                                                   | Needs selected `ReviewRemark`; project/version UUID from store | 🔗 Parent-dependent (remark); 🏪 project/version from store   |
| `ChecklistExecuteDialog`           | Dedicated **"Checklist" toolbar button** — passes `lists` (fetched checklist array) | Needs available checklist data loaded on mount                 | 🔗 Parent-dependent (checklist list passed via `open(lists)`) |

> ⚠️ **Correction from initial map**: `ChecklistExecuteDialog` is **not** triggered from `ReviewRemarksDetailsDialog` bulk-close. It is triggered directly by the "Checklist" toolbar button in `GridReviewRemarks`, which calls `executeDialog?.open(lists)`. The `@close-remark` event from `ReviewRemarksDetailsDialog` calls `openBulkCloseDialog()`, which opens a `ConfirmationDialog`.

---

#### `TabGeneralRemarks`

> Source: `sbom-quality/TabGeneralRemarks.vue`

Static i18n text only. No components, no dialogs.

---

### `TabSourceCode`

> Source: `components/projects/projectsVersions/TabSourceCode.vue`

No portal-specific UI components outside of dialogs.

#### Dialogs inside `TabSourceCode`

| Dialog                    | Trigger                                          | Dependency                                                    | Can use shared state? |
| ------------------------- | ------------------------------------------------ | ------------------------------------------------------------- | --------------------- |
| `ConfirmationDialog`      | "Delete" row action                              | Needs selected source `_key` + URL                            | 🔗 Parent-dependent   |
| `NewExternalSourceDialog` | "Add" button (create) / "Edit" row action (edit) | Create: no item needed; Edit: needs selected `ExternalSource` | ⚠️ Partial            |

---

### `TabOverallReviews`

> Source: `components/projects/projectsVersions/TabOverallReviews.vue`

No portal-specific UI components outside of dialogs.

#### Dialogs inside `TabOverallReviews`

| Dialog                | Trigger                     | Dependency                                                                                                        | Can use shared state?                         |
| --------------------- | --------------------------- | ----------------------------------------------------------------------------------------------------------------- | --------------------------------------------- |
| `OverallReviewDialog` | "Add Overall Review" button | Reads project key, version key, SBOM history, selected SBOM, approvable SBOM key — all from store inside `open()` | ✅ 🏪 Yes — `open()` now takes zero arguments |
| `OverallAuditDialog`  | "Add Audit" button          | Same as above                                                                                                     | ✅ 🏪 Yes — `open()` now takes zero arguments |

> **Note**: `OverallReviewDialog` is instantiated **twice** in the component tree:
>
> 1. At root (`ProjectsVersions.vue`) — via the "Overall Review" slot action in `ProjectMenu`
> 2. Inside `TabOverallReviews` — via the "Add" button within the tab
>
> Since `open()` now reads entirely from the store, the root-level instance is redundant and can be eliminated (see TODO section).

---

### `TabNoticeFile`

> Source: `components/projects/projectsVersions/TabNoticeFile.vue`

| Component     | Description                                          |
| ------------- | ---------------------------------------------------- |
| `JsonViewer3` | JSON preview renderer (`vue-json-viewer`, 3rd-party) |

#### Dialogs inside `TabNoticeFile`

| Dialog            | Trigger                         | Dependency                                                                              | Can use shared state?                   |
| ----------------- | ------------------------------- | --------------------------------------------------------------------------------------- | --------------------------------------- |
| `ProjectSettings` | "Edit 3rd Party Address" button | Calls `showDialog(currentProject.value)` — `currentProject` sourced from `projectStore` | 🏪 Yes — could read from store directly |

---

### `TabAuditLog`

> Source: `components/projects/projectsVersions/TabAuditLog.vue`

Thin wrapper — renders `GridAuditLog` with a `fetchMethod` prop. No portal-specific UI components or dialogs.

---

## Dialog Detail Reference

### `ComponentDetailsDialog`

> Source: `components/dialog/ComponentDetailsDialog.vue`
> Used in: `TabComponentList`

🔗 **Parent-dependent** — needs fetched `ComponentDetails` plus additional row-derived policy/license context passed via `open()`. Not held in any shared store.

| Component              | Description                                    |
| ---------------------- | ---------------------------------------------- |
| `PolicyStatusTableRow` | Renders a policy status row (portal component) |
| `JsonViewer3`          | Raw SPDX data tab (`vue-json-viewer`)          |

#### Nested dialogs inside `ComponentDetailsDialog`

| Dialog                                  | Dependency                                                  | Can use shared state? |
| --------------------------------------- | ----------------------------------------------------------- | --------------------- |
| `ReviewRemarkDialog`                    | Needs selected `ReviewRemark` from the dialog's remark list | 🔗 Parent-dependent   |
| `LicenseRuleDialog`                     | Needs component + license from dialog context               | 🔗 Parent-dependent   |
| `PolicyDecisionDialog`                  | Needs `PolicyRuleStatus` from the policy status table row   | 🔗 Parent-dependent   |
| `ReviewRemarksDetailsDialog`            | Needs selected `ReviewRemark` from the dialog's remark list | 🔗 Parent-dependent   |
| `ConfirmationDialog` (close remark)     | Needs selected remark                                       | 🔗 Parent-dependent   |
| `ConfirmationDialog` (cancel remark)    | Needs selected remark                                       | 🔗 Parent-dependent   |
| `ConfirmationDialog` (mark in-progress) | Needs selected remark                                       | 🔗 Parent-dependent   |

---

## Full Tree

```
ProjectsVersions.vue
├── [Header]
│   └── ProjectMenu
│       ├── RequestFOSSDD             🏪
│       ├── RequestReview             🏪
│       ├── RequestApproval           🏪
│       ├── AddChildrenErrorDialog    🏪
│       └── ConfirmationDialog        🏪
│
├── [Dialogs at root]
│   ├── VersionDialogForm             ⚠️ (projectID removed; edit mode still receives a version argument)
│   ├── ConfirmationDialog            🏪 (delete version)
│   └── OverallReviewDialog           ✅ 🏪 (open() reads store; duplicate of TabOverallReviews instance)
│
└── [v-tabs]
    ├── TabOverview                    ✅ 🏪 (stats fetched/cached via sbomStore)
    │   ├── ChartHeader
    │   ├── Bar                        (vue-chartjs)
    │   └── Doughnut                   (vue-chartjs)
    │
    ├── TabComponentList
    │   ├── PolicyDecisionCell         🔗 (per-row)
    │   ├── ComponentDetailsDialog     🔗 (selected row)
    │   │   ├── PolicyStatusTableRow
    │   │   ├── JsonViewer3            (vue-json-viewer)
    │   │   ├── ReviewRemarkDialog     🔗
    │   │   ├── LicenseRuleDialog      🔗
    │   │   ├── PolicyDecisionDialog   🔗
    │   │   ├── ReviewRemarksDetailsDialog 🔗
    │   │   └── ConfirmationDialog × 3 🔗
    │   ├── LicenseRuleDialog          🔗 (selected row)
    │   ├── PolicyDecisionDialog       🔗 (selected row)
    │   └── BulkPolicyDecisionsDialog  ⚠️
    │
    ├── TabSBOMHistory
    │   └── GridSBOM [channelView]     ✅ (reloadSboms → sbomStore.fetchSBOMHistory)
    │       ├── ReviewRemarkDialog     🔗
    │       ├── ConfirmationDialog     🔗
    │       └── SbomValidationErrorsDialog 🔗
    │
    ├── TabSBOMCompare
    │   └── ComponentCompareDialog     🔗 (selected row)
    │
    ├── TabSBOMQualityMain
    │   ├── TabScanRemarks             🏪 (no dialogs)
    │   ├── TabLicenseRemarks          🏪 (no dialogs)
    │   ├── TabReviewRemarks
    │   │   └── GridReviewRemarks
    │   │       ├── ReviewRemarkDialog 🔗
    │   │       ├── ConfirmationDialog × 5 🔗
    │   │       ├── ReviewRemarksDetailsDialog 🔗
    │   │       └── ChecklistExecuteDialog 🔗 ⚠️ triggered by ReviewRemarks toolbar btn, not by ReviewRemarksDetailsDialog
    │   └── TabGeneralRemarks          (static text only)
    │
    ├── TabSourceCode
    │   ├── ConfirmationDialog         🔗
    │   └── NewExternalSourceDialog    ⚠️
    │
    ├── TabOverallReviews
    │   ├── OverallReviewDialog        ✅ 🏪 (open() reads store)
    │   └── OverallAuditDialog         ✅ 🏪 (open() reads store)
    │
    ├── TabNoticeFile
    │   ├── JsonViewer3                (vue-json-viewer)
    │   └── ProjectSettings            🏪
    │
    └── TabAuditLog
        └── GridAuditLog [fetchMethod]  (no portal-specific components)
```

---

## Summary: Which dialogs can replace prop-passing with shared state?

| Dialog                                                | Current pattern                                                                 | Status                                                                                                                                                               |
| ----------------------------------------------------- | ------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `VersionDialogForm`                                   | `open(config)` — `projectID` removed ✅; `version` still passed for edit mode   | ⚠️ Partial — `ProjectsVersions.vue` could use `sbomStore.currentVersion`, but `GridVersions.vue` still opens create/edit for arbitrary rows                          |
| `OverallReviewDialog` (×2)                            | `open()` — reads entirely from store                                            | ✅ Done; root-level duplicate instance can be removed                                                                                                                |
| `OverallAuditDialog`                                  | `open()` — reads entirely from store                                            | ✅ Done                                                                                                                                                              |
| `ProjectSettings`                                     | `showDialog(currentProject)` — arg is from store                                | ⚠️ Could read from store directly                                                                                                                                    |
| `RequestFOSSDD` / `RequestReview` / `RequestApproval` | No args — reads store internally                                                | 🏪 Already store-based                                                                                                                                               |
| `BulkPolicyDecisionsDialog`                           | Receives computed list from parent                                              | ⚠️ List could be held in a dedicated store slice, but the current transient parent-computed list is also coherent                                                    |
| `NewExternalSourceDialog` (create)                    | No item needed                                                                  | 🏪 Store-based for create mode                                                                                                                                       |
| `NewExternalSourceDialog` (edit)                      | Receives selected `ExternalSource`                                              | 🔗 Needs item; store-based only if selected item is stored                                                                                                           |
| `ComponentDetailsDialog`                              | Receives fetched `ComponentDetails` plus policy/license context from the parent | 🔗 Must remain parent-dependent unless both the selected row and the fetched detail payload move into shared state                                                   |
| `LicenseRuleDialog`                                   | Receives component + license                                                    | 🔗 Must remain parent-dependent                                                                                                                                      |
| `PolicyDecisionDialog`                                | Receives component + policy                                                     | 🔗 Must remain parent-dependent                                                                                                                                      |
| `ComponentCompareDialog`                              | Receives diff data for row                                                      | 🔗 Must remain parent-dependent                                                                                                                                      |
| `ReviewRemarkDialog` (edit)                           | Receives selected remark; SBOM history now from store cache ✅                  | 🔗 Remark item must remain parent-dependent                                                                                                                          |
| `ReviewRemarksDetailsDialog`                          | Receives selected remark                                                        | 🔗 Must remain parent-dependent                                                                                                                                      |
| `SbomValidationErrorsDialog`                          | Receives SBOM validation errors                                                 | 🔗 Must remain parent-dependent                                                                                                                                      |
| `ChecklistExecuteDialog`                              | Receives `lists` from `GridReviewRemarks` toolbar button                        | ⚠️ Lists could move to shared state, but `checklistsStore` would need new project-scoped/applicable-checklist state because it currently models admin checklist CRUD |
| All `ConfirmationDialog` instances                    | Config set by parent per-action                                                 | 🔗 Must remain parent-dependent                                                                                                                                      |

---

## TODO: Next Steps for Clean Architecture

> These items are ordered roughly by impact vs effort. Items marked ✅ were addressed in the latest commit.

### ✅ Already done (latest commit — `feat/sbomrefactor`)

- **`sbomStore` centralization**: `fetchSBOMHistory`, `fetchSBOMStats`, `fetchGeneralVersionStats` moved into `sbomStore`. Stats are now cached and invalidated reactively on version/SBOM selection change.
- **`OverallReviewDialog` + `OverallAuditDialog`**: `open()` signatures changed to zero args; both now read directly from `projectStore` / `sbomStore`.
- **`VersionDialogForm`**: `projectID` removed from the config argument; `projectStore.currentProject._key` is read internally. `version` is still passed for edit mode.
- **`ReviewRemarkDialog`**: Uses cached `sbomStore.channelSpdxs` instead of fetching SBOM history again when the dialog is for the currently-active version.
- **`GridSBOM`**: `reloadSboms()` delegates to `sbomStore.fetchSBOMHistory()`.
- **`TabOverview`**: Stats fetching delegated to `sbomStore`; local reactive vars removed.
- **`SpdxTagDialog` (portal-level)**: Deleted — replaced by `DSpdxTagDialog` from `@shared`.

---

### 🔜 Immediate next steps

#### 1. Eliminate the `OverallReviewDialog` double-instantiation at root

**What**: `ProjectsVersions.vue` holds a root-level `<OverallReviewDialog ref="reviewDia">` just for the "Overall Review" slot action in `ProjectMenu`. `TabOverallReviews` holds its own instance.

**Action**: Remove the root-level instance and the `reviewDia` ref from `ProjectsVersions.vue`. Instead, have `TabOverallReviews` expose an `open()` method (or use a Pinia flag in `sbomStore` / a dedicated `dialogStore` entry) so the "Overall Review" menu action can trigger the tab-level dialog — e.g. by emitting an event or setting a `dialogStore.isOverallReviewOpen` flag.

**Benefit**: Eliminates duplicate component instantiation and the confusing slot-based trigger pattern.

---

#### 2. Narrow `VersionDialogForm` parent input where safe

**What**: `VersionDialogForm.open(config)` still receives `config.version` from the parent. In `ProjectsVersions.vue`, that version is already available as `sbomStore.currentVersion`, but `GridVersions.vue` also uses the same dialog for create mode and for editing arbitrary rows.

**Action**: If the goal is only to simplify `ProjectsVersions.vue`, add a store-backed edit path for the current version there. Do **not** assume `open()` can simply become zero-argument globally without redesigning the create/edit API, because a no-arg `open()` is ambiguous and `GridVersions.vue` still needs to edit arbitrary rows. Keep `DialogVersionFormConfig` unless all call sites can be migrated cleanly.

**Benefit**: Reduces parent coupling for the current-version edit flow without breaking the existing create/edit usage in `GridVersions.vue`.

---

#### 3. Move `ChecklistExecuteDialog` lists to `checklistsStore`

**What**: `GridReviewRemarks` fetches and holds the available checklists in a local `lists` ref, then passes them to `ChecklistExecuteDialog.open(lists)`.

**Action**: Move checklist fetching into shared state only if `checklistsStore` is extended with project-scoped/applicable-checklist state. The store already exists at `stores/checklists.store.ts`, but today it is focused on admin checklist CRUD (`adminService.getChecklist()`), not `ProjectService.getApplicableChecklists(projectKey)`. `ChecklistExecuteDialog.open()` could then read from that new store state.

**Benefit**: Removes an API call from the component; checklist data becomes reusable across contexts.

---

#### 4. Move `BulkPolicyDecisionsDialog` list to a store slice

**What**: `TabComponentList` builds a `DialogBulkPolicyDecisionEntry[]` list from filtered table items and passes it to `BulkPolicyDecisionsDialog`.

**Action**: Add a `selectedComponentsForBulk` slice to `sbomStore` (or a new `policyDecisionStore`). `TabComponentList` writes to it; `BulkPolicyDecisionsDialog` reads from it.

**Benefit**: Removes the ⚠️ Partial dependency, making both components fully decoupled.

---

#### 5. Introduce a `selectedComponent` store slice for `ComponentDetailsDialog`

**What**: `ComponentDetailsDialog` receives the full `ComponentInfo` object on every row click. All nested dialogs (`ReviewRemarkDialog`, `LicenseRuleDialog`, `PolicyDecisionDialog`, etc.) receive further slices of that object.

**Action**: If this is pursued, store more than the clicked row. `TabComponentList` currently fetches `ComponentDetails` and also passes row-derived policy/license context into `ComponentDetailsDialog.open(...)`. A `selectedComponent` slice alone would not remove the dependency chain; the fetched detail payload (or a store action that owns that fetch) would also need to move into shared state.

**Benefit**: Can reduce prop drilling and help route-driven component opening, but only if shared state has a single clear source of truth for both the selected row and the fetched detail payload.

---

### 📁 File & folder structure refactoring

#### 6. Co-locate SBOM-domain components under a `sbom/` module

**Current state**: SBOM-related code is scattered across:

- `components/grids/GridSBOM.vue` (692 lines)
- `components/grids/GridReviewRemarks.vue` (716 lines)
- `components/projects/projectsVersions/Tab*.vue`
- `components/projects/projectsVersions/sbom-quality/Tab*.vue`
- `components/dialog/OverallReviewDialog.vue`, `OverallAuditDialog.vue`
- `components/dialog/ComponentDetailsDialog.vue` (863 lines)
- `components/dialog/ReviewRemarkDialog.vue`
- `components/dialog/SbomValidationErrorsDialog.vue`

**Proposed target structure**:

```
libs/portal/
  sbom/
    store/
      sbom.store.ts           ← already exists; move here
      sbom.store.test.ts
    grids/
      GridSBOM.vue
      GridReviewRemarks.vue
    dialogs/
      OverallReviewDialog.vue
      OverallAuditDialog.vue
      ReviewRemarkDialog.vue
      ReviewRemarksDetailsDialog.vue
      SbomValidationErrorsDialog.vue
      ComponentDetailsDialog.vue
      ComponentCompareDialog.vue
    tabs/
      TabOverview.vue
      TabComponentList.vue
      TabSBOMHistory.vue
      TabSBOMCompare.vue
      TabSBOMQualityMain.vue
      TabSourceCode.vue
      TabOverallReviews.vue
      TabNoticeFile.vue
      TabAuditLog.vue
      sbom-quality/
        TabScanRemarks.vue
        TabLicenseRemarks.vue
        TabReviewRemarks.vue
        TabGeneralRemarks.vue
    charts/
      TabOverview/ (or keep co-located)
```

**Action**: Create the `sbom/` module folder, move files, update import paths, and update `frontend/apps/portal/vite.config.mts` so `unplugin-vue-components` scans the new `libs/portal/sbom/**` directories. `components.d.ts` is generated, so it should be regenerated rather than hand-maintained. Run `tsc` to confirm no broken imports.

---

#### 7. Split oversized single-file components

| File                         | Lines | Suggested split                                                           |
| ---------------------------- | ----- | ------------------------------------------------------------------------- |
| `TabSBOMCompare.vue`         | 999   | Extract `CompareToolbar`, `CompareFilters`, `CompareTable` sub-components |
| `ComponentDetailsDialog.vue` | 863   | Extract `ComponentReviewRemarksPanel`, `ComponentPolicyStatusPanel`       |
| `TabComponentList.vue`       | 782   | Extract `ComponentListFilters`, `ComponentListTable`                      |
| `GridReviewRemarks.vue`      | 716   | Extract `ReviewRemarksToolbar`, `ReviewRemarksTable`                      |
| `GridSBOM.vue`               | 692   | Extract `SbomUploadSection`, `SbomTable`                                  |

---

#### 9. Add component-level tests for the refactored dialogs

**What**: `OverallReviewDialog`, `OverallAuditDialog`, and `VersionDialogForm` now depend on store state instead of props. Existing prop-based tests (if any) need to be updated; new store-mocked tests added.

---

### 🔧 Code quality

#### 10. Migrate `OverallReviewDialog` and `OverallAuditDialog` from Options API to Composition API

**What**: Both dialogs use `defineComponent` with Options-API-style `setup()`. All surrounding components use `<script setup>`.

**Action**: Rewrite using `<script setup lang="ts">` for consistency and to enable better tree-shaking and type inference.

---

#### 11. Standardize SBOM data reload patterns

**What**: Some components call `sbomStore.fetchSBOMHistory()` (✅ correct), while others still call `versionService.getSbomHistory()` directly (e.g. in `ReviewRemarkDialog`'s fallback branch for non-active versions).

**Action**: Audit all direct `versionService` calls within SBOM components. Where the version matches `sbomStore.currentVersion`, always use the store action. Document the intentional fallback for cross-version scenarios.

---

#### 12. Remove stale documentation that still mentions `DialogVersionFormConfig.projectID`

**What**: `DialogVersionFormConfig` in `DialogConfigs.ts` no longer contains a `projectID` field; that cleanup is already done.

**Action**: Remove any remaining comments, notes, or follow-up tasks that still describe `projectID` as if it exists. Keep `DialogVersionFormConfig` itself unless step 2 fully removes the remaining `version` argument across all call sites.
