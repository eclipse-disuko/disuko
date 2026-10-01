<template>
  <v-dialog v-model="show" max-width="500" persistent @after-leave="resetDialog">
    <v-form ref="formRef">
      <DialogLayout :config="dialogConfig" @close="close" @primary-action="submit" @secondary-action="close">
        <v-text-field
          v-model="formData.url"
          label="URL"
          :placeholder="t('PH_EXTERNAL_SOURCE_URL')"
          variant="outlined"
          density="compact"
          :rules="[rules.required, rules.validURL, rules.maxLength]"
          hide-details="auto"
          class="mb-4"
          @keyup.enter="canSubmit && submit()" />
        <v-textarea
          v-model="formData.comment"
          :label="t('COL_DESCRIPTION')"
          :placeholder="t('PH_EXTERNAL_SOURCE_DESCRIPTION')"
          variant="outlined"
          auto-complete="off"
          spellcheck="false"
          density="compact"
          :rules="[rules.required, rules.maxLength]"
          hide-details="auto"
          rows="3" />
      </DialogLayout>
    </v-form>
  </v-dialog>
</template>

<script setup lang="ts">
import {createFormRules} from '@cli/utils/validationRules';
import DialogLayout from '@shared/layouts/DialogLayout.vue';
import {computed, defineExpose, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import type {VForm} from 'vuetify/components';

const {t} = useI18n();
const show = ref(false);
const isSubmitting = ref(false);
const formRef = ref<InstanceType<typeof VForm> | null>(null);
const formData = ref({
  url: '',
  comment: '',
});

const emit = defineEmits<{
  submit: [data: {url: string; comment?: string}];
}>();

// Create rules with config: URL field requires all validations, comment only has maxLength
const rules = createFormRules(t, {
  required: true,
  validURL: true,
  maxLength: 2000,
});

const canSubmit = computed(() => formData.value.url.trim().length > 0);

const dialogConfig = computed(() => ({
  title: `${t('BTN_ADD')} ${t('TITLE_SOURCE_CODE')}`,
  primaryButton: {
    text: t('BTN_ADD'),
  },
  secondaryButton: {
    text: t('BTN_CANCEL'),
  },
}));

const open = () => {
  show.value = true;
};

const close = () => {
  show.value = false;
};

const resetDialog = () => {
  formRef.value?.reset();
  formData.value = {
    url: '',
    comment: '',
  };
  isSubmitting.value = false;
};

const submit = async () => {
  if (!canSubmit.value) {
    return;
  }

  const isValid = await formRef.value?.validate();
  if (!isValid?.valid) {
    return;
  }

  isSubmitting.value = true;
  try {
    emit('submit', {
      url: formData.value.url,
      comment: formData.value.comment || undefined,
    });
    // Parent component will handle closing the dialog
  } catch (error) {
    console.error('Error in dialog submit:', error);
  } finally {
    isSubmitting.value = false;
  }
};

defineExpose({
  open,
  close,
});
</script>
