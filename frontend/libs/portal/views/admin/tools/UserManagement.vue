<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {UserDto} from '@disclosure-portal/model/Users';
import {formatDate} from '@disclosure-portal/utils/View';
import DAutocompleteUser from '@shared/components/disco/DAutocompleteUser.vue';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DIconButton from '@shared/components/disco/DIconButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader} from '@shared/types/table';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import AdminService from '../../../services/admin';

const {t} = useI18n();
const snackbar = useSnackbar();

interface EntityDetail {
  entityID: string;
  entityType: string;
  entitySubType?: string;
  entityStatus?: string;
  entityName?: string;
  projectID?: string;
  projectName?: string;
  disableDeleteReason?: string;
}

const username = ref('');
const loading = ref(false);
const dryRunExecuted = ref(false);
const selectedUserDetails = ref<UserDto | null>(null);
const loadingUserDetails = ref(false);
const affectedEntities = ref<{
  user_tasks_count: number;
  user_roles_count: number;
  data_traces_count: number;
} | null>(null);

const tasksData = ref<EntityDetail[]>([]);
const rolesData = ref<EntityDetail[]>([]);
const logsData = ref<EntityDetail[]>([]);

const loadedPanels = ref<Set<string>>(new Set());

const tasksSearch = ref('');
const rolesSearch = ref('');
const logsSearch = ref('');

const selectedTasksStatusFilter = ref<string[]>([]);
const selectedTasksNameFilter = ref<string[]>([]);
const selectedRolesNameFilter = ref<string[]>([]);
const selectedLogsFilter = ref<string[]>([]);

const filteredTasksData = computed(() => {
  let filtered = tasksData.value;

  if (selectedTasksStatusFilter.value.length > 0) {
    filtered = filtered.filter(
      (task) => task.entityStatus && selectedTasksStatusFilter.value.includes(task.entityStatus),
    );
  }

  if (selectedTasksNameFilter.value.length > 0) {
    filtered = filtered.filter(
      (task) => task.entitySubType && selectedTasksNameFilter.value.includes(task.entitySubType),
    );
  }

  return filtered;
});

const filteredRolesData = computed(() => {
  if (selectedRolesNameFilter.value.length === 0) return rolesData.value;
  return rolesData.value.filter(
    (role) => role.entitySubType && selectedRolesNameFilter.value.includes(role.entitySubType),
  );
});

const filteredLogsData = computed(() => {
  if (selectedLogsFilter.value.length === 0) return logsData.value;
  return logsData.value.filter((log) => log.entityName && selectedLogsFilter.value.includes(log.entityName));
});

const tasksStatusOptions = computed(() => {
  const statuses = tasksData.value.map((task) => task.entityStatus).filter((status) => status);
  return [...new Set(statuses)].map((status) => ({text: status, value: status}));
});

const tasksNameOptions = computed(() => {
  const names = tasksData.value.map((task) => task.entitySubType).filter((n) => n);
  return [...new Set(names)].map((name) => ({text: name, value: name}));
});

const rolesNameOptions = computed(() => {
  const names = rolesData.value.map((role) => role.entitySubType).filter((n) => n);
  return [...new Set(names)].map((name) => ({text: name, value: name}));
});

const logsOptions = computed(() => {
  const logs = logsData.value.map((log) => log.entityName);
  return [...new Set(logs)].map((log) => ({
    text: t('LOG_TYPE_' + log) !== 'LOG_TYPE_' + log ? t('LOG_TYPE_' + log) : log,
    value: log,
  }));
});

const tasksHeaders = computed<DataTableHeader[]>(() => [
  {title: 'Task Type', align: 'start', value: 'entitySubType', sortable: true},
  {title: 'Project', align: 'start', value: 'projectName', sortable: true},
  {title: 'Status', align: 'start', value: 'entityStatus', width: '120', sortable: true},
  {title: t('USER_MANAGEMENT_DISABLE_DELETE_REASON'), align: 'start', value: 'disableDeleteReason', sortable: true},
  {title: t('COL_ACTIONS'), align: 'end', value: 'actions', width: '100', sortable: false},
]);

const rolesHeaders = computed<DataTableHeader[]>(() => [
  {title: 'Role', align: 'start', value: 'entityStatus', sortable: true},

  {title: 'Project', align: 'start', value: 'projectName', sortable: true},
  {title: t('USER_MANAGEMENT_DISABLE_DELETE_REASON'), align: 'start', value: 'disableDeleteReason', sortable: true},
  {title: t('COL_ACTIONS'), align: 'end', value: 'actions', width: '100', sortable: false},
]);

