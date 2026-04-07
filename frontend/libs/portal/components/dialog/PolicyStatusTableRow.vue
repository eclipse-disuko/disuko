<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import Icons from '@disclosure-portal/constants/icons';
import {PolicyDecisionSlim} from '@disclosure-portal/model/PolicyDecision';
import {ComponentDetails, UnmatchedLicense} from '@disclosure-portal/model/Project';
import {PolicyRuleStatus} from '@disclosure-portal/model/VersionDetails';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import {
  getIconColorForPolicyType,
  getIconForPolicyType,
  policyStateToTranslationKey,
} from '@disclosure-portal/utils/View';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

// Props
interface Props {
  item: PolicyRuleStatus | UnmatchedLicense;
  policyDecisionApplied?: PolicyDecisionSlim | null;
  isPolicyDecisionPresent: boolean;
  details: ComponentDetails;
  project: any;
  responsible: boolean;
  isDeprecated: boolean;
  isUnmatched?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  isUnmatched: false,
});

// Emits
const emit = defineEmits<{
  openReviewRemarkDialog: [licenseId: string];
  sendReviewMail: [item: UnmatchedLicense];
  openLicenseRuleDialog: [licenseId: string];
  openPolicyDecisionDialog: [policy: PolicyRuleStatus | null];
  close: [];
}>();

// Composables
const {t} = useI18n();

// Functions
const formatText = (text: string): string => {
  text = escapeHtml(text);
  if (text.includes(' AND ') || text.includes(' OR ')) {
    return text
      .replace(/ AND /g, ' <strong class="db-highlight">AND</strong> ')
      .replace(/ OR /g, ' <strong class="db-highlight">OR</strong> ');
  }
  return text;
};

const findLicenseFamily = (id: string) => {
  const lic = props.details.KnownLicenses?.find((l) => l.License.licenseId === id);
  if (!lic?.License.meta.family) {
    return t('LIC_FAMILY_UNKNOWN');
  }
  return t(`LIC_FAMILY_${lic?.License.meta.family?.toUpperCase().replace(' ', '_')}`) || '';
};

const helpText = (status: PolicyRuleStatus) => {
  return t(`HELP_TT_${status.type.toUpperCase()}`) || '';
};

const getLicenseId = () => {
  if (props.isUnmatched) {
    const unmatchedItem = props.item as UnmatchedLicense;
    return unmatchedItem.known && unmatchedItem.referenced ? unmatchedItem.referenced : unmatchedItem.orig;
  }
  return (props.item as PolicyRuleStatus).licenseMatched;
};

const getLicenseDisplayText = () => {
  if (props.isUnmatched) {
    const unmatchedItem = props.item as UnmatchedLicense;
    return unmatchedItem.orig;
  }
  return (props.item as PolicyRuleStatus).licenseMatched;
};

const getPolicyType = () => {
  if (props.isUnmatched) {
    return 'noassertion';
  }
  return (props.item as PolicyRuleStatus).type;
};

const getPolicyName = () => {
  if (props.isUnmatched) {
    return '';
  }
  return (props.item as PolicyRuleStatus).name;
};

const getPolicyKey = () => {
  if (props.isUnmatched) {
    return '';
  }
  return (props.item as PolicyRuleStatus).key;
};

const getPolicy = () => {
  if (props.isUnmatched) {
    return null;
  }
  return props.item as PolicyRuleStatus;
};

const shouldShowLicenseRuleInfo = () => {
  const licenseId = getLicenseId();
  return props.details.LicenseRuleApplied && licenseId === props.details.LicenseRuleApplied.licenseDecisionId;
};

const isPolicyDecisionMadeAndFound = () => {
  if (props.isUnmatched) {
    return false;
  }
  return (props.item as PolicyRuleStatus).isDecisionMade;
};

const canMakeWarnedPolicyDecision = () => {
  if (props.isUnmatched) {
    return false;
  }
  return (props.item as PolicyRuleStatus).canMakeWarnedDecision;
};

