<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {PolicyLabels} from '@disclosure-portal/constants/policyLabels';
import type {Project} from '@disclosure-portal/model/Project';
import ProjectService from '@disclosure-portal/services/projects';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useDialogStore} from '@disclosure-portal/stores/dialog.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {canDeleteProject, getDeleteTooltip} from '@disclosure-portal/utils/project-deletion-error';
import config from '@shared/utils/config';
import {storeToRefs} from 'pinia';
import {computed, nextTick, type Ref, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const reviewPresetRD = {
  comment: 'PRESET_RD_COMM',
  reviewer: 'CUSTOMER1',
};

const {t} = useI18n();
const router = useRouter();
const appStore = useAppStore();
const projectStore = useProjectStore();
const wizardStore = useWizardStore();
const dialogStore = useDialogStore();

const labelTools = computed(() => appStore.getLabelsTools);
const {hasVehiclePlatformChildren, hasOnlyVehiclePlatformChildren} = storeToRefs(projectStore);
const currentProject = computed((): Project => projectStore.currentProject!);

const confirmDialogVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);

const reqfoss = ref();
const reqreview = ref();
const reqapproval = ref();
const addChildrenErrorDialog = ref();

const isVehiclePlatform = computed(() => {
  if (!currentProject.value) return false;
  for (const lbl of currentProject.value.policyLabels) {
    if (labelTools.value.policyLabelsMap[lbl]?.name === PolicyLabels.VEHICLE_PLATFORM) {
      return true;
    }
  }
  return false;
});
const isVehiclePlatformOrItsChildren = computed(() => {
  return hasVehiclePlatformChildren.value || isVehiclePlatform.value;
});

const showGenerateFoss = computed(
  () =>
    isVehiclePlatformOrItsChildren.value ||
    (currentProject.value.supplierExtraData &&
      currentProject.value.supplierExtraData.external &&
      !currentProject.value.parent),
);
const isRequestReviewAllowed = computed(() => {
  return (
    currentProject.value.accessRights.allowProject.update &&
    currentProject.value.accessRights.allowRequestPlausi &&
    currentProject.value.accessRights.allowRequestPlausi.create &&
    !currentProject.value.parent
  );
});
const showMoreButton = computed(() => {
  return (
    currentProject.value.isApprovalAllowed ||
    isVehiclePlatform.value ||
    hasOnlyVehiclePlatformChildren.value ||
    isRequestReviewAllowed.value ||
    currentProject.value.isProjectOwner ||
    currentProject.value.accessRights.allowProject.delete
  );
});

const handleDialogOpen = (dialogRef: Ref<any, any>, ...openArgs: unknown[]) => {
  if (!projectStore.areMandatoryProjectSettingsSet) {
    addChildrenErrorDialog.value?.open();
  } else {
    dialogRef.value?.open(...openArgs);
  }
};

const generateFoss = () => handleDialogOpen(reqfoss, isVehiclePlatformOrItsChildren.value);
const requestApproval = () => handleDialogOpen(reqapproval, isVehiclePlatformOrItsChildren.value);
const reviewRD = () => handleDialogOpen(reqreview, reviewPresetRD, isVehiclePlatformOrItsChildren.value);
const requestReview = () => handleDialogOpen(reqreview, undefined, isVehiclePlatformOrItsChildren.value);

const showConfirmDialog = (config: Partial<IConfirmationDialogConfig>) => {
  confirmConfig.value = {
    key: currentProject.value._key,
    name: currentProject.value.name,
    ...config,
  } as IConfirmationDialogConfig;
  confirmDialogVisible.value = true;
};

