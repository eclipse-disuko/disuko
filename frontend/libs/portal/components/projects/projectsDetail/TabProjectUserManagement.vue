<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
// TODO: TabProjectUserManagement & TabProjectChildrenUserManagement have a lot of duplicated code. Refactor to a common component.
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {ErrorDialogInterface} from '@disclosure-portal/components/dialog/DialogInterfaces';
import DHTTPError from '@disclosure-portal/model/DHTTPError';
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import {ProjectUser, UserType} from '@disclosure-portal/model/Project';
import projectService from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import eventBus from '@disclosure-portal/utils/eventbus';
import {getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {DataTableHeader, DataTableHeaderFilterItems, SortItem} from '@shared/types/table';
import _ from 'lodash';
import {computed, onBeforeMount, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useTableActionSlider} from '@shared/composables/useTableActionSlider';

const {t} = useI18n();
const projectStore = useProjectStore();
const {info} = useSnackbar();
const {sliderWidth} = useTableActionSlider();

const userDialogVisible = ref(false);
const confirmationDialogVisible = ref(false);
const confirmationDialogConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);

const users = ref<ProjectUser[]>([]);
const dataAreLoaded = ref(false);
const selectedFilterUserType = ref<string[]>([]);
const search = ref('');
const userDialogMode = ref<'create' | 'edit'>('create');
const editingUser = ref<ProjectUser>(new ProjectUser());
const ownerRemaining = ref(false);

const userDialogRef = ref();
const errorDialog = ref<ErrorDialogInterface | null>(null);
const sortItems = ref<SortItem[]>([{key: 'userType', order: 'asc'}]);

const projectModel = computed(() => projectStore.currentProject!);
const userHeaders = computed((): DataTableHeader[] => [
  ...(projectModel.value.allowUserManagementUpdate || projectModel.value.allowUserManagementDelete
    ? [
      {
        title: t('COL_ACTIONS'),
        align: 'start',
        value: 'actions',
        width: sliderWidth.value,
        sortable: false,
      } as DataTableHeader,
    ]
    : []),
  {
    title: t('COL_USER'),
    align: 'start',
    sortable: true,
    value: 'userId',
    width: 420,
  },
  {
    title: t('COL_USER_TYPE'),
    align: 'start',
    value: 'userType',
    width: 160,
    sortable: true,
  },
  {
    title: t('COL_USER_ROLE'),
    align: 'start',
    sortable: true,
    value: 'responsible',
    width: 180,
  },
  {
    title: t('COL_USER_COMMENT'),
    align: 'start',
    sortable: true,
    value: 'comment',
    width: 180,
  },
]);

const filteredList = computed(() => {
  return _.filter(users.value, filterOnType);
});

const filterOnType = (item: ProjectUser): boolean => {
  return selectedFilterUserType.value.length === 0 || selectedFilterUserType.value.includes(item.userType);
};

const possibleUserTypes = computed((): DataTableHeaderFilterItems[] => {
  if (!users.value) {
    return [];
  }

  const uniqueUserTypes = [...new Set(users.value.map(({userType}) => userType))];

  return uniqueUserTypes.map((value: string) => {
    return {
      value,
    } as DataTableHeaderFilterItems;
  });
});

const reload = async () => {
  dataAreLoaded.value = false;
  try {
    const response = await projectService.getUserManagement(projectModel.value._key);
    users.value = response.users;
  } finally {
    dataAreLoaded.value = true;
  }
};

const showCreateUserDialog = () => {
  userDialogMode.value = 'create';
  editingUser.value = new ProjectUser();
  userDialogVisible.value = true;
};

const showEditUserDialog = (user: ProjectUser) => {
  userDialogMode.value = 'edit';
  editingUser.value = user;
  ownerRemaining.value =
    user.userType !== UserType.OWNER || users.value.filter((u) => u.userType === UserType.OWNER).length > 1;
  userDialogVisible.value = true;
};

const showDeleteUserDialog = async (user: ProjectUser) => {
  let userName = user.userId;
  if (user.userProfile.user) {
    userName = `${user.userProfile.lastname}, ${user.userProfile.forename} (${user.userProfile.user})`;
  }

  if (user.responsible) {
    confirmationDialogConfig.value = {
      type: ConfirmationType.DELETE,
      key: user.userId,
      name: userName,
      description: 'DLG_CAN_NOT_DELETE_RESPONSIBLE',
      extendedDetails: '' + t('USER_IS_RESPONSIBLE'),
      okButton: 'Btn_remove',
      okButtonIsDisabled: true,
      title: 'DLG_WARNING_TITLE',
    };
    confirmationDialogVisible.value = true;
    return;
  }
  if (user.userType !== UserType.OWNER || users.value.filter((u) => u.userType === UserType.OWNER).length > 1) {
    const r = await projectService.getPendingApprovalOrReviewUsage(projectModel.value._key, user.userId);
    const isInUse = r.data.success;
    if (isInUse) {
      confirmationDialogConfig.value = {
        type: ConfirmationType.NOT_SET,
        key: user.userId,
        name: userName,
        description: 'DLG_CAN_NOT_DELETE_IN_USE',
        extendedDetails: t('USER_IN_PENDING_APPROVAL'),
        okButton: 'Btn_remove',
        okButtonIsDisabled: true,
        title: 'DLG_WARNING_TITLE',
      };
    } else {
      confirmationDialogConfig.value = {
        type: ConfirmationType.DELETE,
        key: user.userId,
        name: userName,
        description: 'DLG_CONFIRMATION_DESCRIPTION_REMOVE',
        okButton: 'Btn_remove',
        okButtonIsDisabled: false,
      };
    }
    confirmationDialogVisible.value = true;
  } else {
    const dialog = new ErrorDialogConfig();
    dialog.title = t('user_removal_error_title');
    dialog.description = t('user_removal_error_message');
    errorDialog.value?.open(dialog);
  }
};

