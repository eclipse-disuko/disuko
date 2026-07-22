<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {DecisionType} from '@disclosure-portal/components/dialog/DialogConfigs';
import {ComponentInfo} from '@disclosure-portal/model/VersionDetails';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

interface DecisionOption {
  type: DecisionType;
  visible: boolean;
  canMake: boolean;
  dimmed: boolean;
  icon: string;
  color: string;
  label: string;
  tooltip: string;
}

const props = defineProps<{
  item: ComponentInfo;
}>();

const emit = defineEmits<{openPolicyDecision: [type: DecisionType]}>();

const {t} = useI18n();

const canWarn = computed(() => props.item.policyRuleStatus.some((pr) => pr.canMakeWarnedDecision));
const canDeny = computed(() => props.item.policyRuleStatus.some((pr) => pr.canMakeDeniedDecision));

const hasWarnDecisions = computed(() => props.item.policyDecisionsApplied.some((pd) => pd.policyEvaluated === 'warn'));
const hasDenyDecisions = computed(() => props.item.policyDecisionsApplied.some((pd) => pd.policyEvaluated === 'deny'));

const showWarnColumn = computed(() => canWarn.value || hasWarnDecisions.value);
const showDenyColumn = computed(() => canDeny.value || hasDenyDecisions.value);
const showAppliedOnlyColumn = computed(
  () => !canWarn.value && !canDeny.value && props.item.policyDecisionsApplied.length > 0,
);
const showBothColumns = computed(() => showWarnColumn.value && showDenyColumn.value && !showAppliedOnlyColumn.value);

const activeDecisions = computed(() => props.item.policyDecisionsApplied.filter((pd) => !pd.previewMode));
const previewDecisions = computed(() => props.item.policyDecisionsApplied.filter((pd) => pd.previewMode));
const hasAnyDecisions = computed(() => activeDecisions.value.length > 0 || previewDecisions.value.length > 0);

const isWarnDisabled = computed(() => !!props.item.policyDecisionDeniedReason);
const warnTooltip = computed(() =>
  isWarnDisabled.value ? t('TT_' + props.item.policyDecisionDeniedReason) : t('TT_warned_policy_decision'),
);

const deniableRuleStatuses = computed(() => props.item.policyRuleStatus.filter((p) => p.canMakeDeniedDecision));
const isDenyDisabled = computed(
  () =>
    props.item.policyDecisionDeniedReason === 'DECISION_DENIED_COMPONENT_VERSION_NOT_SET' ||
    (deniableRuleStatuses.value.length > 0 && deniableRuleStatuses.value.every((p) => !!p.deniedDecisionDeniedReason)),
);
const denyTooltip = computed(() => {
  if (!isDenyDisabled.value) {
    return t('TT_denied_policy_decision');
  }

  if (props.item.policyDecisionDeniedReason === 'DECISION_DENIED_COMPONENT_VERSION_NOT_SET') {
    return t('TT_' + props.item.policyDecisionDeniedReason);
  }

  if (
    deniableRuleStatuses.value.length > 0 &&
    deniableRuleStatuses.value.every((p) => !!p.deniedDecisionDeniedReason)
  ) {
    return t(`TT_${props.item.policyRuleStatus[0].deniedDecisionDeniedReason}`);
  }

  return t('TT_denied_policy_decision');
});

const getDecisionIcon = (canMake: boolean) => {
  if (canMake) return 'mdi-checkbox-marked-circle-plus-outline';
  if (activeDecisions.value.length > 0) return 'mdi-information-outline';
  if (previewDecisions.value.length > 0) return 'mdi-progress-alert';
  return '';
};

const getDecisionColor = (canMake: boolean, canMakeColor: string) => {
  if (canMake) return canMakeColor;
  if (activeDecisions.value.length > 0) return '';
  if (previewDecisions.value.length > 0) return 'grey';
  return '';
};

const appliedOnlyIcon = computed(() => getDecisionIcon(false));
const appliedOnlyColor = computed(() => getDecisionColor(false, ''));

const decisionOptions = computed<DecisionOption[]>(() => [
  {
    type: 'warn',
    visible: showWarnColumn.value && !showAppliedOnlyColumn.value,
    canMake: canWarn.value,
    dimmed: canWarn.value && isWarnDisabled.value,
    icon: getDecisionIcon(canWarn.value),
    color: getDecisionColor(canWarn.value, 'primary'),
    label: canWarn.value ? t('WARN') : t('INFO'),
    tooltip: warnTooltip.value,
  },
  {
    type: 'deny',
    visible: showDenyColumn.value && !showAppliedOnlyColumn.value,
    canMake: canDeny.value,
    dimmed: canDeny.value && isDenyDisabled.value,
    icon: getDecisionIcon(canDeny.value),
    color: getDecisionColor(canDeny.value, 'orange'),
    label: canDeny.value ? t('DENY') : t('INFO'),
    tooltip: denyTooltip.value,
  },
]);

const onDecisionClick = (option: DecisionOption) => {
  if (option.canMake && !option.dimmed) {
    emit('openPolicyDecision', option.type);
  }
};
</script>

<template>
  <div class="flex flex-col" :class="showBothColumns ? 'gap-y-0.5 py-0.5' : 'justify-center'">
    <template v-for="option in decisionOptions" :key="option.type">
      <span v-if="option.visible">
        <DCActionButton
          size="small"
          density="compact"
          class="h-auto min-h-7 py-0.5"
          :color="option.color"
          :icon="option.icon"
          :text="option.label"
          :style="option.dimmed ? 'opacity: 0.38;' : ''"
          @click.stop="onDecisionClick(option)">
          <template #tooltip>
            <template v-if="option.canMake">
              {{ option.tooltip }}
              <br />
            </template>

            <template v-if="hasAnyDecisions">
              <v-divider></v-divider>
              {{ t('TT_POLICY_DECISIONS_AVAILABLE') }}

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
          </template>
        </DCActionButton>
      </span>
    </template>

    <span v-if="showAppliedOnlyColumn">
      <DCActionButton
        size="small"
        density="compact"
        class="h-auto min-h-7 px-1 py-0.5"
        :color="appliedOnlyColor"
        :icon="appliedOnlyIcon"
        :text="t('INFO')">
        <template #tooltip>
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
      </DCActionButton>
    </span>
  </div>
</template>