const logsHeaders = computed<DataTableHeader[]>(() => [
  {title: 'Log Type', align: 'start', value: 'entityName', sortable: true},
  {title: t('USER_MANAGEMENT_DISABLE_DELETE_REASON'), align: 'start', value: 'disableDeleteReason', sortable: true},
  {title: t('COL_ACTIONS'), align: 'end', value: 'actions', width: '100', sortable: false},
]);

const showConfirmDialog = ref(false);
const confirmDialogConfig = ref<{
  entityType: string;
  entityName: string;
  entityId?: string;
  isDeleteAll: boolean;
}>({
  entityType: '',
  entityName: '',
  entityId: undefined,
  isDeleteAll: false,
});

const onUserChanged = async (user: UserDto) => {
  username.value = user.user || '';
  dryRunExecuted.value = false;
  affectedEntities.value = null;
  tasksData.value = [];
  rolesData.value = [];
  logsData.value = [];
  loadedPanels.value.clear();

  if (user._key) {
    loadingUserDetails.value = true;
    try {
      const response = await AdminService.getUser(user._key);
      selectedUserDetails.value = response.data;
    } catch {
      selectedUserDetails.value = user;
    } finally {
      loadingUserDetails.value = false;
    }
  } else {
    selectedUserDetails.value = user;
  }
};

const resetForm = () => {
  username.value = '';
  dryRunExecuted.value = false;
  affectedEntities.value = null;
  tasksData.value = [];
  rolesData.value = [];
  logsData.value = [];
  loadedPanels.value.clear();
  selectedUserDetails.value = null;
};

const executeDryRun = async (silent: boolean = false) => {
  if (!username.value.trim()) {
    snackbar.error(t('USER_MANAGEMENT_USERNAME_REQUIRED'));
    return;
  }

  loading.value = true;

  try {
    const response = await AdminService.executeDryRun(username.value);
    if (response.data.success) {
      affectedEntities.value = response.data.entities_effected;
      dryRunExecuted.value = true;
      if (!silent) {
        snackbar.info(t(response.data.message));
      }
    } else {
      snackbar.error(t(response.data.message));
    }
  } catch (error) {
    const err = error as {response?: {data?: {message?: string}}};
    const errorMessage = err.response?.data?.message
      ? t(err.response.data.message)
      : t('USER_MANAGEMENT_DRY_RUN_ERROR');
    snackbar.error(errorMessage);
  } finally {
    loading.value = false;
  }
};

const onAccordionChange = async (values: unknown) => {
  const panelValues = Array.isArray(values) ? values : values ? [values] : [];

  for (const panelValue of panelValues) {
    if (panelValue && typeof panelValue === 'string' && !loadedPanels.value.has(panelValue)) {
      await loadEntityDetails(panelValue);
      loadedPanels.value.add(panelValue);
    }
  }
};

const loadEntityDetails = async (entity: string) => {
  if (!username.value.trim()) {
    return;
  }

  loading.value = true;

  try {
    const response = await AdminService.getPersonalDetails(username.value, entity);
    if (response.data.success) {
      switch (entity) {
        case 'tasks':
          tasksData.value = response.data.data;
          break;
        case 'roles':
          rolesData.value = response.data.data;
          break;
        case 'logs':
          logsData.value = response.data.data;
          break;
      }
    } else {
      snackbar.error(t(response.data.message));
    }
  } catch (error) {
    const err = error as {response?: {data?: {message?: string}}};
    const errorMessage = err.response?.data?.message ? t(err.response.data.message) : t('ERROR_LOADING_DATA');
    snackbar.error(errorMessage);
  } finally {
    loading.value = false;
  }
};

const openConfirmDialog = (entityType: string, entityName: string, entityId?: string, isDeleteAll: boolean = false) => {
  confirmDialogConfig.value = {
    entityType,
    entityName,
    entityId,
    isDeleteAll,
  };
  showConfirmDialog.value = true;
};

