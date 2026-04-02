<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import icons from '@disclosure-portal/constants/icons';
import {ChecklistItem, PolicyStatusTypes, TriggerTypes} from '@disclosure-portal/model/Checklist';
import {IObligation} from '@disclosure-portal/model/IObligation';
import {LicenseSlim} from '@disclosure-portal/model/License';
import {NameKeyIdentifier} from '@disclosure-portal/model/ProjectsResponse';
import {ScanRemarkLevel} from '@disclosure-portal/model/Quality';
import {ReviewTemplate} from '@disclosure-portal/model/ReviewTemplate';
import {default as AdminService, default as adminService} from '@disclosure-portal/services/admin';
import licenseService from '@disclosure-portal/services/license';
import {useChecklistsStore} from '@disclosure-portal/stores/checklists.store';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import {SearchOptions} from '@disclosure-portal/utils/Table';
import {debounce} from 'lodash';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {minMax, longText, required} = useRules();
const checklistStore = useChecklistsStore();
const labelStore = useLabelStore();

const isVisible = ref(false);
const isEdit = ref(false);
const item = ref(new ChecklistItem());
const dialog = ref<DiscoForm>();
const title = ref('');
const confirmText = ref('');
const saving = ref(false);
const licensesLoading = ref(false);
const selectedLicenses = ref<(LicenseSlim | null)[]>([]);
const selectedClassification = ref<(IObligation | null)[]>([]);
const allObligations = ref<IObligation[]>([]);
const allReviewTemplates = ref<ReviewTemplate[]>();
const selectedTargetTemplate = ref<ReviewTemplate | null>();
const searchResults = ref<LicenseSlim[]>([]);
const selectedCompNames = ref<(string | null)[]>([]);

const srLevels = computed(
  () =>
    Object.values(ScanRemarkLevel)
      .filter((value) => value !== '')
      .map((value) => ({
        key: value,
        name: t(`SCAN_REMARK_STATUS_${value}`) || value,
      })) as NameKeyIdentifier[],
);

const triggerTypes = computed(
  () =>
    Object.values(TriggerTypes).map((value) => ({
      key: value,
      name: t(`CHECKLIST_TRIGGER_${value}`) || value,
    })) as NameKeyIdentifier[],
);

const psTypes = computed(
  () =>
    Object.values(PolicyStatusTypes).map((value) => ({
      key: value,
      name: t(`CHECKLIST_PC_${value}`) || value,
    })) as NameKeyIdentifier[],
);

const rules = {
  name: minMax(t('NPV_DIALOG_TF_NAME'), 3, 80, false),
  description: longText(t('NP_DIALOG_TF_DESCRIPTION')),
  template: required(t('LBL_CHECKLIST_TEMPLATE')),
};

const dialogConfig = computed(() => ({
  title: t(title.value),
  primaryButton: {text: t(confirmText.value), disabled: saving.value, loading: saving.value},
  secondaryButton: {text: t('BTN_CANCEL'), disabled: saving.value},
}));

const canAddLicense = computed(() => {
  const last = selectedLicenses.value.at(-1);
  return !!last && !!last.licenseId;
});

const canSelectPolicyStatus = computed(
  () =>
    item.value.triggerType === TriggerTypes.POLICY_STATUS ||
    item.value.triggerType === TriggerTypes.CLASS_OR ||
    item.value.triggerType === TriggerTypes.CLASS_AND ||
    item.value.triggerType === TriggerTypes.COMPONENT_NAME ||
    item.value.triggerType === TriggerTypes.LICENSE,
);

const addLicense = () => {
  selectedLicenses.value.push(null);
};

const addClassification = () => {
  selectedClassification.value.push(null);
};

const addComponentName = () => {
  selectedCompNames.value.push(null);
};

