<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import ProjectService, {DocumentDownloadVersion} from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {info: snack} = useSnackbar();
const show = ref(false);
const projectKey = ref('');
const projectName = ref('');
const appId = ref('');
const form = ref(null);
const {t} = useI18n();
const projectStore = useProjectStore();

const open = (approvalGuid: string) => {
  show.value = true;
  projectKey.value = projectStore.currentProject!._key;
  projectName.value = projectStore.currentProject!.name;
  appId.value = approvalGuid;
};

const close = async () => {
  show.value = false;
  if (!projectStore.currentProject?.isGroup) {
    await projectStore.fetchProjectByKey(projectKey.value);
  }
};

const downloadDiscoDocument = (lang: string) => {
  const link = document.createElement('a');
  link.click();
  link.target = '_blank';
  link.rel = 'noopener noreferrer';
  ProjectService.downloadDocumentByTask(projectKey.value, appId.value, 'disclosure', lang, DocumentDownloadVersion.None)
    .then((res) => {
      link.download = `Disclosure_${projectName.value.replaceAll(' ', '_')}-${lang}.pdf`;
      link.href = URL.createObjectURL(new Blob([res.data as BlobPart]));
      link.click();
    })
    .catch((e) => {
      snack(t('Error downloading disco document'));
      console.error('Error downloading disco document', e);
    });
};
defineExpose({open});
</script>

<template>
  <v-form ref="form">
    <v-dialog v-model="show" content-class="large" scrollable width="700" persistent>
      <v-card class="pa-8 dDialog">
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <span class="text-h5">{{ t('DOWNLOAD_DOCUMENT') }}</span>
            </v-col>
            <v-col cols="2" align="right">
              <DCloseButton @click="close" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text class="pt-2">
          <v-row class="py-2">
            <!-- v-html is used to render because DOWNLOAD_INTITIAL_DOCUMENT contains html -->
            <v-col cols="12" xs="12" class="pb-10" v-html="t('DOWNLOAD_INITIAL_DOCUMENT')"> </v-col>
            <v-col xs="6">
              <DCActionButton
                text="Download (de)"
                icon="mdi-download"
                :hint="t('TT_TAD_disco_doc')"
                @click="downloadDiscoDocument('de')"
                class="none-uppercase" />
            </v-col>
            <v-col xs="6">
              <DCActionButton
                text="Download (en)"
                icon="mdi-download"
                :hint="t('TT_TAD_disco_doc')"
                @click="downloadDiscoDocument('en')"
                class="none-uppercase" />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <DCActionButton isDialogButton @click="close" :text="t('BTN_OK')" variant="elevated"> </DCActionButton>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
</template>
