<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ApproverRoles} from '@disclosure-portal/model/Approval';
import DHTTPError from '@disclosure-portal/model/DHTTPError';
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import {UserDto} from '@disclosure-portal/model/Users';
import eventBus from '@disclosure-portal/utils/eventbus';
import DAutocompleteUser from '@shared/components/disco/DAutocompleteUser.vue';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const props = defineProps<{
  customer1User: UserDto;
  customer2User: UserDto;
  projectKey: string;
  title: string;
  confirmText: string;
}>();

const emit = defineEmits<{
  (e: 'confirm', customer1: string, customer2: string): void;
  (e: 'close'): void;
}>();

const {t} = useI18n();

const formRef = ref<VForm | null>(null);
const customer1Ref = ref();
const customer2Ref = ref();

const customer1 = ref<string>('');
const customer2 = ref<string>('');

watch(
  () => props.customer1User,
  (newVal) => {
    if (newVal?.user) {
      customer1.value = newVal.user;
    }
  },
  {immediate: true},
);

watch(
  () => props.customer2User,
  (newVal) => {
    if (newVal?.user) {
      customer2.value = newVal.user;
    }
  },
  {immediate: true},
);

const emitClose = () => {
  emit('close');
};

const confirm = async () => {
  await nextTick();

  const formValidationResult = await formRef.value?.validate();
  const formIsValid = formValidationResult?.valid ?? false;

  const customer1ValidationResult = await customer1Ref.value?.validate();
  const customer1Valid = customer1ValidationResult ?? false;

  const customer2ValidationResult = await customer2Ref.value?.validate();
  const customer2Valid = customer2ValidationResult ?? false;

  if (formIsValid && customer1Valid && customer2Valid) {
    if (customer1.value === customer2.value) {
      const error = new DHTTPError();
      error.title = '' + t('SBOM_REQUEST_INTERNAL_APPROVAL');
      error.message = '' + t('EQUAL_OWNER_APPROVERS_ERROR_MESSAGE');
      eventBus.emit('on-api-error', error);
      return;
    }
    if ((customer1.value !== '' && customer2.value === '') || (customer1.value === '' && customer2.value !== '')) {
      const d = new ErrorDialogConfig();
      d.title = '' + t('SBOM_REQUEST_INTERNAL_APPROVAL');
      d.description = '' + t('BOTH_OR_NONE_OWNER_APPROVERS_ALLOWED_ERROR_MESSAGE');
      eventBus.emit('on-error', {error: d});
      return;
    }
    emit('confirm', customer1.value, customer2.value);
  }
};

const ro = (role: ApproverRoles) => {
  return Boolean(
    (role === ApproverRoles.Customer1 && props.customer1User?.user) ||
    (role === ApproverRoles.Customer2 && props.customer2User?.user),
  );
};

const dialogConfig = computed(() => ({
  title: props.title,
  primaryButton: {text: props.confirmText},
  secondaryButton: {text: t('BTN_CANCEL')},
}));
</script>

<template>
  <DialogLayout :config="dialogConfig" @primary-action="confirm" @secondary-action="emitClose" @close="emitClose">
    <v-form ref="formRef" @submit.prevent="confirm">
      <Stack>
        <DAutocompleteUser
          ref="customer1Ref"
          v-model="customer1"
          :project-key="projectKey"
          :label="t('FIRST_REPORTER_LABEL')"
          :required="true"
          :preselect="customer1User"
          :readonly="ro(ApproverRoles.Customer1)" />
        <DAutocompleteUser
          ref="customer2Ref"
          v-model="customer2"
          :project-key="projectKey"
          :label="t('SECOND_REPORTER_LABEL')"
          :required="true"
          :preselect="customer2User"
          :readonly="ro(ApproverRoles.Customer2)" />
      </Stack>
    </v-form>
  </DialogLayout>
</template>
