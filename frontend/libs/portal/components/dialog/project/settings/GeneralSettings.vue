<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import icons from '@disclosure-portal/constants/icons';
import {Project, ProjectSettingsModel} from '@disclosure-portal/model/Project';
import {useAppStore} from '@disclosure-portal/stores/app';
import useRules from '@disclosure-portal/utils/Rules';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const appStore = useAppStore();
const {minMax, minMaxArray, longText} = useRules();
const applicationSelectorRef = ref();

const item = defineModel<Project>('item', {required: true});
const settingsModel = defineModel<ProjectSettingsModel>('settings', {required: true});

const labelTools = computed(() => appStore.getLabelsTools);

const rules = {
  name: minMax(t('NP_DIALOG_TF_DEVELOPER'), 3, 80, false),
  freeLabels: minMaxArray(t('WIZARD_project_tags'), 0, 20),
  description: longText(t('NP_DIALOG_TF_DESCRIPTION')),
};
const validate = async () => {
  if (item.value.isGroup) {
    return true;
  }
  return (await applicationSelectorRef.value?.validate()) ?? false;
};

defineExpose({validate});
</script>

<template>
  <v-card-text class="p-0 pt-8" v-if="item && labelTools.policyLabelsMap">
    <Stack>
      <div class="grid w-full grid-cols-1 gap-3" :class="{'md:grid-cols-2': !item.isGroup}">
        <v-text-field
          autocomplete="off"
          variant="outlined"
          v-model="item.name"
          class="required mb-auto"
          :rules="rules.name"
          :label="t(item.isGroup ? 'NP_DIALOG_GROUP_NAME' : 'NP_DIALOG_TF_NAME')"
          hide-details="auto"
          autofocus></v-text-field>

        <ApplicationSelector
          v-if="!item.isGroup"
          class="mb-auto"
          ref="applicationSelectorRef"
          v-model="item.applicationMeta!"
          :is-required="false"></ApplicationSelector>
      </div>

      <v-textarea
        variant="outlined"
        v-model="item.description"
        :rules="rules.description"
        :label="t('NP_DIALOG_TF_DESCRIPTION')"
        hide-details
        counter="1000"></v-textarea>

      <div>
        <v-checkbox
          v-model="settingsModel.noFossProject"
          color="primary"
          density="compact"
          :label="t('NO_FOSS_MARKER')"
          hide-details />
        <div v-if="settingsModel?.noFossProject">
          <v-icon size="small" class="m-2 mb-3 ml-5" color="brand"> mdi-alert </v-icon>
          <span>{{ t('NO_FOSS_WARNING') }}</span>
        </div>
      </div>

      <v-combobox
        variant="outlined"
        multiple
        chips
        v-model="item.freeLabels"
        :label="t('NP_DIALOG_SB_FREELABELS')"
        hide-details
        :rules="rules.freeLabels">
        <template v-slot:chip="{item, props}">
          <DLabel :labelName="item.raw" closable :parentProps="props" :iconName="icons.TAG" />
        </template>
      </v-combobox>

      <div class="border-data-table p-2">
        <!-- Labels Display -->
        <div class="editProjectSubheadline bg-used-components -mt-4 mb-3 w-24">{{ t('HEADLINE_LABEL_SET') }}</div>
        <!-- Schema Label -->
        <div class="labelSetGroup">
          <DLabel
            :labelName="labelTools.schemaLabelsMap[item.schemaLabel]?.name || 'UNKNOWN_LABEL'"
            :iconName="icons.SCHEMA" />
        </div>
        <!-- Policy Labels -->
        <div class="labelSetGroup d-flex flex-wrap">
          <DLabel
            v-for="(l, i) in item.policyLabels"
            :key="i"
            :labelName="labelTools.policyLabelsMap[l]?.name || 'UNKNOWN_LABEL'"
            :iconName="icons.POLICY" />
        </div>
        <!-- Project Labels -->
        <div class="labelSetGroup d-flex flex-wrap">
          <DLabel
            v-for="(l, i) in item.projectLabels"
            :key="i"
            :labelName="labelTools.projectLabelsMap[l]?.name || 'UNKNOWN_LABEL'"
            :iconName="icons.PROJECT" />
        </div>
      </div>
    </Stack>
  </v-card-text>
</template>
