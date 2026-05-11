<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {getProjectUserTypes, ProjectKeyName, ProjectUser, UserType} from '@disclosure-portal/model/Project';
import {UserDto} from '@disclosure-portal/model/Users';
import useRules from '@disclosure-portal/utils/Rules';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const showDialog = defineModel<boolean>('showDialog');

const props = withDefaults(
  defineProps<{
    mode: 'create' | 'edit';
    projectKey: string;
    user: ProjectUser;
    ownerRemaining: boolean;
    projectsKeyName?: ProjectKeyName[];
    targetProjectKey?: string;
  }>(),
  {},
);
const emit = defineEmits<{
  (e: 'createUser', user: ProjectUser): void;
  (e: 'editUser', user: ProjectUser, oldUserId: string, targetProjectKey?: string): void;
  (e: 'createMultiUser', user: ProjectUser, selectedProjectKeys: string[]): void;
}>();

const {t} = useI18n();
const title = ref('');
const item = ref<ProjectUser>({} as ProjectUser);
const oldUserId = ref('');
const userTypes = ref(getProjectUserTypes());
const rulesSelectbox = [
  (value: string) => {
    return !!value || t('VALIDATION_required');
  },
];
const rulesMultiSelectbox = [
  (value: string) => {
    return (!!value && value.length > 0) || t('VALIDATION_required');
  },
];
const {minMax} = useRules();
const activeRules = ref({
  comment: minMax(t('COL_USER_COMMENT'), 1, 100, true),
});
const isInternalUser = ref(false);
const ownerRemaining = ref(false);
const isResponsible = ref(false);
const targetProjectKey = ref<string | undefined>(undefined);
const projectsKeyNameIntern = ref<ProjectKeyName[]>([]);
const selectedProjectKeys = ref<string[]>([]);
const formUserDialog = ref<VForm | null>(null);
const autocompleteUserRef = ref();

const isTypeOwnerSelected = computed(() => item.value.userType === UserType.OWNER);

const userChanged = (user: UserDto) => {
  isInternalUser.value = user.isInternal;
};

const open = (projectsKeyNameList: ProjectKeyName[] = []) => {
  ownerRemaining.value = true;
  isResponsible.value = false;
  showDialog.value = true;
  title.value = 'UM_DIALOG_TITLE_NEW_USER';
  item.value = {} as ProjectUser;
  projectsKeyNameIntern.value = projectsKeyNameList;
  selectedProjectKeys.value = [];
  autocompleteUserRef.value?.resetForm();
  formUserDialog.value?.resetValidation();
};

const edit = (user: ProjectUser, ownerRemainingFlag: boolean, targetProjectKeyVal: string | undefined = undefined) => {
  item.value = JSON.parse(JSON.stringify(user));
  isResponsible.value = item.value.responsible;
  showDialog.value = true;
  title.value = 'UM_DIALOG_TITLE_EDIT_USER';
  oldUserId.value = user.userId;
  ownerRemaining.value = ownerRemainingFlag;
  isInternalUser.value = user.userProfile.isInternal;
  targetProjectKey.value = targetProjectKeyVal;
  projectsKeyNameIntern.value = [];
  selectedProjectKeys.value = [];
  autocompleteUserRef.value?.resetForm();
  formUserDialog.value?.resetValidation();
};

const close = () => {
  autocompleteUserRef.value?.resetForm();
  showDialog.value = false;
};

const doDialogAction = () => {
  nextTick(async () => {
    const validForm = (await validate())?.valid;
    const validUser = await autocompleteUserRef.value?.validateOnCreate();
    if (validUser && validForm) {
      if (props.mode === 'create') {
        if (projectsKeyNameIntern.value.length > 0 && selectedProjectKeys.value.length > 0) {
          emit('createMultiUser', item.value, selectedProjectKeys.value);
        } else {
          emit('createUser', item.value);
        }
      }
      if (props.mode === 'edit') {
        if (targetProjectKey.value && targetProjectKey.value.length > 0) {
          emit('editUser', item.value, oldUserId.value, targetProjectKey.value);
        } else {
          emit('editUser', item.value, oldUserId.value);
        }
      }
    }
  });
};

const validate = () => {
  return formUserDialog.value?.validate();
};

function toggleSelectAll() {
  if (allProjectsSelected.value) {
    selectedProjectKeys.value = [];
  } else {
    selectedProjectKeys.value = projectsKeyNameIntern.value.map((e) => e.key);
  }
}

const allProjectsSelected = computed(() => {
  return selectedProjectKeys.value.length === projectsKeyNameIntern.value.length;
});

const someProjectsSelected = computed(() => {
  return selectedProjectKeys.value.length > 0;
});

watch(
  () => showDialog.value,
  (newValue) => {
    if (newValue) {
      if (props.mode === 'create') {
        open(props.projectsKeyName);
      } else {
        edit({...props.user}, props.ownerRemaining, props.targetProjectKey);
      }
    }
  },
);

