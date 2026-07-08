<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {DialogVersionFormConfig} from '@disclosure-portal/components/dialog/DialogConfigs';
import {useHead} from '@unhead/vue';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {SpdxFile} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useIdleStore} from '@shared/stores/idle.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {getStrWithMaxLength} from '@disclosure-portal/utils/View';
import useSnackbar from '@shared/composables/useSnackbar';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import _ from 'lodash';
import {storeToRefs} from 'pinia';
import {computed, nextTick, onUnmounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';

const {t, locale} = useI18n();
const router = useRouter();
const route = useRoute();
const appStore = useAppStore();
const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const {info: snack} = useSnackbar();
const idle = useIdleStore();
const {dashboardCrumbs, projectsCrumb, ...breadcrumbs} = useBreadcrumbsStore();
const pageTitle = ref('');
useHead({title: pageTitle});
const {currentProject} = storeToRefs(projectStore);

const {currentVersion, channelSpdxs} = storeToRefs(sbomStore);

const reviewDia = ref(null);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmVisible = ref(false);
const dataAreLoaded = ref(false);
const selectedTab = ref('');
const editDlg = ref(null);
const sbomMenuOpen = ref(false);

const currentSpdx = computed(() => sbomStore.getSelectedSBOM || spdxFileHistory.value[0]);
const currentProjectEmpty = computed(() => _.isEmpty(currentProject.value));
const versionDetails = computed(() => currentVersion.value);
const versionName = computed(() => currentVersion.value?.name || '');
const spdxFileHistory = computed(() => channelSpdxs.value);
const iconForSelectedSpdx = computed(() =>
  currentProject.value?.approvablespdx.spdxkey === currentSpdx.value?._key ? 'mdi-star' : 'mdi-star-outline',
);
const encodedCurrentProjectParent = computed(() => encodeURIComponent(currentProject.value?.parent));

const iconColorForSelectedSpdx = computed(() =>
  currentProject.value?.approvablespdx.spdxkey === currentSpdx.value?._key ? 'primary' : 'grey',
);

const hintForSelectedSpdx = computed(() =>
  currentProject.value?.approvablespdx.spdxkey === currentSpdx.value?._key
    ? 'TT_approvable_spdx'
    : 'TT_not_approvable_spdx',
);

const hintForDisabledSelectedSpdx = computed(() =>
  currentProject.value?.approvablespdx.spdxkey === currentSpdx.value?._key
    ? 'TT_approvable_spdx'
    : 'TT_not_approvable_spdx',
);
const projectId = computed(() => (Array.isArray(route.params?.uuid) ? route.params.uuid[0] : route.params?.uuid || ''));
const versionKey = computed(() =>
  Array.isArray(route.params?.version) ? route.params.version[0] : route.params?.version || '',
);
const spdxKey = computed(() =>
  Array.isArray(route.params?.currentSbom) ? route.params.currentSbom[0] : route.params?.currentSbom || '',
);
const encodedProjectId = computed(() => encodeURIComponent(projectId.value));
const encodedVersion = computed(() => encodeURIComponent(versionDetails.value?._key || ''));
const encodedSbomKey = computed(() => {
  const sbomKey = spdxKey.value || currentSpdx.value?._key;
  return sbomKey ? encodeURIComponent(sbomKey) : '';
});
const tabUrlPart = computed(() => {
  return `/dashboard/projects/${encodedProjectId.value}/versions/${versionKey.value}`;
});
const userIsOwner = computed(() => currentProject.value?.accessRights?.groups.find((g: string) => g == 'Owner'));
const componentId = computed(() => {
  const componentId = Array.isArray(route.params?.componentId)
    ? route.params.componentId[0]
    : route.params?.componentId;
  if (componentId) {
    return `/${componentId}`;
  } else {
    return '';
  }
});

const resetUrl = async () => {
  const spdx = currentSpdx.value ? `/${currentSpdx.value._key}` : '';
  await router.replace(
    `/dashboard/projects/${encodedProjectId.value}/versions/${encodedVersion.value}/overview${spdx}`,
  );
};

const reload = async () => {
  if (currentProject.value?._key !== projectId.value) {
    await projectStore.fetchProjectByKey(projectId.value);
  }
  if (!versionDetails.value || versionDetails.value._key !== versionKey.value) {
    sbomStore.setCurrentVersion(versionKey.value);
    await sbomStore.fetchAllSBOMsFlat();
  }
  let selectedByRoute = false;
  if (spdxKey.value) {
    const sel = spdxFileHistory.value.find((spdx) => spdx._key === spdxKey.value);
    if (sel) {
      sbomStore.setSelectedSBOMKey(sel._key);
      selectedByRoute = true;
    }
  }
  if (!selectedByRoute) {
    sbomStore.setSelectedSBOMKey(spdxFileHistory.value[0]?._key || '');
    await resetUrl();
  }
  if (route.name === 'VersionSubTap') {
    await resetUrl();
  }
  dataAreLoaded.value = true;
};

const initPage = async () => {
  await nextTick();
  appStore.setDummyDesignMode(currentProject.value?.isDummy ?? false);
  initBreadcrumbs();
};

const editVersion = () => {
  const config = {
    version: versionDetails.value,
  } as unknown as DialogVersionFormConfig;
  (editDlg.value as any)?.open(config);
};

const initBreadcrumbs = () => {
  const currentGroupCrumb = {
    title: currentProject.value?.parentName ?? '',
    href: `/dashboard/groups/${encodedCurrentProjectParent.value}/children`,
  };
  const currentProjectCrumb = {
    title: currentProject.value?.name ?? '',
    href: `/dashboard/projects/${encodedProjectId.value}/overview`,
  };
  const currentversionCrumb = {
    title: versionName.value,
    href: `/dashboard/projects/${encodedProjectId.value}/versions/${encodedVersion.value}/overview`,
  };
  const groupProjectCrumbs = currentProject.value?.parent
    ? [currentGroupCrumb, currentProjectCrumb]
    : [currentProjectCrumb];
  let breadCrumb = [];
  breadCrumb = [...dashboardCrumbs, projectsCrumb, ...groupProjectCrumbs, currentversionCrumb];
  if (currentSpdx.value) {
    breadCrumb[breadCrumb.length - 1] = {
      title: t('BREAD_SBOM_DELIVERIES'),
      href: `/dashboard/projects/${encodedProjectId.value}/sbomlist`,
    };
    breadCrumb.push({
      title: versionDetails.value.name,
      href: `/dashboard/projects/${encodedProjectId.value}/sbomlist`,
    });
  }
  breadcrumbs.setCurrentBreadcrumbs(breadCrumb);
};
const showOverallReviewDialog = () => {
  (reviewDia.value as any)?.open();
};
const showDeletionConfirmationDialog = async () => {
  await versionService.getApprovalOrReviewUsage(projectId.value, versionKey.value).then((r) => {
    const isInUse = r.data.success;
    if (isInUse) {
      confirmConfig.value = {
        type: ConfirmationType.NOT_SET,
        title: 'DLG_WARNING_TITLE',
        key: '',
        name: '',
        description: 'VERSION_IN_APPROVAL',
        okButton: 'Btn_delete',
        okButtonIsDisabled: true,
      };
    } else {
      confirmConfig.value = {
        type: ConfirmationType.DELETE,
        key: versionDetails.value?._key,
        name: versionDetails.value?.name,
        description: 'DLG_CONFIRMATION_DESCRIPTION',
        okButton: 'Btn_delete',
      };
    }
    confirmVisible.value = true;
  });
};

const setApprovable = async (spdxFileKey: string) => {
  const approvableSpdx = {
    spdxkey: '',
    versionkey: '',
  } as ApprovableSPDXDto;
  if (spdxFileKey !== currentProject.value?.approvablespdx.spdxkey) {
    approvableSpdx.spdxkey = spdxFileKey;
    approvableSpdx.versionkey = versionKey.value;
  }
  await projectService.updateApprovableSpdx(approvableSpdx, currentProject.value?._key);
  currentProject.value.approvablespdx = approvableSpdx;
};

const doDeleteVersion = async (config: IConfirmationDialogConfig) => {
  if (config.okButtonIsDisabled) return;
  await versionService.deleteVersion(projectId.value, versionKey.value).then(() => {
    snack(t('DIALOG_version_delete_success'));
    close();
  });
};

const close = () => {
  if (versionKey.value) {
    router.push(`/dashboard/projects/${encodedProjectId.value}/overview`);
  } else {
    router.push('/dashboard/projects');
  }
};

const selectSbom = async (spdx: SpdxFile) => {
  sbomMenuOpen.value = false;
  await selectedSpdxChanged(spdx);
};

const selectedSpdxChanged = async (newSpdx: SpdxFile) => {
  const subTab =
    (route.name === 'licenseRemarks' && '/licenseRemarks') ||
    (route.name === 'reviewRemarks' && '/reviewRemarks') ||
    (route.name === 'generalRemarks' && '/generalRemarks') ||
    '';

  await router.push(
    `/dashboard/projects/${encodedProjectId.value}/versions/${encodedVersion.value}/${selectedTab.value}/${newSpdx._key}${subTab}`,
  );
};

watch(
  currentProject,
  () => {
    idle.hide();
  },
  {deep: true},
);

watch(
  [spdxKey, versionKey, projectId],
  async () => {
    await reload();
  },
  {immediate: true},
);

watch(dataAreLoaded, async (dal) => {
  if (dal) {
    await initPage();
  }
});

const sbomTabNames: Record<string, string> = {
  overview: 'TAB_OVERVIEW',
  component: 'TAB_Components',
  history: 'SBOM_DELIVERIES',
  sbomCompare: 'TAB_SBOM_COMPARE',
  sbomQuality: 'TAB_QUALITY',
  source: 'TAB_SourceCode',
  overallReviews: 'TAB_OVERALL_REVIEWS',
  notice: 'TAB_NoticeFile',
  auditLog: 'TAB_PROJECT_AUDIT',
};

// Set up reactive page title
watch(
  () => [
    currentProject.value?.name,
    currentSpdx.value?.uploaded,
    currentSpdx.value?.metaInfo?.name,
    selectedTab.value,
    locale.value,
  ],
  ([projectName, uploaded, sbomName]) => {
    if (projectName && uploaded) {
      const tabKey = sbomTabNames[selectedTab.value];
      const tabName = tabKey ? t(tabKey) : '';
      pageTitle.value = tabName
        ? `${sbomName} | ${String(projectName)} | ${tabName}`
        : `${sbomName} | ${String(projectName)}`;
    }
  },
  {immediate: true},
);

onUnmounted(() => {
  sbomStore.setCurrentVersion('');
  sbomStore.setSelectedSBOMKey('');
  appStore.unsetDummyDesignMode();
});
</script>

<template>
  <div v-if="currentProject" class="h-full p-4" data-testid="projects-versions">
    <div v-if="!currentProjectEmpty" class="flex flex-row flex-wrap items-center gap-2 pb-3">
      <div class="flex flex-row items-center gap-2">
        <v-chip v-if="currentProject.isDummy" class="dummy-tag mr-2" label>DUMMY</v-chip>
        <div class="text-h6 d-secondary-text">SBOM:</div>
        <v-menu v-if="spdxFileHistory.length > 0" v-model="sbomMenuOpen" location="bottom start">
          <template v-slot:activator="{props}">
            <div v-bind="props" class="sbom-selector flex cursor-pointer items-center gap-1">
              <DDateCellWithTooltip class="d-secondary-text" :value="currentSpdx?.uploaded"></DDateCellWithTooltip>
              <span class="text-h6" v-if="currentSpdx?.metaInfo">
                {{ getStrWithMaxLength(39, currentSpdx.metaInfo.name) }}
              </span>
              <v-chip v-if="currentSpdx?.tag" size="x-small">{{ getStrWithMaxLength(10, currentSpdx.tag) }}</v-chip>
              <v-chip v-if="currentSpdx?.isRecent" size="x-small">{{ t('SBOM_LATEST') }}</v-chip>
              <v-icon size="small" color="primary">mdi-chevron-down</v-icon>
            </div>
          </template>
          <v-list density="compact">
            <v-list-item
              v-for="spdx in spdxFileHistory"
              :key="spdx._key"
              :active="spdx._key === currentSpdx?._key"
              @click="selectSbom(spdx)">
              <span class="align-center text-caption flex flex-row gap-1">
                <span class="d-secondary-text">{{ formatDateAndTime(spdx.uploaded) }}</span>
                <span>-</span>
                <span v-if="spdx.metaInfo">
                  {{ getStrWithMaxLength(39, spdx.metaInfo.name) }}
                </span>
                <v-chip v-if="spdx.tag" size="x-small">{{ spdx.tag }}</v-chip>
                <v-chip v-if="spdx.isRecent" size="x-small">{{ t('SBOM_LATEST') }}</v-chip>
                <DOverallStateIcon v-if="spdx.overallReview" :review="spdx.overallReview" />
                <v-icon v-if="currentProject?.approvablespdx.spdxkey === spdx._key" color="primary" size="x-small">
                  mdi-star
                </v-icon>
              </span>
            </v-list-item>
          </v-list>
        </v-menu>
        <DOverallStateIcon v-if="currentSpdx?.overallReview" :review="currentSpdx.overallReview" class="mx-1" />
        <span v-if="spdxFileHistory.length >= 1 && userIsOwner" class="flex items-center">
          <v-btn
            :icon="iconForSelectedSpdx"
            variant="text"
            density="compact"
            size="small"
            :color="iconColorForSelectedSpdx"
            @click="setApprovable(currentSpdx!._key)">
          </v-btn>
          <Tooltip :text="t(hintForSelectedSpdx)" location="top" />
        </span>
        <span v-else-if="spdxFileHistory.length > 0 && !userIsOwner" class="flex items-center">
          <v-btn :icon="iconForSelectedSpdx" variant="text" density="compact" size="small" disabled></v-btn>
          <Tooltip :text="t(hintForDisabledSelectedSpdx)" location="top" />
        </span>
      </div>
      <v-spacer></v-spacer>
      <DCActionButton
        v-if="currentProject?.accessRights?.allowProjectVersion?.update"
        icon="mdi-pencil"
        :hint="t('TT_edit_project')"
        :text="t('BTN_EDIT')"
        data-testid="edit"
        @click="editVersion"></DCActionButton>
      <ProjectMenu v-if="currentProject">
        <v-divider></v-divider>
        <MenuItem
          v-if="
            !currentProject.isDeprecated &&
            spdxFileHistory.length > 0 &&
            currentProject?.accessRights?.allowProjectVersion?.read
          "
          icon="mdi-message-draw"
          :tooltip="t('TT_overall_review')"
          :text="t('BTN_OVRERALL_REVIEW')"
          @click="showOverallReviewDialog">
        </MenuItem>
        <MenuItem
          v-if="!currentProject.isDeprecated && currentProject?.accessRights?.allowProject?.delete"
          icon="mdi-delete"
          :tooltip="t('TT_delete_version')"
          :text="t('TT_delete_version')"
          @click="showDeletionConfirmationDialog">
        </MenuItem>
      </ProjectMenu>
    </div>
    <div v-if="dataAreLoaded && versionDetails">
      <v-card>
        <v-tabs v-model="selectedTab" slider-color="brand" active-class="active" show-arrows bg-color="tabsHeader">
          <v-tab value="overview" :to="`${tabUrlPart}/overview/${encodedSbomKey}`">
            {{ t('TAB_OVERVIEW') }}
          </v-tab>
          <v-tab
            value="component"
            :to="`${tabUrlPart}/component/${encodedSbomKey}${componentId}`"
            :disabled="!currentSpdx">
            {{ t('TAB_Components') }}
          </v-tab>
          <v-tab value="history" :to="`${tabUrlPart}/history/${encodedSbomKey}`">
            {{ t('TAB_SBOM_DELIVERIES') }}
          </v-tab>
          <v-tab value="sbomCompare" :to="`${tabUrlPart}/sbomCompare/${encodedSbomKey}`" :disabled="!currentSpdx">
            {{ t('TAB_SBOM_COMPARE') }}
          </v-tab>
          <v-tab value="sbomQuality" :to="`${tabUrlPart}/sbomQuality/${encodedSbomKey}`" :disabled="!currentSpdx">
            {{ t('TAB_QUALITY') }}
          </v-tab>
          <v-tab value="source" :to="`${tabUrlPart}/source/${encodedSbomKey}`" :disabled="!currentSpdx">
            {{ t('TAB_SourceCode') }}
          </v-tab>
          <v-tab value="overallReviews" :to="`${tabUrlPart}/overallReviews/${encodedSbomKey}`" :disabled="!currentSpdx">
            {{ t('TAB_OVERALL_REVIEWS') }}
          </v-tab>
          <v-tab value="notice" :to="`${tabUrlPart}/notice/${encodedSbomKey}`" :disabled="!currentSpdx">
            {{ t('TAB_NoticeFile') }}
          </v-tab>
          <v-tab
            v-if="currentProject?.accessRights?.allowProjectAudit?.read"
            value="auditLog"
            :to="`${tabUrlPart}/auditLog/${encodedSbomKey}`"
            :disabled="!currentSpdx">
            {{ t('TAB_PROJECT_AUDIT') }}
          </v-tab>
        </v-tabs>
        <v-tabs-window v-model="selectedTab">
          <v-tabs-window-item value="overview">
            <TabOverview ref="overview"></TabOverview>
          </v-tabs-window-item>
          <v-tabs-window-item value="component">
            <TabComponentList ref="component"></TabComponentList>
          </v-tabs-window-item>
          <v-tabs-window-item value="history">
            <TabSBOMHistory></TabSBOMHistory>
          </v-tabs-window-item>
          <v-tabs-window-item value="sbomCompare">
            <TabSBOMCompare ref="sbomCompare"></TabSBOMCompare>
          </v-tabs-window-item>
          <v-tabs-window-item value="sbomQuality">
            <TabSBOMQualityMain ref="quality"></TabSBOMQualityMain>
          </v-tabs-window-item>
          <v-tabs-window-item value="source">
            <TabSourceCode ref="source"></TabSourceCode>
          </v-tabs-window-item>
          <v-tabs-window-item value="overallReviews">
            <TabOverallReviews ref="overallReviews" @reloadParent="reload"></TabOverallReviews>
          </v-tabs-window-item>
          <v-tabs-window-item value="notice">
            <TabNoticeFile ref="notice"></TabNoticeFile>
          </v-tabs-window-item>
          <v-tabs-window-item v-if="currentProject?.accessRights?.allowProjectAudit?.read" value="auditLog">
            <TabAuditLog ref="auditLog"></TabAuditLog>
          </v-tabs-window-item>
        </v-tabs-window>
      </v-card>
    </div>
    <VersionDialogForm ref="editDlg"></VersionDialogForm>
    <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doDeleteVersion">
    </ConfirmationDialog>
    <OverallReviewDialog ref="reviewDia" @reload="reload"></OverallReviewDialog>
  </div>
</template>

<style scoped lang="scss">
.dummy-tag {
  margin-left: 2px;
  border: 1px solid rgb(var(--v-theme-chartYellow));
}

.sbom-selector {
  border-radius: 4px;
  padding: 4px 8px;

  &:hover {
    background-color: rgba(var(--v-theme-primary), 0.08);
  }
}
</style>
