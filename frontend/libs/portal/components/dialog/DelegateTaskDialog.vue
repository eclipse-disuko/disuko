// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import DAutocompleteUser from '@shared/components/disco/DAutocompleteUser.vue';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const {t} = useI18n();

const showDialog = defineModel<boolean>('showDialog');

const emit = defineEmits<{
  (e: 'confirm', delegateUserId: string): void;
}>();



const formRef = ref<VForm | null>(null);
const delegateUser = ref('');
const delegateUserRef = ref<InstanceType<typeof DAutocompleteUser>>();

const close = () => {
  showDialog.value = false;
  delegateUser.value = '';
};

const confirm = async () => {
  const formValidationResult = await formRef.value?.validate();
  const formIsValid = formValidationResult?.valid ?? false;

  const userValidationResult = await delegateUserRef.value?.validateOnCreate();
  const userValid = userValidationResult ?? false;

  if (formIsValid && userValid) {
    emit('confirm', delegateUser.value);
    close();
  }
};
</script>

<template>
  <v-dialog v-model="showDialog" content-class="small" width="600">
    <v-card class="pa-8 dDialog" flat>
      <v-card-title>
        <v-row>
          <v-col cols="10" class="d-flex align-center">
            <span class="text-h5">{{ t('DELEGATE_TASK_TITLE') }}</span>
          </v-col>
          <v-col cols="2" align="right">
            <DCloseButton @click="close" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-form ref="formRef" @submit.prevent="confirm">
          <v-row>
            <v-col cols="12">
              <p class="text-body-2 mb-4">{{ t('DELEGATE_TASK_DESCRIPTION') }}</p>
              <DAutocompleteUser
                ref="delegateUserRef"
                v-model="delegateUser"
                :project-key="projectKey"
                :label="t('DELEGATE_TO_LABEL')"
                only-internal-users
                required />
            </v-col>
          </v-row>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <DCActionButton
          size="small"
          isDialogButton
          variant="text"
          @click="close"
          class="mr-5"
          :text="t('BTN_CANCEL')" />
        <DCActionButton
          size="small"
          isDialogButton
          variant="flat"
          color="primary"
          @click="confirm"
          :text="t('BTN_DELEGATE')" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
