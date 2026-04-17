<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import ConfirmationDialog from '@disclosure-portal/components/dialog/ConfirmationDialog.vue';
import {ErrorDialogInterface} from '@disclosure-portal/components/dialog/DialogInterfaces';
import ErrorDialog from '@disclosure-portal/components/dialog/ErrorDialog.vue';
import NewUserDialog from '@disclosure-portal/components/dialog/NewUserDialog.vue';
import ProjectChildrenMembersAddedDialog from '@disclosure-portal/components/dialog/ProjectChildrenMembersAddedDialog.vue';
import DHTTPError from '@disclosure-portal/model/DHTTPError';
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import {IDefaultSelectItem} from '@disclosure-portal/model/IObligation';
import {
  IProjectChildrenMembers,
  ProjectChildMemberCombi,
  ProjectChildrenMemberSuccessResponse,
  ProjectKeyName,
  ProjectUser,
  UserType,
} from '@disclosure-portal/model/Project';
import projectService from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import eventBus from '@disclosure-portal/utils/eventbus';
import {getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DIconButton from '@shared/components/disco/DIconButton.vue';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader} from '@shared/types/table';
import _, {indexOf} from 'lodash';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const {t} = useI18n();
const projectStore = useProjectStore();
const {info} = useSnackbar();

const dataAreLoaded = ref(false);
const projectChildrenMembers = ref<IProjectChildrenMembers>({} as IProjectChildrenMembers);
const menu = ref(false);
const selectedFilterUserType = ref<string[]>([]);
const possibleUserTypes = ref<IDefaultSelectItem[]>([]);
const search = ref('');
const userDialogVisible = ref(false);
const projectChildrenMembersAddedVisible = ref(false);
const targetUser = ref('');
const projectChildrenMembersAddedResponseData = ref<ProjectChildrenMemberSuccessResponse[]>([]);
const errorDialog = ref<ErrorDialogInterface | null>(null);
const userDialogRef = ref();
const confirmationDialogVisible = ref(false);
const confirmationDialogConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const userDialogMode = ref<'create' | 'edit'>('create');
const editingUser = ref<ProjectUser>(new ProjectUser());
const ownerRemaining = ref(false);
const projectsKeyName = ref<ProjectKeyName[]>([]);
const targetProjectKey = ref('');
const tableHeight = ref(0);
const route = useRoute();

const projectModel = computed(() => projectStore.currentProject!);

const userHeaders = computed<DataTableHeader[]>(() => {
  const res: DataTableHeader[] = [
    {
      title: t('COL_USER'),
      align: 'start',
      sortable: true,
      class: 'tableHeaderCell',
      value: 'projectMember.userId',
      width: 420,
    },
    {
      title: t('COL_USER_TYPE'),
      align: 'start',
      sortable: true,
      class: 'tableHeaderCell',
      value: 'projectMember.userType',
      width: 160,
    },
    {
      title: t('COL_USER_ROLE'),
      align: 'start',
      sortable: true,
      class: 'tableHeaderCell',
      value: 'projectMember.responsible',
      width: 160,
    },
    {
      title: t('COL_USER_COMMENT'),
      align: 'start',
      sortable: true,
      width: 160,
      class: 'tableHeaderCell',
      value: 'projectMember.comment',
    },
    {
      title: t('PROJECT'),
      align: 'start',
      filterable: true,
      class: 'tableHeaderCell',
      width: 160,
      value: 'projectName',
    },
  ];
  if (projectModel.value.allowUserManagementUpdate || projectModel.value.allowUserManagementDelete) {
    res.unshift({
      title: t('COL_ACTIONS'),
      align: 'center',
      filterable: true,
      sortable: false,
      class: 'tableHeaderCell',
      value: 'actions',
      width: 120,
    });
  }
  return res;
});

const sortItems = () => [{key: 'projectMember.created', order: 'asc' as const}] as const;

const reload = async () => {
  dataAreLoaded.value = false;
  if (projectModel.value._key) {
    projectChildrenMembers.value = await projectService.getChildrenMembers(projectModel.value._key);
    possibleUserTypes.value = getPossibleUserTypes();
  }
  dataAreLoaded.value = true;
};

