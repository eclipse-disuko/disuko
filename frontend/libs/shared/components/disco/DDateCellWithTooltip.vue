<template>
  <v-tooltip
    v-if="value && value !== '0001-01-01T00:00:00Z'"
    :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
    location="bottom"
    content-class="dpTooltip">
    <template v-slot:activator="{props}">
      <span v-bind="props">
        {{ !showTime ? formatDate(value) : formatDateAndTime(value) }}
      </span>
    </template>
    <span>{{ formatDateAndTime(value) }}</span>
  </v-tooltip>
</template>

<script setup lang="ts">
import {formatDate, formatDateAndTime} from '@disclosure-portal/utils/Table';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';

defineProps<{
  value: string | Date | undefined;
  showTime?: boolean;
}>();
</script>