const executeEntityDeletion = async () => {
  showConfirmDialog.value = false;

  if (!username.value.trim()) {
    return;
  }

  loading.value = true;

  try {
    let response;

    if (confirmDialogConfig.value.entityId) {
      response = await AdminService.deletePersonalDataByEntityId(
        confirmDialogConfig.value.entityType,
        confirmDialogConfig.value.entityId,
      );
    } else if (confirmDialogConfig.value.isDeleteAll) {
      snackbar.error(t('FEATURE_NOT_IMPLEMENTED'));
      loading.value = false;
      return;
    } else {
      response = await AdminService.deletePersonalDataByEntity(username.value, confirmDialogConfig.value.entityType);
    }

    if (response.data.success) {
      // Show context-specific success message based on entity type
      let successMessage = '';
      if (confirmDialogConfig.value.entityId) {
        // Single entity deletion
        switch (confirmDialogConfig.value.entityType) {
          case 'tasks':
            successMessage = 'USER_MANAGEMENT_TASK_DELETED';
            break;
          case 'roles':
            successMessage = 'USER_MANAGEMENT_ROLE_DELETED';
            break;
          case 'logs':
            successMessage = 'USER_MANAGEMENT_LOG_DELETED';
            break;
          default:
            successMessage = 'USER_MANAGEMENT_ENTITY_DELETED';
        }
      } else {
        // Bulk entity deletion
        switch (confirmDialogConfig.value.entityType) {
          case 'tasks':
            successMessage = 'USER_MANAGEMENT_ALL_TASKS_DELETED';
            break;
          case 'roles':
            successMessage = 'USER_MANAGEMENT_ALL_ROLES_DELETED';
            break;
          case 'logs':
            successMessage = 'USER_MANAGEMENT_ALL_LOGS_DELETED';
            break;
          default:
            successMessage = response.data.message || 'USER_MANAGEMENT_ENTITY_DELETED';
        }
      }
      snackbar.info(t(successMessage));
      await executeDryRun(true);
      loadedPanels.value.delete(confirmDialogConfig.value.entityType);
    } else {
      snackbar.error(t(response.data.message));
    }
  } catch (error) {
    const err = error as {response?: {data?: {message?: string}}};
    const errorMessage = err.response?.data?.message ? t(err.response.data.message) : t('ERROR_DELETING_DATA');
    snackbar.error(errorMessage);
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <v-card class="pa-4 mb-3 h-auto border-none" flat color="transparent">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h5 mb-3 pr-2">{{ t('USER_MANAGEMENT_DRY_RUN_TITLE') }}</h1>
      </v-col>
    </v-row>

    <v-row align="start" no-gutters>
      <v-col cols="12" md="6" lg="4" class="pr-md-4">
        <p class="mb-4 text-xs">
          {{ t('USER_MANAGEMENT_DRY_RUN_DESCRIPTION') }}
        </p>

        <DAutocompleteUser
          v-model="username"
          :label="t('UM_DIALOG_USER_ID')"
          :disabled="loading || dryRunExecuted"
          :active="null"
          required
          @user-changed="onUserChanged" />

        <div class="d-flex align-center ga-2 mt-6 mb-4">
          <DCActionButton
            :loading="loading"
            :disabled="loading"
            :text="t('USER_MANAGEMENT_DRY_RUN_BUTTON')"
            @click="executeDryRun" />

          <DCActionButton
            v-if="dryRunExecuted"
            variant="outlined"
            class="ml-2"
            :text="t('BTN_RESET')"
            @click="resetForm" />
        </div>
      </v-col>

      <!-- User Details Panel -->
      <v-col cols="12" md="6" lg="8" v-if="selectedUserDetails">
        <v-card class="px-4 pt-0 pb-4" flat color="transparent">
          <div class="d-flex align-center mb-3">
            <v-icon class="mr-2" color="primary">mdi-account-circle</v-icon>
            <span class="text-h6">{{ t('USER_MANAGEMENT_USER_DETAILS') }}</span>
          </div>
          <v-progress-linear v-if="loadingUserDetails" indeterminate class="mb-3"></v-progress-linear>
          <v-row v-else dense>
            <v-col cols="12" sm="6" md="4">
              <span class="d-text d-secondary-text">{{ t('COL_USER_ID') }}</span
              ><br />
              <span class="font-weight-medium">{{ selectedUserDetails.user }}</span>
            </v-col>
            <v-col cols="12" sm="6" md="4">
              <span class="d-text d-secondary-text">{{ t('COL_FORENAME') }}</span
              ><br />
              <span>{{ selectedUserDetails.forename }}</span>
            </v-col>
            <v-col cols="12" sm="6" md="4">
              <span class="d-text d-secondary-text">{{ t('COL_LASTNAME') }}</span
              ><br />
              <span>{{ selectedUserDetails.lastname }}</span>
            </v-col>
            <v-col cols="12" sm="6" md="4" class="mt-2">
              <span class="d-text d-secondary-text">{{ t('COL_EMAIL') }}</span
              ><br />
              <span>{{ selectedUserDetails.email }}</span>
            </v-col>
            <v-col cols="12" sm="6" md="4" class="mt-2">
              <span class="d-text d-secondary-text">{{ t('USER_STATUS') }}</span
              ><br />
              <v-chip :color="selectedUserDetails.active ? 'success' : 'warning'" size="small" variant="flat">
                {{ selectedUserDetails.active ? t('ICON_LABEL_TEXT_ACTIVE') : t('ICON_LABEL_TEXT_INACTIVE') }}
              </v-chip>
            </v-col>
            <v-col cols="12" sm="6" md="4" class="mt-2">
              <span class="d-text d-secondary-text">{{ t('USER_ACCESS_SCOPE') }}</span
              ><br />
              <span>{{
                selectedUserDetails.isInternal ? t('USER_ACCESS_SCOPE_INTERNAL') : t('USER_ACCESS_SCOPE_EXTERNAL')
              }}</span>
            </v-col>
            <v-col cols="12" sm="6" md="4" class="mt-2">
              <span class="d-text d-secondary-text">{{ t('COL_CREATED') }}</span
              ><br />
              <span>{{ formatDate(selectedUserDetails.created) }}</span>
            </v-col>
            <v-col cols="12" sm="6" md="4" class="mt-2">
              <span class="d-text d-secondary-text">{{ t('COL_UPDATED') }}</span
              ><br />
              <span>{{ formatDate(selectedUserDetails.updated) }}</span>
            </v-col>
            <v-col
              cols="12"
              sm="6"
              md="4"
              class="mt-2"
              v-if="selectedUserDetails.roles && selectedUserDetails.roles.length > 0">
              <span class="d-text d-secondary-text">{{ t('ROLES') }}</span
              ><br />
              <span v-for="(role, index) in selectedUserDetails.roles" :key="index">
                <v-chip size="x-small" class="mr-1 mb-1" variant="tonal">{{ t(role) }}</v-chip>
              </span>
            </v-col>
            <v-col cols="12" sm="6" md="4" class="mt-2" v-if="selectedUserDetails.metaData?.department">
              <span class="d-text d-secondary-text">{{ t('DEPARTMENT') }}</span
              ><br />
              <span>{{ selectedUserDetails.metaData.department }}</span>
            </v-col>
            <v-col cols="12" sm="6" md="4" class="mt-2" v-if="selectedUserDetails.metaData?.departmentDescription">
              <span class="d-text d-secondary-text">{{ t('DEPARTMENT_DESCRIPTION') }}</span
              ><br />
              <span>{{ selectedUserDetails.metaData.departmentDescription }}</span>
            </v-col>
            <v-col cols="12" sm="6" md="4" class="mt-2" v-if="selectedUserDetails.metaData?.companyIdentifier">
              <span class="d-text d-secondary-text">{{ t('COMPANY_IDENTIFIER') }}</span
              ><br />
              <span>{{ selectedUserDetails.metaData.companyIdentifier }}</span>
            </v-col>
            <v-col cols="12" sm="6" md="4" class="mt-2">
              <span class="d-text d-secondary-text">{{ t('DEPROVISIONED_DATE') }}</span
              ><br />
              <span>{{
                selectedUserDetails.deprovisioned ? formatDate(selectedUserDetails.deprovisioned) : t('NOT_SET')
              }}</span>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>

    <!-- Detailed Action Plan -->
    <v-row v-if="dryRunExecuted && affectedEntities">
      <v-col cols="12" xs="12">
        <v-card class="pa-3">
          <v-row>
            <v-col cols="12">
              <h2 class="pb-3">Detailed Action Plan</h2>
              <span class="caption">User: {{ username }}</span>
              <div class="text-caption text-medium-emphasis mt-2">
                <v-icon size="small" class="mr-1">mdi-information-outline</v-icon>
                Click on each section below to expand and view details
              </div>
            </v-col>
          </v-row>

          <v-divider class="my-3"></v-divider>

          <v-expansion-panels variant="accordion" class="elevation-0" @update:model-value="onAccordionChange">
            <!-- Tasks Section -->
            <v-expansion-panel class="mb-2 border" value="tasks">
              <v-expansion-panel-title>
                <template v-slot:default>
                  <v-icon class="mr-2" color="primary">mdi-chevron-down</v-icon>
                  <v-icon class="mr-2">mdi-checkbox-marked-circle-outline</v-icon>
                  <span class="font-weight-medium">
                    {{ t('USER_MANAGEMENT_ENTITY_TASKS') }} ({{ affectedEntities.user_tasks_count }})
                  </span>
                </template>
                <template v-slot:actions>
                  <DCActionButton
                    @click.stop="openConfirmDialog('tasks', t('USER_MANAGEMENT_ENTITY_TASKS'), undefined, false)"
                    :disabled="loading"
                    size="small"
                    variant="tonal"
                    icon="mdi-delete"
                    class="mr-2"
                    :text="t('USER_MANAGEMENT_DELETE_ENTITY')" />
                </template>
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <v-progress-linear v-if="loading" indeterminate></v-progress-linear>
                <div v-else-if="tasksData.length > 0" class="pa-2">
                  <div class="d-flex mb-3 justify-end">
                    <v-text-field
                      v-model="tasksSearch"
                      density="compact"
                      :label="t('labelSearch')"
                      append-inner-icon="mdi-magnify"
                      variant="outlined"
                      hide-details
                      style="max-width: 500px"
                      single-line></v-text-field>
                  </div>
                  <v-data-table
                    density="compact"
                    class="striped-table custom-data-table"
                    :search="tasksSearch"
                    :headers="tasksHeaders"
                    :items="filteredTasksData"
                    :items-per-page="10"
                    :footer-props="{'items-per-page-options': [5, 10, 25, 50]}"
                    item-key="entityID">
                    <template v-slot:[`header.entitySubType`]="{column, toggleSort, getSortIcon}">
                      <span class="mr-1">{{ column.title }}</span>
                      <GridHeaderFilterIcon
                        v-model="selectedTasksNameFilter"
                        :column="column"
                        :label="t('Lbl_filter_name')"
                        :allItems="tasksNameOptions" />
                      <v-icon
                        class="v-data-table-header__sort-icon"
                        :icon="getSortIcon(column)"
                        @click="toggleSort(column)" />
                    </template>
                    <template v-slot:[`header.entityStatus`]="{column, toggleSort, getSortIcon}">
                      <span class="mr-1">{{ column.title }}</span>
                      <GridHeaderFilterIcon
                        v-model="selectedTasksStatusFilter"
                        :column="column"
                        :label="t('Lbl_filter_status')"
                        :allItems="tasksStatusOptions" />
                      <v-icon
                        class="v-data-table-header__sort-icon"
                        :icon="getSortIcon(column)"
                        @click="toggleSort(column)" />
                    </template>
                    <template v-slot:[`item.projectName`]="{item}">
                      <a
                        v-if="item.projectID"
                        :href="`#/dashboard/projects/${item.projectID}/approvals`"
                        class="text-primary text-decoration-none">
                        {{ item.projectName }}
                      </a>
                      <span v-else>{{ item.projectName }}</span>
                    </template>
                    <template v-slot:[`item.entityStatus`]="{item}">
                      <span class="text-uppercase">{{ item.entityStatus }}</span>
                    </template>
                    <template v-slot:[`item.disableDeleteReason`]="{item}">
                      {{ item.disableDeleteReason ? t(item.disableDeleteReason) : '' }}
                    </template>
                    <template v-slot:[`item.actions`]="{item}">
                      <DIconButton
                        :disabled="!!item.disableDeleteReason"
                        icon="mdi-delete"
                        :hint="item.disableDeleteReason ? t(item.disableDeleteReason) : t('Delete')"
                        color="error"
                        @clicked="
                          openConfirmDialog('tasks', item.entitySubType || item.entityName || '', item.entityID, false)
                        " />
                    </template>
                  </v-data-table>
                </div>
                <div v-else class="pa-4 text-medium-emphasis text-center">No tasks found</div>
              </v-expansion-panel-text>
            </v-expansion-panel>

            <!-- Roles Section -->
            <v-expansion-panel class="mb-2 border" value="roles">
              <v-expansion-panel-title>
                <template v-slot:default>
                  <v-icon class="mr-2" color="primary">mdi-chevron-down</v-icon>
                  <v-icon class="mr-2">mdi-account-group</v-icon>
                  <span class="font-weight-medium">
                    {{ t('USER_MANAGEMENT_ENTITY_ROLES') }} ({{ affectedEntities.user_roles_count }})
                  </span>
                </template>
                <template v-slot:actions>
                  <DCActionButton
                    @click.stop="openConfirmDialog('roles', t('USER_MANAGEMENT_ENTITY_ROLES'), undefined, false)"
                    :disabled="loading"
                    size="small"
                    variant="tonal"
                    icon="mdi-delete"
                    class="mr-2"
                    :text="t('USER_MANAGEMENT_DELETE_ENTITY')" />
                </template>
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <v-progress-linear v-if="loading" indeterminate></v-progress-linear>
                <div v-else-if="rolesData.length > 0" class="pa-2">
                  <div class="d-flex mb-3 justify-end">
                    <v-text-field
                      v-model="rolesSearch"
                      density="compact"
                      :label="t('labelSearch')"
                      append-inner-icon="mdi-magnify"
                      variant="outlined"
                      hide-details
                      style="max-width: 500px"
                      single-line></v-text-field>
                  </div>
                  <v-data-table
                    density="compact"
                    class="striped-table custom-data-table"
                    :search="rolesSearch"
                    :headers="rolesHeaders"
                    :items="filteredRolesData"
                    :items-per-page="10"
                    :footer-props="{'items-per-page-options': [5, 10, 25, 50]}"
                    item-key="entityID">
                    <template v-slot:[`header.entitySubType`]="{column, toggleSort, getSortIcon}">
                      <span class="mr-1">{{ column.title }}</span>
                      <GridHeaderFilterIcon
                        v-model="selectedRolesNameFilter"
                        :column="column"
                        :label="t('Lbl_filter_role')"
                        :allItems="rolesNameOptions" />
                      <v-icon
                        class="v-data-table-header__sort-icon"
                        :icon="getSortIcon(column)"
                        @click="toggleSort(column)" />
                    </template>
                    <template v-slot:[`item.projectName`]="{item}">
                      <a
                        v-if="item.projectID"
                        :href="`#/dashboard/projects/${item.projectID}/users`"
                        class="text-primary text-decoration-none">
                        {{ item.projectName }}
                      </a>
                      <span v-else>{{ item.projectName }}</span>
                    </template>
                    <template v-slot:[`item.disableDeleteReason`]="{item}">
                      {{ item.disableDeleteReason ? t(item.disableDeleteReason) : '' }}
                    </template>
                    <template v-slot:[`item.actions`]="{item}">
                      <DIconButton
                        :disabled="!!item.disableDeleteReason"
                        icon="mdi-delete"
                        :hint="item.disableDeleteReason ? t(item.disableDeleteReason) : t('Delete')"
                        color="error"
                        @clicked="
                          openConfirmDialog('roles', item.entitySubType || item.entityName || '', item.entityID, false)
                        " />
                    </template>
                  </v-data-table>
                </div>
                <div v-else class="pa-4 text-medium-emphasis text-center">No roles found</div>
              </v-expansion-panel-text>
            </v-expansion-panel>

            <v-expansion-panel class="mb-2 border" value="logs">
              <v-expansion-panel-title>
                <template v-slot:default>
                  <v-icon class="mr-2" color="primary">mdi-chevron-down</v-icon>
                  <v-icon class="mr-2">mdi-database</v-icon>
                  <span class="font-weight-medium">
                    {{ t('USER_MANAGEMENT_ENTITY_TRACES') }} ({{ affectedEntities.data_traces_count }})
                  </span>
                </template>
                <template v-slot:actions>
                  <DCActionButton
                    @click.stop="openConfirmDialog('logs', t('USER_MANAGEMENT_ENTITY_TRACES'), undefined, false)"
                    :disabled="loading"
                    size="small"
                    variant="tonal"
                    icon="mdi-delete"
                    class="mr-2"
                    :text="t('USER_MANAGEMENT_DELETE_ENTITY')" />
                </template>
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <v-progress-linear v-if="loading" indeterminate></v-progress-linear>
                <div v-else-if="logsData.length > 0" class="pa-2">
                  <div class="d-flex mb-3 justify-end">
                    <v-text-field
                      v-model="logsSearch"
                      density="compact"
                      :label="t('labelSearch')"
                      append-inner-icon="mdi-magnify"
                      variant="outlined"
                      hide-details
                      style="max-width: 500px"
                      single-line></v-text-field>
                  </div>
                  <v-data-table
                    density="compact"
                    class="striped-table custom-data-table"
                    :search="logsSearch"
                    :headers="logsHeaders"
                    :items="filteredLogsData"
                    :items-per-page="10"
                    :footer-props="{'items-per-page-options': [5, 10, 25, 50]}"
                    item-key="entityID">
                    <template v-slot:[`header.entityName`]="{column, toggleSort, getSortIcon}">
                      <span class="mr-1">{{ column.title }}</span>
                      <GridHeaderFilterIcon
                        v-model="selectedLogsFilter"
                        :column="column"
                        :label="t('Lbl_filter_log_type')"
                        :allItems="logsOptions" />
                      <v-icon
                        class="v-data-table-header__sort-icon"
                        :icon="getSortIcon(column)"
                        @click="toggleSort(column)" />
                    </template>
                    <template v-slot:[`item.entityName`]="{item}">
                      {{
                        t('LOG_TYPE_' + item.entityName) !== 'LOG_TYPE_' + item.entityName
                          ? t('LOG_TYPE_' + item.entityName)
                          : item.entityName
                      }}
                    </template>
                    <template v-slot:[`item.disableDeleteReason`]="{item}">
                      {{ item.disableDeleteReason ? t(item.disableDeleteReason) : '' }}
                    </template>
                    <template v-slot:[`item.actions`]="{item}">
                      <DIconButton
                        :disabled="!!item.disableDeleteReason"
                        icon="mdi-delete"
                        :hint="item.disableDeleteReason ? t(item.disableDeleteReason) : t('Delete')"
                        color="error"
                        @clicked="openConfirmDialog('logs', item.entityName || '', item.entityID, false)" />
                    </template>
                  </v-data-table>
                </div>
                <div v-else class="pa-4 text-medium-emphasis text-center">No logs found</div>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-card>
      </v-col>
    </v-row>

    <!-- Delete All Section -->
    <v-row v-if="dryRunExecuted && affectedEntities" class="mt-3">
      <v-col cols="12" xs="12">
        <v-divider class="my-3"></v-divider>
        <DCActionButton
          @click="openConfirmDialog('all', '', undefined, true)"
          :disabled="
            loading ||
            affectedEntities.user_tasks_count > 0 ||
            affectedEntities.user_roles_count > 0 ||
            affectedEntities.data_traces_count > 0
          "
          large
          icon="mdi-delete"
          class="pa-3"
          :hint="
            t(
              affectedEntities.user_tasks_count <= 0 &&
                affectedEntities.user_roles_count <= 0 &&
                affectedEntities.data_traces_count <= 0
                ? 'USER_MANAGEMENT_CONFIRM_DELETE_ALL_MESSAGE'
                : 'USER_MANAGEMENT_CONFIRM_DELETE_NOT_ALLOWED',
            )
          "
          :text="t('USER_MANAGEMENT_DELETE_ALL')" />
      </v-col>
    </v-row>

    <!-- Confirmation Dialog -->
    <v-dialog v-model="showConfirmDialog" content-class="small" width="500">
      <v-card class="pa-6 dDialog" flat>
        <v-card-title>
          <v-row>
            <v-col cols="12" class="d-flex align-center">
              <span class="text-h5">{{ t('USER_MANAGEMENT_CONFIRM_DELETE_TITLE') }}</span>
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12">
              <p class="text-body-1">
                {{
                  confirmDialogConfig.isDeleteAll
                    ? t('USER_MANAGEMENT_CONFIRM_DELETE_ALL_MESSAGE')
                    : t('USER_MANAGEMENT_CONFIRM_DELETE_MESSAGE')
                }}
                <strong>{{ username }}</strong
                >?
                <span v-if="!confirmDialogConfig.isDeleteAll"> <br />({{ confirmDialogConfig.entityName }}) </span>
              </p>
              <p class="text-body-2 text-medium-emphasis mt-2">
                {{ t('USER_MANAGEMENT_CONFIRM_DELETE_WARNING') }}
              </p>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" class="text-none" @click="showConfirmDialog = false">
            {{ t('BTN_CANCEL') }}
          </v-btn>
          <v-btn color="primary" variant="flat" class="text-none" @click="executeEntityDeletion">
            {{ t('Btn_confirm') }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>
