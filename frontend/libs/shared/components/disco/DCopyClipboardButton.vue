<script lang="ts" setup>
import {useClipboard} from '@shared/utils/clipboard';
import DIconButton from './DIconButton.vue';

const props = defineProps<{
  icon?: string;
  hint: string;
  content: string;
}>();

const {copyToClipboard} = useClipboard();

const getIconOrDefault = () => {
  return props.icon || 'mdi-content-copy';
};

const onClick = () => {
  copyToClipboard(props.content);
};
</script>

<template>
  <div class="d-inline">
    <DIconButton :icon="getIconOrDefault()" :hint="hint" @clicked="onClick" />
  </div>
</template>
