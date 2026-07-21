<script setup lang="ts">
import type {AugmentedPolicyRule, PolicyRule} from '@cli/models/Sbom';
import {getIconColorForPolicyType, getIconForPolicyType} from '@disclosure-portal/utils/View';
import Tooltip from '@shared/components/disco/Tooltip.vue';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps<{
  policyRules: PolicyRule[];
}>();

const {t} = useI18n();

const truncateText = (text: string | undefined, maxLength: number = 120): string => {
  if (!text) return '';
  return text.length > maxLength ? text.slice(0, maxLength) + '...' : text;
};

const statusPolicyRules = computed<AugmentedPolicyRule[]>(() =>
  props.policyRules.map((p) => ({
    ...p,
    _iconColor: getIconColorForPolicyType(p.type),
    _icon: getIconForPolicyType(p.type),
    _description: truncateText(p.description),
  })),
);

const baseHeaderClass = 'tableHeaderCell';
const policyHeaders = [
  {title: t('STATUS'), key: 'type', width: 130, align: 'center' as const, class: baseHeaderClass},
  {title: t('COL_LICENSE'), key: 'licenseMatched', width: 180, align: 'start' as const, class: baseHeaderClass},
  {title: t('COL_POLICY_NAME'), key: 'name', width: 240, align: 'start' as const, class: baseHeaderClass},
  {title: t('COL_DESCRIPTION'), key: 'description', align: 'start' as const, class: baseHeaderClass},
];
</script>

<template>
  <v-data-table
    height="350"
    :headers="policyHeaders"
    :items="statusPolicyRules"
    item-key="name"
    density="compact"
    :items-per-page="25"
    :footer-props="{'items-per-page-options': [25, 50, 100, -1]}"
    class="striped-table custom-data-table elevation-0 p-2"
    hover>
    <template #item.type="{item}">
      <div class="flex h-full items-center justify-center">
        <v-icon :color="item._iconColor" dense>{{ item._icon }}</v-icon>
        <Tooltip>{{ item.type }}</Tooltip>
      </div>
    </template>
    <template #item.description="{item}">
      <span>{{ item._description }}</span>
      <Tooltip>{{ item.description || '' }}</Tooltip>
    </template>
  </v-data-table>
</template>
