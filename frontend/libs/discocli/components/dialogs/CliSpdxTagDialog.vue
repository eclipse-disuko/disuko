<script setup lang="ts">
import {Tags} from '@disclosure-portal/constants/ruleValidations';
import useRules from '@disclosure-portal/utils/Rules';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';
import {VForm} from 'vuetify/components';
import {projectService as cliProjectService} from '../../../discocli/services/projectService';
import {useAppStore as useCliAppStore} from '../../../discocli/stores/app';

defineOptions({name: 'DSpdxTagDialog'});

interface Props {
  presetTag?: string;
  versionID: string;
  spdxID: string;
  spdxName: string;
  channelView?: boolean;
  isCliApp?: boolean;
}
const props = defineProps<Props>();

const {t} = useI18n();
const isVisible = ref(false);
const tag = ref('');
const cliAppStore = useCliAppStore();
const route = useRoute();
const dialog = ref<VForm | null>(null);
const {info: snack} = useSnackbar();

function showDialog() {
  if (props.presetTag) {
    tag.value = props.presetTag;
  }
  isVisible.value = true;
}

function reset() {
  if (props.presetTag) {
    tag.value = props.presetTag;
  } else {
    dialog.value?.reset();
  }
}

function close() {
  dialog.value?.reset();
  isVisible.value = false;
}

async function updateTag() {
  if (props.isCliApp) {
    const projectUuid = typeof route.params.id === 'string' ? route.params.id : null;
    if (!projectUuid) {
      snack(t('No project selected'));
      return;
    }
    await cliProjectService.updateSbomTag(projectUuid, props.versionID, props.spdxID, tag.value);
    await cliAppStore.fetchCurrentProject(projectUuid);
  }
}

async function doDialogAction() {
  await nextTick();
  const info = await dialog.value?.validate();
  if (!info.valid) {
    return;
  }
  try {
    await updateTag();
    snack(t('DIALOG_SBOM_TAG_UPDATE_SUCCESS'));
    dialog.value?.reset();
    isVisible.value = false;
  } catch (error) {
    console.error('Error updating tag:', error);
    snack(t('Error updating SBOM tag'));
  }
}

const activeRules = ref({
  tag: useRules().minMax(t('COL_SBOM_TAG'), Tags.TAG_MIN_LENGTH, Tags.TAG_MAX_LENGTH, false),
});
</script>

<template>
  <slot :showDialog="showDialog">
    <v-btn text="Replace me" size="small" color="primary" @click.stop="showDialog"></v-btn>
  </slot>
  <v-dialog v-model="isVisible" content-class="msmall" scrollable width="500">
    <v-form ref="dialog" @submit.prevent="doDialogAction">
      <v-card class="pa-8 dDialog" flat>
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <span class="text-h5">{{ t('SBOM_TAG_TITLE') + spdxName }}</span>
            </v-col>
            <v-col cols="2" align="right">
              <DCloseButton @click="close" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text class="pt-2">
          <v-row dense>
            <v-col cols="12" xs="12" class="errorBorder">
              <v-text-field
                autocomplete="off"
                variant="outlined"
                :rules="activeRules.tag"
                v-model="tag"
                :label="t('COL_SBOM_TAG')"
                autofocus />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <DCActionButton size="small" variant="text" @click="close" class="mr-5" :text="t('BTN_CLOSE')" />
          <DCActionButton size="small" variant="flat" @click="doDialogAction" :text="t('NP_DIALOG_BTN_EDIT')" />
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
