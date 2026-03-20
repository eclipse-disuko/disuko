<script setup lang="ts">
import {IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {DialogReviewRemarkConfig} from '@disclosure-portal/components/dialog/DialogConfigs';
import ReviewRemarksDetailsDialog from '@disclosure-portal/components/dialog/ReviewRemarksDetailsDialog.vue';
import {useReviewRemarkActions} from '@disclosure-portal/composables/useReviewRemarkActions';
import Icons from '@disclosure-portal/constants/icons';
import {PolicyDecisionSlim} from '@disclosure-portal/model/PolicyDecision';
import {ComponentDetails, Details, ProjectModel, UnmatchedLicense} from '@disclosure-portal/model/Project';
import {LicenseMeta, ReviewRemark, ReviewRemarkStatus, compareRRLevel} from '@disclosure-portal/model/Quality';
import {ComponentInfoSlim, PolicyRuleStatus} from '@disclosure-portal/model/VersionDetails';
import versionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import {
  getIconColorForPolicyType,
  getIconColorReviewRemarkLevel,
  getIconForPolicyType,
  getIconReviewRemarkLevel,
  policyStateToTranslationKey,
} from '@disclosure-portal/utils/View';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader} from '@shared/types/table';
import _l from 'lodash';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import JsonViewer3 from 'vue-json-viewer';

interface LocalDetails extends Details {
  url?: boolean;
}

interface LocalComponentDetails extends ComponentDetails {
  Attributes: LocalDetails[];
}

const emit = defineEmits(['reloadAfterCreation', 'triggerBulk']);

const userStore = useUserStore();
const snack = useSnackbar();
const {t} = useI18n();

const {
  confirmCloseConfig,
  confirmCancelConfig,
  confirmInProgressConfig,
  closeVisible,
  cancelVisible,
  inProgressVisible,
  isOpen,
  isInProgress,
  openCloseRemarkDialog,
  openCancelRemarkDialog,
  doCloseRemark: doCloseRemarkAction,
  doCancelRemark: doCancelRemarkAction,
  doMarkInProgress: doMarkInProgressAction,
} = useReviewRemarkActions();

const responsible = ref(false);
const show = ref(false);
const name = ref('');
const version = ref('');
const description = ref('');
const selectedTab = ref(0);
const headers = ref<DataTableHeader[]>([
  {
    title: t('COL_ACTIONS'),
    align: 'center',
    width: 150,
    filterable: false,
    class: 'tableHeaderCell',
    value: 'action',
  },
  {
    title: t('KEY'),
    align: 'start',
    filterable: true,
    class: 'tableHeaderCell',
    value: 'Key',
  },
  {
    title: t('COL_VALUES'),
    align: 'start',
    filterable: true,
    class: 'tableHeaderCell',
    value: 'Value',
  },
]);
const project = ref<ProjectModel>({} as ProjectModel);
const projectVersionId = ref('');
const sbomId = ref('');
const noLicenses = ref(false);
const noSbomLicenses = ref(false);
const showTooltipUnmatched = ref<boolean[]>([]);
const showTooltip = ref<boolean[]>([]);
const details = ref<LocalComponentDetails>({} as LocalComponentDetails);
const reviewRemarkDialog = ref();
const licenseRuleDialog = ref();
const policyDecisionDialog = ref();
const viewRemarkDialog = ref();
const reviewRemarks = ref<ReviewRemark[]>([]);
const loadingRemarks = ref(false);

const isDeprecated = computed(() => useProjectStore().currentProject!.isDeprecated);

const fetchReviewRemarks = async (projectKey: string, versionKey: string, sbomUuid: string, spdxId: string) => {
  loadingRemarks.value = true;
  try {
    const response = await versionService.getReviewRemarksForComponent(projectKey, versionKey, sbomUuid, spdxId);
    reviewRemarks.value = response.data;
  } catch (error) {
    console.error('Failed to fetch review remarks:', error);
    reviewRemarks.value = [];
  } finally {
    loadingRemarks.value = false;
  }
};

