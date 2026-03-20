<script setup lang="ts">
import {useForm} from '@shared/composables/useForm';
import {computed} from 'vue';
import Tooltip from './Tooltip.vue';

const props = defineProps<{
  modelValue: string;
  readonly?: boolean;
  label?: string;
  required?: boolean;
  help?: string;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();

const {rules} = useForm();
const value = computed({
  get: () => props.modelValue,
  set: (val: string) => emit('update:modelValue', val),
});
</script>

<template>
  <v-textarea
    class="group"
    autocomplete="off"
    :readonly="readonly ?? false"
    :variant="readonly ? 'solo' : 'outlined'"
    v-model="value"
    :label="label"
    :rules="required ? [rules.required] : []"
    :class="required && !readonly ? `required ${$attrs.class}` : $attrs.class"
    :title="value"
    :rows="1"
    v-bind="$attrs"
    persistent-placeholder
    hide-details="auto">
    <template v-if="help" #append-inner>
      <Tooltip :text="help" as-parent>
        <v-icon
          icon="mdi-help-circle-outline"
          class="cursor-help text-gray-400 opacity-0 group-focus-within:opacity-100 group-hover:opacity-100 transition-opacity duration-250"></v-icon>
      </Tooltip>
    </template>
  </v-textarea>
</template>
