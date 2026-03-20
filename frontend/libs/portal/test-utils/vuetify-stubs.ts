import {vi} from 'vitest';

/**
 * Common Vuetify component stubs for testing
 * Import this in your tests to avoid repeating stub definitions
 */
export const vuetifyStubs = {
  'v-form': {
    template: '<form><slot /></form>',
    methods: {
      validate: vi.fn(() => Promise.resolve({valid: true})),
      reset: vi.fn(),
      resetValidation: vi.fn(),
    },
  },
  'v-text-field': {
    template:
      '<input type="text" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" :autocomplete="autocomplete" :required="required" />',
    props: [
      'modelValue',
      'label',
      'rules',
      'hideDetails',
      'required',
      'autocomplete',
      'variant',
      'type',
      'disabled',
      'readonly',
      'placeholder',
    ],
    emits: ['update:modelValue'],
  },
  'v-textarea': {
    template:
      '<textarea :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" :rows="rows" :counter="counter"></textarea>',
    props: [
      'modelValue',
      'rows',
      'counter',
      'label',
      'rules',
      'hideDetails',
      'variant',
      'noResize',
      'disabled',
      'readonly',
      'placeholder',
    ],
    emits: ['update:modelValue'],
  },
  'v-checkbox': {
    template:
      '<input type="checkbox" :checked="modelValue" @change="$emit(\'update:modelValue\', $event.target.checked)" />',
    props: ['modelValue', 'label', 'hideDetails', 'outlined', 'color', 'disabled', 'readonly'],
    emits: ['update:modelValue'],
  },
  'v-select': {
    template:
      '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><slot /></select>',
    props: [
      'modelValue',
      'items',
      'itemTitle',
      'itemValue',
      'label',
      'rules',
      'hideDetails',
      'variant',
      'multiple',
      'disabled',
      'readonly',
    ],
    emits: ['update:modelValue'],
  },
  'v-autocomplete': {
    template: '<input type="text" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: [
      'modelValue',
      'items',
      'itemTitle',
      'itemValue',
      'label',
      'rules',
      'hideDetails',
      'variant',
      'multiple',
      'disabled',
      'readonly',
    ],
    emits: ['update:modelValue'],
  },
  'v-btn': {
    template: '<button type="button" :disabled="disabled"><slot /></button>',
    props: ['color', 'variant', 'size', 'disabled', 'loading', 'icon', 'to', 'href'],
  },
  'v-icon': {
    template: '<i><slot /></i>',
    props: ['icon', 'size', 'color'],
  },
  'v-card': {
    template: '<div class="v-card"><slot /></div>',
    props: ['variant', 'color', 'elevation'],
  },
  'v-card-title': {
    template: '<div class="v-card-title"><slot /></div>',
  },
  'v-card-text': {
    template: '<div class="v-card-text"><slot /></div>',
  },
  'v-card-actions': {
    template: '<div class="v-card-actions"><slot /></div>',
  },
  'v-dialog': {
    template: '<div v-if="modelValue" class="v-dialog"><slot /></div>',
    props: ['modelValue', 'width', 'maxWidth', 'persistent'],
    emits: ['update:modelValue'],
  },
  'v-menu': {
    template: '<div class="v-menu"><slot /><slot name="activator" /></div>',
    props: ['modelValue', 'openOnClick', 'closeOnContentClick'],
    emits: ['update:modelValue'],
  },
  'v-tooltip': {
    template: '<div class="v-tooltip"><slot /><slot name="activator" /></div>',
    props: ['text', 'location'],
  },
  'v-divider': {
    template: '<hr class="v-divider" />',
  },
  'v-list': {
    template: '<div class="v-list"><slot /></div>',
  },
  'v-list-item': {
    template: '<div class="v-list-item"><slot /></div>',
    props: ['title', 'subtitle', 'value'],
  },
  'v-list-item-title': {
    template: '<div class="v-list-item-title"><slot /></div>',
  },
  'v-list-item-subtitle': {
    template: '<div class="v-list-item-subtitle"><slot /></div>',
  },
  'v-chip': {
    template: '<span class="v-chip"><slot /></span>',
    props: ['color', 'size', 'variant', 'closable'],
  },
  'v-progress-circular': {
    template: '<div class="v-progress-circular"></div>',
    props: ['indeterminate', 'size', 'color'],
  },
  'v-progress-linear': {
    template: '<div class="v-progress-linear"></div>',
    props: ['indeterminate', 'modelValue', 'color'],
  },
  'v-alert': {
    template: '<div class="v-alert" :type="type"><slot /></div>',
    props: ['type', 'variant', 'closable'],
  },
  'v-radio-group': {
    template: '<div class="v-radio-group"><slot /></div>',
    props: ['modelValue', 'label', 'rules', 'hideDetails'],
    emits: ['update:modelValue'],
  },
  'v-radio': {
    template:
      '<input type="radio" :value="value" :checked="modelValue === value" @change="$emit(\'update:modelValue\', value)" />',
    props: ['modelValue', 'value', 'label'],
    emits: ['update:modelValue'],
  },
  'v-switch': {
    template:
      '<input type="checkbox" :checked="modelValue" @change="$emit(\'update:modelValue\', $event.target.checked)" />',
    props: ['modelValue', 'label', 'color', 'disabled'],
    emits: ['update:modelValue'],
  },
  'v-tabs': {
    template: '<div class="v-tabs"><slot /></div>',
    props: ['modelValue', 'color'],
    emits: ['update:modelValue'],
  },
  'v-tab': {
    template: '<button type="button" class="v-tab"><slot /></button>',
    props: ['value'],
  },
  'v-window': {
    template: '<div class="v-window"><slot /></div>',
    props: ['modelValue'],
    emits: ['update:modelValue'],
  },
  'v-window-item': {
    template: '<div class="v-window-item"><slot /></div>',
    props: ['value'],
  },
  'v-container': {
    template: '<div class="v-container"><slot /></div>',
  },
  'v-row': {
    template: '<div class="v-row"><slot /></div>',
  },
  'v-col': {
    template: '<div class="v-col"><slot /></div>',
    props: ['cols', 'sm', 'md', 'lg', 'xl'],
  },
  'v-spacer': {
    template: '<div class="v-spacer"></div>',
  },
  'v-table': {
    template: '<table class="v-table"><slot /></table>',
  },
  'v-data-table': {
    template: '<div class="v-data-table"><slot /></div>',
    props: ['headers', 'items', 'itemsPerPage', 'page'],
  },
};
