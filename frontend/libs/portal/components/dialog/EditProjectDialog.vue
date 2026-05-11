<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {DialogEditProjectConfig} from '@disclosure-portal/components/dialog/DialogConfigs';
import Icons from '@disclosure-portal/constants/icons';
import {PolicyLabels} from '@disclosure-portal/constants/policyLabels';
import ProjectPostRequest from '@disclosure-portal/model/ProjectPostRequest';
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import useRules from '@disclosure-portal/utils/Rules';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const emit = defineEmits<{(e: 'edited'): void}>();

const projectStore = useProjectStore();
const {t} = useI18n();
const {minMax, longText, minMaxArray} = useRules();
const labels = useLabelStore();

const isVisible = ref(false);
const project = ref<ProjectPostRequest>(new ProjectPostRequest());
const projectEditForm = ref<VForm | null>(null);
const config = ref<DialogEditProjectConfig>({} as DialogEditProjectConfig);
const applicationSelectorRef = ref<InstanceType<typeof ApplicationSelector>>();

const isEnterpriseOrMobilePlatform = computed((): boolean => {
  const enterpriseLabelKey = labels.getLabelByNameAndType(PolicyLabels.ENTERPRISE_PLATFORM, 'POLICY')?._key || '';
  const mobileLabelKey = labels.getLabelByNameAndType(PolicyLabels.MOBILE_PLATFORM, 'POLICY')?._key || '';

  if (!enterpriseLabelKey && !mobileLabelKey) {
    return false;
  }

  return project.value.policyLabels.includes(enterpriseLabelKey) || project.value.policyLabels.includes(mobileLabelKey);
});
const rules = {
  name: minMax(t('NP_DIALOG_TF_NAME'), 3, 80, false),
  freeLabels: minMaxArray(t('NP_DIALOG_SB_FREELABELS'), 1, 20),
  description: longText(t('NP_DIALOG_TF_DESCRIPTION')),
};

const open = (newConf: DialogEditProjectConfig) => {
  config.value = newConf;
  if (config.value.project) {
    project.value.fillWithProjectModel(config.value.project);
  } else if (config.value.projectSlim) {
    const d = config.value.projectSlim as unknown as ProjectSlim;
    project.value.fillWithProjectSlim(d);
  } else {
    console.error('insufficient config for dialog');
    return;
  }
  isVisible.value = true;
};

const reset = () => projectEditForm.value?.reset();

const doDialogAction = async () => {
  await nextTick();

  const validation = await projectEditForm.value?.validate();

  if (!validation?.valid) {
    return;
  }

  await projectStore.updateProject(project.value);

  if (config.value.projectSlim) {
    emit('edited');
  }
  isVisible.value = false;
};

defineExpose({open});
</script>

<template>
  <div>
    <v-dialog v-model="isVisible" content-class="large" scrollable width="700">
      <v-form ref="projectEditForm">
        <v-card class="pa-8" data-testid="projects-editor">
          <v-card-title>
            <v-row>
              <v-col cols="10" class="d-flex align-center">
                <span class="text-h5">
                  {{ t('NP_DIALOG_TITLE_EDIT') }}
                </span>
              </v-col>
              <v-col cols="2" class="px-0 text-right">
                <DCloseButton
                  @click="
                    () => {
                      reset();
                      isVisible = false;
                    }
                  " />
              </v-col>
            </v-row>
          </v-card-title>
          <v-card-text class="pt-2" v-if="project">
            <!-- Form Fields -->
            <v-row dense>
              <v-col cols="12" sm="6" class="errorBorder pb-3">
                <v-text-field
                  autocomplete="off"
                  variant="outlined"
                  v-model="project.name"
                  class="required"
                  :rules="rules.name"
                  :label="t('NP_DIALOG_TF_NAME')"
                  hide-details="auto"
                  autofocus />
              </v-col>
              <v-col cols="12" sm="6">
                <ApplicationSelector
                  ref="applicationSelectorRef"
                  v-model="project.applicationMeta!"
                  :is-required="isEnterpriseOrMobilePlatform" />
              </v-col>
              <v-col cols="12">
                <v-textarea
                  variant="outlined"
                  v-model="project.description"
                  :rules="rules.description"
                  :label="t('NP_DIALOG_TF_DESCRIPTION')"
                  hide-details="auto"
                  counter="1000" />
              </v-col>
              <v-col cols="12">
                <v-combobox
                  variant="outlined"
                  multiple
                  chips
                  v-model="project.freeLabels"
                  :label="t('NP_DIALOG_SB_FREELABELS')"
                  hide-details="auto"
                  :rules="rules.freeLabels">
                  <template v-slot:chip="{item, props}">
                    <DLabel :labelName="item.raw" closable :parentProps="props" :iconName="Icons.TAG" />
                  </template>
                </v-combobox>
              </v-col>
              <v-col cols="12" class="ma-1 border-data-table mt-3">
                <!-- Labels Display -->
                <div class="editProjectSubheadline mt-n4 bg-used-components mb-3">{{ t('HEADLINE_LABEL_SET') }}</div>
                <Stack class="mb-3 ml-3 gap-2">
                  <ProjectLabel :label="labels.getLabelByKey(project.schemaLabel)"></ProjectLabel>
                  <div class="flex flex-wrap gap-1">
                    <ProjectLabel
                      v-for="labelKey in project.policyLabels"
                      :key="labelKey"
                      :label="labels.getLabelByKey(labelKey)" />
                  </div>
                  <div class="flex flex-wrap gap-1">
                    <ProjectLabel
                      v-for="labelKey in project.projectLabels"
                      :key="labelKey"
                      :label="labels.getLabelByKey(labelKey)" />
                  </div>
                </Stack>
              </v-col>
            </v-row>
          </v-card-text>
          <!-- Actions -->
          <v-card-actions class="justify-end">
            <DCActionButton
              size="small"
              variant="text"
              @click="
                () => {
                  reset();
                  isVisible = false;
                }
              "
              class="mr-5"
              :text="t('BTN_CANCEL')" />

            <DCActionButton
              isDialogButton
              :loading="projectStore.loading"
              :disabled="projectStore.loading"
              size="small"
              variant="flat"
              @click="doDialogAction()"
              :text="t('NP_DIALOG_BTN_EDIT')" />
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </div>
</template>
