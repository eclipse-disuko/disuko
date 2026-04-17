<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <v-dialog v-model="show" width="700" scrollable max-height="500" persistent>
    <DialogLayout :config="dialogConfig" @close="closeDialog" @primary-action="closeDialog">
      <div class="position-relative">
        <v-list class="pa-0 noneBorder">
          <v-list-item v-for="(error, index) in formattedErrors" :key="index">
            {{ error }}
          </v-list-item>
        </v-list>
        <DCopyClipboardButton
          :tableButton="true"
          class="position-absolute top-0 right-0 mr-4"
          :hint="t('TT_COPY_TO_CLIPBOARD')"
          :content="rawErrors" />
      </div>
    </DialogLayout>
  </v-dialog>
</template>

<script setup lang="ts">
import DCopyClipboardButton from '@shared/components/disco/DCopyClipboardButton.vue';

import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const show = ref(false);
const rawErrors = ref('');

const formattedErrors = computed(() => (rawErrors.value ? rawErrors.value.split('\n').filter(Boolean) : []));

const dialogConfig = computed(() => ({
  title: t('VALIDATE_SCHEMA'),
  primaryButton: {text: t('BTN_OK')},
}));

const open = (errors: string) => {
  show.value = true;
  rawErrors.value = errors;
};

const closeDialog = () => {
  show.value = false;
};
defineExpose({open});
</script>
