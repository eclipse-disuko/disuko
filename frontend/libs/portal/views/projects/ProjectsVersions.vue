<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {DialogVersionFormConfig} from '@disclosure-portal/components/dialog/DialogConfigs';
import {usePageTitle} from '@disclosure-portal/composables/usePageTitle';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {SpdxFile} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
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
const {useReactiveTitle} = usePageTitle();
const {currentProject} = storeToRefs(projectStore);

const {selectedSpdx, currentVersion, channelSpdxs} = storeToRefs(sbomStore);

const reviewDia = ref(null);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmVisible = ref(false);
const dataAreLoaded = ref(false);
const selectedTab = ref('');
const editDlg = ref(null);

const currentSpdx = computed(() => selectedSpdx.value || spdxFileHistory.value[0]);
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
      sbomStore.setSelectedSpdx(sel);
      selectedByRoute = true;
    }
  }
  if (!selectedByRoute) {
    sbomStore.setSelectedSpdx(spdxFileHistory.value[0] || null);
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

// Function to get tab display name from route
const getTabDisplayName = () => {
  const path = route.path;
  if (path.includes('/overview')) return t('TAB_OVERVIEW');
  if (path.includes('/component')) return t('TAB_Components');
  if (path.includes('/history')) return t('SBOM_DELIVERIES');
  if (path.includes('/sbomCompare')) return t('TAB_SBOM_COMPARE');
  if (path.includes('/sbomQuality')) return t('TAB_QUALITY');
  if (path.includes('/source')) return t('TAB_SourceCode');
  if (path.includes('/overallReviews')) return t('TAB_OVERALL_REVIEWS');
  if (path.includes('/notice')) return t('TAB_NoticeFile');
  if (path.includes('/auditLog')) return t('TAB_PROJECT_AUDIT');
  return t('TAB_OVERVIEW');
};

// Set up reactive page title
watch(
  () => [currentProject.value?.name, versionDetails.value?.name, route.name, locale.value],
  ([projectName, versionName]) => {
    if (projectName && versionName) {
      const tabName = getTabDisplayName();
      useReactiveTitle(`${String(projectName)} | ${String(versionName)} | ${tabName}`);
    }
  },
  {immediate: true},
);

onUnmounted(() => {
  sbomStore.setCurrentVersion('');
  sbomStore.setSelectedSpdx({} as SpdxFile);
  appStore.unsetDummyDesignMode();
});
</script>

<template>
  <div v-if="currentProject" class="h-full p-4" data-testid="projects-versions">
    <div v-if="!currentProjectEmpty" class="d-flex align-center ga-2 flex-row flex-wrap pb-3">
      <div>
        <v-chip v-if="currentProject.isDummy" class="dummy-tag mr-2" label>DUMMY</v-chip>
        <span class="text-h5" style="display: inline-block">{{ t('PROJECT') }}</span>
        <span class="text-h5 project-name px-2" :title="currentProject.name">
          <q>
            <span>{{ currentProject.name }}</span>
          </q>
        </span>
      </div>
      <v-spacer></v-spacer>
      <div class="d-flex align-center flex-row">
        <v-select
          v-if="spdxFileHistory.length > 0"
          :model-value="currentSpdx"
          :label="t('VERSION') + ' &quot;' + versionName + '&quot; > ' + t('SBOM_COMPARE_CURRENT')"
          :items="spdxFileHistory"
          style="display: inline-block; max-width: 550px"
          variant="outlined"
          density="compact"
          hide-details
          return-object
          color="inputActiveBorderColor"
          location="offsetY"
          @update:modelValue="selectedSpdxChanged">
          <template v-slot:item="{item, props}">
            <v-list-item v-bind="props" title="">
              <v-icon
                v-if="currentProject?.approvablespdx.spdxkey === item.raw._key"
                color="primary"
                size="x-small"
                class="pb-1">
                mdi-star
              </v-icon>
              <span class="text-caption ml-5">{{ formatDateAndTime(item.raw.Uploaded) }}&nbsp;</span>
              <span class="text-caption d-secondary-text" v-if="item.raw.MetaInfo"
                >&nbsp;-&nbsp;{{ getStrWithMaxLength(39, item.raw.MetaInfo.Name) }}</span
              >
              <span class="text-caption d-secondary-text" v-if="item.raw.Tag">&nbsp;({{ item.raw.Tag }})</span>
              <span class="text-caption d-secondary-text mr-1" v-if="item.raw.isRecent"
                >&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span
              >
              <span class="text-caption d-secondary-text mr-1" v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }}</span>
              <DOverallStateIcon v-if="item.raw.OverallReview" :review="item.raw.OverallReview" />
            </v-list-item>
          </template>
          <template v-slot:selection="{item}">
            <div class="d-inline py-1">
              <v-icon
                v-if="currentProject?.approvablespdx.spdxkey == item.raw._key"
                color="primary"
                size="x-small"
                class="pr-2 pb-1">
                mdi-star
              </v-icon>
              <span v-else class="placeholder-icon"></span>
            </div>
            <span class="text-caption">{{ formatDateAndTime(item.raw.Uploaded) }}</span>
            <span class="text-caption d-secondary-text" v-if="item.raw.MetaInfo">
              - {{ getStrWithMaxLength(39, item.raw.MetaInfo.Name) }}
            </span>
            <span class="text-caption d-secondary-text" v-if="item.raw.Tag"
              >&nbsp;({{ getStrWithMaxLength(10, item.raw.Tag) }})
            </span>
            <span class="text-caption" v-if="item.raw.isRecent">&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span>
            <span class="text-caption d-secondary-text" v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }} </span>&nbsp;
            <DOverallStateIcon v-if="item.raw.OverallReview" :review="item.raw.OverallReview" />
          </template>
        </v-select>
        <span v-if="spdxFileHistory.length >= 1 && userIsOwner">
          <DCActionButton
            variant="text"
            :tableButton="true"
            :icon="iconForSelectedSpdx"
            :color="iconColorForSelectedSpdx"
            :hint="t(hintForSelectedSpdx)"
            @click="setApprovable(currentSpdx!._key)"></DCActionButton>
        </span>
        <span v-else-if="spdxFileHistory.length > 0 && !userIsOwner">
          <DCActionButton
            variant="text"
            :tableButton="true"
            :disabled="true"
            :icon="iconForSelectedSpdx"
            :hint="t(hintForDisabledSelectedSpdx)"></DCActionButton>
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
          @click="showOverallReviewDialog"></MenuItem>
        <MenuItem
          v-if="!currentProject.isDeprecated && currentProject?.accessRights?.allowProject?.delete"
          icon="mdi-delete"
          :tooltip="t('TT_delete_version')"
          :text="t('TT_delete_version')"
          @click="showDeletionConfirmationDialog"></MenuItem>
      </ProjectMenu>
    </div>
    <div v-if="dataAreLoaded && versionDetails">
      <v-card>
        <v-tabs v-model="selectedTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
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
          <v-tabs-window-item value="overallReviews" class="pa-3">
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
    <ConfirmationDialog
      v-model:showDialog="confirmVisible"
      :config="confirmConfig"
      @confirm="doDeleteVersion"></ConfirmationDialog>
    <OverallReviewDialog ref="reviewDia" @reload="reload"></OverallReviewDialog>
  </div>
</template>

<style scoped lang="scss">
.dummy-tag {
  margin-left: 2px;
  border: 1px solid rgb(var(--v-theme-chartYellow));
}
.project-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
  position: relative;
  transition: max-width 0.3s ease;
}

.project-name:hover {
  max-width: none;
  white-space: normal;
  z-index: 1;
  padding: 5px;
}
</style>
