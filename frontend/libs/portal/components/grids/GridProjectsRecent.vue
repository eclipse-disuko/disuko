<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <v-main tag="div" class="pt-0">
    <v-row class="shrink pb-4" v-if="(items && items.length > 0) || !loaded">
      <v-col>
        <v-row class="shrink">
          <v-col cols="12" xs="12">
            <v-data-table
              :headers="headers"
              :items="items"
              v-model:sort-by="sortItems"
              density="compact"
              class="striped-table custom-data-table"
              :loading="loading"
              :loading-text="t('PROJECTS_LOADING')"
              hide-default-footer
              @click:row="(event, dataItem) => onClickRow(dataItem.item)">
              <template #[`body.append`]>
                <tr>
                  <td></td>
                  <td></td>
                  <td class="d-flex mt-2 justify-end">
                    <span class="font-weight-light">
                      {{ t('TABLE_ITEMS') }}
                      <span class="font-weight-light"> {{ items.length }}</span>
                    </span>
                  </td>
                </tr>
              </template>
              <template #[`item.status`]="{item}">
                <span :class="getItemClass(item)">
                  {{ !item.status ? 'new' : item.status }}
                </span>
              </template>
              <template #[`item.name`]="{item}">
                <v-tooltip
                  :text="`${itemTooltip(item)}`"
                  :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
                  location="bottom"
                  content-class="dpTooltip">
                  <template v-slot:activator="{props}">
                    <span v-bind="props">{{ item.name }}</span>
                  </template>
                </v-tooltip>
              </template>
              <template #[`item.updated`]="{item}">
                <DDateCellWithTooltip :value="item.updated" />
              </template>
            </v-data-table>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-row class="ma-5 d-flex shrink justify-center" v-else>
      <h1 class="d-headline">{{ t('GRID_PROJECTS_RECENT_EMPTY') }}</h1>
    </v-row>
  </v-main>
</template>

<script setup lang="ts">
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import ProjectService from '@disclosure-portal/services/projects';
import {openUrl} from '@disclosure-portal/utils/url';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {onMounted, reactive, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const router = useRouter();
const {t} = useI18n();

const items = reactive<ProjectSlim[]>([]);
const loading = ref(false);
const loaded = ref(false);
const headers: any[] = [
  {
    title: t('COL_STATUS'),
    text: 'COL_STATUS',
    sortable: true,
    filterable: true,
    value: 'status',
    width: '140',
  },
  {
    title: t('COL_NAME'),
    text: 'COL_NAME',
    align: 'start',
    sortable: true,
    value: 'name',
    width: '240',
  },
  {
    title: t('COL_UPDATED'),
    text: 'COL_UPDATED',
    align: 'start',
    sortable: true,
    width: 160,
    value: 'updated',
  },
];

const itemTooltip = (item: ProjectSlim) => {
  return item.description === '' ? item.name : item.description;
};

const sortItems = ref([{key: 'updated', order: 'desc'}]);

const onClickRow = (item: ProjectSlim) => {
  if (item.isGroup) {
    openUrl('/dashboard/groups/' + encodeURIComponent(item._key), router);
  } else {
    openUrl('/dashboard/projects/' + encodeURIComponent(item._key), router);
  }
};

const reload = async () => {
  loading.value = true;
  items.push(...(await ProjectService.getRecent()).data.projects);
  loading.value = false;
  loaded.value = true;
};

onMounted(async () => {
  await reload();
});

const getItemClass = (item: ProjectSlim) => {
  return {
    [`pStatus${item.status ? item.status : 'new'}`]: true,
  };
};
</script>
<style></style>
