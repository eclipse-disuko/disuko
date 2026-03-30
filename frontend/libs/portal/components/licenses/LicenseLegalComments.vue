// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {useI18n} from 'vue-i18n';

defineProps<{
  license: {
    meta?: {
      legalComments?: string;
    };
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
      <div class="position-relative fill-height overflow-auto">
        <div class="p-4 pt-12 text-justify [&_*]:text-[14px]">
          <Markdown
            v-if="license && license.meta && license.meta.legalComments"
            :text="license.meta?.legalComments"
          ></Markdown>
          <div v-else>{{ t('NO_COMMENT') }}</div>
        </div>
        <div class="position-absolute top-0 right-0 d-flex justify-start">
          <DCopyClipboardButton
            v-if="license && license.meta && license.meta.legalComments"
            class="ma-2"
            :tableButton="true"
            :hint="t('TT_COPY_LEGAL_COMMENTS')"
            :content="license.meta.legalComments"
          />
          <DIconButton
            v-if="license && license.meta && license.meta.legalComments"
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
