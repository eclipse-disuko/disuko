<script setup lang="ts">
import type {License, PolicyCat, PolicyRule} from '@cli/models/PolicyRule';
import {POLICY_STATUS_DEFS} from '@cli/models/PolicyRule';
import DCopyClipboardButton from '@shared/components/disco/DCopyClipboardButton.vue';
import DialogLayout, {type DialogLayoutConfig} from '@shared/layouts/DialogLayout.vue';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import type {DataTableHeader} from 'vuetify';

const props = defineProps<{
  policyRule: PolicyRule | null;
  show: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'close'): void;
}>();

const {t} = useI18n();
const search = ref('');
const loading = ref(false);
const sortItems = ref([{key: 'name', order: false}]);

const selectedCategory = ref<'all' | PolicyCat>('all');

// Replace STATUS_DEFS usage
function toggleCategory(cat: PolicyCat) {
  selectedCategory.value = selectedCategory.value === cat ? 'all' : cat;
}
const filteredLicenses = computed(() => {
  if (!props.policyRule?.licenses) return [] as License[];
  if (selectedCategory.value === 'all') return props.policyRule.licenses;
  return props.policyRule.licenses.filter(
    (l: License) => (l.policyType || '').toLowerCase() === selectedCategory.value,
  );
});

const statusButtons = computed(() => {
  const licenses: License[] = props.policyRule?.licenses || [];
  const totalLicenses = licenses.length;
  const currentCategory = selectedCategory.value;
  return POLICY_STATUS_DEFS.map((def) => {
    const defKey = def.key;
    const count =
      defKey === 'all' ? totalLicenses : licenses.filter((l) => (l.policyType || '').toLowerCase() === defKey).length;
    const active = currentCategory === defKey;
    const displayLabel = t(def.labelKey);
    return {
      ...def,
      count,
      active,
      ariaLabel: `${count} ${displayLabel}`,
      displayLabel,
    };
  });
});

const statusButtonsWithClasses = computed(() => {
  const borderClassByKey: Record<string, string> = {
    allow: 'allow-border',
    warn: 'warning-border',
    deny: 'deny-border',
    all: 'all-border',
  };
  const twBorderColorByKey: Record<string, string> = {
    allow: 'border-green-500',
    warn: 'border-yellow-500',
    deny: 'border-red-600',
    all: 'border-slate-500',
  };
  const activeBgByKey: Record<string, string> = {
    allow: 'bg-green-100',
    warn: 'bg-yellow-100',
    deny: 'bg-red-100',
    all: 'bg-slate-200',
  };
  return statusButtons.value.map((btn) => {
    const classes: string[] = ['policy-filter-btn', 'status-btn'];
    if (!btn.active) {
      classes.push('border-slate-400');
    } else {
      classes.push(twBorderColorByKey[btn.key] || 'border-slate-500');
      classes.push(borderClassByKey[btn.key] || 'all-border', 'active');
    }
    // Common static utility classes
    const staticUtils =
      'min-w-[150px] font-medium normal-case transition duration-150 focus-visible:shadow-[0_0_0_2px_rgba(var(--v-theme-primary),0.5)]';
    const activeUtils = btn.active
      ? `${activeBgByKey[btn.key] || 'bg-slate-200'} text-slate-900 shadow-sm`
      : 'hover:bg-slate-50';
    const renderClasses = [...classes, staticUtils, activeUtils, 'border'].filter(Boolean).join(' ');
    // Style object for thin borders (0.5px inactive, 1px active)
    const style = btn.active ? {borderWidth: '1px'} : {borderWidth: '0.5px'};
    return {...btn, classes, renderClasses, style};
  });
});

const show = computed({
  get: () => props.show,
  set: (value: boolean) => {
    emit('update:show', value);
    if (!value) {
      emit('close');
    }
  },
});

//table header should be computed otherwise not translated on language change
const licenseHeaders = computed<DataTableHeader[]>(() => [
  {key: 'name', title: t('COL_LICENSE_NAME'), align: 'start', sortable: true, width: 250},
  {key: 'identifier', title: t('COL_LICENSE_ID'), align: 'start', sortable: true, width: 150},
  {key: 'policyType', title: t('COL_POLICY_TYPE'), align: 'center', sortable: true, width: 150},
  {key: 'key', title: t('COL_KEY'), align: 'start', sortable: true, width: 200},
]);

const getLicenseInfoForClipboard = (item: License): string => {
  return (
    `License Reference\n` +
    `Policy Rule: ${props.policyRule?.name}\n` +
    `License Name: ${item.name}\n` +
    `SPDX ID: ${item.identifier}\n` +
    `License Key: ${item.key}`
  );
};

function resetDialog() {
  show.value = false;
}

function close() {
  show.value = false;
}

const dialogConfig = computed<DialogLayoutConfig>(() => ({
  title: `${t('TITLE_POLICY_RULES')} - ${props.policyRule?.name || ''}`,
}));
</script>

<template>
  <v-dialog v-model="show" content-class="large overflow-hidden" width="1200" @after-leave="resetDialog">
    <DialogLayout v-if="policyRule" :config="dialogConfig" @close="close">
      <!-- Category filter buttons -->
      <div class="mb-4 flex flex-wrap gap-3" role="group" :aria-label="t('FILTER_GROUP_POLICY_TYPES')">
        <v-btn
          v-for="btn in statusButtonsWithClasses"
          :key="btn.key"
          :prepend-icon="btn.active ? btn.outlineIcon || btn.icon : btn.icon"
          density="comfortable"
          variant="outlined"
          :color="btn.color"
          :class="btn.renderClasses"
          :style="btn.style"
          :aria-pressed="btn.active"
          :aria-label="btn.ariaLabel"
          :aria-current="btn.active ? 'true' : undefined"
          @click="btn.key === 'all' ? (selectedCategory = 'all') : toggleCategory(btn.key as any)">
          <span :class="['mr-1 font-semibold', btn.key !== 'all' && btn.count === 0 ? 'opacity-50' : '']">{{
            btn.count
          }}</span>
          <span>{{ btn.displayLabel }}</span>
        </v-btn>
      </div>

      <!-- Search field -->
      <div class="d-flex align-center mb-4">
        <v-spacer />
        <v-text-field
          v-model="search"
          append-inner-icon="mdi-magnify"
          :label="t('labelSearch')"
          variant="outlined"
          clearable
          density="compact"
          hide-details
          class="max-w-[400px]" />
      </div>

      <v-data-table
        fixed-header
        height="400"
        :headers="licenseHeaders"
        :items="filteredLicenses"
        :search="search"
        :sort-by="sortItems"
        :loading="loading"
        :items-per-page="100"
        :footer-props="{'items-per-page-options': [10, 25, 50, 100, -1]}"
        density="compact"
        class="striped-table">
        <template v-slot:item.actions="{item}">
          <DCopyClipboardButton
            icon="mdi-content-copy"
            :hint="t('TT_COPY_LICENSE_INFO')"
            :content="getLicenseInfoForClipboard(item)" />
        </template>
      </v-data-table>
    </DialogLayout>
  </v-dialog>
</template>
