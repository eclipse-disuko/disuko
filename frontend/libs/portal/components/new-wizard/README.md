# New Wizard - Business Rules & Edge Cases

## Architecture Overview

```
NewWizard.vue (Container)
    ├── WizardLayout.vue (Contains logic for rendering)
    │   ├── Stepper.vue
    │   │   └── Step.vue
    │   ├── Dynamic Step Components
    │   │   ├── WizardStepPlatform
    │   │   ├── WizardStepDetails
    │   │   ├── WizardStepArchitecture
    │   │   ├── WizardStepTargetUsers
    │   │   ├── WizardStepDistributionTarget
    │   │   ├── WizardStepDevelopment
    │   │   ├── WizardStepOwner
    │   │   ├── WizardStepDeveloper
    │   │   └── WizardStepSummary
    │   └── Navigation Buttons
    │
    └── State Management
        ├── wizard.store.ts (Centralized logic)
        ├── useNewWizard.ts (Configuration & Reusables between Pinia and Vue)
        └── NewWizard.ts (Models)
```

### Key Files & Responsibilities

**`wizard.store.ts`** - Central state management

- Manages wizard state (open/closed, current step, mode)
- Stores project data being created/edited
- Handles step navigation and validation
- Dynamically adjusts available steps based on selections
- Watches project changes to mark steps complete/incomplete

**`useNewWizard.ts`** - Configuration layer

- Defines step configurations with IDs and i18n keys
- Builds step arrays for different project types (Enterprise/Product/Other)
- Provides validation rules for form fields
- Maps step IDs to their Vue components for dynamic rendering
- Utility functions for step management (merge, remove, build)

**`NewWizard.ts`** - Type definitions and models

- Defines all TypeScript types and interfaces
- Step types, platform types, architecture types, etc.
- Project model structure for wizard
- Constants for step IDs, platforms, architectures, etc.

**Step Components** - Individual wizard pages

- Each component handles one step of the wizard
- Validates and auto-advances when complete
- Uses wizard store to read/write project data
- Selection-based steps use WizardCard.vue
- Form-based steps use OwnerSettings/DeveloperSettings components

## Dynamic Step Flow

### Platform Step Variations

**Enterprise IT / Mobile:**

1. `Platform → 2. Details → 3. Architecture → 4. Target Users → 5. Distribution Target → 6. Development → 7. Owner → 8. Developer → 9. Summary`

**Product:**

1. `Platform → 2. Details → 3. Architecture → 4. Development → 5. Owner → 6. Developer → 7. Summary`

- Skips: Target Users, Distribution Target

**Other:**

1. `Platform → 2. Details → 3. Development → 4. Owner → 5. Developer → 6. Summary`

- Skips: Architecture, Target Users, Distribution Target

### Development Step Variations

**In-house Development:**

- Developer step is removed from the flow
- Owner department data auto-fills developer fields
- Developer is company/department

**Internal Developer:**

- Supplier marked as non-external
- Developer is company/department

**External Developer:**

- Supplier marked as external
- Only developer name is required instead of a company/department

## Edge Cases & Special Behaviors

### Application Selector Visibility (Details Step)

**Shown when:**

- Platform is Enterprise IT OR Mobile
- AND Application Connector capability is enabled

### Architecture Options

**Enterprise IT / Mobile platforms show:**

- Frontend or Client
- Backend

**Product platform shows:**

- Product onboard
- Product offboard

### Distribution Target Pre-selection

**Only applies when Architecture is Frontend:**

- **Company** target user → Auto-selects **Company** distribution target, skips Distribution Target step
- **Business Partner** target user → Auto-selects **Business Partner** distribution target, skips Distribution Target step
- **End Customer** target user → Auto-selects **Business Partner** distribution target, skips Distribution Target step

**When Architecture is NOT Frontend (Backend or Product architectures):**

- User proceeds to Distribution Target step for manual selection
- No auto-selection occurs

### Step Completion Logic

A step is marked complete when:

- **Platform:** Target platform is selected
- **Details:** Project name passes validation (3-80 chars)
- **Architecture:** Architecture option is selected
- **Target Users:** Target user type is selected
- **Distribution Target:** Distribution target is selected
- **Development:** Development type is selected
- **Owner:** Department selected AND both addresses (customer + notice contact) are valid
- **Developer:** Department selected OR (supplier name, address, and number are all valid)
- **Summary:** All previous steps are completed

## Navigation Rules

- Users can navigate to any previously seen step
- Cannot skip forward to unseen steps unless current step is complete
- Back button always available (except first step)
- Next button disabled if current step incomplete
- Summary step requires all previous steps complete
