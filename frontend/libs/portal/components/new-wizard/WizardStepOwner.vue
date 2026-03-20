<script setup lang="ts">
import {useNewWizard} from '@disclosure-portal/composables/useNewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const wizardStore = useWizardStore();
const {validationRules} = useNewWizard();

const stepOwnerForm = ref(null);

const validateSelf = async () => {
  const info = await (stepOwnerForm.value as any)?.validate();
  return info?.valid;
};

onMounted(async () => {
  if (
    wizardStore.project.projectSettings?.customerMeta?.dept ||
    wizardStore.project.projectSettings?.customerMeta?.address ||
    wizardStore.project.projectSettings?.noticeContactMeta?.address
  ) {
    await validateSelf();
  }
});
</script>

<template>
  <v-form ref="stepOwnerForm">
    <Stack>
      <h2 class="text-body-1 py-0">{{ t('WIZARD_OWNER_description') }}</h2>

      <DAutocompleteCompany
        id="owner-company"
        v-if="RightsUtils.rights().isInternal"
        v-model="wizardStore.project.projectSettings.customerMeta.dept"
        :label="t('COMPANY')"
        required
        aria="owner company" />

      <v-textarea
        variant="outlined"
        no-resize
        rows="4"
        v-model="wizardStore.project.projectSettings.customerMeta.address"
        :label="t('PROJECT_SETTINGS_ADDRESS')"
        hide-details="auto"
        data-testid="OwnerSettings__Address"
        :rules="validationRules.address" />

      <v-textarea
        v-if="!wizardStore.isVehicleOnboardArchitecture"
        id="thirdparty-address"
        rows="5"
        autocomplete="off"
        :placeholder="t('PLACEHOLDER_NOTICE_CONTACT_ADDRESS')"
        persistent-placeholder
        variant="outlined"
        v-model="wizardStore.project.projectSettings.noticeContactMeta.address"
        :label="t('NOTICE_CONTACT_ADDRESS')"
        hide-details="auto"
        :rules="validationRules.address" />
    </Stack>
  </v-form>
</template>