const open = async (existing?: ChecklistItem) => {
  const [ob, tp] = await Promise.all([AdminService.getAllObligations(), adminService.getReviewTemplates()]);
  allObligations.value = ob.data.items;
  allReviewTemplates.value = tp.data;

  if (existing) {
    isEdit.value = true;
    item.value = {
      _key: existing._key,
      name: existing.name,
      triggerType: existing.triggerType,
      classifications: existing.classifications,
      policyLabels: existing.policyLabels,
      policyStatus: existing.policyStatus,
      scanRemarks: existing.scanRemarks,
      licenseIds: existing.licenseIds,
      componentNames: existing.componentNames,
      targetTemplateKey: existing.targetTemplateKey,
      targetTemplateName: existing.targetTemplateName,
    };
    selectedLicenses.value = item.value.licenseIds
      ? (await licenseService.lookup(existing.licenseIds)).data.items
      : [null];
    selectedClassification.value = item.value.classifications
      ? item.value.classifications.map((id) => allObligations.value.find((obl) => obl?._key === id) || null)
      : [null];
    selectedCompNames.value = item.value.componentNames ? item.value.componentNames : [null];
    selectedTargetTemplate.value = allReviewTemplates.value.find((t) => t._key == existing.targetTemplateKey);
  } else {
    isEdit.value = false;
    selectedLicenses.value = [null];
    selectedClassification.value = [null];
    selectedCompNames.value = [null];
    selectedTargetTemplate.value = null;
    item.value = {
      triggerType: TriggerTypes.DEFAULT,
    } as ChecklistItem;
  }
  dialog.value?.reset();
  title.value = existing ? 'CHECKLIST_ITEM_DIALOG_EDIT_TITLE' : 'CHECKLIST_ITEM_DIALOG_ADD_TITLE';
  confirmText.value = existing ? 'NP_DIALOG_BTN_EDIT' : 'NP_DIALOG_BTN_CREATE';
  isVisible.value = true;
};

const doDialogAction = async () => {
  await nextTick();
  const info = await dialog.value?.validate();
  if (!info?.valid) {
    return;
  }
  saving.value = true;
  item.value.licenseIds = selectedLicenses.value.filter((item) => !!item).map((item) => item.licenseId);
  item.value.classifications = selectedClassification.value.filter((item) => !!item).map((item) => item._key);
  item.value.componentNames = selectedCompNames.value.filter((item) => item !== null && item.length > 0) as string[];
  item.value.targetTemplateKey = selectedTargetTemplate.value?._key!;
  item.value.targetTemplateName = selectedTargetTemplate.value?.title!;

  if (isEdit.value) {
    await checklistStore.editItem(item.value);
  } else {
    await checklistStore.addItem(item.value);
  }

  saving.value = false;
  isVisible.value = false;
};

const debouncedSearchLicenses = debounce(async (query: string) => {
  if (!query) {
    searchResults.value = [];
    licensesLoading.value = false;
    return;
  }
  licensesLoading.value = true;
  const options = {} as SearchOptions;
  options.filterString = query;
  options.itemsPerPage = 50;
  options.page = 1;
  options.sortBy = [{key: 'name', order: 'asc'}];
  const response = await licenseService.search(options);
  searchResults.value = response.data.licenses;
  licensesLoading.value = false;
}, 300);

const close = () => {
  isVisible.value = false;
};

defineExpose({
  open,
});
</script>

