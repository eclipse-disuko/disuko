<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import Icons from '@disclosure-portal/constants/icons';
import SimpleProfileData from '@disclosure-portal/model/ProfileData';
import {UserDto, UserRequestDto} from '@disclosure-portal/model/Users';
import AdminService from '@disclosure-portal/services/admin';
import profileService from '@disclosure-portal/services/profile';
import {useAppStore} from '@disclosure-portal/stores/app';
import {createNavItemsGroup, useUserStore} from '@disclosure-portal/stores/user';
import {formatDate} from '@disclosure-portal/utils/View';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps<{
  userProfile: UserDto;
  rolesSwitchable: boolean;
  hasUsersAccess: boolean;
}>();
const emit = defineEmits(['reloadUserProfile']);
const {t} = useI18n();
const userStore = useUserStore();
const appStore = useAppStore();
const icons = Icons;
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const selectedTab = ref('userOverview');
const confirmVisible = ref(false);

const enableDisable = () => {
  if (!props.userProfile.active) {
    confirmConfig.value = confirmConfig.value = {
      type: ConfirmationType.NOT_SET, // allows to separate in onConfirm Callback see below if not need set to ConfirmationType.NOT_SET
      key: '',
      name: '',
      description: 'DLG_CONFIRMATION_DESCRIPTION_ENABLE_USER',
      extendedDetails: '',
      okButton: 'BTN_ENABLE_USER',
      okButtonIsDisabled: false,
    } as IConfirmationDialogConfig;
  } else {
    confirmConfig.value = confirmConfig.value = {
      type: ConfirmationType.NOT_SET, // allows to separate in onConfirm Callback see below if not need set to ConfirmationType.NOT_SET
      key: '',
      name: '',
      description: 'DLG_CONFIRMATION_DESCRIPTION_DISABLE_USER',
      extendedDetails: '',
      okButton: 'BTN_DISABLE_USER',
      okButtonIsDisabled: false,
    } as IConfirmationDialogConfig;
  }
  confirmVisible.value = true;
};

const doEnableDisable = async (config: IConfirmationDialogConfig) => {
  await AdminService.enableDisableUser(props.userProfile._key, {
    active: !props.userProfile.active,
  } as UserRequestDto);
  emit('reloadUserProfile', props.userProfile._key);
};

const fetchRoles = async () => {
  return props.hasUsersAccess
    ? AdminService.getUserProjectRoles(props.userProfile._key)
    : profileService.getProjectRoles();
};

const fetchTasks = async () => {
  return props.hasUsersAccess ? AdminService.getUserTasks(props.userProfile._key) : profileService.getTasks();
};

const applyNewUserRoles = async (user: UserDto, forceNonInternal: boolean) => {
  const userDtoResponse = (
    await AdminService.updateUserRoles(user._key, {
      roles: user.roles,
    })
  ).data;

  if (user.user === userStore.getProfile.user) {
    let oAuthTokenResponse: SimpleProfileData;
    if (forceNonInternal) {
      oAuthTokenResponse = (await AdminService.getUserTokensForNonInternal(userDtoResponse._key)).data;
    } else {
      oAuthTokenResponse = (await AdminService.getUserTokens(userDtoResponse._key)).data;
    }
    userStore.setSimpleProfileData(oAuthTokenResponse);
    createNavItemsGroup();
  }

  emit('reloadUserProfile', userDtoResponse._key);
};
</script>