const open = async (
  projectData: ProjectModel,
  versionKey: string,
  sbomIdData: string,
  data: ComponentDetails,
  policyStatus?: PolicyRuleStatus[],
  unmatched?: UnmatchedLicense[],
  policyDecisionsApplied?: PolicyDecisionSlim[],
  policyDecisionDeniedReason?: string,
) => {
  responsible.value = userStore.getProfile.user === projectData.responsible;
  details.value = data;

  details.value.Attributes = details.value.Attributes?.map((attr) => {
    if (attr.Key === 'homepage') {
      attr.url = true;
    }

    return attr;
  });

  if (policyStatus) {
    details.value.PolicyStatus = policyStatus;
  }
  if (unmatched) {
    details.value.UnmatchedLicenses = unmatched;
  }
  if (policyDecisionsApplied) {
    details.value.PolicyDecisionsApplied = policyDecisionsApplied;
  }
  if (policyDecisionDeniedReason) {
    details.value.PolicyDecisionDeniedReason = policyDecisionDeniedReason;
  }

  project.value = projectData;
  projectVersionId.value = versionKey;
  sbomId.value = sbomIdData;

  noLicenses.value = data.UnknownLicenses?.length === 0 && data.KnownLicenses?.length === 0;
  noSbomLicenses.value = data.ExtractedLicenses?.length === 0 && data.IdentifiedViaAlias?.length === 0;
  name.value = details.value.Attributes?.find((a) => a.Key === 'name')?.Value ?? '';
  description.value = details.value.Attributes?.find((a) => a.Key === 'description')?.Value ?? '';
  version.value = details.value.Attributes?.find((a) => a.Key === 'versionInfo')?.Value ?? '';

  show.value = true;
  selectedTab.value = 0;

  showTooltip.value = [];
  showTooltipUnmatched.value = [];

  // Fetch review remarks for this component
  const spdxId = details.value.Attributes?.find((a) => a.Key === 'SPDXID')?.Value ?? '';
  if (spdxId) {
    await fetchReviewRemarks(projectData._key, versionKey, sbomIdData, spdxId);
  }
};

const reloadReviewRemarks = async () => {
  const spdxId = details.value.Attributes?.find((a) => a.Key === 'SPDXID')?.Value ?? '';
  if (spdxId && project.value._key && projectVersionId.value && sbomId.value) {
    await fetchReviewRemarks(project.value._key, projectVersionId.value, sbomId.value, spdxId);
  }
};

const doCloseRemark = async (config: IConfirmationDialogConfig) => {
  if (!project.value._key || !projectVersionId.value) return;
  await doCloseRemarkAction(config, project.value._key, projectVersionId.value, reloadReviewRemarks);
};

const doCancelRemark = async (config: IConfirmationDialogConfig) => {
  if (!project.value._key || !projectVersionId.value) return;
  await doCancelRemarkAction(config, project.value._key, projectVersionId.value, reloadReviewRemarks);
};

const doMarkInProgress = async (config: IConfirmationDialogConfig) => {
  if (!project.value._key || !projectVersionId.value) return;
  await doMarkInProgressAction(config, project.value._key, projectVersionId.value, reloadReviewRemarks);
};

const close = () => {
  show.value = false;
};

const closeAndReload = () => {
  show.value = false;
  emit('reloadAfterCreation');
};

const closeAndTriggerBulk = () => {
  show.value = false;
  emit('triggerBulk');
};

const searchForSimilarLicense = (licenseText: string) => {
  localStorage.setItem('licenseText', licenseText);
  openLinkToLicenses();
};

const openLinkToLicenses = () => {
  const url = '/#/dashboard/licenses/compare';
  window.open(url);
};

const viewRemark = (remark: ReviewRemark) => {
  if (remark) {
    viewRemarkDialog.value?.open(remark);
  }
};

const openReviewRemarkDialog = (licenseId: string | undefined = undefined) => {
  const comp = {
    spdxId: details.value.Attributes?.find((a) => a.Key === 'SPDXID')?.Value ?? '',
    name: name.value,
    version: version.value,
    licenseExpression: '',
    componentInfo: [],
  };
  const licenses: LicenseMeta[] = [];
  if (licenseId) {
    const license = new LicenseMeta();
    license.licenseId = licenseId;
    license.licenseName = '';
    licenses.push(license);
  }
  reviewRemarkDialog.value?.open({
    spdxID: sbomId.value,
    components: [comp],
    licenses: licenses,
  } as DialogReviewRemarkConfig);
};

