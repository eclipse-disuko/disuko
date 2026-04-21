<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import ProjectService, {DocumentDownloadVersion} from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {info: snack} = useSnackbar();
const {t} = useI18n();
const projectStore = useProjectStore();

const currentProject = computed(() => projectStore.currentProject!);

const show = ref(false);
const appId = ref('');

const open = (approvalGuid: string) => {
  show.value = true;
  appId.value = approvalGuid;
};

const close = async () => {
  show.value = false;
  if (!currentProject.value.isGroup) {
    await projectStore.fetchProjectByKey(currentProject.value._key);
  }
};

const downloadDiscoDocument = (lang: string) => {
  const link = document.createElement('a');
  link.click();
  link.target = '_blank';
  link.rel = 'noopener noreferrer';
  ProjectService.downloadDocumentByTask(
    currentProject.value._key,
    appId.value,
    'disclosure',
    lang,
    DocumentDownloadVersion.None,
  )
    .then((res) => {
      link.download = `Disclosure_${currentProject.value.name.replaceAll(' ', '_')}-${lang}.pdf`;
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
  <v-dialog v-model="show" content-class="large" scrollable width="700" persistent>
    <DialogLayout :config="{title: t('DOWNLOAD_DOCUMENT')}" @close="close">
      <Stack>
        <div v-html="t('DOWNLOAD_INITIAL_DOCUMENT')"></div>
        <Stack direction="row">
          <DCActionButton
            text="Download (de)"
            icon="mdi-download"
            :hint="t('TT_TAD_disco_doc')"
            @click="downloadDiscoDocument('de')"
            class="none-uppercase" />
          <DCActionButton
            text="Download (en)"
            icon="mdi-download"
            :hint="t('TT_TAD_disco_doc')"
            @click="downloadDiscoDocument('en')"
            class="none-uppercase" />
        </Stack>
      </Stack>
    </DialogLayout>
  </v-dialog>
</template>
