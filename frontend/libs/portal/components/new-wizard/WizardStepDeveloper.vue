<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useNewWizard} from '@disclosure-portal/composables/useNewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const wizardStore = useWizardStore();
const {validationRules} = useNewWizard();

const stepDeveloperForm = ref(null);

const validateSelf = async () => {
  const info = await (stepDeveloperForm.value as any)?.validate();
  return info?.valid;
};

onMounted(async () => {
  if (
    wizardStore.project.projectSettings?.documentMeta?.supplierName ||
    wizardStore.project.projectSettings?.documentMeta?.supplierDept ||
    wizardStore.project.projectSettings?.documentMeta?.supplierAddress ||
    wizardStore.project.projectSettings?.documentMeta?.supplierNr
  ) {
    await validateSelf();
  }
});
</script>

<template>
  <v-form ref="stepDeveloperForm">
    <Stack>
      <h2 class="text-body-1 py-0">{{ t('WIZARD_DEVELOPER_description') }}</h2>

      <v-text-field
        v-if="wizardStore.project.projectSettings.supplierExtraData.external"
        autocomplete="off"
        v-model="wizardStore.project.projectSettings.documentMeta.supplierName"
        :label="t('DEVELOPER_NAME')"
        variant="outlined"
        :rules="validationRules.name"
        class="required -mt-3 pt-3"
        hide-details="auto" />

      <DAutocompleteCompany
        id="developer-company"
        v-if="!wizardStore.project.projectSettings.supplierExtraData.external && RightsUtils.rights().isInternal"
        v-model="wizardStore.project.projectSettings.documentMeta.supplierDept"
        :label="t('COMPANY')"
        :required="true"
        class="-mt-3 pt-3"
        aria="supplier company" />

      <v-textarea
        no-resize
        rows="4"
        v-model="wizardStore.project.projectSettings.documentMeta.supplierAddress"
        :label="t('PROJECT_SETTINGS_ADDRESS')"
        hide-details="auto"
        :rules="validationRules.address"
        variant="outlined"
        data-testid="DeveloperSettings__Address" />

      <v-text-field
        autocomplete="off"
        variant="outlined"
        v-model="wizardStore.project.projectSettings.documentMeta.supplierNr"
        :label="t('SUPPLIER_NR')"
        hide-details="auto"
        :rules="validationRules.supplierNr" />
    </Stack>
  </v-form>
</template>