const showDeletionConfirmationDialog = async () => {
  if (currentProject.value.isDummy) {
    showConfirmDialog({
      type: ConfirmationType.DELETE,
      description: 'DLG_CONFIRMATION_DESCRIPTION_DUMMY',
      okButton: 'Btn_delete',
    });
  } else {
    try {
      const approvalOrReviewUsage = await ProjectService.getApprovalOrReviewUsage(currentProject.value._key);
      const isInUse = approvalOrReviewUsage.data.success;

      if (isInUse) {
        showConfirmDialog({
          type: ConfirmationType.DEPRECATE,
          title: 'DLG_WARNING_TITLE',
          name: '',
          okButton: 'BTN_DEPRECATE',
          description: 'PROJECT_IN_APPROVAL_DEPRECATION',
          emphasiseText: 'PROJECT_DEPRECATION_UNREVERTABLE',
          emphasiseConfirmationText: 'PROJECT_DEPRECATION_UNREVERTABLE_CONFIRM',
        });
      } else {
        showConfirmDialog({
          type: ConfirmationType.DELETE,
          description: 'DLG_CONFIRMATION_DESCRIPTION',
          okButton: 'Btn_delete',
        });
      }
    } catch (e) {
      console.error('Error checking approval or review usage:', e);
    }
  }
};

const showDeprecationConfirmationDialog = () => {
  showConfirmDialog({
    type: ConfirmationType.DEPRECATE,
    description: 'DLG_DEPRECATION_CONFIRMATION_DESCRIPTION',
    emphasiseText: 'PROJECT_DEPRECATION_UNREVERTABLE',
    emphasiseConfirmationText: 'PROJECT_DEPRECATION_UNREVERTABLE_CONFIRM',
    okButton: 'BTN_DEPRECATE',
  });
};
const onConfirm = async (config: IConfirmationDialogConfig) => {
  if (config.okButtonIsDisabled) return;
  if (config.type === ConfirmationType.DELETE) {
    await projectStore.deleteProject(config.key);
    await router.push('/dashboard/projects');
  } else if (config.type === ConfirmationType.DEPRECATE) {
    await projectStore.deprecateProject(config.key);
    await projectStore.fetchProjectByKey(config.key);
  } else if (config.type === ConfirmationType.NOT_SET) {
    dialogStore.isSettingsDialogOpen = true;
  }
};

const openNewWizard = () => {
  wizardStore.openWizard({project: currentProject.value, mode: 'edit'});
};

const openSettingsDialog = () => {
  dialogStore.isSettingsDialogOpen = true;
};

const isProjectResponsible = computed(() => {
  return projectStore.currentProject?.responsible.toLowerCase() === useUserStore().getProfile.user.toLowerCase();
});

const isApprovalOrFossGenerationDisabled = computed(() => {
  if (config.useFutureIt && !isVehiclePlatformOrItsChildren.value) {
    return currentProject.value.isApprovalDisabled || !isProjectResponsible.value;
  } else if (config.useFutureProduct && isVehiclePlatformOrItsChildren.value) {
    return currentProject.value.isApprovalDisabled || !isProjectResponsible.value;
  } else {
    return currentProject.value.isApprovalDisabled;
  }
});

watch(
  currentProject,
  async () => {
    await nextTick();
    if (currentProject.value.deptMissing) {
      confirmConfig.value = {
        type: ConfirmationType.NOT_SET,
        key: currentProject.value._key,
        name: currentProject.value.name,
        description: 'DLG_CONFIRMATION_DESCRIPTION_DEPT_MISSING',
        okButton: 'Btn_edit_settings',
      };
      confirmDialogVisible.value = true;
    }
  },
  {immediate: true},
);
</script>

