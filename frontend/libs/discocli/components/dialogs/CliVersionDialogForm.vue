<script setup lang="ts">
import type {VersionRequest} from '@cli/models/Version';
import {projectService} from '@cli/services/projectService';
import {useAppStore} from '@cli/stores/app';
import type {DialogVersionFormConfig} from '@disclosure-portal/components/dialog/DialogConfigs';
import type {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import useSnackbar from '@shared/composables/useSnackbar';
import type {DialogLayoutConfig} from '@shared/layouts/DialogLayout.vue';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const appStore = useAppStore();
const {info: snack} = useSnackbar();
const {minMax, longText} = useRules();

const isVisible = ref(false);
const req = ref<VersionRequest>({name: '', description: ''});
const versionDialog = ref<DiscoForm | null>(null);
const title = ref('');
const confirmText = ref('');
const config = ref<DialogVersionFormConfig>({} as DialogVersionFormConfig);

const rules = {
  name: minMax(t('NPV_DIALOG_TF_NAME'), 3, 80, false),
  description: longText(t('NP_DIALOG_TF_DESCRIPTION')),
};

const open = (newConf: DialogVersionFormConfig) => {
  config.value = newConf;
  versionDialog.value?.reset();
  title.value = newConf.version ? 'NPV_DIALOG_EDIT_TITLE' : 'NPV_DIALOG_TITLE';
  confirmText.value = newConf.version ? 'NP_DIALOG_BTN_EDIT' : 'NP_DIALOG_BTN_CREATE';
  if (newConf.version) {
    req.value = {
      name: newConf.version.name,
      description: newConf.version.description || '',
    };
  } else {
    req.value = {name: '', description: ''};
  }
  isVisible.value = true;
};

const closeDialog = () => {
  versionDialog.value?.reset();
  isVisible.value = false;
};

const doDialogAction = async () => {
  const formValidation = await versionDialog.value?.validate();
  if (!formValidation?.valid) {
    return;
  }

  const projectUuid = config.value.projectID || appStore.getCurrentProject()?.uuid;
  if (!projectUuid) {
    snack(t('DIALOG_token_missing_error'));
    return;
  }

  try {
    await projectService.createProjectVersion(projectUuid, req.value);
    snack(t(config.value.version ? 'DIALOG_CHANNEL_EDIT_SUCCESS' : 'DIALOG_CHANNEL_CREATE_SUCCESS'));
    await appStore.refetchCurrentProject(projectUuid);
    closeDialog();
  } catch (error) {
    snack(t('DIALOG_version_error'));
  }
};

defineExpose({open});

const dialogConfig = computed(
  (): DialogLayoutConfig => ({
    title: t(title.value),
    secondaryButton: {text: t('BTN_CANCEL')},
    primaryButton: {text: t(confirmText.value)},
  }),
);
</script>

<template>
  <v-dialog v-model="isVisible" content-class="large" scrollable width="500">
    <v-form ref="versionDialog">
      <DialogLayout
        :config="dialogConfig"
        @close="closeDialog"
        @secondary-action="closeDialog"
        @primary-action="doDialogAction">
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
      </DialogLayout>
    </v-form>
  </v-dialog>
</template>
