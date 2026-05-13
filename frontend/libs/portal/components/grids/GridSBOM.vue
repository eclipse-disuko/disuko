<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import ErrorDialogConfig from '@shared/types/ErrorDialogConfig';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {NameKeyIdentifier, VersionSbomsFlat} from '@disclosure-portal/model/ProjectsResponse';
import {Group} from '@disclosure-portal/model/Rights';
import {SpdxFile} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useIdleStore} from '@shared/stores/idle.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import eventBus from '@shared/utils/eventbus';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {formatDateTime, formatDateTimeShort, originShort, originTooltip} from '@disclosure-portal/utils/View';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {
  DataTabelIndex,
  DataTableHeader,
  DataTableHeaderFilterItems,
  DataTableItem,
  SortItem,
} from '@shared/types/table';
import {useClipboard} from '@shared/utils/clipboard';
import config from '@shared/utils/config';
import dayjs from 'dayjs';
import _ from 'lodash';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';
import {useLanguageStore} from '@shared/stores/language.store';
import {storeToRefs} from 'pinia';

type DataTableItems = DataTabelIndex & VersionSbomsFlat;

interface Props {
  channelView?: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits(['openVersion']);

const {t} = useI18n();
const appStore = useAppStore();
const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const route = useRoute();
const router = useRouter();
const idle = useIdleStore();
const languageStore = useLanguageStore();
const {appLanguage} = storeToRefs(languageStore);
const {copyToClipboard} = useClipboard();

const projectModel = computed(() => projectStore.currentProject!);
const versionDetails = computed(() => sbomStore.getCurrentVersion);
const spdxFileHistory = computed(() => sbomStore.getChannelSpdxs);
const labelTools = computed(() => appStore.getLabelsTools);

const search = ref('');
const uploadURL = ref('');
const isBranchSelectionEnabled = ref(true);
const selectedFilterChannel = ref<string[]>([]);
const selectedBranch = ref<NameKeyIdentifier>({} as NameKeyIdentifier);
const sortItems = ref<SortItem[]>([{key: 'uploaded', order: 'desc'}]);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmVisible = ref(false);
const {info: snack} = useSnackbar();
const branches = computed(() => sbomStore.allVersions);
const reviewRemarkDialog = ref();
const dlgSbomValidationErrors = ref();
const helpText = ref('');
const upload = ref();

const sortByName = (a: SpdxFile, b: SpdxFile): number => {
  return b.metaInfo.name.localeCompare(a.metaInfo.name);
};

const headers = computed((): DataTableHeader[] => [
  {
    title: t('COL_ACTIONS'),
    sortable: false,
    align: 'center',
    value: 'actions',
    width: 100,
  },
  {
    title: t('COL_SPDX_FILENAME'),
    align: 'start',
    value: 'searchIndex',
    width: 380,
    sortable: true,
    sortRaw: sortByName,
  },
  {
    title: t('COL_REVIEW_STATUS'),
    align: 'center',
    value: 'overallReview',
    width: 80,
  },
  ...(!props.channelView
    ? [
        {
          title: t('COL_SBOM_BRANCH'),
          align: 'start',
          value: 'versionName',
          sortable: true,
          width: 110,
        } as DataTableHeader,
      ]
    : []),
  {
    title: t('COL_SBOM_TAG'),
    align: 'start',
    value: 'tag',
    sortable: true,
    width: 130,
  },
  {
    title: t('COL_SBOM_FORMAT'),
    align: 'start',
    value: 'metaInfo.spdxVersion',
    sortable: true,
    width: 100,
  },
  {
    title: t('COL_SBOM_ORIGIN'),
    align: 'start',
    value: 'origin',
    sortable: true,
    width: 100,
  },
  {
    title: t('COL_SBOM_UPLOADER'),
    align: 'start',
    value: 'uploader',
    sortable: true,
    width: 110,
  },
  {
    title: t('COL_UPLOADED'),
    align: 'start',
    value: 'uploaded',
    sortable: true,
    width: 110,
  },
]);

const items = computed((): DataTableItems[] => {
  const getSearchIndex = (file: VersionSbomsFlat) => {
    const approvalInfo = file.approvalInfo
      ? ` ${t(`SBOM_STATUS_${file.approvalInfo.status}`)} ${file.approvalInfo.comment}`
      : '';

    return `${file._key} ${file.metaInfo.name} ${formatDateTime(file.uploaded)}${approvalInfo}`;
  };

  if (!props.channelView) {
    return sbomStore.getAllSBOMsFlat.map((file) => ({
      ...file,
      searchIndex: getSearchIndex(file),
    }));
  }
  return spdxFileHistory.value.map((file) => ({
    ...(file as VersionSbomsFlat),
    versionName: versionDetails.value.name,
    versionKey: versionDetails.value._key,
    searchIndex: getSearchIndex(file as VersionSbomsFlat),
  }));
});

const filteredList = computed((): DataTableItems[] => {
  return items.value.filter(filterOnChannel);
});

const possibleChannels = computed((): DataTableHeaderFilterItems[] => {
  if (!items.value) {
    return [];
  }

  const uniqueVersionNames = [...new Set(items.value.map((item: VersionSbomsFlat) => item.versionName))];

  return uniqueVersionNames.map((value: string) => {
    return {
      value,
    } as DataTableHeaderFilterItems;
  });
});

const isOwnerOrDomainAdmin = computed(
  (): boolean =>
    projectModel.value &&
    projectModel.value.accessRights &&
    (projectModel.value.accessRights.groups.includes(Group.ProjectOwner) ||
      projectModel.value.accessRights.groups.includes(Group.UserDomainAdmin)),
);

const filterOnChannel = (item: VersionSbomsFlat) => {
  return selectedFilterChannel.value.length === 0 || selectedFilterChannel.value.includes(item.versionName);
};

const getReferenceInfoForClipboard = (item: VersionSbomsFlat): string => {
  const schemaLabelName = labelTools.value.schemaLabelsMap[projectModel.value.schemaLabel]
    ? labelTools.value.schemaLabelsMap[projectModel.value.schemaLabel].name
    : 'UNKNOWN_LABEL';
  const policyLabelNames = projectModel.value.policyLabels
    .map((l: string) =>
      labelTools.value.policyLabelsMap[l] ? labelTools.value.policyLabelsMap[l].name : 'UNKNOWN_LABEL',
    )
    .join(', ');
  const tabName = 'component';
  const defaultPoicyFilter = 'NOT_SET';
  const deleviryLink = `https://${window.location.host}/#/dashboard/projects/${encodeURIComponent(projectModel.value._key)}/versions/${encodeURIComponent(item.versionKey)}/${tabName}/${defaultPoicyFilter}/${item._key}`;

  return `Disclosure Portal SBOM Reference

Project Name: ${projectModel.value.name}
Project Identifier: ${projectModel.value._key}
Project Schema Label: ${schemaLabelName}
Project Policy Labels: ${policyLabelNames}
Project Version: ${item.versionName}
Version Identifier:  ${item.versionKey}
Reference Timestamp: ${formatDateAndTime(dayjs().toISOString())} (UTC)
SBOM Name: ${item.metaInfo.name}
SBOM Identifier: ${item._key}
Origin: ${item.origin}
Uploader: ${item.uploader}
Upload Date: ${formatDateTimeShort(item.uploaded, true)} (UTC)
SBOM SHA-256: ${item.hash}
Deliveries Link: ${deleviryLink}`;
};

const reloadSboms = async () => {
  await sbomStore.fetchAllSBOMsFlat(true);
};
const toggleLock = async (item: VersionSbomsFlat) => {
  await projectService.toggleSpdxLock(projectModel.value._key, item.versionKey, item._key);
  await reloadSboms();
};

const setApprovable = async (item: VersionSbomsFlat) => {
  const approvableSpdx = {
    spdxkey: '',
    versionkey: '',
  } as ApprovableSPDXDto;
  if (item._key !== projectModel.value.approvablespdx.spdxkey) {
    approvableSpdx.spdxkey = item._key;
    approvableSpdx.versionkey = item.versionKey;
  }
  await projectService
    .updateApprovableSpdx(approvableSpdx, projectModel.value._key)
    .then(() => (projectModel.value.approvablespdx = approvableSpdx));
  await sbomStore.fetchAllSBOMsFlat(true);
};

const downloadFile = (item: VersionSbomsFlat) => {
  const link = document.createElement('a');
  link.click();
  link.target = '_blank';
  projectService
    .downloadSpdxHistoryFile(projectModel.value._key, item.versionKey, item._key)
    .then((res) => {
      const spdxFiles = items.value.filter((sbomFile) => sbomFile._key === item._key);
      if (spdxFiles && spdxFiles.length > 0) {
        const updated = formatDateTimeShort(spdxFiles[0].updated);
        link.download = item.versionName + '_' + updated + '.json';
        link.href = URL.createObjectURL(new Blob([res.data as unknown as BlobPart]));
        link.click();
      }
    })
    .catch((e) => {
      console.error('cannot find spdxFile ' + e);
    });
};

const uploadProgress = (file: File, progress: number) => {
  idle.show(t('PROGRESS_UPLOADING') + ' (' + file.name + ')', progress);
};

const fileUploaded = (_file: File, response: any) => {
  if (response.docIsValid) {
    snack(t('upload_spdx_description'));
    reloadSboms();
  } else {
    if (response.validationFailedMessage === '') {
      const d = new ErrorDialogConfig();
      d.description = t('upload_error_message');
      d.title = '' + t('VALIDATE_SCHEMA');
      d.copyDesc = true;
      d.description += response.message + ' ' + response.raw;
      d.reqId = response.reqID;
      eventBus.emit('on-error', {error: d});
    } else {
      dlgSbomValidationErrors.value?.open(response.validationFailedMessage, helpText.value);
    }
  }
  idle.hide();
};

const fileUploadFailed = () => {
  idle.hide();
};

const updateContextHelp = () => {
  const ht = route.meta?.helpText as Record<string, string>;
  if (ht?.[appLanguage.value]) {
    helpText.value = ht?.[appLanguage.value];
  } else {
    helpText.value = '';
  }
};

const doDelete = async (config: IConfirmationDialogConfig) => {
  await projectService.deleteSpdx(projectModel.value._key, config.contextKey!, config.key);
  snack(t('SBOM_DELETED'));
  await reloadSboms();
};
const onUploadUrlChangedAndShowFileDialog = (url: string) => {
  uploadURL.value = url;
  upload.value?.uploadClick();
};

const openSBOM = (event: Event, item: DataTableItem<VersionSbomsFlat>) => {
  const version: VersionSbomsFlat = item.item;
  emit('openVersion', [version.versionKey]);
  const url = `/dashboard/projects/${encodeURIComponent(projectModel.value._key)}/versions/${encodeURIComponent(version.versionKey)}/overview/${encodeURIComponent(version._key)}`;
  router.push(url);
};

const uploadSPDXFile = () => {
  if (!selectedBranch.value || !projectModel.value) {
    snack(t('SBOM_UPLOAD_DISABLED'));
    return;
  }
  uploadURL.value =
    config.SERVER_URL +
    '/api/v1/projects/' +
    encodeURIComponent(projectModel.value._key) +
    '/versions/' +
    encodeURIComponent(selectedBranch.value.key) +
    '/spdx';
  onUploadUrlChangedAndShowFileDialog(uploadURL.value);
};

const showConfirm = (item: VersionSbomsFlat) => {
  confirmConfig.value = {
    type: ConfirmationType.NOT_SET,
    key: item._key,
    contextKey: item.versionKey,
    name: item.metaInfo.name,
    okButtonIsDisabled: false,
    okButton: 'BTN_DELETE',
    description: 'DLG_CONFIRMATION_DESCRIPTION',
  } as IConfirmationDialogConfig;
  confirmVisible.value = true;
};

const openReviewRemarkDialog = (sbom: VersionSbomsFlat) => {
  reviewRemarkDialog.value?.open({
    versionID: sbom.versionKey,
    spdxID: sbom._key,
  });
};

const copySbomToClipboard = (item: VersionSbomsFlat) => {
  const content = getReferenceInfoForClipboard(item);
  copyToClipboard(content);
};

const getActionButtons = (item: VersionSbomsFlat): TableActionButtonsProps['buttons'] => {
  const isApprovable = projectModel.value && projectModel.value.approvablespdx.spdxkey == item._key;
  const canSetApprovable =
    projectModel.value &&
    projectModel.value.accessRights &&
    projectModel.value.accessRights.groups.find((g: string) => g == 'Owner');

  return [
    {
      icon: isApprovable ? 'mdi-star' : 'mdi-star-outline',
      hint: isApprovable ? t('TT_approvable_spdx') : t('TT_not_approvable_spdx'),
      event: 'setApprovable',
      show: !!canSetApprovable,
      disabled: projectModel.value.isDeprecated,
    },
    {
      icon: item.isLocked ? 'mdi-lock-outline' : 'mdi-lock-open-variant-outline',
      hint: item.isLocked ? t('TT_unlock_spdx') : t('TT_lock_spdx'),
      event: 'toggleLock',
      show: isOwnerOrDomainAdmin.value,
      disabled: projectModel.value.isDeprecated,
    },

    {
      icon: 'mdi-message-plus-outline',
      hint: t('TT_add_review_remark'),
      event: 'addRemark',
      show: true,
      disabled: projectModel.value.isDeprecated,
    },
    {
      icon: 'mdi-content-copy',
      hint: t('TT_COPY_REFERENCE_INFO'),
      event: 'copy',
      show: true,
    },
    {
      icon: 'mdi-download',
      hint: t('TT_download_spdx'),
      event: 'download',
      show: projectModel.value?.accessRights?.allowSBOMAction?.download,
    },
    {
      icon: 'mdi-delete',
      hint: t('TT_delete_spdx'),
      event: 'delete',
      show: isOwnerOrDomainAdmin.value,
      disabled: item.isInUse || item.isLocked || item.isToRetain || projectModel.value.isDeprecated,
    },
  ];
};

watch(appLanguage, () => {
  updateContextHelp();
});

onMounted(async () => {
  updateContextHelp();
  if (!props.channelView) {
    sbomStore.fetchAllSBOMsFlat().then(() => {
      selectedBranch.value = branches.value[0];
      if (versionDetails.value) {
        const branchFromVersion = branches.value.find((g) => g.key == versionDetails.value._key);
        if (branchFromVersion) {
          selectedBranch.value = branchFromVersion;
          isBranchSelectionEnabled.value = false;
        }
      }
    });
  } else {
    selectedBranch.value = {
      name: versionDetails.value.name,
      key: versionDetails.value._key,
    };
    isBranchSelectionEnabled.value = false;
  }
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #description>
      <!-- v-html is used to render because DOWNLOAD_INTITIAL_DOCUMENT contains html -->
      <span class="text-caption" v-html="t('SBOM_DELIVERIES_DISCLAIMER_TEXT')"> </span>
    </template>
    <template #buttons>
      <div>
        <DiscoFileUpload
          ref="upload"
          :uploadTargetUrl="uploadURL"
          acceptTypes=".json,.spdx"
          @reqFailed="fileUploadFailed"
          @reqFinished="fileUploaded"
          @reqProgress="uploadProgress" />
        <DCActionButton
          :text="t('BTN_UPLOAD')"
          icon="mdi-upload"
          :hint="selectedBranch?.name ? t('BTN_UPLOAD') : t('SBOM_UPLOAD_DISABLED')"
          @clicked="uploadSPDXFile"
          v-if="projectModel && projectModel.accessRights && projectModel.accessRights.allowSBOMAction.upload" />
      </div>
      <v-select
        v-model="selectedBranch"
        density="compact"
        variant="outlined"
        :disabled="!isBranchSelectionEnabled"
        item-title="name"
        return-object
        :items="branches"
        :label="t('LBL_UPLOAD_CHANNEL')"
        v-if="projectModel && projectModel.accessRights && projectModel.accessRights.allowSBOMAction.upload"
        hide-details />
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <div ref="tableSbomDeliveries" class="fill-height">
        <v-data-table
          density="compact"
          fixed-header
          :sort-by="sortItems"
          :search="search"
          :headers="headers"
          :items="filteredList"
          @click:row="openSBOM"
          :footer-props="{
            'items-per-page-options': [10, 50, 100, -1],
          }"
          class="striped-table fill-height">
          <template #[`header.versionName`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterChannel"
                  :column="column"
                  :label="t('COL_SBOM_BRANCH')"
                  :allItems="possibleChannels">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`item.searchIndex`]="{item}">
            {{ formatDateTime(item.uploaded) }} -&nbsp;{{ item.metaInfo.name }}
            <br />
            <span class="font-weight-bold">UUID: </span>
            <span>{{ item._key }}</span>
            <br v-if="item.approvalInfo && item.approvalInfo.status" />
            <span class="font-weight-bold" v-if="item.approvalInfo && item.approvalInfo.status">{{
              t(`SBOM_STATUS_${item.approvalInfo.status}`)
            }}</span>
            <span v-if="item.approvalInfo && item.approvalInfo.status">
              <v-icon
                small
                v-if="item.approvalInfo && item.approvalInfo.comment && item.approvalInfo.comment.length > 0"
                >chevron_right</v-icon
              >
              {{ item.approvalInfo.comment }}</span
            >
            <br v-if="item.isToDelete" />
            <span v-if="item.isToDelete" class="font-weight-bold text-[rgb(var(--v-theme-error))]">{{
              t('SBOM_ABOUT_DELETION_NOTE')
            }}</span>
            <br v-if="item.isToRetain" />
            <span v-if="item.isToRetain" class="font-weight-bold text-[rgb(var(--v-theme-success))]">{{
              t('SBOM_MARKED_FOR_RETENTION')
            }}</span>
          </template>
          <template #[`item.overallReview`]="{item}">
            <DOverallStateIcon v-if="item.overallReview" :review="item.overallReview" />
          </template>
          <template #[`item.uploaded`]="{item}">
            <DDateCellWithTooltip :value="item.uploaded" />
          </template>
          <template #[`item.tag`]="{item}">
            <v-chip
              v-if="
                item.approvalInfo.isInApproval ||
                projectModel.isDeprecated ||
                !projectModel.accessRights.allowSBOMAction.upload ||
                !projectModel.accessRights.allowSBOMAction.delete
              "
              color="labelBackgroundColor"
              class="mr-1 mb-1 px-2 py-2"
              label>
              <v-icon class="pr-2" small color="labelIconColor" left>mdi-label</v-icon>
              <span v-if="!item.tag" class="letterSpacing">{{ t('SPDX_TAG_UNSET') }}</span>
              <span v-else class="letterSpacing">{{ item.tag }}</span>
            </v-chip>
            <DSpdxTagDialog
              :presetTag="item.tag"
              :versionID="item.versionKey"
              :spdxID="item._key"
              :spdxName="item.metaInfo.name"
              :channel-view="channelView"
              v-slot="{showDialog}"
              v-else>
              <v-chip color="labelBackgroundColor" class="mr-1 mb-1 px-2 py-2" label link @click.stop="showDialog">
                <v-icon class="pr-2" small color="primary" left>mdi-label</v-icon>
                <span v-if="!item.tag" class="letterSpacing">{{ t('SPDX_TAG_UNSET') }}</span>
                <span v-else class="letterSpacing">{{ item.tag }}</span>
              </v-chip>
            </DSpdxTagDialog>
          </template>
          <template v-slot:[`item.origin`]="{item}">
            <Tooltip v-if="originTooltip(item.origin)" location="bottom" :text="originTooltip(item.origin)" as-parent>
              {{ originShort(item.origin) }}
            </Tooltip>
            <span v-else>{{ item.origin }}</span>
          </template>
          <template #[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
              @toggleLock="isOwnerOrDomainAdmin ? toggleLock(item) : undefined"
              @setApprovable="setApprovable(item)"
              @addRemark="openReviewRemarkDialog(item)"
              @copy="copySbomToClipboard(item)"
              @download="downloadFile(item)"
              @delete="showConfirm(item)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <ReviewRemarkDialog ref="reviewRemarkDialog" />
  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doDelete" />
  <SbomValidationErrorsDialog ref="dlgSbomValidationErrors" />
</template>
