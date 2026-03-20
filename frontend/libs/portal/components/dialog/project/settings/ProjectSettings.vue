<script lang="ts" setup>
import {PolicyLabels} from '@disclosure-portal/constants/policyLabels';
import {Project, ProjectSettingsModel} from '@disclosure-portal/model/Project';
import ProjectPostRequest from '@disclosure-portal/model/ProjectPostRequest';
import {Group} from '@disclosure-portal/model/Rights';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useDialogStore} from '@disclosure-portal/stores/dialog.store';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import useRules from '@disclosure-portal/utils/Rules';
import useSnackbar from '@shared/composables/useSnackbar';
import {storeToRefs} from 'pinia';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info: snack} = useSnackbar();
const appStore = useAppStore();
const projectStore = useProjectStore();
const {minMax, minMaxArray} = useRules();
const dialogStore = useDialogStore();
const idle = useIdleStore();

const {isSettingsDialogOpen} = storeToRefs(dialogStore);

const activeTab = ref('general');
const settingsModel = ref<ProjectSettingsModel>({} as ProjectSettingsModel);
const formRef = ref();
const generalForm = ref();

const labelTools = computed(() => appStore.getLabelsTools);
const projectModel = ref<Project>(Object.assign({}, projectStore.currentProject!));

const isVehicleOnboard = computed(() => {
  const onboardLabel = labelTools.value.policyLabelItems.find((label) => label.name === PolicyLabels.ONBOARD);
  if (!onboardLabel) {
    return false;
  }
  return projectModel.value.policyLabels.includes(onboardLabel._key);
});

const rules = {
  name: minMax(t('NP_DIALOG_TF_DEVELOPER'), 3, 80, false),
  address: minMax(t('NP_DIALOG_TF_ADDRESS'), 3, 300, true),
  supplierNr: minMax(t('NP_DIALOG_TF_SUPPLIER_NR'), 1, 25, true),
  freeLabels: minMaxArray(t('WIZARD_project_tags'), 1, 20),
};

const showDialog = async () => {
  isSettingsDialogOpen.value = true;
  projectModel.value = Object.assign({}, projectStore.currentProject!);
};

watch(isSettingsDialogOpen, async (newVal) => {
  if (newVal) {
    if (projectModel.value.hasParent) {
      settingsModel.value = new ProjectSettingsModel();
      settingsModel.value.customerMeta = projectModel.value.parentProjectSettings.customerMeta;
      settingsModel.value.noticeContactMeta = projectModel.value.parentProjectSettings.noticeContactMeta;
      settingsModel.value.supplierExtraData = projectModel.value.parentProjectSettings.supplierExtraData;
      settingsModel.value.documentMeta = projectModel.value.parentProjectSettings.documentMeta;
      settingsModel.value.customIds = [...projectModel.value.customIds];
      settingsModel.value.noFossProject = projectModel.value.isNoFoss;
    } else {
      settingsModel.value = new ProjectSettingsModel();
      settingsModel.value.fill(projectModel.value);
    }

    await nextTick();
    formRef.value?.resetValidation();
  }
});

const doDialogAction = async () => {
  const validationResult = await formRef.value?.validate();
  const generalFormValid = await generalForm.value?.validate();
  if (!validationResult?.valid || !generalFormValid) {
    const errors = validationResult?.errors as Array<{id: string; errorMessages: string[]}> | undefined;
    if (errors && errors.length) {
      if (errors.some((e) => e.id === 'owner-company')) {
        activeTab.value = 'owner';
      } else if (errors.some((e) => e.id === 'developer-company')) {
        activeTab.value = 'developer';
      }
    }
    snack(t('INVALID_PROJECT_SETTINGS'));
    return;
  }
  idle.show(t('PROJECT_IS_UPDATING'));
  const projectUpdate = new ProjectPostRequest();
  Object.assign(projectUpdate, {
    ...projectModel.value,
    projectSettings: settingsModel.value,
    id: projectModel.value._key,
  });

  await projectStore.updateProject(projectUpdate);

  isSettingsDialogOpen.value = false;
};

defineExpose({
  showDialog,
  activeTab,
});
</script>

<template>
  <slot :showDialog="showDialog"> </slot>
  <v-dialog v-model="isSettingsDialogOpen" width="800" persistent scrollable height="800">
    <v-form ref="formRef">
      <v-card scrollable class="p-8" data-testid="projects-editor">
        <v-card-title>
          <Stack direction="row" align="center">
            <span class="text-h5">
              {{ t('PVD_DIALOG_TITLE') }}
            </span>
            <v-spacer></v-spacer>
            <DCloseButton @click="isSettingsDialogOpen = false" />
          </Stack>
        </v-card-title>
        <v-card-text class="pt-2">
          <v-tabs v-model="activeTab" slider-color="mbti" show-arrows bg-color="tabsHeader">
            <v-tab value="general">{{ t('TAB_GENERAL') }}</v-tab>
            <v-tab value="owner">{{ t('TAB_OWNER') }}</v-tab>
            <v-tab value="developer">{{ t('TAB_DEVELOPER') }}</v-tab>
            <v-tab value="customids">{{ t('TAB_CUSTOM_IDs') }}</v-tab>
          </v-tabs>
          <v-tabs-window v-model="activeTab">
            <v-tabs-window-item value="general" eager>
              <GeneralSettings
                v-model:item="projectModel"
                v-model:settings="settingsModel"
                ref="generalForm"></GeneralSettings>
            </v-tabs-window-item>
            <v-tabs-window-item value="owner" eager>
              <OwnerSettings
                v-model:customer-meta="settingsModel.customerMeta"
                v-model:notice-meta="settingsModel.noticeContactMeta"
                :vehicle-onboard="isVehicleOnboard"
                :active-rules="rules"
                :rights="projectModel.accessRights"
                :has-parent="projectModel.hasParent" />
            </v-tabs-window-item>
            <v-tabs-window-item value="developer" eager>
              <DeveloperSettings
                v-model:settings="settingsModel"
                v-model:project="projectModel"
                :active-rules="rules"
                :rights="projectModel.accessRights"
                :has-parent="projectModel.hasParent" />
            </v-tabs-window-item>

            <v-tabs-window-item value="customids">
              <CustomIdSettings
                v-model="settingsModel.customIds"
                :readonly="
                  projectModel.accessRights && !projectModel.accessRights.groups?.includes(Group.ProjectOwner)
                " />
            </v-tabs-window-item>
          </v-tabs-window>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <DCActionButton
            v-if="projectModel.accessRights.groups?.includes(Group.ProjectOwner)"
            size="small"
            variant="text"
            class="mr-5"
            is-dialog-button
            :text="t('BTN_CANCEL')"
            @click="isSettingsDialogOpen = false"></DCActionButton>
          <DCActionButton
            v-if="projectModel.accessRights.groups?.includes(Group.ProjectOwner)"
            size="small"
            variant="flat"
            color="primary"
            :text="t('Btn_save')"
            is-dialog-button
            :loading="projectStore.loading"
            @click="doDialogAction"></DCActionButton>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