const openLicenseRuleDialog = (licenseId: string) => {
  const componentSpdxId = details.value.Attributes?.find((a) => a.Key === 'SPDXID')?.Value ?? '';
  const licenseExpression = details.value.Attributes?.find((a) => a.Key === 'licenseEffective')?.Value ?? '';

  const component = new ComponentInfoSlim();
  component.spdxId = componentSpdxId;
  component.name = name.value;
  component.version = version.value;
  component.licenseExpression = licenseExpression;

  licenseRuleDialog.value?.open({
    licenseId,
    component,
    policyStatus: details.value.PolicyStatus,
  });
};

const openPolicyDecisionDialog = (policy: PolicyRuleStatus | null) => {
  if (policy === null) {
    return;
  }

  const componentSpdxId = details.value.Attributes?.find((a) => a.Key === 'SPDXID')?.Value ?? '';
  const licenseExpression = details.value.Attributes?.find((a) => a.Key === 'licenseEffective')?.Value ?? '';

  const component = new ComponentInfoSlim();
  component.spdxId = componentSpdxId;
  component.name = name.value;
  component.version = version.value;

  component.licenseExpression = licenseExpression;

  policyDecisionDialog.value?.open({
    component,
    policies: [policy],
    type: policy.type,
  });
};

const isInExtractedLicenses = (id: string) => {
  const res = _l.some(
    details.value.ExtractedLicenses,
    (license) => license.LicenseId.toLowerCase() === id.toLowerCase(),
  );
  if (!res) {
    return {ExtractedText: ''};
  }
  return res;
};

const getExtractedLicenseById = (id: string) => {
  const res = _l.find(
    details.value.ExtractedLicenses,
    (license) => license.LicenseId.toLowerCase() === id.toLowerCase(),
  );
  if (!res) {
    return {ExtractedText: id};
  }
  return res;
};

const mailtoLink = (item: UnmatchedLicense) => {
  if (item.known) {
    const lic = _l.find(
      details.value.KnownLicenses,
      (license) => license.License.licenseId.toLowerCase() === item.referenced.toLowerCase(),
    );
    if (!lic) {
      return '';
    }
    const profile = userStore.getProfile;
    const projectLink = `${window.location.origin}/#/dashboard/projects/${encodeURIComponent(project.value._key)}`;
    const versionLink = `${projectLink}/versions/${encodeURIComponent(projectVersionId.value)}`;
    const sbomLink = `${versionLink}/component/NOT_SET/${sbomId.value}`;
    const licenseLink = `${window.location.origin}/#/dashboard/licenses/${lic.License.licenseId}`;
    return t('REQ_REVIEW_MAIL_CONTENT')
      .replace(/%name/g, lic.License.name)
      .replace(/%id/g, lic.OrigName)
      .replace(/%requestor/, `${profile.forename} ${profile.lastname}`)
      .replace(/%projectName/, project.value.name)
      .replace(/%projectLink/, projectLink)
      .replace(/%versionLink/, versionLink)
      .replace(/%sbom/, sbomLink)
      .replace(/%comp/, name.value)
      .replace(/%licenseLink/, licenseLink)
      .replace(/%user/g, `${profile.forename} ${profile.lastname}`);
  }
  const extracted = getExtractedLicenseById(item.orig);
  let text = '';
  if (extracted) {
    text = extracted.ExtractedText.substring(0, 200);
  }
  const profile = userStore.getProfile;
  const projectLink = `${window.location.origin}/#/dashboard/projects/${encodeURIComponent(project.value._key)}`;
  const versionLink = `${projectLink}/versions/${encodeURIComponent(projectVersionId.value)}`;
  const sbomLink = `${versionLink}/component/NOT_SET/${sbomId.value}`;
  return t('REQ_REVIEW_UNKNOWN_MAIL_CONTENT')
    .replace(/%id/g, item.orig)
    .replace(/%requestor/, `${profile.forename} ${profile.lastname}`)
    .replace(/%projectName/, project.value.name)
    .replace(/%projectLink/, projectLink)
    .replace(/%versionLink/, versionLink)
    .replace(/%sbom/, sbomLink)
    .replace(/%comp/, name.value)
    .replace(/%text/, text)
    .replace(/%user/g, `${profile.forename} ${profile.lastname}`);
};

