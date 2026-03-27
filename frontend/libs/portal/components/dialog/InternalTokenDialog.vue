// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {capabilityMap, type InternalToken, type InternalTokenRequest} from '@disclosure-portal/model/InternalToken';
import adminService from '@disclosure-portal/services/admin';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import Select from '@shared/components/disco/Select.vue';
import TextArea from '@shared/components/disco/TextArea.vue';
import TextField from '@shared/components/disco/TextField.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {useClipboard} from '@shared/utils/clipboard';
import dayjs from 'dayjs';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info: snack} = useSnackbar();
const {longText} = useRules();
const {copyToClipboard} = useClipboard();

const emit = defineEmits(['reload']);

const isVisible = ref(false);
const isEdit = ref(false);
const isTokenDisplay = ref(false);
const internalToken = ref({} as InternalTokenRequest);
const dialog = ref<DiscoForm | null>(null);
const title = ref('');
const confirmText = ref('');

const isDatePickerVisible = ref(false);
const selectedDate = ref('');

// Format the displayed date
const formattedDate = computed(() => {
  if (!selectedDate.value) return '';
  return dayjs(selectedDate.value).format('DD.MM.YYYY');
});

const minDate = computed(() => dayjs().format('YYYY-MM-DD'));
const maxDate = computed(() => dayjs().add(7, 'day').format('YYYY-MM-DD'));

// Capability options for multi-select
const capabilityOptions = computed(() => [
  {name: t('CAPABILITY_STATISTICS_CSV'), value: capabilityMap.StatisticsCSV},
  {name: t('CAPABILITY_CUSTOM_LICENSES'), value: capabilityMap.CustomLicenses},
]);

const rules = {
  name: [(v: string) => !!v || t('FIELD_REQUIRED')],
  description: longText(t('NP_DIALOG_TF_DESCRIPTION')),
  expiry: [(v: string) => !!v || t('FIELD_REQUIRED')],
};

const open = (existing?: InternalToken) => {
  isTokenDisplay.value = false;

  if (existing) {
    isEdit.value = true;
    internalToken.value = {
      _key: existing._key,
      name: existing.name,
      description: existing.description,
      expiry: existing.expiry ? existing.expiry.split('T')[0] : '',
      capabilities: existing.capabilities || [],
    } as InternalTokenRequest;
  } else {
    isEdit.value = false;
    selectedDate.value = dayjs().format('YYYY-MM-DD');

    // Only reset values if they're empty (preserve existing input)
    if (!internalToken.value.name && !internalToken.value.description && !internalToken.value.expiry) {
      internalToken.value = {
        name: '',
        description: '',
        expiry: '',
        capabilities: [],
      } as InternalTokenRequest;
    }
  }

  title.value = existing ? 'INTERNAL_TOKEN_RENEW_TITLE' : 'INTERNAL_TOKEN_ADD_TITLE';
  confirmText.value = existing ? 'BTN_RENEW' : 'NP_DIALOG_BTN_CREATE';
  isVisible.value = true;
};

const showToken = (tokenData: InternalToken, isRenewal: boolean = false) => {
  isTokenDisplay.value = true;
  isEdit.value = true; // Use edit mode for read-only display
  internalToken.value = {
    _key: tokenData._key,
    name: tokenData.name,
    description: tokenData.description,
    expiry: tokenData.expiry ? tokenData.expiry.split('T')[0] : '',
    capabilities: tokenData.capabilities || [],
    token: tokenData.token, // Add the token to display
  } as InternalTokenRequest & {token: string};

  title.value = isRenewal ? 'DLG_RENEWAL_TITLE' : 'UM_DIALOG_TITLE_TOKEN_ISSUED';
  confirmText.value = 'BTN_CLOSE';
  isVisible.value = true;
};

const doDialogAction = async () => {
  if (isTokenDisplay.value) {
    // If in token display mode, just close the dialog
    isVisible.value = false;
    return;
  }

  await nextTick(async () => {
    dialog.value?.validate().then(async (info) => {
      if (!info.valid) {
        return;
      }

      if (isEdit.value) {
        const response = await adminService.renewInternalToken(internalToken.value._key!);
        const renewedToken = response.data;

        snack(t('DIALOG_basicauth_renew_success'));
        emit('reload');

        // Show the renewed token
        setTimeout(() => {
          showToken(renewedToken, true);
        }, 100);
      } else {
        // Create request data for new tokens
        const requestData: InternalTokenRequest = {
          name: internalToken.value.name,
          description: internalToken.value.description,
          capabilities: internalToken.value.capabilities || [],
          expiry: selectedDate.value ? dayjs(selectedDate.value).toISOString() : '',
          revoked: false as boolean,
        };

        const response = await adminService.createInternalToken(requestData as InternalToken);
        const createdToken = response.data;

        snack(t('DIALOG_basicauth_create_success'));
        emit('reload');

        // Clear form after successful creation
        clearForm();

        // Show the created token
        setTimeout(() => {
          showToken(createdToken, false);
        }, 100);
      }

      if (!isTokenDisplay.value) {
        isVisible.value = false;
      }
    });
  });
};

