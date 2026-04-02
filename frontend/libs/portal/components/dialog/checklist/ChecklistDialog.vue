<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import icons from '@disclosure-portal/constants/icons';
import {Checklist} from '@disclosure-portal/model/Checklist';
import {useChecklistsStore} from '@disclosure-portal/stores/checklists.store';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info: snack} = useSnackbar();
const {minMax, longText} = useRules();
const checklistsStore = useChecklistsStore();
const labelStore = useLabelStore();

const isVisible = ref(false);
const isEdit = ref(false);
const item = ref(new Checklist());
const dialog = ref<DiscoForm>();
const saving = ref(false);

const rules = {
  name: minMax(t('NPV_DIALOG_TF_NAME'), 3, 80, false),
  description: longText(t('NP_DIALOG_TF_DESCRIPTION')),
  labels: [(values: string[]) => (!!values && values.length > 0) || t('IS_REQUIRED', {fieldName: 'Labels'})],
};

const dialogConfig = computed(() => ({
  title: t(isEdit.value ? 'CHECKLIST_DIALOG_EDIT_TITLE' : 'CHECKLIST_DIALOG_ADD_TITLE'),
  primaryButton: {
    text: t(isEdit.value ? 'NP_DIALOG_BTN_EDIT' : 'NP_DIALOG_BTN_CREATE'),
    disabled: saving.value,
    loading: saving.value,
  },
  secondaryButton: {text: t('BTN_CANCEL'), disabled: saving.value},
}));

const open = async (existing?: Checklist) => {
  if (existing) {
    isEdit.value = true;
    item.value = {
      ...existing,
    };
  } else {
    isEdit.value = false;
    item.value = new Checklist();
  }
  dialog.value?.reset();
  isVisible.value = true;
};

const doDialogAction = async () => {
  await nextTick();
  const info = await dialog.value?.validate();
  if (!info?.valid) {
    return;
  }
  saving.value = true;
  if (isEdit.value) {
    await checklistsStore.editChecklist(item.value);
    snack(t('DIALOG_CHECKLIST_UPDATE_SUCCESS'));
  } else {
    await checklistsStore.createChecklist(item.value);
    snack(t('DIALOG_CHECKLIST_CREATE_SUCCESS'));
  }
  saving.value = false;
  isVisible.value = false;
};

const close = () => {
  isVisible.value = false;
};

defineExpose({
  open,
});
</script>

<template>
  <v-dialog v-model="isVisible" content-class="large" persistent scrollable width="840">
    <DialogLayout :config="dialogConfig" @primary-action="doDialogAction" @secondary-action="close" @close="close">
      <v-form ref="dialog" @submit.prevent="doDialogAction">
        <Stack class="sm:flex-row justify-space-between gap-4">
          <Stack class="flex-grow-1 gap-4">
            <div>{{ t('HEADER_ENGLISH') }}</div>
            <TextField
              required
              v-model="item.name"
              :rules="rules.name"
              :label="t('LBL_NAME')"
              :persistent-placeholder="false"
              autofocus
              tabindex="1" />
            <TextArea
              no-resize
              tabindex="3"
              v-model="item.description"
              :rules="rules.description"
              :label="t('DESCRIPTION')"
              :persistent-placeholder="false"
              :rows="5"
              :counter="1000" />
          </Stack>
          <v-divider vertical class="mb-5 mt-2 hidden sm:block"></v-divider>
          <Stack class="flex-grow-1 gap-4">
            <div>{{ t('HEADER_GERMAN') }}</div>
            <TextField
              required
              v-model="item.nameDE"
              :rules="rules.name"
              :label="t('LBL_NAME')"
              :persistent-placeholder="false"
              autofocus
              tabindex="1" />
            <TextArea
              no-resize
              tabindex="3"
              v-model="item.descriptionDE"
              :rules="rules.description"
              :label="t('DESCRIPTION')"
              :persistent-placeholder="false"
              :rows="5"
              :counter="1000" />
          </Stack>
        </Stack>
        <v-select
          variant="outlined"
          hide-details="auto"
          v-model="item.policyLabels"
          item-title="name"
          item-value="_key"
          tabindex="5"
          clearable
          multiple
          :items="labelStore.policyLabels"
          :label="t('AL_DIALOG_SB_LABELS_OR')"
          :rules="rules.labels"
          v-bind:menu-props="{location: 'bottom'}">
          <template v-slot:chip="{item, props}">
            <DLabel closable :parentProps="props" :labelName="item.title" :iconName="icons.TAG" />
          </template>
        </v-select>
        <v-checkbox
          v-model="item.active"
          hide-details
          color="primary"
          :label="t('ACTIVE_FLAG')"
          class="shrink mt-0 pt-0" />
      </v-form>
    </DialogLayout>
  </v-dialog>
</template>
