<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import useDimensions from '@disclosure-portal/composables/useDimensions';
import Icons from '@disclosure-portal/constants/icons';
import Label from '@disclosure-portal/model/Label';
import {Rights} from '@shared/types/Rights';
import SchemaModel from '@disclosure-portal/model/Schema';
import AdminService from '@disclosure-portal/services/admin';
import {useUserStore} from '@disclosure-portal/stores/user';
import eventBus from '@shared/utils/eventbus';
import {formatDateAndTime, getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import {IMap} from '@disclosure-portal/utils/View';
import useSnackbar from '@shared/composables/useSnackbar';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
import {computed, nextTick, onMounted, onUnmounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {useUrls} from '@shared/composables/useUrls';

const {t} = useI18n();
const router = useRouter();
const breadcrumbs = useBreadcrumbsStore();
const {openUrl} = useUrls();

const items = ref<SchemaModel[]>([]);
const headers = computed<DataTableHeader[]>(() => {
  return [
    {
      title: t('COL_STATUS'),
      align: 'center',
      value: 'active',
      width: 100,
      maxWidth: 120,
    },
    {
      title: t('COL_NAME'),
      align: 'start',
      value: 'name',
      width: 240,
      minWidth: 100,
      sortable: true,
    },
    {
      title: t('COL_DESCRIPTION'),
      align: 'start',
      width: 200,
      minWidth: 200,
      value: 'description',
    },
    {
      title: t('COL_SCHEMA_VERSION'),
      align: 'start',
      value: 'version',
      width: 100,
      sortable: true,
    },
    {
      title: t('COL_LABEL'),
      align: 'start',
      width: 180,
      maxWidth: 180,
      value: 'label',
    },
    {
      title: t('CREATED'),
      align: 'start',
      value: 'created',
      width: 120,
      maxWidth: 120,
      sortable: true,
    },
    {
      title: t('UPDATED'),
      align: 'start',
      value: 'updated',
      width: 120,
      maxWidth: 120,
      sortable: true,
    },
  ];
});
const search = ref('');
const schemaLabels = ref<Label[]>([]);
const labelsMap = ref<IMap<Label>>({});
const rights = ref<Rights>({} as Rights);
const icons = Icons;

const customFilterTable = (value: string, search: string, item: SchemaModel) => {
  if (value) {
    const dateTime = formatDateAndTime(value);
    if (dateTime && dateTime !== 'Invalid date') {
      return dateTime.indexOf(search) > -1;
    }
    let found = false;
    if (labelsMap.value[item.label]) {
      found = found || labelsMap.value[item.label].name.toLowerCase().indexOf(search.toLowerCase()) > -1;
    }
    return found || ('' + value).toLowerCase().indexOf(search.toLowerCase()) > -1;
  }
  return false;
};

const sortItems = () => {
  return [{key: 'active', order: 'desc'} as SortItem];
};

const refresh = async () => {
  items.value = (await AdminService.getAllSchemas()).data;
  if (!items.value) {
    items.value = [];
  }
  await reloadLabels();
};

const reloadLabels = async () => {
  schemaLabels.value = (await AdminService.getSchemaLabels()).data;
  createLabelsMap();
};

const createLabelsMap = () => {
  labelsMap.value = {};
  for (const lbl of schemaLabels.value) {
    labelsMap.value[lbl._key] = lbl;
  }
};

const initBreadcrumbs = () => {
  // EventBus-Implementierung für Breadcrumbs
  breadcrumbs.setCurrentBreadcrumbs([
    {
      title: t('BC_Dashboard'),
      href: '/dashboard/home',
    },
    {
      title: t('BC_ADMIN'),
      href: '/dashboard/admin',
    },
    {
      title: t('BC_SBOM_Schemes'),
      href: '/dashboard/admin/schemas/',
    },
  ]);
};

const isCreateSchemaDialogOpen = ref(false);

const showCreateSchemaDialog = () => {
  isCreateSchemaDialogOpen.value = true;
};

const createSchema = async (schemaFormData: FormData) => {
  AdminService.createSchema(schemaFormData)
    .then(async () => {
      useSnackbar().info(t('DIALOG_schemas_successfully'));
      await refresh();
      isCreateSchemaDialogOpen.value = false;
    })
    .catch((e) => {
      console.error('got error ' + JSON.stringify(e.response.data.code));
    });
};

const onClickRow = (event: Event, item: DataTableItem<SchemaModel>) => {
  if (rights.value.allowSchema && rights.value.allowSchema.read) {
    openUrl('/dashboard/admin/schemas/' + item.item._key, router);
  }
};

const getNameOfType = (type: number) => {
  if (type === 0) {
    return t('TYPE_JSON');
  } else if (type === 1) {
    return t('TYPE_XML');
  } else if (type === 2) {
    return t('TYPE_OCTET');
  }
  return '';
};

const {calculateHeight} = useDimensions();
const tableHeight = ref(0);
const dataTableAsElement = ref<HTMLElement | null>(null);
const updateTableHeight = () => {
  nextTick(() => {
    if (dataTableAsElement.value) {
      tableHeight.value = calculateHeight(dataTableAsElement.value, false);
    }
  });
};

const userStore = useUserStore();
onMounted(async () => {
  rights.value = userStore.getRights;
  initBreadcrumbs();
  updateTableHeight();
  eventBus.on('window-resize', updateTableHeight);

  await refresh();
});

onUnmounted(() => {
  eventBus.off('window-resize', updateTableHeight);
});
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h1 class="text-h5">{{ t('Schemes') }}</h1>
      <DCActionButton
        large
        :text="t('BTN_ADD')"
        icon="mdi-plus"
        :hint="t('TT_add_schema')"
        @click="showCreateSchemaDialog"
        v-if="rights.allowSchema && rights.allowSchema.create" />
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <div ref="dataTableAsElement" class="fill-height">
        <v-data-table
          density="compact"
          hide-default-footer
          class="striped-table fill-height"
          fixed-header
          @click:row="onClickRow"
          :headers="headers"
          :items="items"
          :search="search"
          :height="tableHeight"
          :sort-by="sortItems()"
          :sort-desc="[1]"
          :custom-filter="customFilterTable"
          :items-per-page="50"
          :item-class="getCssClassForTableRow">
          <template #[`item.active`]="{item}">
            <div class="flex justify-center">
              <v-icon size="x-small" :color="item.active ? 'success' : 'warning'">{{ icons.CIRCLE_FILLED }}</v-icon>
            </div>
          </template>
          <template #[`item.created`]="{item}">
            <DDateCellWithTooltip :value="String(item.created)" />
          </template>
          <template #[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="String(item.updated)" />
          </template>
          <template #[`item.name`]="{item}">
            <span>{{ item.name }}</span>
          </template>
          <template #[`item.version`]="{item}">
            {{ item.version }}
          </template>
          <template #[`item.type`]="{item}">
            {{ getNameOfType(item.type) }}
          </template>
          <template #[`item.label`]="{item}">
            <DLabel
              v-if="item.label"
              :labelName="labelsMap[item.label] ? labelsMap[item.label].name : 'UNKNOWN_LABEL'"
              :iconName="icons.SCHEMA"
              class="mt-1" />
            <Tooltip location="bottom">
              {{ labelsMap[item.label] ? labelsMap[item.label].description : '' }}
            </Tooltip>
          </template>
          <template #[`item.description`]="{item}">
            <Truncated>{{ item.description }}</Truncated>
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <NewSchemaDialog
    :isOpen="isCreateSchemaDialogOpen"
    :labels="schemaLabels"
    @update:isOpen="isCreateSchemaDialogOpen = $event"
    @onSave="createSchema" />
</template>
