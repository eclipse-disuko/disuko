// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ApprovalFormHooks, useApprovalFormBase} from '@disclosure-portal/composables/approval/useApprovalFormBase';
import {DocumentMeta, ExternalApprovalRequest} from '@disclosure-portal/model/ApprovalRequest';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {OverallReviewState, SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import config from '@shared/utils/config';
import {computed, nextTick, ref, watch} from 'vue';

export function useFossDdForm(externalHooks: Partial<ApprovalFormHooks> = {}) {
  const vehicle = ref(false);
  const radioGroup = ref(0);
  const childProjectChannels = ref<Map<string, VersionSlim>>(new Map());
  const allChannelSboms = ref<Map<string, SpdxFile[]>>(new Map());
  const activePanel = ref<number | null>(null);
  const dd = ref();

  const defaultC1 = () => {
    if (externalHooks.defaultC1) return externalHooks.defaultC1();
    return noFOSS.value ? false : vehicle.value;
  };
  const defaultC2 = (countApprovables: number, selectedSbom: SpdxFile | null) => {
    if (externalHooks.defaultC2) return externalHooks.defaultC2(countApprovables, selectedSbom);
    if (noFOSS.value) return false;
    if (vehicle.value) return false;
    return countApprovables > 0 || selectedSbom != null;
  };
  const defaultC3 = (countApprovables: number) => {
    if (externalHooks.defaultC3) return externalHooks.defaultC3(countApprovables);
    if (noFOSS.value) return false;
    if (vehicle.value) return false;
    return !(countApprovables > 0);
  };
  const defaultC4 = () => {
    if (externalHooks.defaultC4) return externalHooks.defaultC4();
    return noFOSS.value ? false : !vehicle.value;
  };
  const defaultRadioGroupFn = () => {
    if (externalHooks.defaultRadioGroup) return externalHooks.defaultRadioGroup();
    return noFOSS.value ? 3 : 1;
  };

  const hooks: ApprovalFormHooks = {
    ...externalHooks,
    defaultC1,
    defaultC2,
    defaultC3,
    defaultC4,
    defaultRadioGroup: defaultRadioGroupFn,
    afterSetDefaultFlags: () => {
      radioGroup.value = defaultRadioGroupFn();
      externalHooks.afterSetDefaultFlags?.();
    },
    afterResetFormState: () => {
      activePanel.value = null;
      radioGroup.value = 0;
      childProjectChannels.value.clear();
      allChannelSboms.value.clear();
      externalHooks.afterResetFormState?.();
    },
    loadSBOMHistPreHook: async () => {
      await base.sbomStore.fetchAllSBOMsFlat();
      await externalHooks.loadSBOMHistPreHook?.();
    },
  };

  const base = useApprovalFormBase(hooks);

  const {noFOSS, selectedProjects, approvableInfo, channels, stats} = base;

  watch(radioGroup, () => {
    if (radioGroup.value == 3) {
      noFOSS.value = true;
    }
  });

  const isRdConfirmationMissing = computed(() => {
    if (!vehicle.value) {
      return false;
    }

    if (!base.projectModel.value.isGroup) {
      const approvableSpdx = approvableInfo.value.projects?.[0]?.approvablespdx;

      if (!approvableSpdx?.spdxkey || !approvableSpdx?.versionkey) {
        return false;
      }

      const channel = channels.value.find((c) => c._key === approvableSpdx.versionkey);

      if (!channel) {
        return false;
      }

      const hasAuditedReview = channel.overallReviews?.some(
        (review) => review.sbomId === approvableSpdx.spdxkey && review.state === OverallReviewState.AUDITED,
      );

      return !hasAuditedReview;
    }

    if (base.projectModel.value.isGroup && approvableInfo.value.projects) {
      const selectedProjectsSet = new Set(selectedProjects.value);

      for (const project of approvableInfo.value.projects) {
        if (selectedProjectsSet.size > 0 && !selectedProjectsSet.has(project.projectKey)) {
          continue;
        }

        if (!project.approvablespdx?.spdxkey || !project.approvablespdx?.versionkey) {
          continue;
        }

        const channel = childProjectChannels.value.get(project.approvablespdx.versionkey);

        if (!channel) {
          return true;
        }

        const hasAuditedReview = channel.overallReviews?.some(
          (review) => review.sbomId === project.approvablespdx.spdxkey && review.state === OverallReviewState.AUDITED,
        );

        if (!hasAuditedReview) {
          return true;
        }
      }
    }

    return false;
  });

  const isDeniedOrUnasserted = computed(() => {
    return vehicle.value && (stats.value.denied > 0 || stats.value.noAssertion > 0);
  });

  const isEnterpriseOrMobileOrOther = computed(() => {
    return !vehicle.value && (stats.value.denied > 0 || stats.value.noAssertion > 0);
  });

  const showRedWarnDeniedDecisionsMessage = computed(
    () => !isDeniedOrUnasserted.value && approvableInfo.value.hasDeniedDecisions,
  );

  const open = async (isVehicle: boolean) => {
    base.idle.showIdle = true;
    await base.sbomStore.fetchAllSBOMsFlat();
    approvableInfo.value = await projectService.getApprovableInfo(base.projectModel.value._key);

    vehicle.value = isVehicle;
    if (vehicle.value) {
      base.withZip.value = true;
    }
    noFOSS.value = base.projectModel.value.isNoFoss;
    base.setDefaultFlags();
    await base.autoSelect();

    if (!base.projectModel.value.isGroup) {
      allChannelSboms.value.clear();
      for (const channel of channels.value) {
        const versionEntry = base.sbomStore.getAllSBOMs.find((v) => v.versionKey === channel._key);
        allChannelSboms.value.set(channel._key, versionEntry?.spdxFileHistory ?? []);
      }
    }

    if (base.projectModel.value.isGroup && approvableInfo.value.projects) {
      childProjectChannels.value.clear();

      const versionFetchPromises = approvableInfo.value.projects
        .filter((p) => p.approvablespdx.versionkey)
        .map(async (project) => {
          try {
            const versionDetails = await versionService.getVersion(
              project.projectKey,
              project.approvablespdx.versionkey,
            );
            childProjectChannels.value.set(project.approvablespdx.versionkey, versionDetails.data);
          } catch (error) {
            console.error(`Failed to fetch version details for project ${project.projectKey}:`, error);
          }
        });

      await Promise.all(versionFetchPromises);
    }

    base.idle.showIdle = false;
    base.isVisible.value = true;
  };

  const doDialogAction = async () => {
    await nextTick();
    const info = await base.form.value?.validate();
    if (!info?.valid) {
      return;
    }

    if (isRdConfirmationMissing.value && config.enforceFOSSOfficeConfirmation) {
      return;
    }

    const metaDoc: DocumentMeta = new DocumentMeta();
    if (vehicle.value) {
      metaDoc.c1 = radioGroup.value == 1;
      metaDoc.c2 = radioGroup.value == 2;
      metaDoc.c3 = false;
      metaDoc.c4 = false;
      metaDoc.c5 = false;
    } else {
      metaDoc.c1 = base.c1.value;
      metaDoc.c2 = base.c2.value;
      metaDoc.c3 = base.c3.value;
      metaDoc.c4 = base.c4.value;
      metaDoc.c5 = base.c5.value;
    }
    metaDoc.c6 = noFOSS.value;

    const req: ExternalApprovalRequest = {
      comment: base.comment.value,
      guidProject: base.projectModel.value._key,
      metaDoc: metaDoc,
      withZip: base.withZip.value,
      fossVersion: 'vanilla',
      selectedProjects: selectedProjects.value,
    };

    base.idle.showIdle = true;
    if (!base.projectModel.value.isGroup) {
      const approvableSpdx = {
        spdxkey: '',
        versionkey: '',
      } as ApprovableSPDXDto;
      approvableSpdx.spdxkey = base.selectedSbom.value?._key ?? '';
      approvableSpdx.versionkey = base.selectedChannel.value?._key ?? '';
      await projectService.updateApprovableSpdx(approvableSpdx, base.projectModel.value._key);
    }

    const response = await (vehicle.value
      ? projectService.createVehicleApproval(req, base.projectModel.value._key)
      : projectService.createExternalApproval(req, base.projectModel.value._key));

    if (response) {
      await base.jobStore.pollJobStatus(base.projectModel.value._key, response.jobKey);
      base.isVisible.value = false;
      dd.value?.open(response.approvalGuid);
    } else {
      base.idle.showIdle = false;
    }
  };

  return {
    ...base,
    vehicle,
    radioGroup,
    childProjectChannels,
    allChannelSboms,
    activePanel,
    dd,
    isRdConfirmationMissing,
    isDeniedOrUnasserted,
    isEnterpriseOrMobileOrOther,
    showRedWarnDeniedDecisionsMessage,
    open,
    doDialogAction,
  };
}