<template>
  <v-container fluid>
    <v-row class="header">
      <v-col md="auto">
        <h1 class="d-headline">{{ t('BC_Profile') }}</h1>
      </v-col>
    </v-row>
    <v-row class="expand" v-if="userProfile">
      <v-col>
        <v-card>
          <v-tabs v-model="selectedTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
            <v-tab value="userOverview">
              {{ t('TAB_USER_OVERVIEW') }}
            </v-tab>
            <v-tab value="roles">
              {{ t('TITLE_ROLES') }}
            </v-tab>
            <v-tab value="tasks">
              {{ t('TASKS') }}
            </v-tab>
            <v-tab value="auditLog" v-if="hasUsersAccess">
              {{ t('TITLE_AUDIT_LOG') }}
            </v-tab>
          </v-tabs>
          <v-tabs-window v-model="selectedTab">
            <v-tabs-window-item value="userOverview">
              <v-row class="pa-4">
                <v-col cols="12" xs="6" sm="6" md="3" lg="1">
                  <span class="d-text d-secondary-text">{{ t('COL_USER_ID') }}</span
                  ><br />
                  <span>{{ userProfile.user }}</span>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="1">
                  <span class="d-text d-secondary-text">{{ t('COL_FORENAME') }}</span
                  ><br />
                  <span>{{ userProfile.forename }}</span>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="1">
                  <span class="d-text d-secondary-text">{{ t('COL_LASTNAME') }}</span
                  ><br />
                  <span>{{ userProfile.lastname }}</span>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="3">
                  <span class="d-text d-secondary-text">{{ t('COL_EMAIL') }}</span
                  ><br />
                  <span>{{ userProfile.email }}</span>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="1">
                  <span class="d-text d-secondary-text">{{ t('COL_CREATED') }}</span
                  ><br />
                  <span>{{ formatDate(userProfile.created) }}</span>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="1">
                  <span class="d-text d-secondary-text">{{ t('COL_UPDATED') }}</span
                  ><br />
                  <span>{{ formatDate(userProfile.updated) }}</span>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="2">
                  <span class="d-text d-secondary-text">{{ t('USER_STATUS') }}</span
                  ><br />
                  <div v-if="userProfile.active">
                    <v-icon size="small" color="success">{{ icons.CIRCLE_FILLED }}</v-icon>
                    <DCActionButton
                      large
                      :text="t('DISABLE_USER')"
                      :hint="t('TT_disable_user')"
                      @click="enableDisable"
                      class="mx-2"
                      v-if="hasUsersAccess" />
                    <span v-else class="d-subtitle-2 pl-1">{{ t('ICON_LABEL_TEXT_ACTIVE') }}</span>
                  </div>
                  <div v-else>
                    <v-icon size="small" color="warning">{{ icons.CIRCLE_FILLED }}</v-icon>
                    <DCActionButton
                      large
                      :text="t('ENABLE_USER')"
                      :hint="t('TT_enable_user')"
                      @click="enableDisable"
                      class="mx-2"
                      v-if="hasUsersAccess" />
                    <span v-else class="d-subtitle-2 pl-1">{{ t('ICON_LABEL_TEXT_INACTIVE') }}</span>
                  </div>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="2">
                  <span class="d-text d-secondary-text">{{ t('ROLES') }}</span
                  ><br />
                  <span v-if="!userProfile.roles || userProfile.roles.length === 0">{{ t('NO_ROLES') }}</span>
                  <span v-else v-for="(item, index) in userProfile.roles" :key="index">
                    <span>{{ t(item) }}</span
                    ><br />
                  </span>
                  <SwitchUserRolesDialog
                    ref="dlgUserRoles"
                    v-if="props && rolesSwitchable && userProfile.roles && userProfile.roles.length > 0"
                    @applyNewUserRoles="applyNewUserRoles"
                    v-slot="{showDialog}">
                    <DCActionButton
                      large
                      :text="t('BTN_RESTRICT_ROLES')"
                      icon="mdi-shield-off"
                      :hint="t('TT_RESTRICT_ROLES')"
                      @click="showDialog(userProfile)" />
                  </SwitchUserRolesDialog>
                </v-col>
              </v-row>
              <v-row class="pa-4">
                <v-col cols="12" xs="6" sm="6" md="3" lg="2">
                  <span class="d-text d-secondary-text">{{ t('COL_TERMS_DATE') }}</span
                  ><br />
                  <span>{{ formatDate(userProfile.termsOfUseDate) }}</span>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="2">
                  <span class="d-text d-secondary-text">{{ t('COL_TERMS_VERSION') }}</span
                  ><br />
                  <span>{{ userProfile.termsOfUseVersion }}</span>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="2">
                  <span class="d-text d-secondary-text">{{ t('COL_TERMS_ACCEPTANCE') }}</span
                  ><br />
                  <v-icon v-if="userProfile.termsOfUse" size="small" color="primary">mdi-check</v-icon>
                  <v-icon v-else size="small" class="greyCheck">mdi-check</v-icon>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="2">
                  <span class="d-text d-secondary-text">{{ t('USER_ACCESS_SCOPE') }}</span
                  ><br />
                  <span v-if="userProfile.isInternal">{{ t('USER_ACCESS_SCOPE_INTERNAL') }}</span>
                  <span v-else>{{ t('USER_ACCESS_SCOPE_EXTERNAL') }}</span>
                </v-col>
              </v-row>
              <v-row v-if="userProfile.metaData" class="pa-4">
                <v-col v-if="userProfile.metaData" cols="12" xs="6" sm="6" md="3" lg="2">
                  <span class="d-text d-secondary-text">{{ t('DEPARTMENT') }}</span
                  ><br />
                  <span>{{ userProfile.metaData.department }}</span>
                </v-col>
                <v-col v-if="userProfile.metaData" cols="12" xs="6" sm="6" md="3" lg="2">
                  <span class="d-text d-secondary-text">{{ t('DEPARTMENT_DESCRIPTION') }}</span
                  ><br />
                  <span>{{ userProfile.metaData.departmentDescription }}</span>
                </v-col>
                <v-col v-if="userProfile.metaData" cols="12" xs="6" sm="6" md="3" lg="2">
                  <span class="d-text d-secondary-text">{{ t('COMPANY_IDENTIFIER') }}</span
                  ><br />
                  <span>{{ userProfile.metaData.companyIdentifier }}</span>
                </v-col>
                <v-col cols="12" xs="6" sm="6" md="3" lg="2">
                  <span class="d-text d-secondary-text">{{ t('DEPROVISIONED_DATE') }}</span
                  ><br />
                  <span>{{ userProfile.deprovisioned ? formatDate(userProfile.deprovisioned) : t('NOT_SET') }}</span>
                </v-col>
                <v-col>
                  <v-switch
                    color="primary"
                    v-model="appStore.alternateRender"
                    :label="t('SWITCH_ALTERNATE_RENDER')"></v-switch>
                </v-col>
              </v-row>
            </v-tabs-window-item>

            <v-tabs-window-item value="roles">
              <GridRoles :fetch-method="fetchRoles" />
            </v-tabs-window-item>
            <v-tabs-window-item value="tasks">
              <GridTask
                :fetch-method="fetchTasks"
                :hideEditAction="true"
                :readOnly="hasUsersAccess"
                :in-own-view="false" />
            </v-tabs-window-item>
            <v-tabs-window-item value="auditLog" v-if="hasUsersAccess">
              <GridAuditLog :fetch-method="() => AdminService.getUserAuditTrail(userProfile._key)" />
            </v-tabs-window-item>
          </v-tabs-window>
        </v-card>
      </v-col>
    </v-row>

    <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doEnableDisable" />
  </v-container>
</template>