<template>
  <v-btn :hidden="!showMoreButton" id="projectMenu" class="text-none" width="100px" variant="tonal" color="primary">
    <v-icon size="large">mdi-menu-down</v-icon>
    <span class="font-bold">{{ t('MENU_MORE') }}</span>
  </v-btn>
  <v-menu activator="#projectMenu">
    <v-list>
      <MenuItem
        v-if="currentProject.isProjectOwner && !currentProject.isGroup"
        icon="mdi-checkbox-marked-circle-plus-outline"
        :text="t('BTN_OPEN_PROJECT_WIZARD')"
        :tooltip="t('BTN_OPEN_PROJECT_WIZARD')"
        @click="openNewWizard"></MenuItem>

      <MenuItem
        v-if="currentProject.isApprovalAllowed && !showGenerateFoss"
        icon="mdi-checkbox-marked-circle-plus-outline"
        :text="t('BTN_REQUEST_APPROVAL')"
        :disabled="isApprovalOrFossGenerationDisabled"
        @click="requestApproval">
        <template #tooltip>
          <span v-if="currentProject.isDummy">{{ t('TT_ONLY_FOR_REAL_PROJECTS') }}</span>
          <span v-else-if="!isProjectResponsible && isApprovalOrFossGenerationDisabled">{{
            t('TT_REQUEST_APPROVAL_RESPONSIBLE_ONLY')
          }}</span>
          <span v-else-if="currentProject.hasParent">{{ t('TT_REQUEST_APPROVAL_PARENT') }}</span>
          <span v-else-if="currentProject.deptMissing">{{ t('TT_REQUEST_APPROVAL_MISSING_DEPT') }}</span>
          <span v-else-if="!currentProject.isApprovalDisabled">{{ t('TT_REQUEST_APPROVAL') }}</span>
        </template>
      </MenuItem>

      <MenuItem
        v-if="currentProject.isApprovalAllowed && showGenerateFoss"
        icon="mdi-checkbox-marked-circle-plus-outline"
        :text="t('MENU_BTN_GENERATE_FOSS_DD')"
        :disabled="isApprovalOrFossGenerationDisabled"
        @click="generateFoss">
        <template #tooltip>
          <span v-if="currentProject.isDummy">{{ t('TT_ONLY_FOR_REAL_PROJECTS') }}</span>
          <span v-else-if="!isProjectResponsible && isApprovalOrFossGenerationDisabled">{{
            t('TT_GENERATE_FOSS_DD_RESPONSIBLE_ONLY')
          }}</span>
          <span v-else-if="currentProject.hasParent && isVehiclePlatformOrItsChildren">
            {{ t('TT_GENERATE_FOSS_DD_PARENT') }}
          </span>
          <span v-else-if="!currentProject.isApprovalDisabled">{{ t('TT_GENERATE_FOSS_DD') }}</span>
        </template>
      </MenuItem>

      <!-- RD FOSS Review @see disclosure-portal#4675 & disclosure-portal#5064 -->
      <MenuItem
        v-if="isVehiclePlatform || hasOnlyVehiclePlatformChildren"
        icon="mdi-checkbox-marked-circle-plus-outline"
        :text="t('PRESET_RD_TITLE')"
        :disabled="currentProject.isDummy"
        :tooltip="currentProject.isDummy ? t('TT_ONLY_FOR_REAL_PROJECTS') : t('PRESET_RD_TT')"
        @click="reviewRD"></MenuItem>

      <MenuItem
        v-if="isRequestReviewAllowed"
        icon="mdi-checkbox-marked-circle-plus-outline"
        :text="t('BTN_REQUEST_PLAUSI')"
        :tooltip="t('TT_REQUEST_PLAUSI')"
        @click="requestReview"></MenuItem>

      <MenuItem
        v-if="currentProject.isProjectOwner"
        icon="mdi-archive-outline"
        :text="t('BTN_DEPRECATE')"
        :tooltip="t('TT_deprecate_project')"
        @click="showDeprecationConfirmationDialog"></MenuItem>

      <MenuItem
        v-if="currentProject.accessRights.allowProject.delete"
        icon="mdi-delete"
        :text="t('TT_delete_project')"
        :tooltip="t(getDeleteTooltip(currentProject))"
        :disabled="!canDeleteProject(currentProject)"
        @click="showDeletionConfirmationDialog"></MenuItem>

      <slot></slot>
    </v-list>
  </v-menu>

  <template>
    <AddChildrenErrorDialog ref="addChildrenErrorDialog" @open-settings="openSettingsDialog"></AddChildrenErrorDialog>
    <RequestFOSSDD ref="reqfoss"></RequestFOSSDD>
    <RequestReview ref="reqreview"></RequestReview>
    <RequestApproval ref="reqapproval"></RequestApproval>
    <ConfirmationDialog
      v-model:showDialog="confirmDialogVisible"
      :config="confirmConfig"
      @confirm="onConfirm"></ConfirmationDialog>
  </template>
</template>
