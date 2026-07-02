<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {DecisionType} from '@disclosure-portal/components/dialog/DialogConfigs';
import {ComponentInfo} from '@disclosure-portal/model/VersionDetails';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps<{
  item: ComponentInfo;
}>();

const emit = defineEmits<{openPolicyDecision: [type: DecisionType]}>();

const {t} = useI18n();

const canMakeWarnedPolicyDecision = computed(() => props.item.policyRuleStatus.some((pr) => pr.canMakeWarnedDecision));
const hasWarnedDecisions = computed(() =>
  props.item.policyDecisionsApplied.some((pd) => pd.policyEvaluated === 'warn'),
);
const isWarnedPolicyDecisionInfosAvailable = computed(
  () => canMakeWarnedPolicyDecision.value || hasWarnedDecisions.value,
);

const canMakeDeniedPolicyDecision = computed(() => props.item.policyRuleStatus.some((pr) => pr.canMakeDeniedDecision));
const hasDeniedDecisions = computed(() =>
  props.item.policyDecisionsApplied.some((pd) => pd.policyEvaluated === 'deny'),
);
const isDeniedPolicyDecisionInfosAvailable = computed(
  () => canMakeDeniedPolicyDecision.value || hasDeniedDecisions.value,
);

const isOnlyPolicyDecisionsAppliedPresent = computed(
  () =>
    !canMakeWarnedPolicyDecision.value &&
    !canMakeDeniedPolicyDecision.value &&
    props.item.policyDecisionsApplied.length > 0,
);

const isWarnedPolicyDecisionDisabled = computed(() => !!props.item.policyDecisionDeniedReason);
const warnedPolicyDecisionTooltip = computed(() =>
  isWarnedPolicyDecisionDisabled.value
    ? t('TT_' + props.item.policyDecisionDeniedReason)
    : t('TT_warned_policy_decision'),
);

const isDeniedPolicyDecisionDisabled = computed(
  () =>
    props.item.policyDecisionDeniedReason === 'DECISION_DENIED_COMPONENT_VERSION_NOT_SET' ||
    (props.item.policyRuleStatus.filter((p) => p.canMakeDeniedDecision).length > 0 &&
      props.item.policyRuleStatus.filter((p) => p.canMakeDeniedDecision).every((p) => !!p.deniedDecisionDeniedReason)),
);
const deniedPolicyDecisionTooltip = computed(() => {
  if (!isDeniedPolicyDecisionDisabled.value) {
    return t('TT_denied_policy_decision');
  }

  if (props.item.policyDecisionDeniedReason === 'DECISION_DENIED_COMPONENT_VERSION_NOT_SET') {
    return t('TT_' + props.item.policyDecisionDeniedReason);
  }

  if (
    props.item.policyRuleStatus.filter((p) => p.canMakeDeniedDecision).length > 0 &&
    props.item.policyRuleStatus.filter((p) => p.canMakeDeniedDecision).every((p) => !!p.deniedDecisionDeniedReason)
  ) {
    return t(`TT_${props.item.policyRuleStatus[0].deniedDecisionDeniedReason}`);
  }

  return t('TT_denied_policy_decision');
});

const activeDecisions = computed(() => props.item.policyDecisionsApplied.filter((pd) => !pd.previewMode));
const previewDecisions = computed(() => props.item.policyDecisionsApplied.filter((pd) => pd.previewMode));
const hasAnyDecisions = computed(() => activeDecisions.value.length > 0 || previewDecisions.value.length > 0);

const warnedIcon = computed(() => {
  if (canMakeWarnedPolicyDecision.value) return 'mdi-checkbox-marked-circle-plus-outline';
  if (activeDecisions.value.length > 0) return 'mdi-information-outline';
  if (previewDecisions.value.length > 0) return 'mdi-progress-alert';
  return '';
});

const warnedIconColor = computed(() => {
  if (canMakeWarnedPolicyDecision.value) return 'primary';
  if (activeDecisions.value.length > 0) return '';
  if (previewDecisions.value.length > 0) return 'grey';
  return '';
});

const deniedIcon = computed(() => {
  if (canMakeDeniedPolicyDecision.value) return 'mdi-checkbox-marked-circle-plus-outline';
  if (activeDecisions.value.length > 0) return 'mdi-information-outline';
  if (previewDecisions.value.length > 0) return 'mdi-progress-alert';
  return '';
});

const deniedIconColor = computed(() => {
  if (canMakeDeniedPolicyDecision.value) return 'orange';
  if (activeDecisions.value.length > 0) return '';
  if (previewDecisions.value.length > 0) return 'grey';
  return '';
});

const onlyAppliedDecisionsIcon = computed(() => {
  if (activeDecisions.value.length > 0) return 'mdi-information-outline';
  if (previewDecisions.value.length > 0) return 'mdi-progress-alert';
  return '';
});
const onlyAppliedDecisionsIconColor = computed(() => {
  if (activeDecisions.value.length > 0) return '';
  if (previewDecisions.value.length > 0) return 'grey';
  return '';
});

