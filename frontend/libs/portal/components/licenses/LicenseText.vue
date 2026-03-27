// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {useI18n} from 'vue-i18n';

defineProps<{
  license: {
    text?: string;
  };
}>();
const {t} = useI18n();
const translateContentInNewTab = () => {
  const translationUrl = t('TRANSLATION_URL_FORMAT');
  window.open(translationUrl, '_blank');
};
</script>

<template>
  <TableLayout has-title has-tab>
    <template #table>
      <div v-if="license.text" class="position-relative fill-height overflow-auto">
        <div class="pa-4" style="max-width: 1280px">
          <p class="license-text text-justify">{{ license.text }}</p>
        </div>
        <div class="position-absolute top-0 right-0 mr-8">
          <DCopyClipboardButton
            class="ma-2"
            :tableButton="true"
            :hint="t('TT_COPY_LICENSE_TEXT')"
            :content="license.text"
          />
          <DIconButton
            class="ma-2"
            icon="mdi-translate"
            :hint="t('TT_TRANSLATE_IN_SEPARATE_TAB')"
            @clicked="translateContentInNewTab"
          />
        </div>
      </div>
    </template>
  </TableLayout>
</template>
