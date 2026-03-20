<script setup lang="ts">
import Icons from '@disclosure-portal/constants/icons';
import type {Label} from '@disclosure-portal/model/Label';
import Tooltip from '@shared/components/disco/Tooltip.vue';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

interface Props {
  label: Partial<Label>;
}

const props = defineProps<Props>();

const {t} = useI18n();

const labelIcons: Record<Label['type'], string> = {
  SCHEMA: Icons.SCHEMA,
  POLICY: Icons.POLICY,
  PROJECT: Icons.PROJECT_LABEL,
};

const labelTooltipTitles: Record<Label['type'], string> = {
  SCHEMA: t('TT_schema_label_with_description'),
  POLICY: t('TT_policy_label_with_description'),
  PROJECT: t('TT_project_label_with_description'),
};

const iconName = computed(() => (props.label.type ? labelIcons[props.label.type] : Icons.TAG));
const tooltipTitle = computed(() => (props.label.type ? labelTooltipTitles[props.label.type] : ''));
</script>

<template>
  <span>
    <DLabel css-clazzes="m-0" :labelName="label.name" :iconName="iconName"></DLabel>
    <Tooltip v-if="label.description">
      <span v-if="tooltipTitle">{{ tooltipTitle }}</span>
      <span v-if="label.description">{{ label.description }}</span>
    </Tooltip>
  </span>
</template>
