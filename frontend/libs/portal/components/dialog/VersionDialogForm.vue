<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {DialogVersionFormConfig} from '@disclosure-portal/components/dialog/DialogConfigs';
import ProjectVersionPostRequest from '@disclosure-portal/model/ProjectVersionPostRequest';
import versionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const sbomStore = useSbomStore();
const projectStore = useProjectStore();
const {info: snack} = useSnackbar();
const {minMax, longText} = useRules();

const isVisible = ref(false);
const req = ref(new ProjectVersionPostRequest());
const versionDialog = ref<DiscoForm | null>(null);
const title = ref('');
const confirmText = ref('');
const config = ref<DialogVersionFormConfig>({} as DialogVersionFormConfig);

const rules = {
  name: minMax(t('NPV_DIALOG_TF_NAME'), 3, 80, false),
  description: longText(t('NP_DIALOG_TF_DESCRIPTION')),
};

const dialogConfig = computed(() => ({
  title: t(title.value),
  primaryButton: {text: t(confirmText.value)},
  secondaryButton: {text: t('BTN_CANCEL')},
}));

const open = (newConf: DialogVersionFormConfig) => {
  config.value = newConf;
  versionDialog.value?.reset();
  title.value = newConf.version ? 'NPV_DIALOG_EDIT_TITLE' : 'NPV_DIALOG_TITLE';
  confirmText.value = newConf.version ? 'NP_DIALOG_BTN_EDIT' : 'NP_DIALOG_BTN_CREATE';
  if (newConf.version) {
    req.value = ProjectVersionPostRequest.toProjectVersionPostRequest(newConf.version);
  }
  isVisible.value = true;
};

const doDialogAction = async () => {
  await nextTick();
  const info = await versionDialog.value?.validate();
  if (!info?.valid) {
    return;
  }
  if (config.value.version) {
    await versionService.updateProjectVersion(projectStore.currentProject!._key, config.value.version._key, req.value);
    snack(t('DIALOG_version_edit_success'));
  } else {
    await versionService.createVersion(projectStore.currentProject!._key, req.value);
    snack(t('DIALOG_version_create_success'));
  }
  await sbomStore.fetchAllSBOMsFlat(true);
  versionDialog.value?.reset();
  isVisible.value = false;
};

const close = () => {
  isVisible.value = false;
};

defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" content-class="large" scrollable width="500">
    <DialogLayout :config="dialogConfig" @primary-action="doDialogAction" @secondary-action="close" @close="close">
      <v-form ref="versionDialog" @submit.prevent="doDialogAction">
        <Stack>
          <v-text-field
            autocomplete="off"
            class="required"
            v-model="req.name"
            :rules="rules.name"
            :label="t('NPV_DIALOG_TF_NAME')"
            autofocus
            hide-details="auto"
            variant="outlined" />
          <v-textarea
            no-resize
            v-model="req.description"
            :rules="rules.description"
            :label="t('NP_DIALOG_TF_DESCRIPTION')"
            :counter="1000"
            hide-details="auto"
            variant="outlined" />
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
</template>
