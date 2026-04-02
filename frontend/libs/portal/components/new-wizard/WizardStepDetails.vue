<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useNewWizard} from '@disclosure-portal/composables/useNewWizard';
import {useCapabilitiesStore} from '@disclosure-portal/stores/capabilities';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const wizardStore = useWizardStore();
const {validationRules} = useNewWizard();

const {t} = useI18n();
const capabilitiesStore = useCapabilitiesStore();

const stepDetailsForm = ref(null);
const applicationSelectorRef = ref(null);

const validateSelf = async () => {
  const info = await (stepDetailsForm.value as any)?.validate();
  let isAppValid = true;
  if (wizardStore.isEnterpriseOrMobilePlatform && capabilitiesStore.applicationConnector) {
    isAppValid = applicationSelectorRef.value ? await (applicationSelectorRef.value as any)?.validate() : false;
  }
  return info?.valid && isAppValid;
};

onMounted(async () => {
  if (wizardStore.project.name || wizardStore.project.description) {
    await validateSelf();
  }
});

const onEnter = async () => {
  if (await validateSelf()) {
    wizardStore.nextStep();
  }
};
</script>

<template>
  <v-form ref="stepDetailsForm" class="overflow-x-hidden" @submit.prevent="onEnter">
    <Stack>
      <h2 class="text-body-1 py-0">{{ t('WIZARD_page_details_hint') }}</h2>
      <div
        class="grid grid-cols-1 gap-3 w-full"
        :class="{'md:grid-cols-2': wizardStore.isEnterpriseOrMobilePlatform && capabilitiesStore.applicationConnector}">
        <v-text-field
          autocomplete="off"
          class="required mb-auto"
          required
          :rules="validationRules.name"
          v-model="wizardStore.project.name"
          variant="outlined"
          :label="wizardStore.project.isGroup ? t('NP_DIALOG_GROUP_NAME') : t('WIZARD_project_name')"
          hide-details="auto"></v-text-field>

        <ApplicationSelector
          v-if="wizardStore.isEnterpriseOrMobilePlatform && capabilitiesStore.applicationConnector"
          ref="applicationSelectorRef"
          v-model="wizardStore.project.applicationMeta!"
          :is-required="false"></ApplicationSelector>
      </div>
      <v-textarea
        v-model="wizardStore.project.description"
        class="expand"
        variant="outlined"
        rows="8"
        :label="t('WIZARD_project_description')"
        no-resize
        hide-details
        counter="1000"
        :rules="validationRules.description"></v-textarea>
      <v-checkbox
        v-model="wizardStore.project.projectSettings.noFossProject"
        outlined
        color="primary"
        hide-details
        :label="t('NO_FOSS_MARKER')"></v-checkbox>
      <div v-if="wizardStore.project.projectSettings?.noFossProject">
        <v-icon small class="mx-2">warning</v-icon>
        <span>{{ t('NO_FOSS_WARNING') }}</span>
      </div>
    </Stack>
  </v-form>
</template>
