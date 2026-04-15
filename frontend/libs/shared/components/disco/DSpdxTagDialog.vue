<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {Tags} from '@disclosure-portal/constants/ruleValidations';
import projectService from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import useRules from '@disclosure-portal/utils/Rules';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const props = defineProps<{
  presetTag?: string;
  versionID: string;
  spdxID: string;
  spdxName: string;
  channelView?: boolean;
}>();

const {t} = useI18n();
const isVisible = ref(false);
const tag = ref('');
const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const form = ref<VForm | null>(null);
const {info: snack} = useSnackbar();

const showDialog = () => {
  if (props.presetTag) {
    tag.value = props.presetTag;
  }
  isVisible.value = true;
};

const close = () => {
  form.value?.reset();
  isVisible.value = false;
};

const doDialogAction = async () => {
  await nextTick();
  const info = await form.value?.validate();
  if (!info?.valid) {
    return;
  }
  try {
    await projectService.updateSpdxTag(projectModel.value._key, props.versionID, props.spdxID, tag.value);
    snack(t('DIALOG_SPDX_TAG_UPDATE_SUCCESS'));
    await sbomStore.fetchAllSBOMsFlat(true);
    form.value?.reset();
    isVisible.value = false;
  } catch (error) {
    console.error('Error updating SPDX tag:', error);
  }
};

const projectModel = computed(() => projectStore.currentProject!);

const activeRules = ref({
  tag: useRules().minMax(t('COL_SBOM_TAG'), Tags.TAG_MIN_LENGTH, Tags.TAG_MAX_LENGTH, false),
});

const dialogConfig = computed(() => ({
  title: t('SPDX_TAG_TITLE') + props.spdxName,
  secondaryButton: {text: t('BTN_CLOSE')},
  primaryButton: {text: t('NP_DIALOG_BTN_EDIT')},
}));
</script>

<template>
  <slot :showDialog="showDialog">
    <v-btn text="Replace me" size="small" color="primary" @click.stop="showDialog"></v-btn>
  </slot>
  <v-dialog v-model="isVisible" content-class="msmall" scrollable width="500">
    <v-form ref="form" @submit.prevent="doDialogAction">
      <DialogLayout :config="dialogConfig" @close="close" @secondaryAction="close" @primaryAction="doDialogAction">
        <Stack class="errorBorder">
          <v-text-field
            autocomplete="off"
            variant="outlined"
            :rules="activeRules.tag"
            v-model="tag"
            :label="t('COL_SBOM_TAG')"
            autofocus />
        </Stack>
      </DialogLayout>
    </v-form>
  </v-dialog>
</template>