const sendReviewMail = (item: UnmatchedLicense) => {
  const mailToUri = mailtoLink(item);
  if (!mailToUri || mailToUri === '') {
    return;
  }
  navigator.clipboard
    .writeText(mailToUri)
    .then(() => {
      snack.info(t('SNACK_INVITATION_COPIED'));
    })
    .catch(() => {
      snack.info(t('SNACK_WENT_WRONG'));
    });
  const link = encodeURI(mailToUri);
  if (!link) {
    return;
  }
  window.open(link, '_blank');
};

const getLicenseEffectiveAttribute = (attributes: Details[]) => {
  const foundAttribute = _l.find(attributes, {Key: 'licenseEffective'});
  return foundAttribute ? foundAttribute.Value : 'NOASSERTION';
};

const dialogLayoutConfig = computed(() => {
  const desc =
    description.value !== 'NOASSERTION'
      ? description.value.length > 75
        ? description.value.substring(0, 75) + '...'
        : description.value
      : '';

  return {
    title: name.value,
    titleTooltip: `${name.value} ${version.value} ${desc}`,
    primaryButton: {text: t('BTN_CLOSE')},
  };
});

const findPolicyDecisionApplied = (item: PolicyRuleStatus): PolicyDecisionSlim | null => {
  if (!item.isDecisionMade) {
    return null;
  }
  return (
    details.value.PolicyDecisionsApplied.find(
      (pd) => pd.policyId === item.key && pd.licenseId.toLowerCase() === item.licenseMatched.toLowerCase(),
    ) ?? null
  );
};

const isPolicyDecisionPresent = computed(() => details.value.PolicyDecisionsApplied.length > 0);

const sortedReviewRemarks = computed(() => {
  const statusOrder = new Map<ReviewRemarkStatus, number>([
    [ReviewRemarkStatus.OPEN, 1],
    [ReviewRemarkStatus.IN_PROGRESS, 2],
    [ReviewRemarkStatus.CLOSED, 3],
    [ReviewRemarkStatus.CANCELLED, 4],
    [ReviewRemarkStatus.NOT_SET, 5],
  ]);

  return [...reviewRemarks.value].sort((a, b) => {
    const statusCompare = statusOrder.get(a.status)! - statusOrder.get(b.status)!;
    if (statusCompare !== 0) return statusCompare;

    return compareRRLevel(b.level, a.level);
  });
});

defineExpose({
  open,
});
</script>

