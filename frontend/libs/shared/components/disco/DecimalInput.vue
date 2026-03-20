<template>
  <v-text-field
    variant="outlined"
    :model-value="displayValue"
    @input="onInput"
    @blur="onBlur"
    v-bind="$attrs"
    :readonly="readonly"
    hide-details
    return-object
  ></v-text-field>
</template>

<script setup lang="ts">
import {computed, ref, watch} from 'vue';

const props = defineProps({
  modelValue: {
    type: [Number, String],
    default: '0',
  },
  reset: {
    type: Boolean,
    default: false,
  },
  readonly: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:modelValue']);

const internalValue = ref(props.modelValue);
const tempValue = ref('');

watch(
  () => props.modelValue,
  (newValue) => {
    internalValue.value = newValue;
    tempValue.value = formatNumber(newValue);
  },
);

const formatNumber = (value) => {
  if (value === '' || value === null || isNaN(value)) return '';

  const parts = Number(value).toFixed(2).split('.');
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, '.');
  return parts.join(',');
};

const parseNumber = (value) => {
  if (!value) return '';
  const cleanValue = value.replace(/\./g, '').replace(',', '.');
  return parseFloat(cleanValue).toFixed(2);
};

const displayValue = computed(() => formatNumber(internalValue.value));

const onInput = (event) => {
  tempValue.value = event.target.value;
};

const onBlur = () => {
  const numericValue = parseNumber(tempValue.value) || '0.00';
  internalValue.value = numericValue;
  emit('update:modelValue', numericValue);
};

watch(
  () => props.reset,
  (newValue) => {
    if (newValue) {
      internalValue.value = '0.00';
      tempValue.value = '';
      emit('update:modelValue', '0.00');
    }
  },
);
</script>
