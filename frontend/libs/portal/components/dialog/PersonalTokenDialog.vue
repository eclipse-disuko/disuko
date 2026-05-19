<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {CreatePersonalTokenRequest} from '@disclosure-portal/model/PersonalToken';
import profileService from '@disclosure-portal/services/profile';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import TextArea from '@shared/components/disco/TextArea.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {useClipboard} from '@shared/utils/clipboard';
import dayjs from 'dayjs';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info: snack} = useSnackbar();
const {copyToClipboard} = useClipboard();

const emit = defineEmits(['reload']);

const isVisible = ref(false);
const isTokenDisplay = ref(false);
const description = ref('');
const createdToken = ref('');
const dialog = ref<DiscoForm | null>(null);

const isDatePickerVisible = ref(false);
const selectedDate = ref('');

const formattedDate = computed(() => {
  if (!selectedDate.value) return '';
  return dayjs(selectedDate.value).format('DD.MM.YYYY');
});

const minDate = computed(() => dayjs().add(1, 'day').format('YYYY-MM-DD'));
const maxDate = computed(() => dayjs().add(2, 'year').format('YYYY-MM-DD'));

const rules = {
  expiry: [(v: string) => !!v || t('FIELD_REQUIRED')],
};

const open = () => {
  isTokenDisplay.value = false;
  createdToken.value = '';
  description.value = '';
  selectedDate.value = dayjs().add(30, 'day').format('YYYY-MM-DD');
  isVisible.value = true;
};

const doCreate = async () => {
  if (isTokenDisplay.value) {
    isVisible.value = false;
    return;
  }

  await nextTick(async () => {
    dialog.value?.validate().then(async (info) => {
      if (!info.valid) return;

      const requestData: CreatePersonalTokenRequest = {
        description: description.value,
        expiry: dayjs(selectedDate.value).toISOString(),
      };

      const response = await profileService.createToken(requestData);
      createdToken.value = response.token;
      isTokenDisplay.value = true;
      snack(t('PERSONAL_TOKEN_CREATED'));
      emit('reload');
    });
  });
};

const copyTokenToClipboard = () => {
  if (createdToken.value) {
    copyToClipboard(createdToken.value);
  }
};

const dialogConfig = computed(() => {
  if (isTokenDisplay.value) {
    return {title: t('PERSONAL_TOKEN_ISSUED_TITLE'), primaryButton: {text: t('BTN_CLOSE')}};
  }
  return {
    title: t('PERSONAL_TOKEN_ADD_TITLE'),
    primaryButton: {text: t('NP_DIALOG_BTN_CREATE')},
    secondaryButton: {text: t('BTN_CANCEL')},
  };
});

defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" persistent :width="isTokenDisplay ? '600' : '500'">
    <DialogLayout
      :config="dialogConfig"
      @primary-action="doCreate"
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
      <v-form ref="dialog" @submit.prevent="doCreate">
        <Stack>
          <div v-if="isTokenDisplay && createdToken">
            <v-text-field
              autocomplete="off"
              :label="t('PERSONAL_TOKEN_LABEL')"
              :model-value="createdToken"
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
            <v-alert type="warning" color="warning" class="mt-2">
              {{ t('PERSONAL_TOKEN_SECURITY_WARNING') }}
            </v-alert>
          </div>
          <div v-if="!isTokenDisplay">
            <TextArea
              v-model="description"
              :label="t('NP_DIALOG_TF_DESCRIPTION')"
              variant="outlined"
              auto-grow
              max-rows="4"
              class="no-resize-textarea" />
          </div>
          <div v-if="!isTokenDisplay">
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
                :min="minDate"
                :max="maxDate"
                @update:model-value="isDatePickerVisible = false" />
            </v-menu>
          </div>
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
</template>
