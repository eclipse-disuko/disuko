// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import DCopyClipboardButton from '@shared/components/disco/DCopyClipboardButton.vue';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const emit = defineEmits(['close']);
const show = ref(false);
const config = ref(new ErrorDialogConfig());
const {t} = useI18n();

const open = (newConfig: ErrorDialogConfig) => {
  if (newConfig) {
    config.value = newConfig;
  }
  show.value = true;
};

const close = () => {
  show.value = false;
  emit('close', config.value);
};

const dialogConfig = computed(() => ({
  title: '⚠️ ' + config.value.title || t('DLG_ERROR_TITLE'),
  primaryButton: {text: t('BTN_OK')},
}));

defineExpose({open, close});
</script>
<template>
  <v-dialog v-model="show" width="500" scrollable max-height="500">
    <DialogLayout :config="dialogConfig" @primary-action="close" @close="close">
      <Stack>
        <Stack direction="row" align="center" justify="between" class="w-full">
          <span class="flex-1" v-html="config.description"></span>
          <DCopyClipboardButton
            v-if="config.copyDesc"
            :tableButton="true"
            :hint="t('TT_COPY_TO_CLIPBOARD_ERROR')"
            :content="config.description" />
        </Stack>

        <Stack v-if="config.reqId" direction="row" align="center" justify="between" class="w-full">
          <Stack class="flex-1 gap-0">
            <span v-html="t('DESC_ERROR')"></span>
            <span class="font-italic">{{ config.reqId }}</span>
          </Stack>
          <DCopyClipboardButton :tableButton="true" :hint="t('TT_COPY_TO_CLIPBOARD_REQID')" :content="config.reqId" />
        </Stack>
      </Stack>
    </DialogLayout>
  </v-dialog>
</template>
