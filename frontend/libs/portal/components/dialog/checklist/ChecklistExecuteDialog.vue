// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {Checklist} from '@disclosure-portal/model/Checklist';
import projectService from '@disclosure-portal/services/projects';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const {t} = useI18n();
const {info: snack} = useSnackbar();
const route = useRoute();
const appStore = useAppStore();

const currentProject = computed(() => useProjectStore().currentProject!);
const version = computed(() => appStore.getCurrentVersion);
const sbom = computed(() =>
  Array.isArray(route.params.currentSbom) ? route.params.currentSbom[0] : route.params.currentSbom,
);

const isVisible = ref(false);
const dialog = ref<DiscoForm | null>(null);
const saving = ref(false);
const lists = ref<Checklist[]>([]);
const selected = ref<string | null>(null);

const emits = defineEmits(['reload']);

const dialogConfig = computed(() => ({
  title: t('CHECKLIST_EXECUTE_DT'),
  primaryButton: {text: t('CHECKLIST_EXECUTE'), disabled: saving.value, loading: saving.value},
  secondaryButton: {text: t('BTN_CANCEL'), disabled: saving.value},
}));

const open = (l: Checklist[]) => {
  lists.value = [...l].sort((a, b) => a.name.localeCompare(b.name));
  if (lists.value.length > 0) {
    selected.value = lists.value[0]._key;
  }
  isVisible.value = true;
};

const doDialogAction = async () => {
  await nextTick();
  const info = await dialog.value?.validate();
  if (!info?.valid) {
    return;
  }
  saving.value = true;
  await projectService.executeChecklistChecks(currentProject.value._key, version.value._key, sbom.value, {
    ids: [selected.value as string],
  });
  snack(t('CHECKLIST_EXECUTED'));
  emits('reload');
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
  <v-dialog v-model="isVisible" content-class="large" persistent scrollable width="500">
    <DialogLayout :config="dialogConfig" @primary-action="doDialogAction" @secondary-action="close" @close="close">
      <v-form ref="dialog" @submit.prevent="doDialogAction">
        <Stack>
          <v-select
            v-model="selected"
            :items="lists"
            autofocus
            item-title="name"
            item-value="_key"
            variant="outlined"
            density="compact"
            transition="scale-transition">
          </v-select>
        </Stack>
      </v-form>
      <template #left>
        <div class="text-sm italic">
          <p>{{ t('RR_WARN_TEXT') }}</p>
        </div>
      </template>
    </DialogLayout>
  </v-dialog>
</template>
