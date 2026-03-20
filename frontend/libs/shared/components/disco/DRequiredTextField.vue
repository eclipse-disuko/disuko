<script setup>
import {computed} from 'vue';

const fieldProps = defineProps(['label', 'rules']);

const computedLabel = computed(() => {
  return fieldProps.rules && fieldProps.rules.includes(rules.required) ? `${fieldProps.label} *` : fieldProps.label;
});
</script>
<template>
  <v-text-field autocomplete="off" v-bind="fieldProps" :label="computedLabel"></v-text-field>
</template>