const showCreateUserDialog = () => {
  userDialogMode.value = 'create';
  projectsKeyName.value = [];
  editingUser.value = new ProjectUser();
  userDialogVisible.value = true;
};

const createUser = async (user: ProjectUser) => {
  if (
    projectChildrenMembers.value.list.filter(
      (combi) => combi.projectMember.userId === user.userId && combi.projectKey === projectModel.value._key,
    ).length === 0
  ) {
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

const showCreateMultiUserDialog = () => {
  userDialogMode.value = 'create';
  const uniqueList = _.uniqBy(projectChildrenMembers.value.list, (item) => `${item.projectKey}-${item.projectName}`);
  projectsKeyName.value = uniqueList.map((combi) => new ProjectKeyName(combi.projectKey, combi.projectName));
  userDialogVisible.value = true;
};

const createMultiUser = async (user: ProjectUser, selectedProjectKeys: string[]) => {
  const response = (
    await projectService.addProjectChildrenMember(
      projectModel.value._key,
      user,
      user.comment,
      user.responsible,
      selectedProjectKeys,
    )
  ).data;
  info(t('DIALOG_project_member_create_success'));
  closeUserDialog();
  await reload();
  if (response) {
    targetUser.value = user.userId;
    projectChildrenMembersAddedResponseData.value = response;
    projectChildrenMembersAddedVisible.value = true;
  }
};

const showEditUserDialog = (projectChildMember: ProjectChildMemberCombi) => {
  userDialogMode.value = 'edit';
  const currentOwnerRemaining =
    projectChildMember.projectMember.userType !== UserType.OWNER ||
    projectChildrenMembers.value.list.filter(
      (combi) => combi.projectMember.userType === UserType.OWNER && combi.projectKey === projectChildMember.projectKey,
    ).length > 1;
  editingUser.value = projectChildMember.projectMember;
  ownerRemaining.value = currentOwnerRemaining;
  targetProjectKey.value = projectChildMember.projectKey;
  userDialogVisible.value = true;
};

const editUser = async (user: ProjectUser, oldUserId: string, targetProjectKey: string | undefined) => {
  if (
    (oldUserId === user.userId ||
      projectChildrenMembers.value.list.filter(
        (combi) => combi.projectMember.userId === user.userId && combi.projectKey === targetProjectKey,
      ).length === 0) &&
    targetProjectKey
  ) {
    await projectService.editProjectMember(targetProjectKey, user, oldUserId, user.comment, user.responsible);
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

const showDeleteUserDialog = async (projectChildMember: ProjectChildMemberCombi) => {
  let userName = projectChildMember.projectMember.userId;
  if (projectChildMember.projectMember.userProfile.user) {
    userName = `${projectChildMember.projectMember.userProfile.lastname}, ${projectChildMember.projectMember.userProfile.forename} (${projectChildMember.projectMember.userProfile.user})`;
  }
  if (projectChildMember.projectMember.responsible) {
    confirmationDialogConfig.value = {
      type: ConfirmationType.NOT_SET,
      contextKey: projectChildMember.projectKey,
      key: projectChildMember.projectMember.userId,
      name: userName,
      description: 'DLG_CAN_NOT_DELETE_RESPONSIBLE',
      extendedDetails: t('USER_IS_RESPONSIBLE'),
      okButton: 'Btn_remove',
      okButtonIsDisabled: true,
      title: 'DLG_WARNING_TITLE',
    };
    confirmationDialogVisible.value = true;
    return;
  }
  if (
    projectChildMember.projectMember.userType !== UserType.OWNER ||
    projectChildrenMembers.value.list.filter(
      (combi) => combi.projectMember.userType === UserType.OWNER && combi.projectKey === projectChildMember.projectKey,
    ).length > 1
  ) {
    const r = await projectService.getPendingApprovalOrReviewUsage(
      projectChildMember.projectKey,
      projectChildMember.projectMember.userId,
    );
    const isInUse = r.data.success;
    if (isInUse) {
      confirmationDialogConfig.value = {
        type: ConfirmationType.NOT_SET,
        contextKey: projectChildMember.projectKey,
        key: projectChildMember.projectMember.userId,
        name: userName,
        description: 'DLG_CAN_NOT_DELETE_IN_USE',
        extendedDetails: t('USER_IN_PENDING_APPROVAL'),
        okButton: 'Btn_remove',
        okButtonIsDisabled: true,
        title: 'DLG_WARNING_TITLE',
      };
      confirmationDialogVisible.value = true;
    } else {
      confirmationDialogConfig.value = {
        type: ConfirmationType.NOT_SET,
        contextKey: projectChildMember.projectKey,
        key: projectChildMember.projectMember.userId,
        name: userName,
        description: 'DLG_CONFIRMATION_DESCRIPTION_REMOVE',
        okButton: 'Btn_remove',
        okButtonIsDisabled: false,
      };
      confirmationDialogVisible.value = true;
    }
  } else {
    const dialog = new ErrorDialogConfig();
    dialog.title = t('user_removal_error_title');
    dialog.description = t('user_removal_error_message');
    errorDialog.value?.open(dialog);
  }
};

const deleteUser = async (config: IConfirmationDialogConfig) => {
  if (config.okButtonIsDisabled) return;
  const childProjectKey = config.contextKey;
  const userId = config.key;
  if (childProjectKey && userId) {
    await projectService.deleteProjectMember(childProjectKey, userId);
    info(t('DIALOG_project_member_delete_success'));
    await reload();
  }
};

const closeUserDialog = () => {
  userDialogRef.value?.close();
};

const filteredList = computed(() => {
  if (!projectChildrenMembers.value) {
    return [];
  }
  return _.filter(projectChildrenMembers.value.list, filterOnType);
});

const filterOnType = (item: ProjectChildMemberCombi): boolean => {
  return (
    selectedFilterUserType.value.length === 0 ||
    indexOf(selectedFilterUserType.value, item.projectMember.userType) !== -1
  );
};

const getPossibleUserTypes = (): IDefaultSelectItem[] => {
  if (!projectChildrenMembers.value) {
    return [];
  }
  return _.chain(projectChildrenMembers.value.list)
    .uniqBy('projectMember.userType')
    .map((item: ProjectChildMemberCombi) => {
      return {
        text: item.projectMember.userType,
        value: item.projectMember.userType,
      } as IDefaultSelectItem;
    })
    .value();
};

const customFilterTable = (value: any, search: string, internalItem: any) => {
  const item = internalItem.raw;
  const lowerSearch = search.toLowerCase();
  if (value === item.projectMember.userId) {
    const forename = item.projectMember.userProfile.forename.toLowerCase();
    const lastname = item.projectMember.userProfile.lastname.toLowerCase();
    return (
      forename.indexOf(lowerSearch) !== -1 ||
      lastname.indexOf(lowerSearch) !== -1 ||
      item.projectMember.userId.toLowerCase().indexOf(lowerSearch) !== -1
    );
  }
  if (value === item.projectMember?.responsible && 'project responsible'.includes(lowerSearch)) {
    return item.projectMember.responsible;
  }

  if (typeof value !== 'string') {
    return false;
  }

  return value.toLowerCase().indexOf(lowerSearch) !== -1;
};

const getActionButtons = (item: ProjectChildMemberCombi): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-pencil',
      hint: t('TT_edit_user'),
      event: 'edit',
      show: projectModel.value.allowUserManagementUpdate,
    },
    {
      icon: 'mdi-close',
      hint: t('TT_remove_user'),
      event: 'remove',
      show: projectModel.value.allowUserManagementDelete,
    },
  ];
};

onMounted(async () => {
  await reload();
});

watch(
  () => route.path,
  (_newPath, _oldPath) => {
    if (_newPath.includes('childrenUsers')) {
      reload();
    }
  },
);
</script>

<template>
  <TableLayout has-title has-tab>
    <template #buttons>
      <DCActionButton
        large
        :text="t('BTN_ADD')"
        icon="mdi-plus"
        :hint="t('TT_new_user')"
        @click="showCreateUserDialog"
        v-if="projectModel.allowUserManagementCreate && !projectModel.isDeprecated" />
      <DCActionButton
        large
        :text="t('BTN_ADD_MULTI')"
        icon="mdi-plus"
        :hint="t('TT_new_user')"
        @click="showCreateMultiUserDialog"
        v-if="projectModel.allowUserManagementCreate && !projectModel.isDeprecated" />
      <v-spacer />
      <v-text-field
        autocomplete="off"
        max-width="450px"
        variant="outlined"
        density="compact"
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        clearable
        single-line
        hide-details></v-text-field>
    </template>
    <template #table>
      <v-data-table
        density="compact"
        :loading="!dataAreLoaded"
        fixed-header
        :height="tableHeight"
        class="striped-table custom-data-table h-full"
        :headers="userHeaders"
        :items="filteredList"
        :sort-by="sortItems()"
        :item-class="getCssClassForTableRow"
        :search="search"
        :custom-filter="customFilterTable">
        <template v-slot:[`header.projectMember.userType`]="{column, getSortIcon, toggleSort}">
          <div class="v-data-table-header__content">
            <span>{{ column.title }}</span>
            <v-menu :close-on-content-click="false" v-model="menu">
              <template v-slot:activator="{props}">
                <DIconButton
                  :parentProps="props"
                  icon="mdi-filter-variant"
                  :hint="t('TT_SHOW_FILTER')"
                  :color="selectedFilterUserType.length > 0 ? 'primary' : 'default'" />
              </template>
              <div class="bg-background" style="width: 280px">
                <v-row class="d-flex ma-1 mr-2 justify-end">
                  <DIconButton icon="mdi-close" @clicked="menu = false" color="default" />
                </v-row>
                <v-select
                  v-model="selectedFilterUserType"
                  :items="possibleUserTypes"
                  class="pa-2 mx-2 pb-4"
                  :label="t('Lbl_filter_userType')"
                  clearable
                  multiple
                  item-title="text"
                  item-value="value"
                  variant="outlined"
                  density="compact"
                  menu
                  transition="scale-transition"
                  persistent-clear
                  :list-props="{class: 'striped-filter-dd py-0'}">
                  <template v-slot:item="{props}">
                    <v-list-item v-bind="props" class="px-2 py-0">
                      <template v-slot:prepend="{isSelected}">
                        <v-checkbox hide-details :model-value="isSelected" />
                      </template>
                      <template v-slot:title="{title}">
                        <span class="pStatusFilterEntry"> {{ title }}</span>
                      </template>
                    </v-list-item>
                  </template>
                  <template v-slot:selection="{item, index}">
                    <div v-if="index === 0" class="d-flex align-center">
                      <span class="pStatusFilterEntry">{{ item.title }}</span>
                    </div>
                    <span v-if="index === 1" class="pAdditionalFilter">
                      +{{ selectedFilterUserType.length - 1 }} others
                    </span>
                  </template>
                </v-select>
              </div>
            </v-menu>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </div>
        </template>
        <template v-slot:[`item.projectMember.userId`]="{item}">
          <span v-if="item.projectMember.userProfile.user">
            {{ item.projectMember.userProfile.lastname }}, {{ item.projectMember.userProfile.forename }} ({{
              item.projectMember.userProfile.user
            }})
          </span>
          <span v-else>{{ item.projectMember.userId }}</span>
        </template>
        <template v-slot:[`item.projectMember.responsible`]="{item}">
          <span v-if="item.projectMember.responsible">{{ t('COL_USER_ROLE_RESPONSIBLE') }}</span>
        </template>
        <template v-slot:[`item.actions`]="{item}">
          <TableActionButtons
            variant="compact"
            :buttons="getActionButtons(item)"
            @edit="showEditUserDialog(item)"
            @remove="showDeleteUserDialog(item)" />
        </template>
      </v-data-table>
    </template>
  </TableLayout>
  <template>
    <NewUserDialog
      ref="userDialogRef"
      v-model:showDialog="userDialogVisible"
      :mode="userDialogMode"
      :projectKey="projectModel._key"
      :user="editingUser"
      :ownerRemaining="ownerRemaining"
      :projectsKeyName="projectsKeyName"
      :targetProjectKey="targetProjectKey"
      @createUser="createUser"
      @createMultiUser="createMultiUser"
      @editUser="editUser" />
    <ProjectChildrenMembersAddedDialog
      v-model:showDialog="projectChildrenMembersAddedVisible"
      :targetUser="targetUser"
      :items="projectChildrenMembersAddedResponseData" />
    <ConfirmationDialog
      v-model:showDialog="confirmationDialogVisible"
      :config="confirmationDialogConfig"
      @confirm="deleteUser" />
    <ErrorDialog ref="errorDialog" />
  </template>
</template>