<template>
  <v-dialog v-model="isVisible" content-class="large" persistent scrollable width="640" min-height="640">
    <DialogLayout :config="dialogConfig" @primary-action="doDialogAction" @secondary-action="close" @close="close">
      <v-form ref="dialog" @submit.prevent="doDialogAction">
        <Stack>
          <TextField
            required
            v-model="item.name"
            :rules="rules.name"
            :label="t('AL_DIALOG_TF_NAME')"
            autofocus
            hide-details="auto"
            tabindex="1" />
          <v-select
            v-model="item.triggerType"
            :label="t('LBL_TRIGGER_TYPE')"
            :items="triggerTypes"
            item-title="name"
            item-value="key"
            variant="outlined"></v-select>
          <v-select
            variant="outlined"
            hide-details="auto"
            v-model="item.policyLabels"
            item-title="name"
            item-value="_key"
            clearable
            multiple
            :items="labelStore.policyLabels"
            :label="t('AL_DIALOG_SB_LABELS_OR')"
            v-bind:menu-props="{location: 'bottom'}">
            <template v-slot:chip="{item, props}">
              <DLabel closable :parentProps="props" :labelName="item.title" :iconName="icons.TAG" />
            </template>
          </v-select>
          <v-autocomplete
            variant="outlined"
            hide-details="auto"
            v-model="selectedTargetTemplate"
            item-title="title"
            item-value="_key"
            clearable
            required
            return-object
            :items="allReviewTemplates"
            :label="t('LBL_CHECKLIST_TEMPLATE')"
            :rules="rules.template"
            v-bind:menu-props="{location: 'bottom'}">
          </v-autocomplete>
          <v-select
            class=""
            multiple
            v-if="canSelectPolicyStatus"
            v-model="item.policyStatus"
            :label="t('LBL_POLICY_STATUS')"
            :items="psTypes"
            item-title="name"
            item-value="key"
            variant="outlined"></v-select>
          <v-select
            v-if="item.triggerType === TriggerTypes.SCAN_REMARK"
            variant="outlined"
            hide-details="auto"
            v-model="item.scanRemarks"
            item-title="name"
            item-value="key"
            clearable
            :items="srLevels"
            :label="t('LBL_SCAN_REMARKS')"
            v-bind:menu-props="{location: 'bottom'}">
          </v-select>
          <v-autocomplete
            v-if="item.triggerType === TriggerTypes.LICENSE"
            v-for="(_, index) in selectedLicenses"
            :key="index"
            v-model="selectedLicenses[index]"
            clearable
            :label="t('labelSearchLicense')"
            :items="searchResults"
            return-object
            single-line
            :loading="licensesLoading"
            variant="outlined"
            hide-details="auto"
            @update:search="debouncedSearchLicenses"
            :no-filter="true">
            <template v-slot:item="{item, props}">
              <v-list-item v-bind="props" title="">
                <span class="d-subtitle-2">{{ item.raw.name }}</span>
                <span class="d-text d-secondary-text">&nbsp;({{ item.raw.licenseId }})</span>
              </v-list-item>
            </template>
            <template v-slot:selection="{item}">
              <div class="d-inline">
                <span class="d-subtitle-2">{{ item.raw.name }}</span>
                <span class="d-text d-secondary-text">&nbsp;({{ item.raw.licenseId }})</span>
              </div>
            </template>
          </v-autocomplete>
          <div
            v-if="canAddLicense && item.triggerType === TriggerTypes.LICENSE"
            class="d-flex align-center border-md border-dashed border-opacity-25 p-3 mb-6"
            @click="addLicense">
            <v-icon color="primary">mdi-plus</v-icon>
            <span class="font-weight-light pl-1">{{ t('RR_DIALOG_MORE_LICENSE') }}</span>
          </div>
          <v-autocomplete
            v-if="item.triggerType === TriggerTypes.CLASS_AND || item.triggerType === TriggerTypes.CLASS_OR"
            v-for="(_, index) in selectedClassification"
            :key="index"
            v-model="selectedClassification[index]"
            clearable
            item-title="name"
            :label="t('LBL_OBLIGATION')"
            :items="allObligations"
            return-object
            single-line
            variant="outlined"
            hide-details="auto">
            <template v-slot:item="{item, props}">
              <v-list-item v-bind="props" title="">
                <span class="d-subtitle-2">{{ item.raw?.name }}</span>
                <span class="d-text d-secondary-text">&nbsp;({{ item.raw?.warnLevel }})</span>
              </v-list-item>
            </template>
            <template v-slot:selection="{item}">
              <div class="d-inline">
                <span class="d-subtitle-2">{{ item.raw?.name }}</span>
                <span class="d-text d-secondary-text">&nbsp;({{ item.raw?.warnLevel }})</span>
              </div>
            </template>
          </v-autocomplete>
          <div
            v-if="item.triggerType === TriggerTypes.CLASS_AND || item.triggerType === TriggerTypes.CLASS_OR"
            class="d-flex align-center border-md border-dashed border-opacity-25 p-3 mb-6"
            @click="addClassification">
            <v-icon color="primary">mdi-plus</v-icon>
            <span class="font-weight-light pl-1">{{ t('RR_DIALOG_MORE_OBLIGATION') }}</span>
          </div>
          <TextField
            v-if="item.triggerType === TriggerTypes.COMPONENT_NAME"
            v-for="(_, index) in selectedCompNames"
            :key="index"
            v-model="selectedCompNames[index]"
            :label="t('LBL_COMPONENT_NAME')"
            dynamic-placeholder />
          <div
            v-if="item.triggerType === TriggerTypes.COMPONENT_NAME"
            class="d-flex align-center border-md border-dashed border-opacity-25 p-3 mb-6"
            @click="addComponentName">
            <v-icon color="primary">mdi-plus</v-icon>
            <span class="font-weight-light pl-1">{{ t('LBL_COMPONENT_NAME') }}</span>
          </div>
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
</template>