const canMakeDeniedPolicyDecision = () => {
  if (props.isUnmatched) {
    return false;
  }
  return (props.item as PolicyRuleStatus).canMakeDeniedDecision;
};

const shouldShowInternalLink = () => {
  if (props.isUnmatched) {
    const unmatchedItem = props.item as UnmatchedLicense;
    return props.project?.accessRights?.isInternal && unmatchedItem.known;
  }
  return props.project?.accessRights?.isInternal;
};

const getInternalLinkText = () => {
  if (props.isUnmatched) {
    const unmatchedItem = props.item as UnmatchedLicense;
    return unmatchedItem.referenced;
  }
  return (props.item as PolicyRuleStatus).licenseMatched;
};

const getPolicyRuleHref = () => {
  const key = getPolicyKey();
  if (props.isUnmatched || !key) {
    return '';
  }
  return `#/dashboard/policyrules/${key}`;
};

const licenseChoiceHint = computed(() =>
  props.details.ChoiceDeniedReason ? t('TT_' + props.details.ChoiceDeniedReason) : t('TT_license_rule'),
);

const isLicenseChoiceDisabled = computed(() => !!props.details.ChoiceDeniedReason);

const isWarnedPolicyDecisionDisabled = computed(() => !!props.details.PolicyDecisionDeniedReason);
const warnedPolicyDecisionTooltip = computed(() =>
  isWarnedPolicyDecisionDisabled.value
    ? t('TT_' + props.details.PolicyDecisionDeniedReason)
    : t('TT_warned_policy_decision'),
);

const isDeniedPolicyDecisionDisabled = computed(
  () =>
    props.details.PolicyDecisionDeniedReason === 'DECISION_DENIED_COMPONENT_VERSION_NOT_SET' ||
    !!getPolicy()?.deniedDecisionDeniedReason,
);
const deniedPolicyDecisionTooltip = computed(() => {
  if (!isDeniedPolicyDecisionDisabled.value) {
    return t('TT_denied_policy_decision');
  }

  if (props.details.PolicyDecisionDeniedReason === 'DECISION_DENIED_COMPONENT_VERSION_NOT_SET') {
    return t('TT_' + props.details.PolicyDecisionDeniedReason);
  }

  if (getPolicy()?.deniedDecisionDeniedReason) {
    return t('TT_' + getPolicy()?.deniedDecisionDeniedReason);
  }

  return t('TT_denied_policy_decision');
});

const getActionButtons = computed((): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: Icons.QUESTIONED,
      hint: props.isUnmatched ? t('HELP_TT_UNASSERTED') : helpText(props.item as PolicyRuleStatus),
      event: 'help',
      show: true,
    },
    {
      icon: 'mdi-message-plus-outline',
      hint: t('UM_DIALOG_TITLE_NEW_REVIEW_REMARK'),
      event: 'reviewRemark',
      disabled: props.isDeprecated,
      show: true,
    },
    {
      icon: 'mdi-email-fast-outline',
      hint: props.isUnmatched ? t('TT_request_review') : t('TT_not_request_review'),
      event: 'sendMail',
      disabled: !props.isUnmatched,
      show: true,
    },
    {
      icon: 'mdi-text-box-edit-outline',
      hint: licenseChoiceHint.value,
      event: 'licenseChoice',
      disabled: isLicenseChoiceDisabled.value,
      show: props.details.CanChooseLicense,
    },
    {
      icon: 'mdi-checkbox-marked-circle-plus-outline',
      hint: warnedPolicyDecisionTooltip.value,
      event: 'warnedPolicyDecision',
      disabled: isWarnedPolicyDecisionDisabled.value,
      show: canMakeWarnedPolicyDecision(),
    },
    {
      icon: 'mdi-checkbox-marked-circle-plus-outline',
      hint: deniedPolicyDecisionTooltip.value,
      event: 'deniedPolicyDecision',
      disabled: isDeniedPolicyDecisionDisabled.value,
      show: canMakeDeniedPolicyDecision(),
      color: 'orange',
    },
  ];
});
</script>