<template>
  <v-dialog v-model="show" scrollable width="1000" max-height="700" min-width="630" height="700">
    <DialogLayout :config="dialogLayoutConfig" @close="close" @primaryAction="close">
      <template #title-right>
        <DCActionButton
          class="ml-2"
          large
          icon="mdi-message-plus-outline"
          :hint="t('TT_add_review_remark')"
          :text="t('TT_add_review_remark')"
          @click="openReviewRemarkDialog()"
          v-if="!isDeprecated" />
      </template>

      <v-card class="card-border" min-height="394">
        <v-tabs v-model="selectedTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
          <v-tab value="rules">
            {{ t('TAB_TITLE_POLICY_RULES') }}
          </v-tab>
          <v-tab value="commonLicenses">
            {{ t('TAB_TITLE_LICENSES') }}
          </v-tab>
          <v-tab value="sbomLicenses">
            {{ t('TAB_TITLE_SBOM_LICENSES') }}
          </v-tab>
          <v-tab value="attributes">
            {{ t('TAB_TITLE_ATTRIBUTES') }}
          </v-tab>
          <v-tab value="additional_info">
            {{ t('TAB_TITLE_REMARKS') }}
            <v-badge v-if="reviewRemarks.length > 0" :content="reviewRemarks.length" color="mbti" inline></v-badge>
          </v-tab>
          <v-tab value="raw">
            {{ t('TAB_TITLE_RAW') }}
          </v-tab>
        </v-tabs>

        <v-tabs-window v-model="selectedTab">
          <v-tabs-window-item class="dialogContent" value="rules">
            <Stack>
              <v-table fixed-header density="compact" height="342">
                <template v-slot:default>
                  <thead>
                    <tr>
                      <th class="text-center">
                        {{ t('COL_ACTIONS') }}
                      </th>
                      <th>
                        {{ t('COL_STATUS') }}
                      </th>
                      <th class="text-left">
                        {{ t('COL_LICENSE') }}
                      </th>
                      <th class="text-left">
                        {{ t('COL_POLICY_NAME') }}
                      </th>
                      <th class="text-left">
                        {{ t('COL_LICENSE_FAMILY') }}
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    <template v-if="details.UnmatchedLicenses?.length === 0 && details.PolicyStatus?.length === 0">
                      <tr>
                        <td>
                          <div class="d-flex justify-center flex-row">
                            <DIconButton :icon="Icons.QUESTIONED" :hint="t('HELP_TT_UNASSERTED')" />
                            <DIconButton
                              icon="mdi-message-plus-outline"
                              :hint="t('UM_DIALOG_TITLE_NEW_REVIEW_REMARK')"
                              @clicked="openReviewRemarkDialog()"
                              :disabled="isDeprecated" />
                            <DIconButton
                              icon="mdi-email-fast-outline"
                              :hint="t('TT_not_request_review_no_licence')"
                              disabled />
                          </div>
                        </td>
                        <td>
                          <v-icon :color="getIconColorForPolicyType('noassertion')">
                            {{ getIconForPolicyType('noassertion') }}
                            <Tooltip location="bottom" :text="policyStateToTranslationKey('noassertion')" />
                          </v-icon>
                        </td>
                        <td>
                          <span>{{ getLicenseEffectiveAttribute(details.Attributes) }}</span>
                        </td>
                        <td></td>
                        <td></td>
                      </tr>
                    </template>
                    <PolicyStatusTableRow
                      v-for="(item, index) in details.PolicyStatus.filter((p) => p.type === 'deny')"
                      :key="`policy-status-deny-${index}`"
                      :item="item"
                      :policyDecisionApplied="findPolicyDecisionApplied(item)"
                      :isPolicyDecisionPresent="isPolicyDecisionPresent"
                      :details="details"
                      :project="project"
                      :responsible="responsible"
                      :is-deprecated="isDeprecated"
                      :is-unmatched="false"
                      @close="close"
                      @openReviewRemarkDialog="openReviewRemarkDialog"
                      @sendReviewMail="sendReviewMail"
                      @openLicenseRuleDialog="openLicenseRuleDialog"
                      @openPolicyDecisionDialog="openPolicyDecisionDialog" />
                    <PolicyStatusTableRow
                      v-for="(item, index) in details.UnmatchedLicenses"
                      :key="`unmatched-${index}-${item.referenced}`"
                      :item="item"
                      :isPolicyDecisionPresent="isPolicyDecisionPresent"
                      :details="details"
                      :project="project"
                      :responsible="responsible"
                      :is-deprecated="isDeprecated"
                      :is-unmatched="true"
                      @close="close"
                      @openReviewRemarkDialog="openReviewRemarkDialog"
                      @sendReviewMail="sendReviewMail"
                      @openLicenseRuleDialog="openLicenseRuleDialog"
                      @openPolicyDecisionDialog="openPolicyDecisionDialog" />
                    <PolicyStatusTableRow
                      v-for="(item, index) in details.PolicyStatus.filter((p) => p.type !== 'deny')"
                      :key="`policy-status-other-${index}`"
                      :item="item"
                      :policyDecisionApplied="findPolicyDecisionApplied(item)"
                      :isPolicyDecisionPresent="isPolicyDecisionPresent"
                      :details="details"
                      :project="project"
                      :responsible="responsible"
                      :is-deprecated="isDeprecated"
                      :is-unmatched="false"
                      @close="close"
                      @openReviewRemarkDialog="openReviewRemarkDialog"
                      @sendReviewMail="sendReviewMail"
                      @openLicenseRuleDialog="openLicenseRuleDialog"
                      @openPolicyDecisionDialog="openPolicyDecisionDialog" />
                  </tbody>
                </template>
              </v-table>
              <div v-if="details.Problems && details.Problems.length > 0" class="d-flex flex-column align-start">
                <div>
                  <v-icon class="pl-2" color="warning" size="small">mdi mdi-alert</v-icon>
                  <span class="text-caption pl-2 pr-3">{{ t('COMP_PROBLEM_FOUND') }}</span>
                </div>
                <span class="text-caption pl-4" v-for="(p, i) in details.Problems" :key="`0-problem-${p}-${i}`">
                  - {{ t(p) }}
                </span>
              </div>
              <span v-if="details.ContainsOr" class="text-caption pl-2 pr-3">{{ t('POLICY_RULES_DISCLAIMER') }}</span>
            </Stack>
          </v-tabs-window-item>

          <v-tabs-window-item class="dialogContent" value="commonLicenses">
            <div class="h-[344px] overflow-x-auto">
              <Stack>
                <div v-if="noLicenses">
                  {{ t(noLicenses ? 'INFO_NO_LICENSES_FOUND' : 'IDENTIFIED_LICENSES_ON_COMPONENT') }}
                </div>
                <v-expansion-panels focusable>
                  <v-expansion-panel
                    v-for="(item, i) in details.KnownLicenses"
                    :key="`known-${i}-${item.License.licenseId}`">
                    <v-expansion-panel-title>
                      <DInternalLink
                        :text="item.License.name + ' (' + item.ReferencedName + ')'"
                        :url="'/#/dashboard/licenses/' + item.ReferencedName"
                        class=""
                        v-if="project && project.accessRights && project.accessRights.isInternal"></DInternalLink>
                      <span v-else> {{ item.License.name }}</span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text class="mt-2">
                      <DCopyClipboardButton
                        class="d-flex justify-end"
                        :table-button="true"
                        :hint="t('TT_COPY_LICENSE_TEXT')"
                        :content="item.License.text" />
                      <div class="licenseText text-caption" v-text="item.License.text" />
                    </v-expansion-panel-text>
                  </v-expansion-panel>
                </v-expansion-panels>
                <div
                  class="text-body-2 d-secondary-text pa-4 headline-unknown-license"
                  v-if="details.UnknownLicenses && details.UnknownLicenses.length > 0">
                  {{ t('UNIDENTIFIED_LICENSES_ON_COMPONENT') }}
                </div>
                <div v-for="(item, i) in details.UnknownLicenses" :key="`unknown-${i}-${item}`">
                  <template v-if="isInExtractedLicenses(item)">
                    <v-expansion-panels focusable>
                      <v-expansion-panel :key="`extracted-${i}-${item}`">
                        <v-expansion-panel-title>
                          {{ item }}
                        </v-expansion-panel-title>
                        <v-expansion-panel-text class="mt-2">
                          <div class="d-flex justify-end">
                            <DIconButton
                              icon="mdi-shield-search"
                              :hint="t('TT_search_license_text')"
                              @clicked="searchForSimilarLicense(getExtractedLicenseById(item).ExtractedText)" />
                            <DCopyClipboardButton
                              :table-button="true"
                              :hint="t('TT_COPY_LICENSE_TEXT')"
                              :content="getExtractedLicenseById(item).ExtractedText" />
                          </div>
                          <div class="licenseText text-caption" v-text="getExtractedLicenseById(item).ExtractedText" />
                        </v-expansion-panel-text>
                      </v-expansion-panel>
                    </v-expansion-panels>
                  </template>
                  <template v-else>
                    <div class="unknown-license pa-3 pl-4">
                      {{ item }}
                    </div>
                  </template>
                </div>
                <div class="d-flex flex-column align-start pa-4" v-if="details.Problems && details.Problems.length > 0">
                  <div>
                    <v-icon class="pl-2" color="warning" size="small">mdi mdi-alert</v-icon>
                    <span class="d-subtitle-2 pl-2 pr-3">{{ t('COMP_PROBLEM_FOUND') }}</span>
                  </div>
                  <span class="d-subtitle-2 pl-4" v-for="(p, i) in details.Problems" :key="`1-problem-${p}-${i}`">
                    - {{ t(p) }}
                  </span>
                </div>
              </Stack>
            </div>
          </v-tabs-window-item>

          <v-tabs-window-item class="dialogContent" value="sbomLicenses">
            <div class="h-[344px] overflow-x-auto">
              <Stack class="pa-4">
                <div>{{ t('INFO_NO_LICENSES_FOUND') }}</div>
                <div
                  class="text-body-2 pa-4 headline-unknown-license"
                  v-if="details.IdentifiedViaAlias && details.IdentifiedViaAlias.length > 0">
                  {{ t('IDENTIFIED_LICENSES_ON_SBOM_VIA_ALIAS') }}
                </div>
                <v-expansion-panels focusable>
                  <v-expansion-panel
                    v-for="(item, i) in details.IdentifiedViaAlias"
                    :key="`extracted-on-sbom-${i}-${item.License.LicenseId}`">
                    <v-expansion-panel-title>
                      <DInternalLink
                        :text="`${item.License.LicenseId} (${item.AliasTargetId})`"
                        :url="'/#/dashboard/licenses/' + item.AliasTargetId"
                        class=""
                        v-if="project && project.accessRights && project.accessRights.isInternal"></DInternalLink>
                      <span v-else>{{ `${item.License.LicenseId} (${item.AliasTargetId})` }}</span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text class="mt-2">
                      <div class="d-flex justify-end">
                        <DIconButton
                          icon="mdi-shield-search"
                          :hint="t('TT_search_license_text')"
                          @clicked="searchForSimilarLicense(item.License.ExtractedText)" />
                        <DCopyClipboardButton
                          :table-button="true"
                          :hint="t('TT_COPY_LICENSE_TEXT')"
                          :content="item.License.ExtractedText" />
                      </div>
                      <div class="licenseText text-caption" v-text="item.License.ExtractedText" />
                    </v-expansion-panel-text>
                  </v-expansion-panel>
                </v-expansion-panels>
                <div
                  class="text-body-2 pa-4 headline-unknown-license"
                  v-if="details.ExtractedLicenses && details.ExtractedLicenses.length > 0">
                  {{ t('UNIDENTIFIED_LICENSES_ON_SBOM') }}
                </div>
                <v-expansion-panels focusable>
                  <v-expansion-panel
                    v-for="(item, i) in details.ExtractedLicenses"
                    :key="`extracted-on-sbom-${i}-${item.LicenseId}`">
                    <v-expansion-panel-title>
                      {{ item.LicenseId }}
                    </v-expansion-panel-title>
                    <v-expansion-panel-text class="mt-2">
                      <div class="d-flex justify-end">
                        <DIconButton
                          icon="mdi-shield-search"
                          :hint="t('TT_search_license_text')"
                          @clicked="searchForSimilarLicense(item.ExtractedText)" />
                        <DCopyClipboardButton
                          :table-button="true"
                          :hint="t('TT_COPY_LICENSE_TEXT')"
                          :content="item.ExtractedText" />
                      </div>
                      <div class="licenseText text-caption" v-text="item.ExtractedText" />
                    </v-expansion-panel-text>
                  </v-expansion-panel>
                </v-expansion-panels>
              </Stack>
            </div>
          </v-tabs-window-item>

          <v-tabs-window-item class="dialogContent" value="attributes">
            <v-data-table
              height="342"
              density="compact"
              class="striped-table custom-data-table"
              hide-default-footer
              :items-per-page="100000"
              fixed-header
              :headers="headers"
              :items="details.Attributes">
              <template #[`item.Key`]="{item}">
                <span>{{ t('ATTR_' + item.Key.toUpperCase()) }}</span>
              </template>

              <template #[`item.Value`]="{item}">
                <DExternalLink v-if="item?.url" :url="item.Value" :text="item.Value" />
                <span v-else>{{ item.Value }}</span>
              </template>

              <template #[`item.action`]="{item}">
                <div class="opacity-40 hover:opacity-100">
                  <DCopyClipboardButton :hint="t('TT_COPY_TO_CLIPBOARD')" :content="item.Value" />
                </div>
              </template>
            </v-data-table>
          </v-tabs-window-item>

          <v-tabs-window-item class="dialogContent" value="additional_info">
            <div class="h-[344px] overflow-y-auto pa-4">
              <v-progress-linear v-if="loadingRemarks" indeterminate color="primary"></v-progress-linear>
              <div v-else-if="reviewRemarks.length > 0">
                <div class="text-h6 mb-2">{{ t('TAB_REVIEW_REMARKS') }}</div>
                <v-card
                  v-for="remark in sortedReviewRemarks"
                  :key="remark.key"
                  variant="outlined"
                  class="pa-3 mb-3 cursor-pointer"
                  @click="viewRemark(remark)">
                  <Stack class="gap-2">
                    <Stack direction="row" justify="between" align="center">
                      <div class="d-flex align-center flex-wrap gap-2">
                        <v-icon :color="getIconColorReviewRemarkLevel(remark.level)" size="small">
                          {{ getIconReviewRemarkLevel(remark.level) }}
                        </v-icon>
                        <span class="font-weight-bold">{{ remark.title }}</span>
                        <v-chip size="x-small" :color="getIconColorReviewRemarkLevel(remark.level)" variant="flat">
                          {{ t('REMARK_LEVEL_' + remark.level) }}
                        </v-chip>
                        <v-chip
                          size="small"
                          :color="
                            remark.status === 'OPEN' ? 'error' : remark.status === 'IN_PROGRESS' ? 'warning' : 'success'
                          ">
                          {{ t('REMARK_STATUS_' + remark.status) }}
                        </v-chip>
                      </div>
                      <Stack direction="row" align="center" class="gap-1">
                        <DIconButton
                          v-if="isOpen(remark) || isInProgress(remark)"
                          icon="mdi-cancel"
                          :hint="t('TT_cancel_review_remark')"
                          @clicked="openCancelRemarkDialog(remark)"
                          :disabled="isDeprecated" />
                      </Stack>
                    </Stack>
                    <Stack direction="row" justify="between" align="start" class="text-body-2">
                      <span class="cursor-pointer">{{ remark.description }}</span>
                      <Stack align="end" class="text-caption ml-4">
                        <span>{{ remark.author }}</span>
                        <span class="text-grey">{{ new Date(remark.created).toLocaleDateString() }}</span>
                      </Stack>
                    </Stack>
                    <div v-if="remark.licenses?.length > 0" class="text-caption">
                      <strong>{{ t('LICENSES') }}:</strong> {{ remark.licenses.map((l) => l.licenseId).join(', ') }}
                    </div>
                  </Stack>
                </v-card>
              </div>
              <Stack v-else align="center" justify="center" class="h-[300px]">
                <v-icon size="64" color="grey-lighten-1">mdi-information-outline</v-icon>
                <span class="text-grey">{{ t('NO_ADDITIONAL_INFO_FOR_COMPONENT') }}</span>
              </Stack>
            </div>
          </v-tabs-window-item>

          <v-tabs-window-item class="dialogContent" value="raw">
            <div class="h-[344px] overflow-x-auto">
              <Stack class="pa-4">
                <Stack direction="row" class="sticky top-0 z-1 align-center bg-[rgba(var(--v-theme-surface))]">
                  <h3 class="text-h6">
                    {{ t('CAPTION_SCHEMA_DETAILS') }}
                  </h3>
                  <DCopyClipboardButton
                    :hint="t('TT_COPY_RAW_DETAILS_CONTENT')"
                    :tableButton="true"
                    :content="JSON.stringify(details.RawInfo, null, 2)" />
                </Stack>
                <JsonViewer3
                  class="text-caption"
                  :value="details.RawInfo"
                  :expand-depth="2"
                  aria-expanded="true"
                  theme="jv-dark"
                  sort />
              </Stack>
            </div>
          </v-tabs-window-item>
        </v-tabs-window>
      </v-card>
    </DialogLayout>
  </v-dialog>

  <ReviewRemarkDialog ref="reviewRemarkDialog" @reload="reloadReviewRemarks"></ReviewRemarkDialog>
  <LicenseRuleDialog ref="licenseRuleDialog" @reload="closeAndReload"></LicenseRuleDialog>
  <PolicyDecisionDialog
    ref="policyDecisionDialog"
    @reload="closeAndReload"
    @triggerBulk="closeAndTriggerBulk"></PolicyDecisionDialog>
  <ReviewRemarksDetailsDialog
    ref="viewRemarkDialog"
    :project-uuid="project?._key || ''"
    :version-uuid="projectVersionId"
    @reload="reloadReviewRemarks"
    @close-remark="openCloseRemarkDialog"></ReviewRemarksDetailsDialog>
  <ConfirmationDialog v-model:showDialog="closeVisible" :config="confirmCloseConfig" @confirm="doCloseRemark">
  </ConfirmationDialog>
  <ConfirmationDialog v-model:showDialog="cancelVisible" :config="confirmCancelConfig" @confirm="doCancelRemark">
  </ConfirmationDialog>
  <ConfirmationDialog v-model:showDialog="inProgressVisible" :config="confirmInProgressConfig" @confirm="doMarkInProgress">
  </ConfirmationDialog>
</template>