const handleWarnedClick = () => {
  if (canMakeWarnedPolicyDecision.value && !isWarnedPolicyDecisionDisabled.value) {
    emit('openPolicyDecision', 'warn');
  }
};

const handleDeniedClick = () => {
  if (canMakeDeniedPolicyDecision.value && !isDeniedPolicyDecisionDisabled.value) {
    emit('openPolicyDecision', 'deny');
  }
};
</script>

<template>
  <span v-if="isWarnedPolicyDecisionInfosAvailable && !isOnlyPolicyDecisionsAppliedPresent">
    <span v-if="canMakeWarnedPolicyDecision" @click.stop="handleWarnedClick">
      <v-btn size="small" variant="text" disabled>
      {{ t('no_warned_decision_text') || 'No Warned Decision' }}
    </v-btn>
    </span>
    <v-icon v-else size="small" :color="warnedIconColor">
      {{ warnedIcon }}
    </v-icon>

    <Tooltip>
      <span v-if="canMakeWarnedPolicyDecision">
        {{ warnedPolicyDecisionTooltip }}
        <br />
      </span>

      <template v-if="hasAnyDecisions">
        <v-divider></v-divider>
        {{ t('TT_POLICY_DECISIONS_AVAILABLE') }}
        <br />

        <PolicyDecisionList
          v-if="activeDecisions.length > 0"
          :decisions="activeDecisions"
          :title="t('TT_POLICY_DECISION_APPLIED')"
          icon="mdi-information-outline"
          arrow="→" />

        <PolicyDecisionList
          v-if="previewDecisions.length > 0"
          :decisions="previewDecisions"
          :title="t('TT_POLICY_DECISION_APPLIED_PREVIEW')"
          icon="mdi-progress-alert"
          icon-color="grey"
          arrow="⇢" />
      </template>
    </Tooltip>
    &nbsp;
  </span>


  <span v-if="isDeniedPolicyDecisionInfosAvailable && !isOnlyPolicyDecisionsAppliedPresent">
      <v-btn
        v-if="canMakeDeniedPolicyDecision"
        size="small"
        variant="tonal"
        color="primary"
        @click.stop="handleDeniedClick"
        :disabled="isDeniedPolicyDecisionDisabled"
        class="text-none"
        density="comfortable"
        prepend-icon="mdi-checkbox-marked-circle-plus-outline ">{{ t('Decide')}}
      </v-btn>
      <v-btn
        v-else
        size="small"
        :color="deniedIconColor"
        prepend-icon="mdi-block-helper"
        variant="text">
      {{ t('denied_decision_button_text') || 'Denied Decision' }}
      </v-btn>
    <Tooltip>
      <span v-if="canMakeDeniedPolicyDecision">
        {{ deniedPolicyDecisionTooltip }}
        <br />
      </span>

      <template v-if="hasAnyDecisions">
        <v-divider></v-divider>
        {{ t('TT_POLICY_DECISIONS_AVAILABLE') }}
        <br />

        <PolicyDecisionList
          v-if="activeDecisions.length > 0"
          :decisions="activeDecisions"
          :title="t('TT_POLICY_DECISION_APPLIED')"
          icon="mdi-information-outline"
          arrow="→" />

        <PolicyDecisionList
          v-if="previewDecisions.length > 0"
          :decisions="previewDecisions"
          :title="t('TT_POLICY_DECISION_APPLIED_PREVIEW')"
          icon="mdi-progress-alert"
          icon-color="grey"
          arrow="⇢" />
      </template>
    </Tooltip>
  </span>


  <span v-if="isOnlyPolicyDecisionsAppliedPresent">
    <v-btn
      size="small"
      variant="tonal"
      color="primary"
      @click.stop="handleDeniedClick"
      :disabled="isDeniedPolicyDecisionDisabled"
      class="text-none"
      density="comfortable"
      prepend-icon="mdi-information-outline">
      {{ t('Info') }}
    </v-btn>

    <Tooltip>
      <span v-if="canMakeDeniedPolicyDecision">
        {{ deniedPolicyDecisionTooltip }}
        <br />
      </span>

      <template v-if="hasAnyDecisions">
        <PolicyDecisionList
          v-if="activeDecisions.length > 0"
          :decisions="activeDecisions"
          :title="t('TT_POLICY_DECISION_APPLIED')"
          icon="mdi-information-outline"
          arrow="→" />

        <PolicyDecisionList
          v-if="previewDecisions.length > 0"
          :decisions="previewDecisions"
          :title="t('TT_POLICY_DECISION_APPLIED_PREVIEW')"
          icon="mdi-progress-alert"
          icon-color="grey"
          arrow="⇢" />
      </template>
    </Tooltip>
    <v-icon size="small">mdi-blank</v-icon>
  </span>
</template>
