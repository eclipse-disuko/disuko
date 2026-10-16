# Project Groups

## Overview

- <i aria-hidden="true" class="v-icon notranslate pr-1 material-icons labelIconColor--text mdi mdi-card-multiple-outline" style="font-size: 16px;"></i> A **Project Group** summarizes a list of projects (children)
- <i aria-hidden="true" class="v-icon notranslate pr-1 material-icons labelIconColor--text mdi mdi-account-supervisor-circle" style="font-size: 16px;"></i> A **Project Child** has a parent relation to the group
- <i aria-hidden="true" class="v-icon notranslate pr-1 material-icons labelIconColor--text mdi mdi-layers" style="font-size: 16px;"></i> A **Project SBOM** lists a child project's software components. The table shows:
  - **Policy Status:** Component counts: total, denied, unasserted, warned, questioned, and allowed.
  - **Supplier:** The child's supplier user ID.
  - **Has Approval:** Blue check: child project has an approval. Gray check: none.
- The user needs to be owner of the group **and** owner of the child to create a child relation
- A project can only be child of one single group
- Deleted projects can not be added as child to a group

## Tokens

- Token of the project group are passed on to children
- It is therefore possible to use one technical token for all children

## User Management

- The user management of a group is **not** passed on to children
- A user on a group only has rights on the group itself (exception: tokens)

## Deletion of Groups

- When a group is deleted, the parent entry for all children will be removed
- When a group is deleted, the child projects remain as stand-alone project

## Deletion of Children

- When a child project is deleted, the project is marked as deleted on the project group
- The group will then display the deleted project with a notifier that it is deleted
- The dialog he edit children shows the deleted project child, but **not** all deleted projects.

## Removal of Children

- When removing a child the parent attribute will be cleared and the group projects gets a new children list set
- Tokens on group level can no longer be used to access the removed child project
