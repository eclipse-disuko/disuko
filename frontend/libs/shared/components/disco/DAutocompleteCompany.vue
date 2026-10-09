<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {Department} from '@shared/model/Department';
import Tooltip from '@shared/components/disco/Tooltip.vue';
import {UserMetaData} from '@shared/types/Users';
import _ from 'lodash';
import {computed, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

defineOptions({inheritAttrs: false});

const props = withDefaults(
  defineProps<{
    fetchOptions: (query: string) => Promise<Department[]>;
    required?: boolean;
    label?: string;
    readonly?: boolean;
    help?: string;
    matchingDepartment?: UserMetaData;
  }>(),
  {
    required: false,
    label: '',
    readonly: false,
    help: undefined,
    matchingDepartment: undefined,
  },
);

const emit = defineEmits(['depChanged']);
const dept = defineModel<Department | null>({required: true});

const {t} = useI18n();
const currentQuery = ref('');
const suggestions = ref<Department[]>([]);
const matchingSuggestion = ref<Department>();
const normalize = (value?: string) => value?.trim().toLowerCase() ?? '';
const isMatching = (department: Department) => department.deptId === matchingSuggestion.value?.deptId;
const defaultSuggestions = () => (matchingSuggestion.value ? [matchingSuggestion.value] : []);

const sortSuggestions = (departments: Department[]) =>
  departments.sort((a, b) => {
    if (isMatching(a)) return -1;
    if (isMatching(b)) return 1;
    if (a.level === 0 && b.level !== 0) return -1;
    if (a.level !== 0 && b.level === 0) return 1;
    return b.level - a.level;
  });

const loadMatchingDepartment = async () => {
  matchingSuggestion.value = undefined;
  const companyCode = normalize(props.matchingDepartment?.companyIdentifier);
  const description = normalize(props.matchingDepartment?.departmentDescription);
  if (companyCode && description) {
    const departments = await props.fetchOptions(companyCode);
    matchingSuggestion.value = departments.find(
      (d) => normalize(d.companyCode) === companyCode && normalize(d.descriptionEnglish) === description,
    );
  }
  if (currentQuery.value.length < 3) {
    suggestions.value = defaultSuggestions();
  }
};

watch(() => props.matchingDepartment, loadMatchingDepartment, {deep: true, immediate: true});

const selected = computed({
  get: () => (dept.value?.deptId ? dept.value : null),
  set: (val: Department | null) => {
    dept.value = val ?? ({} as Department);
    if (val) {
      emit('depChanged', val);
    }
  },
});

const rules = computed(() =>
  props.required ? [(v: Department | null) => !!v?.deptId || t('IS_REQUIRED', {fieldName: t('COMPANY_CODE')})] : [],
);

const noDataText = computed(() => {
  if (currentQuery.value.length < 3) {
    return t('TYPE_AT_LEAST_3');
  }
  return t('NO_RESULTS');
});

const searchChanged = (query: string) => {
  currentQuery.value = query;
  if (!query || query.length < 3) {
    suggestions.value = defaultSuggestions();
    return;
  }
  props.fetchOptions(query.toLowerCase().trim()).then((res) => {
    suggestions.value = sortSuggestions(res);
  });
};

const debouncedSearchChanged = _.debounce(searchChanged, 300);
</script>

<template>
  <v-autocomplete
    class="group"
    autocomplete="off"
    v-model="selected"
    :label="label"
    :items="suggestions"
    @update:search="debouncedSearchChanged"
    :no-data-text="noDataText"
    :item-title="() => ''"
    item-value="deptId"
    return-object
    :custom-filter="() => true"
    :required="required"
    :clearable="!readonly"
    :class="{required: required}"
    variant="outlined"
    :rules="rules"
    :readonly="readonly"
    persistent-placeholder
    hide-details="auto">
    <template v-if="help" #append-inner>
      <Tooltip :text="help" as-parent>
        <v-icon
          icon="mdi-help-circle-outline"
          class="cursor-help text-gray-400 opacity-0 transition-opacity duration-250 group-focus-within:opacity-100 group-hover:opacity-100" />
      </Tooltip>
    </template>
    <template #item="{item, props: itemProps}">
      <v-list-item
        v-bind="itemProps"
        :class="{'dep-level-1': item.raw.level === 1 && !isMatching(item.raw)}"
        class="px-2">
        <v-list-item-title>
          <span class="font-weight-bold">
            {{ `[${item.raw.companyCode}] ${item.raw.companyName}` }}
          </span>
          <v-chip
            v-if="isMatching(item.raw)"
            class="ml-2"
            size="x-small"
            color="success"
            variant="tonal"
            prepend-icon="mdi-account-check-outline"
            data-testid="matching-department">
            {{ t('MATCHES_YOUR_PROFILE') }}
          </v-chip>
        </v-list-item-title>
        <v-list-item-subtitle>
          {{ `[${item.raw.deptId}, ${item.raw.orgAbbreviation}, ${item.raw.level}] ${item.raw.descriptionEnglish}` }}
        </v-list-item-subtitle>
      </v-list-item>
    </template>
    <template #selection="{item}">
      <v-list-item class="px-0">
        <v-list-item-title>
          <span class="font-weight-bold">
            {{ `[${item.raw.companyCode}] ${item.raw.companyName}` }}
          </span>
        </v-list-item-title>
        <v-list-item-subtitle>
          {{ `[${item.raw.deptId}, ${item.raw.orgAbbreviation}] ${item.raw.descriptionEnglish}` }}
        </v-list-item-subtitle>
      </v-list-item>
    </template>
    <template v-if="matchingSuggestion && currentQuery.length < 3" #append-item>
      <div class="text-caption text-medium-emphasis px-2 py-1" data-testid="search-other-department-hint">
        <v-icon icon="mdi-magnify" size="small" class="mr-1" />
        {{ t('SEARCH_OTHER_DEPARTMENT_HINT') }}
      </div>
    </template>
  </v-autocomplete>
</template>

<style>
.dep-level-1 {
  color: #808080;
  background-color: rgba(0, 0, 0, 0.2) !important;
}
</style>
