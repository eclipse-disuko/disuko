<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {MailTemplate, UpdateMailTemplate} from '@disclosure-portal/model/MailTemplate';
import mailTemplatesService from '@disclosure-portal/services/mailtemplates.service';
import {DataTableHeader} from '@shared/types/table';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info} = useSnackbar();

const emit = defineEmits(['reload']);

const isVisible = ref(false);
const saving = ref(false);
const testing = ref(false);
const key = ref('');
const form = ref<UpdateMailTemplate>({subject: '', message: '', bcc: '', cc: ''});
const values = ref<{key: string; value: string}[]>([]);
const dialogRef = ref();

const valuesHeaders = computed((): DataTableHeader[] => [
  {title: t('MAIL_TEMPLATE_VALUE_KEY'), value: 'key', align: 'start'},
  {title: t('MAIL_TEMPLATE_VALUE_DESCRIPTION'), value: 'value', align: 'start'},
]);

const dialogConfig = ref({
  title: '',
  loading: false,
  primaryButton: '',
  secondaryButton: '',
});

const open = (item: MailTemplate) => {
  key.value = item._key;
  form.value = {subject: item.subject, message: item.message, bcc: item.bcc, cc: item.cc};
  values.value = Object.entries(item.values ?? {}).map(([k, v]) => ({key: k, value: v}));
  dialogConfig.value = {
    title: t('MAIL_TEMPLATE_DIALOG_TITLE'),
    loading: false,
    primaryButton: t('BTN_SAVE'),
    secondaryButton: t('BTN_CANCEL'),
  };
  dialogRef.value?.reset();
  isVisible.value = true;
};

const close = () => {
  isVisible.value = false;
};

const save = async () => {
  const valid = (await dialogRef.value?.validate())?.valid;
  if (!valid) return;
  saving.value = true;
  dialogConfig.value.loading = true;
  try {
    await mailTemplatesService.update(key.value, form.value);
    info(t('MAIL_TEMPLATE_SAVE_SUCCESS'));
    emit('reload');
    isVisible.value = false;
  } finally {
    saving.value = false;
    dialogConfig.value.loading = false;
  }
};

const test = async () => {
  const valid = (await dialogRef.value?.validate())?.valid;
  if (!valid) return;
  testing.value = true;
  try {
    await mailTemplatesService.test(key.value, form.value.message);
    info(t('MAIL_TEMPLATE_TEST_SUCCESS'));
  } finally {
    testing.value = false;
  }
};

defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" scrollable width="1100" height="700">
    <ReactiveDialogLayout :config="dialogConfig" @primary-action="save" @secondary-action="close" @close="close">
      <v-form ref="dialogRef" @submit.prevent="save">
        <div class="mb-3 flex gap-2">
          <v-text-field
            v-model="form.bcc"
            :label="t('MAIL_TEMPLATE_BCC')"
            :hint="t('MAIL_TEMPLATE_BCC_HINT')"
            variant="outlined"
            class="flex-1"
            persistent-hint
            hide-details="auto" />
          <v-text-field
            v-model="form.cc"
            :label="t('MAIL_TEMPLATE_CC')"
            :hint="t('MAIL_TEMPLATE_CC_HINT')"
            variant="outlined"
            class="flex-1"
            persistent-hint
            hide-details="auto" />
        </div>
        <v-text-field
          v-model="form.subject"
          :label="t('MAIL_TEMPLATE_SUBJECT')"
          variant="outlined"
          class="mb-6"
          hide-details="auto" />
        <div class="flex gap-3">
          <Stack class="flex-1">
            <v-textarea
              v-model="form.message"
              :label="t('MAIL_TEMPLATE_MESSAGE')"
              variant="outlined"
              no-resize
              rows="8" />
          </Stack>
          <div v-if="values.length" class="pa-3 w-90">
            <div class="mb-2 flex items-center gap-1">
              <v-icon icon="mdi-help-circle-outline" color="primary" size="small" />
              <span class="text-body-2 font-weight-medium">{{ t('MAIL_TEMPLATE_VALUES') }}</span>
            </div>
            <v-data-table
              density="compact"
              :items="values"
              :headers="valuesHeaders"
              hide-default-footer
              height="240"
              disable-sort
              item-value="key">
              <template #[`item.key`]="{item}">
                <code>{{ item.key }}</code>
              </template>
            </v-data-table>
          </div>
        </div>
      </v-form>
      <template #leftButtons>
        <DCActionButton
          is-dialog-button
          size="small"
          variant="text"
          :loading="testing"
          :text="t('MAIL_TEMPLATE_BTN_TEST')"
          @click="test" />
      </template>
    </ReactiveDialogLayout>
  </v-dialog>
</template>
