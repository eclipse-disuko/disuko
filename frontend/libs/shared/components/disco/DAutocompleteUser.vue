<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {UserDto} from '@shared/types/Users';
import profileService from '@shared/services/profile';
import projectService from '@disclosure-portal/services/projects';
import {useUserStore} from '@disclosure-portal/stores/user';
import {RuleFunction} from '@disclosure-portal/types/rules';
import config from '@shared/utils/config';
import {debounce} from 'lodash';
import {computed, nextTick, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const props = withDefaults(
  defineProps<{
    preselect?: UserDto;
    projectKey?: string;
    onlyAlphanumeric?: boolean;
    required?: boolean;
    label?: string;
    readonly?: boolean;
    onlyInternalUsers?: boolean;
    active?: boolean | null;
  }>(),
  {
    active: true,
  },
);

const modelValue = defineModel<string>();

const emit = defineEmits<{
  (e: 'userChanged', user: UserDto): void;
}>();

const {t} = useI18n();
const userStore = useUserStore();
const showErrors = ref<boolean>(false);
const formRef = ref<VForm | null>(null);
const selectedUser = ref<UserDto | null>(null);
const searchFieldInput = ref('');
const items = ref<UserDto[]>([]);
const idle = ref(true);
const modeIsOnCreate = ref(false);
const activeRules = ref<RuleFunction[]>([]);
const searchedAndNotFound = ref(false);

const requiredClass = computed(() => (props.required ? 'required' : ''));
const searchValidatorRegex = /^[0-9A-Za-z-_. äöüÄÖÜß]+$/;

const debouncedSearchForUser = debounce(() => searchForUser(), 300);

onMounted(() => {
  // Initialize selected user if preselected
  if (props.preselect && props.preselect.user) {
    items.value.push(props.preselect);
    selectedUser.value = props.preselect;
  }
});

watch(
  () => props.preselect,
  (newVal) => {
    items.value = [];
    if (newVal && newVal.user) {
      items.value.push(newVal);
      selectedUser.value = newVal;
    } else {
      selectedUser.value = null;
    }
  },
);

watch(selectedUser, () => {
  if (selectedUser.value && selectedUser.value.user) {
    modelValue.value = selectedUser.value.user;
    emit('userChanged', selectedUser.value);
  } else {
    modelValue.value = '';
  }
});

const validate = async (): Promise<boolean> => {
  activeRules.value = [];

  if (props.required) {
    activeRules.value.push(() => {
      return (
        (selectedUser.value && !!selectedUser.value.user) ||
        (searchFieldInput.value && searchFieldInput.value.trim().length > 0) ||
        t('user_error_message')
      );
    });
  }

  if (props.onlyAlphanumeric) {
    activeRules.value.push(() => {
      return (
        !searchFieldInput.value ||
        (searchFieldInput.value && searchValidatorRegex.test(searchFieldInput.value)) ||
        t('user_search_error_message')
      );
    });
    activeRules.value.push(() => {
      if (!modeIsOnCreate.value) {
        if (!idle.value) {
          return true;
        }
        if (!searchFieldInput.value) {
          return true;
        }
      }
      if (selectedUser.value && selectedUser.value._key) {
        return true;
      }
      return t('user_error_message');
    });
  }

  await nextTick();
  const validationResult = await formRef.value?.validate();
  return validationResult?.valid ?? false;
};

const prepareForCreate = () => {
  selectedUser.value = null;
  items.value = [];
};

const validateOnCreate = async (): Promise<boolean> => {
  showErrors.value = true;
  modeIsOnCreate.value = true;
  const result = validate();
  modeIsOnCreate.value = false;
  return result;
};

const resetForm = () => {
  idle.value = true;
  searchFieldInput.value = '';
  items.value = [];
  formRef.value?.resetValidation();
};

// Expose method for external use on component's ref
defineExpose({
  validate,
  prepareForCreate,
  validateOnCreate,
  resetForm,
});

const onClickClearHandler = () => {
  searchFieldInput.value = '';
  items.value = [];
  searchedAndNotFound.value = false;
  resetForm();
};

const searchForUser = async () => {
  searchedAndNotFound.value = false;
  if (!searchFieldInput.value || searchFieldInput.value.length < 3) {
    return;
  }
  idle.value = false;
  await validate();

  const hasSpecialChars = searchFieldInput.value && !searchValidatorRegex.test(searchFieldInput.value);
  if (!hasSpecialChars) {
    let response;
    const activeFilter = props.active === null ? undefined : props.active;
    if (props.projectKey) {
      response = await projectService.getUsersBySearchFragment(props.projectKey, searchFieldInput.value, activeFilter);
    } else {
      response = await profileService.getUsersBySearchFragment(searchFieldInput.value, activeFilter);
    }
    items.value = response.data.filter((user: UserDto) => (props.onlyInternalUsers ? user.isInternal === true : true));
    if (items.value.length === 0) {
      searchedAndNotFound.value = true;
    }
    idle.value = true;
  }
};

const openMailto = () => {
  const profile = userStore.getProfile;
  const mailContent = t('INVITATION_MAIL_CONTENT')
    .replace(/%url/, config.SERVER_URL)
    .replace(/%user/, `${profile.forename} ${profile.lastname}`);
  window.open(encodeURI(mailContent));
};

function itemTitleFunction(item: UserDto) {
  return `${item.lastname}, ${item.forename} (${item.user})`;
}
</script>

<template>
  <v-form ref="formRef">
    <v-autocomplete
      v-model="selectedUser"
      v-model:search="searchFieldInput"
      :rules="showErrors ? activeRules : []"
      @update:search="debouncedSearchForUser"
      autocomplete="off"
      item-value="user"
      :item-title="itemTitleFunction"
      @click:clear="onClickClearHandler"
      append-icon=""
      return-object
      variant="outlined"
      :placeholder="label ? label : t('UM_DIALOG_USER_ID')"
      :label="label ? label : t('UM_DIALOG_USER_ID')"
      clearable
      hide-details="auto"
      :hide-no-data="!searchedAndNotFound"
      :no-filter="true"
      :required="required"
      :items="items"
      :class="requiredClass"
      :disabled="readonly">
      <template v-if="searchedAndNotFound" v-slot:no-data>
        <v-list-item>
          <v-list-item-title>
            <span>{{ t('USER_AUTOCOMPLETE_NO_DATA_TEXT') }}</span>
          </v-list-item-title>
          <v-list-item-action>
            <v-spacer />
            <DCActionButton
              class="my-2"
              icon="mdi-mail"
              :hint="t('TT_SEND_INVITE')"
              :text="t('BTN_SEND_MAIL')"
              @click.stop="openMailto" />
            <v-spacer />
          </v-list-item-action>
        </v-list-item>
      </template>
    </v-autocomplete>
  </v-form>
</template>
