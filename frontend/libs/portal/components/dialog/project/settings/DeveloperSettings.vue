<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <Stack class="gap-4">
    <Stack direction="row" v-if="hasParent" class="pt-4">
      <v-icon color="warning">mdi-alert</v-icon>
      <span>{{ t('DEVELOPER_SETTINGS_FROM_PARENT') }}</span>
    </Stack>

    <Stack v-if="!inWizard">
      <v-checkbox
        v-model="projectSettings.supplierExtraData.external"
        hide-details
        color="primary"
        class="mt-0 mr-0 shrink"
        :readonly="hasParent || isNotProjectOwner"
        :disabled="hasParent || isNotProjectOwner">
        <template v-slot:label>
          {{ t('EXTERNAL') }} <span class="small-text ml-2">{{ t('INFO_TEXT_ETERNAL') }}</span>
        </template>
      </v-checkbox>
      <div class="mt-n5 mb-4 ml-1" v-if="!inWizard && projectSettings.supplierExtraData.external">
        <v-icon size="x-small" class="mx-2 ml-8" color="brand">mdi-alert</v-icon>
        <span class="text-body-2">{{ t('EXTERNAL_REMARK') }}</span>
      </div>
    </Stack>

    <v-text-field
      v-if="projectSettings.supplierExtraData.external"
      autocomplete="off"
      v-model="projectSettings.documentMeta.supplierName"
      :label="t('DEVELOPER_NAME')"
      variant="outlined"
      :readonly="hasParent || isNotProjectOwner"
      :disabled="hasParent || isNotProjectOwner"
      :rules="activeRules.name"
      class="required"
      hide-details />

    <DAutocompleteCompany
      id="developer-company"
      v-if="!projectSettings.supplierExtraData.external && RightsUtils.rights().isInternal"
      v-model="projectSettings.documentMeta.supplierDept"
      :readonly="hasParent || isNotProjectOwner"
      :disabled="hasParent || isNotProjectOwner"
      :label="t('COMPANY')"
      :required="true"
      aria="supplier company"></DAutocompleteCompany>

    <v-textarea
      no-resize
      rows="5"
      v-model="projectSettings.documentMeta.supplierAddress"
      :label="t('PROJECT_SETTINGS_ADDRESS')"
      hide-details="auto"
      :rules="activeRules.address"
      variant="outlined"
      data-testid="DeveloperSettings__Address"
      :readonly="hasParent || isNotProjectOwner"
      :disabled="hasParent || isNotProjectOwner" />

    <v-text-field
      autocomplete="off"
      variant="outlined"
      v-model="projectSettings.documentMeta.supplierNr"
      :label="t('SUPPLIER_NR')"
      hide-details="auto"
      :readonly="hasParent || isNotProjectOwner"
      :disabled="hasParent || isNotProjectOwner"
      :rules="activeRules.supplierNr" />
  </Stack>
</template>

<script lang="ts" setup>
import {ProjectLabels} from '@disclosure-portal/constants/policyLabels';
import {ProjectModel, ProjectSettingsModel} from '@disclosure-portal/model/Project';
import {Group, Rights} from '@disclosure-portal/model/Rights';
import {WizardProjectPostRequest} from '@disclosure-portal/model/Wizard';
import {useAppStore} from '@disclosure-portal/stores/app';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {computed, watch} from 'vue';
import {useI18n} from 'vue-i18n';

interface Props {
  activeRules: Record<string, any>;
  rights?: Rights;
  hasParent?: boolean;
  inWizard?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  hasParent: false,
  inWizard: false,
});

const projectSettings = defineModel<ProjectSettingsModel>('settings', {required: true});
const project = defineModel<ProjectModel | WizardProjectPostRequest>('project', {required: true});

const {t} = useI18n();
const appStore = useAppStore();

const isNotProjectOwner = computed(() => props.rights && !props.rights.groups?.includes(Group.ProjectOwner));

const labelTools = appStore.getLabelsTools;

const updateDevelopmentLabels = async () => {
  if (!project.value || props.hasParent || isNotProjectOwner.value) {
    return;
  }

  const labels = project.value?.projectLabels.filter((label) => {
    const labelObj = labelTools.projectLabelsMap[label];
    if (!labelObj) return true; // Keep labels we don't recognize
    const name = labelObj.name;
    // Manual check instead of includes to satisfy TypeScript
    return (
      name !== ProjectLabels.DEVELOPMENT_EXTERNAL &&
      name !== ProjectLabels.DEVELOPMENT_INTERNAL &&
      name !== ProjectLabels.DEVELOPMENT_INHOUSE
    );
  });

  if (projectSettings.value.supplierExtraData.external) {
    const label = labelTools.projectLabelsMapByName[ProjectLabels.DEVELOPMENT_EXTERNAL];
    if (label) {
      labels.push(label._key);
    }
  } else {
    const selectedCompany = projectSettings.value.documentMeta.supplierDept;
    const customerCompany = projectSettings.value.customerMeta.dept;
    const supplierCompanyCode = selectedCompany?.companyCode;
    const customerCompanyCode = customerCompany?.companyCode;

    if (!supplierCompanyCode || !customerCompanyCode) {
      const label = labelTools.projectLabelsMapByName[ProjectLabels.DEVELOPMENT_INTERNAL];
      if (label) {
        labels.push(label._key);
      }
    } else if (supplierCompanyCode === customerCompanyCode) {
      const label = labelTools.projectLabelsMapByName[ProjectLabels.DEVELOPMENT_INHOUSE];
      if (label) {
        labels.push(label._key);
      }
    } else {
      const label = labelTools.projectLabelsMapByName[ProjectLabels.DEVELOPMENT_INTERNAL];
      if (label) {
        labels.push(label._key);
      }
    }
  }

  project.value.projectLabels = labels;
};

watch(
  [() => projectSettings.value.supplierExtraData.external, () => projectSettings.value.documentMeta.supplierDept],
  () => {
    if (!props.inWizard) {
      updateDevelopmentLabels();
    }
  },
  {immediate: true},
);
</script>
