<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import ConfirmationDialog from '@disclosure-portal/components/dialog/ConfirmationDialog.vue';
import ReviewRemarkDialog from '@disclosure-portal/components/dialog/ReviewRemarkDialog.vue';
import SbomValidationErrorsDialog from '@disclosure-portal/components/dialog/SbomValidationErrorsDialog.vue';
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import {IDefaultSelectItem} from '@disclosure-portal/model/IObligation';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {NameKeyIdentifier, VersionSbomsFlat} from '@disclosure-portal/model/ProjectsResponse';
import {Group} from '@disclosure-portal/model/Rights';
import {SpdxFile} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import eventBus from '@disclosure-portal/utils/eventbus';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {formatDateTime, formatDateTimeShort, originShort, originTooltip} from '@disclosure-portal/utils/View';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import DIconButton from '@shared/components/disco/DIconButton.vue';
import DSpdxTagDialog from '@shared/components/disco/DSpdxTagDialog.vue';
import Tooltip from '@shared/components/disco/Tooltip.vue';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import DiscoFileUpload from '@shared/components/widgets/DiscoFileUpload.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {DataTabelIndex, DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
import {useClipboard} from '@shared/utils/clipboard';
import config from '@shared/utils/config';
import dayjs from 'dayjs';
import _ from 'lodash';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';

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
const {copyToClipboard} = useClipboard();

const projectModel = computed(() => projectStore.currentProject!);
const versionDetails = computed(() => sbomStore.getCurrentVersion);
const spdxFileHistory = computed(() => sbomStore.getChannelSpdxs);
const labelTools = computed(() => appStore.getLabelsTools);

const search = ref('');
const uploadURL = ref('');
const isBranchSelectionEnabled = ref(true);
const selectedFilterChannel = ref<string[]>([]);
const statusFilterOpened = ref(false);
const selectedBranch = ref<NameKeyIdentifier>({} as NameKeyIdentifier);
const sortItems = ref<SortItem[]>([{key: 'Uploaded', order: 'desc'}]);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmVisible = ref(false);
const {info: snack} = useSnackbar();
const branches = ref<NameKeyIdentifier[]>([]);
const tableSbomDeliveries = ref<HTMLElement | null>(null);
const reviewRemarkDialog = ref<InstanceType<typeof ReviewRemarkDialog>>();
const dlgSbomValidationErrors = ref<InstanceType<typeof SbomValidationErrorsDialog>>();
const helpText = ref('');
const upload = ref<InstanceType<typeof DiscoFileUpload>>();

const sortByName = (a: SpdxFile, b: SpdxFile): number => {
  return b.MetaInfo.Name.localeCompare(a.MetaInfo.Name);
};

