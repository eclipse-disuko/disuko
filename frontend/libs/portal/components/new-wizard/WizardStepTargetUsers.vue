<script setup lang="ts">
import {WizardCardProps} from '@disclosure-portal/components/new-wizard/WizardCard.vue';
import {architectures, distributionTargets, targetUsers, TargetUsers} from '@disclosure-portal/model/NewWizard';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const wizardStore = useWizardStore();

const cardList = ref<WizardCardProps[]>([
  {
    key: targetUsers.company,
    image: new URL('@disclosure-portal/assets/wizard/projectUser/company.svg', import.meta.url).href,
    title: t('NEW_WIZARD_TARGET_USERS_MB'),
    subtitle: t('NEW_WIZARD_TARGET_USERS_MB_DESC'),
    helptext: t('NEW_WIZARD_TARGET_USERS_MB_HELP'),
    isFlipped: false,
    isActive: wizardStore.project?.targetUsers === targetUsers.company,
  },
  {
    key: targetUsers.businessPartner,
    image: new URL('@disclosure-portal/assets/wizard/projectUser/business-partner.svg', import.meta.url).href,
    title: t('NEW_WIZARD_TARGET_USERS_BUSINESS'),
    subtitle: t('NEW_WIZARD_TARGET_USERS_BUSINESS_DESC'),
    helptext: t('NEW_WIZARD_TARGET_USERS_BUSINESS_HELP'),
    isFlipped: false,
    isActive: wizardStore.project?.targetUsers === targetUsers.businessPartner,
  },
  {
    key: targetUsers.customer,
    image: new URL('@disclosure-portal/assets/wizard/projectUser/end-customer.svg', import.meta.url).href,
    title: t('NEW_WIZARD_TARGET_USERS_CUSTOMER'),
    subtitle: t('NEW_WIZARD_TARGET_USERS_CUSTOMER_DESC'),
    helptext: t('NEW_WIZARD_TARGET_USERS_CUSTOMER_HELP'),
    isFlipped: false,
    isActive: wizardStore.project?.targetUsers === targetUsers.customer,
  },
]);

/**
 * Pre-select distribution target based on selected target users
 * When architecture is frontend, automatically choose distribution target
 * company -> company
 * business partner -> business partner
 * customer -> business partner
 * @param selectedTargetUsers The selected target users
 */
const preSelectDistributionTarget = (selectedTargetUsers: TargetUsers) => {
  const project = wizardStore.project;
  if (project.architecture !== architectures.frontend) {
    wizardStore.nextStep();
    return;
  }

  if (selectedTargetUsers === targetUsers.company) {
    wizardStore.project.distributionTarget = distributionTargets.company;
  } else if (selectedTargetUsers === targetUsers.businessPartner) {
    wizardStore.project.distributionTarget = distributionTargets.businessPartner;
  } else if (selectedTargetUsers === targetUsers.customer) {
    wizardStore.project.distributionTarget = distributionTargets.businessPartner;
  }
  wizardStore.nextTwoSteps();
};
const onCardSelect = (selectedTargetUsers: TargetUsers) => {
  wizardStore.project.targetUsers = selectedTargetUsers;
  preSelectDistributionTarget(selectedTargetUsers);
};
</script>

<template>
  <Stack class="overflow-hidden">
    <h2 class="text-body-1 py-0">{{ t('WIZARD_project_target_Users') }}</h2>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 mb-1">
      <WizardCard v-for="card in cardList" :key="card.key" :card="card" @update="onCardSelect($event)" />
    </div>
  </Stack>
</template>
