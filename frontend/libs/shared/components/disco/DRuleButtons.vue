<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import Icons from '@disclosure-portal/constants/icons';
import {PolicyState} from '@disclosure-portal/model/PolicyRule';
import {
  getIconColorForPolicyType,
  getIconColorForPolicyTypeHighlighted,
  getIconForPolicyType,
  openUrl,
  policyStateToTranslationKey,
} from '@disclosure-portal/utils/View';
import {IRuleBtnCallbacks} from '@shared/components/disco/interfaces';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {getCurrentInstance, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

interface Props {
  policies: PolicyState[];
  callbacks: IRuleBtnCallbacks;
  minWidth?: string;
  noClickable?: boolean;
  forceClickable?: boolean;
}

const props = defineProps<Props>();

const {t} = useI18n();
const router = useRouter();

const selectedFilterPolicyTypes = ref<PolicyState[]>([]);

const isSelected = (policy: PolicyState) => {
  if (!policy || policy.length < 1) {
    return selectedFilterPolicyTypes.value.length === 0;
  }
  return selectedFilterPolicyTypes.value.includes(policy);
};

const getBtnIconColorForPolicyFilterBtn = (policy: PolicyState) => {
  if (!props.forceClickable && props.callbacks.getCountForPolicyFilterBtn(policy) === 0) {
    return 'disabledColor';
  }
  return isSelected(policy) ? getIconColorForPolicyTypeHighlighted(policy) : getIconColorForPolicyType(policy);
};

const getIconPolicyType = (policy: PolicyState) => {
  const icon = getIconForPolicyType(policy);
  if (icon === Icons.WARNING && props.callbacks.getCountForPolicyFilterBtn(policy) === 0) {
    return 'mdi-alert';
  }
  return icon;
};

const getStyle = (policy: PolicyState) => {
  if (!props.forceClickable && props.callbacks.getCountForPolicyFilterBtn(policy) === 0) {
    return 'color: var(--v-disabledColor-base) !important;';
  }
  return '';
};

const getTextKeyForPolicyFilterBtn = (policy: PolicyState): string => {
  switch (policy) {
    case PolicyState.NOT_SET:
      return 'COMPONENTS';
    case PolicyState.QUESTIONED:
    case PolicyState.NOASSERTION:
    case PolicyState.WARN:
    case PolicyState.ALLOW:
    case PolicyState.DENY:
      return policyStateToTranslationKey(policy);
    default:
      return 'unknown_policy';
  }
};

const handlePolicySelect = (policy: PolicyState) => {
  openUrl(props.callbacks.getUrlToComponents(policy) as string, router, () => {});
  if (!props.forceClickable && props.callbacks.getCountForPolicyFilterBtn(policy) === 0) {
    return;
  }
  selectedFilterPolicyTypes.value = [policy];
  props.callbacks.handlePolicySelect(policy, selectedFilterPolicyTypes.value);
};

onMounted(() => {
  if (typeof props.callbacks.getUrlToComponents !== 'function') {
    console.error('getUrlToComponents is not a function');
    return;
  }
  const initPolicy = props.callbacks.getInitSelectedPolicy();
  selectedFilterPolicyTypes.value = initPolicy ? [initPolicy] : [];
  const instance = getCurrentInstance();
  if (instance) {
    props.callbacks.setRuleButtons(instance.proxy as any);
  }
});
</script>

<template>
  <span v-if="policies" class="d-flex" data-testid="ruleButtons">
    <v-tooltip
      :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
      location="bottom"
      v-for="(policy, i) in policies"
      :key="i"
      content-class="dpTooltip">
      <template v-slot:activator>
        <v-btn
          :variant="isSelected(policy) ? 'tonal' : 'text'"
          xcolor="ruleButton"
          size="small"
          v-if="policy"
          :class="[
            'ma-2 text-none card-border mx-2 my-2',
            {active: isSelected(policy)},
            isSelected(policy) ? callbacks.getActiveClassForPolicyFilterBtn(policy) : '',
          ]"
          :plain="!isSelected(policy)"
          :style="getStyle(policy)"
          @click.stop="handlePolicySelect(policy)"
          :min-width="minWidth"
          :disabled="!forceClickable && callbacks.getCountForPolicyFilterBtn(policy) === 0">
          <v-icon
            :color="getBtnIconColorForPolicyFilterBtn(policy)"
            :style="getStyle(policy)"
            :icon="getIconPolicyType(policy)"
            class="mr-1" />
          <span v-if="policy">
            {{ callbacks.getCountForPolicyFilterBtn(policy) }}
            {{ t(getTextKeyForPolicyFilterBtn(policy)) }}
          </span>
        </v-btn>
      </template>
      <span v-if="policy">{{ t(callbacks.getToolTipKeyForPolicyFilterBtn(policy)) }}</span>
    </v-tooltip>
  </span>
</template>