const headers = (): DataTableHeader[] => {
  const res: DataTableHeader[] = [
    {
      title: t('COL_ACTIONS'),
      sortable: false,
      align: 'center',
      class: 'tableHeaderCell',
      value: 'Actions',
      width: 100,
    },
    {
      title: t('COL_SPDX_FILENAME'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'searchIndex',
      width: 380,
      sortable: true,
      sortRaw: sortByName,
    },
    {
      title: t('COL_REVIEW_STATUS'),
      align: 'center',
      class: 'tableHeaderCell',
      value: 'OverallReview',
      width: 80,
    },
  ];
  if (!props.channelView) {
    res.push({
      title: t('COL_SBOM_BRANCH'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'versionName',
      sortable: true,
      width: 110,
    });
  }
  res.push(
    {
      title: t('COL_SBOM_TAG'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'Tag',
      sortable: true,
      width: 100,
    },
    {
      title: t('COL_SBOM_FORMAT'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'MetaInfo.SpdxVersion',
      sortable: true,
      width: 100,
    },
    {
      title: t('COL_SBOM_ORIGIN'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'Origin',
      sortable: true,
      width: 100,
    },
    {
      title: t('COL_SBOM_UPLOADER'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'Uploader',
      sortable: true,
      width: 110,
    },
    {
      title: t('COL_UPLOADED'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'Uploaded',
      sortable: true,
      width: 110,
    },
  );

  return res;
};

const items = computed((): DataTableItems[] => {
  const getSearchIndex = (file: VersionSbomsFlat) => {
    const approvalInfo = file.ApprovalInfo
      ? ` ${t(`SBOM_STATUS_${file.ApprovalInfo.Status}`)} ${file.ApprovalInfo.Comment}`
      : '';

    return `${file._key} ${file.MetaInfo.Name} ${formatDateTime(file.Uploaded)}${approvalInfo}`;
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

const possibleChannels = computed((): IDefaultSelectItem[] => {
  if (!items.value) {
    return [];
  }

  return _.chain(items.value)
    .uniqBy((item: VersionSbomsFlat) => item.versionName)
    .map((item: VersionSbomsFlat) => {
      return {
        text: item.versionName,
        value: item.versionName,
      } as IDefaultSelectItem;
    })
    .value();
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
SBOM Name: ${item.MetaInfo.Name}
SBOM Identifier: ${item._key}
Origin: ${item.Origin}
Uploader: ${item.Uploader}
Upload Date: ${formatDateTimeShort(item.Uploaded, true)} (UTC)
SBOM SHA-256: ${item.Hash}
Deliveries Link: ${deleviryLink}`;
};

const reloadSboms = async () => {
  await sbomStore.fetchAllSBOMsFlat();
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
        const updated = dayjs(spdxFiles[0].Updated.toString()).format(t('DATETIME_FORMAT_SHORT'));
        link.download = item.versionName + '_' + updated + '.json';
        link.href = URL.createObjectURL(new Blob([res.data as unknown as BlobPart]));
        link.click();
      }
    })
    .catch((e) => {
      console.error('cannot find spdxFile ' + e);
    });
};

onMounted(async () => {
  if (!props.channelView) {
    sbomStore.fetchAllSBOMsFlat().then(() => {
      branches.value = sbomStore.allVersions;
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

watch(
  () => appStore.getAppLanguage,
  () => {
    updateContextHelp();
  },
);

onMounted(() => {
  updateContextHelp();
});

const updateContextHelp = () => {
  const lang = appStore.getAppLanguage as 'de' | 'en';
  const ht = route.meta?.helpText as {de: string; en: string};
  if (ht?.[lang]) {
    helpText.value = ht?.[lang];
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
    name: item.MetaInfo.Name,
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
      icon: item.IsLocked ? 'mdi-lock-outline' : 'mdi-lock-open-variant-outline',
      hint: item.IsLocked ? t('TT_unlock_spdx') : t('TT_lock_spdx'),
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
      disabled: item.IsInUse || item.IsLocked || item.IsToRetain || projectModel.value.isDeprecated,
    },
  ];
};
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
      <v-text-field
        autocomplete="off"
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        clearable
        density="compact"
        variant="outlined"
        hide-details></v-text-field>
    </template>
    <template #table>
      <div ref="tableSbomDeliveries" class="fill-height">
        <v-data-table
          density="compact"
          fixed-header
          :sort-by="sortItems"
          :search="search"
          :headers="headers()"
          :items="filteredList"
          @click:row="openSBOM"
          :footer-props="{
            'items-per-page-options': [10, 50, 100, -1],
          }"
          class="striped-table fill-height">
          <template v-slot:[`header.versionName`]="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" v-model="statusFilterOpened">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="selectedFilterChannel.length > 0 ? 'primary' : 'default'"
                    location="top" />
                </template>
                <div class="bg-background" style="width: 280px">
                  <v-row class="d-flex ma-1 mr-2 justify-end">
                    <DIconButton icon="mdi-close" @clicked="statusFilterOpened = false" color="default" />
                  </v-row>
                  <v-select
                    v-model="selectedFilterChannel"
                    :items="possibleChannels"
                    class="pa-2 mx-2 pb-4"
                    :label="t('Lbl_filter_branches')"
                    clearable
                    multiple
                    item-title="text"
                    item-value="value"
                    variant="outlined"
                    density="compact"
                    menu
                    transition="scale-transition"
                    persistent-clear
                    :list-props="{class: 'striped-filter-dd py-0'}">
                    <template v-slot:item="{props}">
                      <v-list-item v-bind="props" class="px-2 py-0">
                        <template v-slot:prepend="{isSelected}">
                          <v-checkbox hide-details :model-value="isSelected" />
                        </template>
                        <template v-slot:title="{title}">
                          <span class="pFilterEntry">
                            {{ title }}
                          </span>
                        </template>
                      </v-list-item>
                    </template>
                    <template v-slot:selection="{item, index}">
                      <div v-if="index === 0" class="d-flex align-center">
                        <span class="pFilterEntry">{{ item.title }}</span>
                      </div>
                      <span v-if="index === 1" class="pAdditionalFilter">
                        +{{ selectedFilterChannel.length - 1 }} others
                      </span>
                    </template>
                  </v-select>
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>
          <template v-slot:[`item.searchIndex`]="{item}">
            {{ formatDateTime(item.Uploaded) }} -&nbsp;{{ item.MetaInfo.Name }}
            <br />
            <span class="font-weight-bold">UUID: </span>
            <span>{{ item._key }}</span>
            <br v-if="item.ApprovalInfo && item.ApprovalInfo.Status" />
            <span class="font-weight-bold" v-if="item.ApprovalInfo && item.ApprovalInfo.Status">{{
              t(`SBOM_STATUS_${item.ApprovalInfo.Status}`)
            }}</span>
            <span v-if="item.ApprovalInfo && item.ApprovalInfo.Status">
              <v-icon
                small
                v-if="item.ApprovalInfo && item.ApprovalInfo.Comment && item.ApprovalInfo.Comment.length > 0"
                >chevron_right</v-icon
              >
              {{ item.ApprovalInfo.Comment }}</span
            >
            <br v-if="item.IsToDelete" />
            <span v-if="item.IsToDelete" class="font-weight-bold text-[rgb(var(--v-theme-error))]">{{
              t('SBOM_ABOUT_DELETION_NOTE')
            }}</span>
            <br v-if="item.IsToRetain" />
            <span v-if="item.IsToRetain" class="font-weight-bold text-[rgb(var(--v-theme-success))]">{{
              t('SBOM_MARKED_FOR_RETENTION')
            }}</span>
          </template>
          <template v-slot:[`item.OverallReview`]="{item}">
            <DOverallStateIcon v-if="item.OverallReview" :review="item.OverallReview" />
          </template>
          <template v-slot:[`item.Uploaded`]="{item}">
            <DDateCellWithTooltip :value="item.Uploaded" />
          </template>
          <template v-slot:[`item.Tag`]="{item}">
            <v-chip
              v-if="
                item.ApprovalInfo.IsInApproval ||
                projectModel.isDeprecated ||
                !projectModel.accessRights.allowSBOMAction.upload ||
                !projectModel.accessRights.allowSBOMAction.delete
              "
              color="labelBackgroundColor"
              class="mr-1 mb-1 px-2 py-2"
              label>
              <v-icon class="pr-2" small color="labelIconColor" left>mdi-label</v-icon>
              <span v-if="!item.Tag" class="letterSpacing">{{ t('SPDX_TAG_UNSET') }}</span>
              <span v-else class="letterSpacing">{{ item.Tag }}</span>
            </v-chip>
            <DSpdxTagDialog
              :presetTag="item.Tag"
              :versionID="item.versionKey"
              :spdxID="item._key"
              :spdxName="item.versionName"
              :channel-view="channelView"
              v-slot="{showDialog}"
              v-else>
              <v-chip color="labelBackgroundColor" class="mr-1 mb-1 px-2 py-2" label link @click.stop="showDialog">
                <v-icon class="pr-2" small color="primary" left>mdi-label</v-icon>
                <span v-if="!item.Tag" class="letterSpacing">{{ t('SPDX_TAG_UNSET') }}</span>
                <span v-else class="letterSpacing">{{ item.Tag }}</span>
              </v-chip>
            </DSpdxTagDialog>
          </template>
          <template v-slot:[`item.Origin`]="{item}">
            <Tooltip v-if="originTooltip(item.Origin)" location="bottom" :text="originTooltip(item.Origin)" as-parent>
              {{ originShort(item.Origin) }}
            </Tooltip>
            <span v-else>{{ item.Origin }}</span>
          </template>
          <template v-slot:[`item.Actions`]="{item}">
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