watch(
  () => item.value.userType,
  (newUserType) => {
    if (newUserType !== UserType.OWNER) {
      item.value.responsible = false;
    }
  },
);
defineExpose({close});
</script>

<template>
  <v-form ref="formUserDialog">
    <v-dialog v-model="showDialog" content-class="small" persistent width="600">
      <v-card class="pa-8">
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <span class="text-h5">{{ t(title) }}</span>
            </v-col>
            <v-col cols="2" align="right">
              <DCloseButton @click="close" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-row align="center">
            <v-col
              cols="12"
              xs="12"
              class="errorBorder"
              v-if="projectsKeyNameIntern && projectsKeyNameIntern.length > 0">
              <v-select
                variant="outlined"
                :items="projectsKeyNameIntern"
                v-model="selectedProjectKeys"
                v-bind:menu-props="{location: 'bottom'}"
                item-value="key"
                item-title="name"
                multiple
                chips
                closable-chips
                :label="t('PROJECTS')"
                class="required pb-2"
                hide-details="auto"
                :rules="rulesMultiSelectbox">
                <template v-slot:prepend-item>
                  <v-list-item :title="t('SELECT_ALL')" @click="toggleSelectAll">
                    <template v-slot:prepend>
                      <v-checkbox-btn
                        :indeterminate="someProjectsSelected && !allProjectsSelected"
                        :model-value="allProjectsSelected"></v-checkbox-btn>
                    </template>
                  </v-list-item>
                  <v-divider class="mt-2"></v-divider>
                </template>
              </v-select>
            </v-col>
            <v-col cols="12" xs="12" class="errorBorder px-2">
              <DAutocompleteUser
                v-model="item.userId"
                :project-key="projectKey"
                :preselect="item.userProfile"
                ref="autocompleteUserRef"
                @userChanged="userChanged"
                :onlyInternalUsers="isTypeOwnerSelected"
                :readonly="!ownerRemaining"
                required />
            </v-col>
            <v-col cols="12" xs="12" class="errorBorder px-2">
              <v-select
                :items="userTypes"
                v-model="item.userType"
                class="required"
                :rules="rulesSelectbox"
                v-bind:menu-props="{location: 'bottom'}"
                :label="t('COL_USER_TYPE')"
                variant="outlined"
                hide-details="auto"
                :disabled="!ownerRemaining || isResponsible"
                required>
                <template v-slot:item="{item, props}">
                  <v-list-item
                    v-bind="{...props, title: undefined}"
                    :disabled="!isInternalUser && item.raw === UserType.OWNER">
                    <template v-if="!isInternalUser && item.raw === UserType.OWNER">
                      <span>{{ item.raw }} {{ t('USER_DIALOG_OWNER_ONLY_FOR_INTERNAL') }}</span>
                    </template>
                    <template v-else>
                      {{
                        item.raw === UserType.OWNER
                          ? `${item.raw} ${t('USER_DIALOG_OWNER_ONLY_FOR_INTERNAL')}`
                          : item.raw
                      }}
                      <tooltip :text="t(`HELP_USER_TYPE_${item.raw}`)"></tooltip>
                    </template>
                  </v-list-item>
                </template>
              </v-select>
            </v-col>
            <v-col cols="12" xs="12" class="px-0">
              <v-text-field
                autocomplete="off"
                variant="outlined"
                hide-details="auto"
                v-model="item.comment"
                class="pr-2 pl-2"
                :label="t('COL_USER_COMMENT')"
                :rules="activeRules.comment" />
            </v-col>
            <v-col cols="12" xs="12">
              <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip">
                <template v-slot:activator="{props}">
                  <div v-bind="props">
                    <v-checkbox
                      v-model="item.responsible"
                      :disabled="isResponsible || item.userType != UserType.OWNER"
                      :label="t('COL_USER_ROLE_RESPONSIBLE')"
                      hide-details
                      class="mt-0 pt-1" />
                  </div>
                </template>
                <span>{{ t('TT_RESPONSIBLE') }}</span>
              </v-tooltip>
            </v-col>
          </v-row>
          <v-row v-if="item.userType" align="start" class="pa-1">
            <v-col cols="1" xs="1">
              <v-icon>mdi-information-outline</v-icon>
            </v-col>
            <v-col cols="11" xs="11">
              <span class="text-body-2">{{ t(`HELP_USER_TYPE_${item.userType}`) }}</span>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="justify-end pr-7">
          <DCActionButton
            isDialogButton
            size="small"
            variant="text"
            @click="close"
            class="mr-5"
            :text="t('BTN_CANCEL')" />
          <v-btn @click="doDialogAction" color="primary" variant="flat">
            <span v-if="mode === 'create'" class="font-weight-bold">{{ t('NP_DIALOG_BTN_CREATE') }}</span>
            <span v-else class="font-weight-bold">{{ t('NP_DIALOG_BTN_SAVE') }}</span>
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
</template>