const copyTokenToClipboard = async () => {
  const token = internalToken.value.token;
  if (token) {
    copyToClipboard(token);
  }
};

const clearForm = () => {
  selectedDate.value = dayjs().format('YYYY-MM-DD');
  isDatePickerVisible.value = false;

  internalToken.value = {
    name: '',
    description: '',
    expiry: '',
    capabilities: [],
  } as InternalTokenRequest;
  dialog.value?.reset();
};

const dialogConfig = computed(() => {
  if (isTokenDisplay.value) {
    return {title: t(title.value), primaryButton: {text: t(confirmText.value)}};
  }
  return {
    title: t(title.value),
    primaryButton: {text: t(confirmText.value)},
    secondaryButton: {text: t('BTN_CANCEL')},
  };
});

defineExpose({
  open,
  showToken,
  clearForm,
});
</script>

<template>
  <v-dialog v-model="isVisible" content-class="large" persistent scrollable :width="isTokenDisplay ? '600' : '500'">
    <DialogLayout
      :config="dialogConfig"
      @primary-action="doDialogAction"
      @secondary-action="isVisible = false"
      @close="isVisible = false">
      <template #left v-if="isTokenDisplay">
        <DCActionButton
          size="small"
          is-dialog-button
          variant="outlined"
          @click="copyTokenToClipboard"
          :text="t('BTN_COPY_TOKEN_CONTENT')"
          icon="mdi-content-copy" />
      </template>
      <v-form ref="dialog" @submit.prevent="doDialogAction">
        <Stack>
          <div v-if="isTokenDisplay && internalToken.token">
            <v-text-field
              autocomplete="off"
              :label="t('NP_DIALOG_HINT_TOKEN')"
              :model-value="internalToken.token"
              readonly
              variant="outlined"
              hide-details
              type="text"
              class="mb-2">
              <template v-slot:append-inner>
                <v-btn
                  icon="mdi-content-copy"
                  variant="text"
                  size="small"
                  @click="copyTokenToClipboard"
                  :title="t('TT_COPY_TOKEN')" />
              </template>
            </v-text-field>
            <div class="text-caption text-medium-emphasis">
              {{ t('NP_DIALOG_HINT_TOKEN') }}
            </div>
          </div>
          <TextField
            v-model="internalToken.name"
            variant="outlined"
            :rules="isTokenDisplay ? [] : rules.name"
            :readonly="isEdit"
            :label="t('NPV_DIALOG_TF_NAME')"
            required />
          <div v-if="isTokenDisplay || isEdit">
            <TextField :model-value="internalToken.expiry" variant="outlined" readonly :label="t('TKN_EXPIRATION')" />
          </div>
          <div v-else>
            <v-menu v-model="isDatePickerVisible" :close-on-content-click="false" :offset-y="true" min-width="auto">
              <template v-slot:activator="{props}">
                <v-text-field
                  autocomplete="off"
                  variant="outlined"
                  class="cursor-pointer"
                  v-model="formattedDate"
                  :label="t('TKN_EXPIRATION')"
                  readonly
                  hide-details="auto"
                  :rules="rules.expiry"
                  required
                  v-bind="props">
                  <template v-slot:append>
                    <v-icon color="primary">mdi-calendar</v-icon>
                  </template>
                </v-text-field>
              </template>
              <v-date-picker
                v-model="selectedDate"
                first-day-of-week="1"
                color="primary"
                show-current
                border="0"
                :elevation="0"
                :min="minDate"
                :max="maxDate"
                @update:model-value="isDatePickerVisible = false" />
            </v-menu>
          </div>
          <Select
            v-model="internalToken.capabilities"
            :items="capabilityOptions"
            :label="t('CAPABILITIES')"
            :readonly="isEdit"
            variant="outlined"
            multiple
            chips
            :closable-chips="!isEdit" />
          <TextArea
            v-model="internalToken.description"
            :rules="isTokenDisplay ? [] : rules.description"
            :label="t('NP_DIALOG_TF_DESCRIPTION')"
            :readonly="isEdit"
            :counter="isTokenDisplay ? undefined : 1000"
            variant="outlined"
            auto-grow
            max-rows="4"
            class="no-resize-textarea" />
          <v-alert v-if="isTokenDisplay" type="warning" color="warning">
            {{
              t(
                'NP_DIALOG_TOKEN_SECURITY_WARNING',
                'This token will only be shown once. Please copy and store it securely.',
              )
            }}
          </v-alert>
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
</template>
