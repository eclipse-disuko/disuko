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

defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" scrollable width="1100">
    <ReactiveDialogLayout :config="dialogConfig" @primary-action="save" @secondary-action="close" @close="close">
      <v-form ref="dialogRef" @submit.prevent="save">
        <div class="flex gap-6">
          <Stack class="flex-1">
            <v-text-field
              v-model="form.subject"
              :label="t('MAIL_TEMPLATE_SUBJECT')"
              variant="outlined"
              hide-details="auto" />
            <v-textarea
              v-model="form.message"
              :label="t('MAIL_TEMPLATE_MESSAGE')"
              variant="outlined"
              no-resize
              rows="8" />
            <v-text-field v-model="form.bcc" :label="t('MAIL_TEMPLATE_BCC')" variant="outlined" hide-details="auto" />
            <v-text-field v-model="form.cc" :label="t('MAIL_TEMPLATE_CC')" variant="outlined" hide-details="auto" />
          </Stack>
          <div v-if="values.length" class="values-table pa-3 w-90 shrink-0">
            <div class="mb-2 flex items-center gap-1">
              <v-icon icon="mdi-help-circle-outline" color="primary" size="small" />
              <span class="text-body-2 font-weight-medium">{{ t('MAIL_TEMPLATE_VALUES') }}</span>
            </div>
            <v-data-table
              density="compact"
              :items="values"
              :headers="valuesHeaders"
              hide-default-footer
              disable-sort
              item-value="key">
              <template #[`item.key`]="{item}">
                <code>{{ item.key }}</code>
              </template>
            </v-data-table>
          </div>
        </div>
      </v-form>
    </ReactiveDialogLayout>
  </v-dialog>
</template>

<style scoped>
.values-table {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.38) !important;
  border-radius: 4px;
}
</style>
