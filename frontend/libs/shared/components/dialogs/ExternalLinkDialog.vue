// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

interface Props {
  url: string;
}

const props = defineProps<Props>();

const emits = defineEmits(['close']);

const isDialogVisible = defineModel<boolean>('isDialogVisible', {required: true});

const {t} = useI18n();

const externalLinkDialogConfig = computed(() => ({
  title: t('EXT_LINK_DIALOG_TITLE'),
  primaryButton: {text: t('EXT_LINK_DIALOG_CONTINUE_BTN')},
  secondaryButton: {text: t('BTN_CANCEL')},
  icon: 'mdi mdi-open-in-new',
}));

const openUrl = () => {
  window.open(props.url, '_blank');
  closeAction();
};

const closeAction = () => {
  isDialogVisible.value = false;
  emits('close');
};
</script>

<template>
  <v-dialog v-model="isDialogVisible" content-class="small" width="540">
    <DialogLayout
      :config="externalLinkDialogConfig"
      @close="closeAction"
      @secondaryAction="closeAction"
      @primaryAction="openUrl">
      <p>
        {{ t('EXT_LINK_DIALOG_TEXT') }}
        <br /><br />
        <i>{{ url }}</i>
      </p>
    </DialogLayout>
  </v-dialog>
</template>