const deleteUser = async (config: IConfirmationDialogConfig) => {
  if (config.okButtonIsDisabled) return;
  if (config.key) {
    await projectService.deleteProjectMember(projectModel.value._key, config.key);
    info(t('DIALOG_project_member_delete_success'));
    await reload();
  }
};

const createUser = async (user: ProjectUser) => {
  if (!users.value.some((u) => u.userId === user.userId)) {
    await projectService.addProjectMember(projectModel.value._key, user, user.comment, user.responsible);
    info(t('DIALOG_project_member_create_success'));
    closeUserDialog();
    await reload();
  } else {
    const error = new DHTTPError();
    error.title = t('user_create_error_title');
    error.message = t('user_create_error_message') + ' ' + user.userId;
    eventBus.emit('on-api-error', error);
  }
};

const editUser = async (user: ProjectUser, oldUserId: string) => {
  if (oldUserId === user.userId || !users.value.some((u) => u.userId === user.userId)) {
    await projectService.editProjectMember(projectModel.value._key, user, oldUserId, user.comment, user.responsible);
    info(t('DIALOG_project_member_edit_success'));
    closeUserDialog();
    await reload();
  } else {
    const error = new DHTTPError();
    error.title = t('user_create_error_title');
    error.message = t('user_create_error_message') + ' ' + user.userId;
    eventBus.emit('on-api-error', error);
  }
};
const customFilterTable = (value: any, search: string, internalItem: any) => {
  const item = internalItem.raw;
  const lowerSearch = search.toLowerCase();
  if (value === item.userId) {
    const forename = item.userProfile.forename.toLowerCase();
    const lastname = item.userProfile.lastname.toLowerCase();
    return (
      forename.indexOf(lowerSearch) !== -1 ||
      lastname.indexOf(lowerSearch) !== -1 ||
      item.userId.toLowerCase().indexOf(lowerSearch) !== -1
    );
  }
  if (value === item.responsible && 'project responsible'.includes(lowerSearch)) {
    return item.responsible;
  }

  if (typeof value !== 'string') {
    return false;
  }

  return value.toLowerCase().indexOf(lowerSearch) !== -1;
};

const closeUserDialog = () => {
  userDialogRef.value?.close();
};

const getActionButtons = (): TableActionButtonsProps['buttons'] => {
  const canModify = !projectModel.value.isDeprecated;

  return [
    {
      icon: 'mdi-pencil',
      hint: t('TT_edit_user'),
      event: 'edit',
      show: projectModel.value.allowUserManagementUpdate && canModify,
    },
    {
      icon: 'mdi-close',
      hint: t('TT_remove_user'),
      event: 'remove',
      show: projectModel.value.allowUserManagementDelete && canModify,
    },
  ];
};

onBeforeMount(async () => {
  await reload();
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <DCActionButton
        :text="t('BTN_ADD')"
        icon="mdi-plus"
        :hint="t('TT_new_user')"
        @click="showCreateUserDialog"
        v-if="projectModel.allowUserManagementCreate && !projectModel.isDeprecated" />
      <v-spacer></v-spacer>
      <v-text-field
        autocomplete="off"
        max-width="450px"
        variant="outlined"
        density="compact"
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        single-line
        clearable
        hide-details></v-text-field>
    </template>
    <template #table>
      <div ref="tableUserManagement" class="fill-height action-slider-table">
        <v-data-table
          density="compact"
          :items-per-page="100"
          :loading="!dataAreLoaded"
          :search="search"
          fixed-header
          class="striped-table custom-data-table fill-height"
          :headers="userHeaders"
          :items="filteredList"
          :item-class="getCssClassForTableRow"
          :sort-by="sortItems"
          :custom-filter="customFilterTable">
          <template #[`header.userType`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterUserType"
                  :column="column"
                  :label="t('COL_USER_TYPE')"
                  :allItems="possibleUserTypes">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`item.userId`]="{item}">
            <span v-if="item.userProfile.user">
              {{ item.userProfile.lastname }}, {{ item.userProfile.forename }} ({{ item.userProfile.user }})
            </span>
            <span v-else>{{ item.userId }}</span>
          </template>
          <template #[`item.responsible`]="{item}">
            <div v-if="item.responsible">{{ t('COL_USER_ROLE_RESPONSIBLE') }}</div>
          </template>
          <template #[`item.actions`]="{item}">
            <TableActionButtons
              variant="slider"
              :buttons="getActionButtons()"
              @edit="showEditUserDialog(item)"
              @remove="showDeleteUserDialog(item)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <NewUserDialog
    ref="userDialogRef"
    v-model:showDialog="userDialogVisible"
    :mode="userDialogMode"
    :projectKey="projectModel._key"
    :user="editingUser"
    :ownerRemaining="ownerRemaining"
    @createUser="createUser"
    @editUser="editUser" />
  <ConfirmationDialog
    v-model:showDialog="confirmationDialogVisible"
    :config="confirmationDialogConfig"
    @confirm="deleteUser" />
  <ErrorDialog ref="errorDialog" />
</template>