<template>
  <tr>
    <td>
      <TableActionButtons
        variant="minimal"
        :buttons="getActionButtons"
        @help="() => {}"
        @reviewRemark="emit('openReviewRemarkDialog', getLicenseId())"
        @sendMail="emit('sendReviewMail', item as UnmatchedLicense)"
        @licenseChoice="emit('openLicenseRuleDialog', getLicenseId())"
        @warnedPolicyDecision="emit('openPolicyDecisionDialog', getPolicy())"
        @deniedPolicyDecision="emit('openPolicyDecisionDialog', getPolicy())" />
    </td>
    <td>
      <span v-if="policyDecisionApplied" class="d-inline-block w-[18px] h-[16px]">
        <v-icon size="small" :color="policyDecisionApplied.previewMode ? 'grey' : ''">
          {{ policyDecisionApplied.previewMode ? 'mdi-progress-alert' : 'mdi-information-outline' }}
        </v-icon>
        <Tooltip>
          <PolicyDecisionItem
            v-if="isPolicyDecisionMadeAndFound()"
            class="ml-4"
            :previewMode="policyDecisionApplied?.previewMode"
            :decision="policyDecisionApplied"
            :license-id="policyDecisionApplied.licenseId" />
          <template v-if="policyDecisionApplied">
            {{
              ` (${t('TT_POLICY_DECISION_BY_AT', {
                creator: policyDecisionApplied.creator,
                created: formatDateAndTime(policyDecisionApplied.created),
              })})`
            }}
          </template>
        </Tooltip>
      </span>
      <span v-else-if="isPolicyDecisionPresent">
        <v-icon>mdi-blank</v-icon>
      </span>
      &nbsp;
      <span>
        <v-icon :color="getIconColorForPolicyType(getPolicyType())">
          {{ getIconForPolicyType(getPolicyType()) }}
        </v-icon>
        <tooltip :text="policyStateToTranslationKey(getPolicyType())"></tooltip>
      </span>
    </td>
    <td>
      <span v-if="shouldShowLicenseRuleInfo()">
        <v-icon size="small" :color="details?.LicenseRuleApplied?.previewMode ? 'grey' : ''">
          {{ details?.LicenseRuleApplied?.previewMode ? 'mdi-progress-alert' : 'mdi-information-outline' }}
        </v-icon>
        <tooltip>
          <span class="text-subtitle-1">{{
            details?.LicenseRuleApplied?.previewMode
              ? t('TT_LICENSE_RULE_APPLIED_PREVIEW')
              : t('TT_LICENSE_RULE_APPLIED')
          }}</span>
          <br />
          <span class="d-text d-secondary-text">{{ t('TT_LICENSE_RULE_EXPRESSION') }}</span>
          <br />
          <span
            class="d-text d-secondary-text"
            v-html="formatText(details?.LicenseRuleApplied?.licenseExpression || '')"></span>
          <br />
          <span class="d-text d-secondary-text">{{
            t('TT_LICENSE_RULE_DECISION', {
              decision: details?.LicenseRuleApplied?.licenseDecisionName || '',
              decisionId: details?.LicenseRuleApplied?.licenseDecisionId || '',
            })
          }}</span>
          <br />
          <span class="d-text d-secondary-text">{{
            t('TT_LICENSE_RULE_BY_AT', {
              creator: details?.LicenseRuleApplied?.creator || '',
              created: formatDateAndTime(details?.LicenseRuleApplied?.created || ''),
            })
          }}</span>
        </tooltip>
        &nbsp;
      </span>
      <DInternalLink
        v-if="shouldShowInternalLink()"
        :text="getInternalLinkText()"
        :url="'/#/dashboard/licenses/' + getInternalLinkText()"
        class="" />
      <span v-else>{{ getLicenseDisplayText() }}</span>
    </td>
    <td>
      <a
        v-if="!isUnmatched && getPolicyKey()"
        :href="getPolicyRuleHref()"
        class="discotext"
        @click="emit('close')"
        target="_blank">
        {{ getPolicyName() }}
      </a>
    </td>
    <td>{{ isUnmatched ? '' : findLicenseFamily(getLicenseId()) }}</td>
  </tr>
</template>
