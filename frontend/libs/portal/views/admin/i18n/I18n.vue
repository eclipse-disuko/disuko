<script lang="ts" setup>
import {getApi} from '@disclosure-portal/api';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const {api} = getApi();

const {t} = useI18n();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();
const router = useRouter();

interface LocaleRegistryItem {
  code: string;
  isDefault: boolean;
  entryCount: number;
  displayName?: string;
  nativeDisplayName?: string;
}

interface I18nLocaleListItem {
  localeCode: string;
  displayName: string;
  nativeName: string;
  isDefault: boolean;
  scope: string;
  entryCount: number;
}

interface LanguageRow {
  code: string;
  name: string;
  nativeName: string;
  keyCount: number;
  missingCount: number;
  extraCount: number;
  isDefault: boolean;
}

const localeRegistry = ref<LocaleRegistryItem[]>([]);
const isLoadingLocales = ref(false);

const fetchLocales = async () => {
  isLoadingLocales.value = true;
  try {
    const res = await api.get<I18nLocaleListItem[]>('/api/v1/i18n');
    localeRegistry.value = (res.data || []).map((item) => ({
      code: item.localeCode,
      isDefault: item.isDefault,
      entryCount: item.entryCount,
      displayName: item.displayName || undefined,
      nativeDisplayName: item.nativeName || undefined,
    }));
  } finally {
    isLoadingLocales.value = false;
  }
};

const languageRows = computed<LanguageRow[]>(() => {
  const defaultLocale = localeRegistry.value.find((item) => item.isDefault) ?? localeRegistry.value[0];
  const defaultCount = defaultLocale?.entryCount ?? 0;

  return localeRegistry.value.map((item) => {
    const diff = item.entryCount - defaultCount;
    return {
      code: item.code,
      name: item.displayName || t(`LANG_${item.code}`),
      nativeName: item.nativeDisplayName || t(`LANG_NATIVE_${item.code}`),
      keyCount: item.entryCount,
      missingCount: diff < 0 ? Math.abs(diff) : 0,
      extraCount: diff > 0 ? diff : 0,
      isDefault: item.isDefault,
    };
  });
});

const search = ref('');
const sortItems = ref<SortItem[]>([{key: 'code', order: 'asc'}]);

const pageTitle = computed(() => t('ADMIN_I18N_PAGE_TITLE'));
const pageDescription = computed(() => t('ADMIN_I18N_PAGE_DESCRIPTION'));

const labels = computed(() => ({
  sectionTitle: t('ADMIN_I18N_SECTION_TITLE'),
  sectionSubtitle: t('ADMIN_I18N_SECTION_SUBTITLE'),
  colLocale: t('ADMIN_I18N_COL_LOCALE'),
  colLanguage: t('COL_LANG').trim(),
  colKeys: t('ADMIN_I18N_COL_KEYS'),
  colMissing: t('ADMIN_I18N_COL_MISSING'),
  colExtra: t('ADMIN_I18N_COL_EXTRA'),
}));

const headers = computed((): DataTableHeader[] => [
  {
    title: labels.value.colLocale,
    align: 'start',
    value: 'code',
    sortable: true,
    width: 120,
  },
  {
    title: labels.value.colLanguage,
    align: 'start',
    value: 'name',
    sortable: true,
    width: 260,
    minWidth: 220,
  },
  {
    title: labels.value.colKeys,
    align: 'start',
    value: 'keyCount',
    sortable: true,
    width: 120,
  },
  {
    title: labels.value.colMissing,
    align: 'start',
    value: 'missingCount',
    sortable: true,
    width: 120,
  },
  {
    title: labels.value.colExtra,
    align: 'start',
    value: 'extraCount',
    sortable: true,
    width: 120,
  },
]);

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    ...dashboardCrumbs,
    {
      title: pageTitle.value,
    },
  ]);
};

const onRowClick = (_event: Event, row: {item: LanguageRow}) => {
  router.push({name: 'I18nLocaleDetails', params: {localeCode: row.item.code}});
};

onMounted(async () => {
  await fetchLocales();
  initBreadcrumbs();
});
</script>

<template>
  <TableLayout>
    <template #description>
      <h1 class="text-h5">{{ pageTitle }}</h1>
      <p class="text-body-2 text-medium-emphasis mt-1">{{ pageDescription }}</p>
    </template>
    <template #buttons>
      <h2 class="text-h6">{{ labels.sectionTitle }}</h2>
      <v-spacer></v-spacer>
      <v-text-field
        autocomplete="off"
        :max-width="500"
        append-inner-icon="mdi-magnify"
        variant="outlined"
        density="compact"
        v-model="search"
        :label="t('labelSearch')"
        single-line
        hide-details></v-text-field>
    </template>
    <template #table>
      <div class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table fill-height"
          fixed-header
          :headers="headers"
          :items="languageRows"
          :search="search"
          :sort-by="sortItems"
          @click:row="onRowClick"
          :items-per-page="10">
          <template #[`item.code`]="{item}">
            <v-chip size="small" color="primary" variant="tonal">{{ item.code.toUpperCase() }}</v-chip>
          </template>
          <template #[`item.name`]="{item}">
            <div class="font-weight-medium">{{ item.name }}</div>
            <div v-if="item.nativeName.toLowerCase() !== item.name.toLowerCase()" class="text-caption text-medium-emphasis">{{ item.nativeName }}</div>
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
</template>
