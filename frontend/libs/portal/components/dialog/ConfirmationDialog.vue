<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps({
  config: {
    type: Object as () => IConfirmationDialogConfig,
    required: false,
    default: null,
  },
});

const emit = defineEmits<{
  confirm: [config: IConfirmationDialogConfig];
}>();

const showDialog = defineModel<boolean>('showDialog', {required: false});

const {t} = useI18n();
const title = ref<string | undefined>(undefined);
const contextKey = ref<string | undefined>(undefined);
const key = ref('');
const name = ref('');
const description = ref('');
const extendedDetails = ref('');
const okButton = ref('');
const okButtonIsDisabled = ref(false);
const emphasiseText = ref('');
const emphasiseConfirmationText = ref('');
const emphasiseConfirmation = ref(true);
const methodConfig = ref<IConfirmationDialogConfig | null>(null);

watch(
  () => props.config,
  (newVal) => {
    if (newVal) {
      open(props.config);
    }
  },
);

const open = (config: IConfirmationDialogConfig | null) => {
  reset();
  if (config) {
    title.value = config.title;
    contextKey.value = config.contextKey;
    key.value = config.key;
    name.value = config.name;
    description.value = config.description;
    extendedDetails.value = config.extendedDetails ?? '';
    okButton.value = config.okButton ?? 'BTN_OK';
    okButtonIsDisabled.value = config.okButtonIsDisabled ?? false;
    emphasiseText.value = config.emphasiseText ?? '';
    emphasiseConfirmationText.value = config.emphasiseConfirmationText ?? '';
    if (emphasiseConfirmationText.value !== '') {
      emphasiseConfirmation.value = false;
    }
    showDialog.value = true;
    nextTick(() => {
      const buttonRef = document.querySelector('.confirmButton') as HTMLElement;
      if (buttonRef) {
        buttonRef.focus();
      }
    });
  }
};
const makeVisible = () => {
  showDialog.value = true;
};
const openWithoutDetails = (newDescription: string, okButtonText: string) => {
  methodConfig.value = {
    description: newDescription,
    okButton: okButtonText,
  } as IConfirmationDialogConfig;
  open(methodConfig.value);
};

const close = () => {
  showDialog.value = false;
  extendedDetails.value = '';
};

const confirm = () => {
  if (methodConfig.value) {
    emit('confirm', methodConfig.value);
    methodConfig.value = null;
    close();
    return;
  }
  if (contextKey.value) {
    emit('confirm', props.config!);
  } else {
    emit('confirm', props.config!);
  }
  close();
};

const reset = () => {
  extendedDetails.value = '';
  emphasiseConfirmation.value = true;
};

const dialogConfig = computed(() => {
  const cfg: any = {
    title: title.value ? t(title.value) : t('DLG_CONFIRMATION_TITLE'),
    secondaryButton: {text: t('BTN_CANCEL')},
  };
  if (okButton.value && !okButtonIsDisabled.value) {
    cfg.primaryButton = {text: t(okButton.value), disabled: !emphasiseConfirmation.value};
  }
  return cfg;
});

defineExpose({openWithoutDetails, makeVisible});
</script>
<template>
  <slot name="default" :open="open"> </slot>
  <v-dialog v-model="showDialog" content-class="small" width="800" max-width="500">
    <DialogLayout
      v-if="okButton"
      :config="dialogConfig"
      @primary-action="confirm"
      @secondary-action="close"
      @close="close">
      <v-card-text class="pa-0">
        <v-row>
          <v-col cols="12" v-if="okButtonIsDisabled"> {{ extendedDetails }}<br />{{ t(description) }}</v-col>
          <v-col cols="12" v-else-if="extendedDetails">
            {{ extendedDetails }}<br />{{ t(description) }}"{{ name }}"?
          </v-col>
          <v-col cols="12" v-else-if="name"> {{ t(description) }}"{{ name }}"?</v-col>
          <v-col cols="12" v-else>
            {{ t(description) }}
          </v-col>
          <v-col cols="12" v-if="emphasiseText">
            <span class="font-weight-bold">{{ t(emphasiseText) }}</span>
          </v-col>
        </v-row>
        <v-row>
          <div class="f-modal-alert">
            <div class="f-modal-icon f-modal-warning scaleWarning">
              <span class="f-modal-body pulseWarningIns"></span>
              <span class="f-modal-dot pulseWarningIns"></span>
            </div>
          </div>
        </v-row>
        <v-row v-if="emphasiseConfirmationText">
          <v-col cols="12" xs="12" class="px-0">
            <v-checkbox
              v-model="emphasiseConfirmation"
              :label="t(emphasiseConfirmationText)"
              hide-details
              class="mt-5 pt-1" />
          </v-col>
        </v-row>
      </v-card-text>
    </DialogLayout>
  </v-dialog>
</template>
