<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {PolicyRuleDto} from '@disclosure-portal/model/PolicyRule';
import {ProjectModel} from '@disclosure-portal/model/Project';
import projectService from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {getCssClassForReadonlyRow} from '@disclosure-portal/utils/Table';
import {openUrlInNewTab} from '@disclosure-portal/utils/url';
import {formatDateTimeShort} from '@disclosure-portal/utils/View';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {useClipboard} from '@shared/utils/clipboard';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {DataTableHeader} from '@shared/types/table';

const {t} = useI18n();
const projectStore = useProjectStore();
const router = useRouter();
const breadcrumbs = useBreadcrumbsStore();
const {copyToClipboard} = useClipboard();

const search = ref('');
const tablePolicyRules = ref<HTMLElement | null>(null);
const rules = ref<PolicyRuleDto[]>([]);

const headers = computed((): DataTableHeader[] => {
  return [
    {
      title: t('COL_ACTIONS'),
      sortable: false,
      align: 'center',
      width: 120,
      value: 'actions',
    },
    {
      title: t('COL_NAME'),
      sortable: true,
      value: 'name',
      width: 200,
    },
    {
      title: t('COL_DESCRIPTION'),
      sortable: false,
      value: 'description',
      align: 'start',
      width: 180,
    },
  ];
});

const projectModel = computed(() => projectStore.currentProject!);

const reloadInternal = async () => {
  if (projectModel.value && projectModel.value._key) {
    if (projectModel.value.isGroup) {
      restartWithCorrectView(projectModel.value);
    }
    rules.value = (await projectService.getPolicyRules(projectModel.value._key)).data;
  }
  initBreadcrumbs();
};

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    {
      title: t('BC_Dashboard'),
      href: '/dashboard/home',
    },
    {
      title: t('BC_Projects'),
      href: '/dashboard/projects/',
    },
    {
      title: '' + projectModel.value.name,
      href: '/dashboard/projects/' + encodeURIComponent(projectModel.value._key),
    },
  ]);
};

const restartWithCorrectView = (item: ProjectModel) => {
  if (item.isGroup) {
    router.push('/dashboard/groups/' + encodeURIComponent(item._key));
  } else {
    router.push('/dashboard/projects/' + encodeURIComponent(item._key));
  }
};

const openRule = (item: PolicyRuleDto) => {
  openUrlInNewTab(`/dashboard/policyrules/${encodeURIComponent(item.key)}`);
};

const getReferenceInfoForClipboard = (item: PolicyRuleDto) => {
  return `Disclosure Portal Policy Rule Reference

Rule Name: ${item.name}
Rule Identifier: ${item.key}
Rule Description: ${item.description}
Rule Created: ${formatDateTimeShort(item.created)}
Rule Updated: ${formatDateTimeShort(item.updated)}`;
};

const copyRuleToClipboard = (item: PolicyRuleDto) => {
  const content = getReferenceInfoForClipboard(item);
  copyToClipboard(content);
};

const actionButtons = computed((): TableActionButtonsProps['buttons'] => [
  {
    icon: 'mdi-content-copy',
    hint: t('TT_COPY_REFERENCE_INFO'),
    event: 'copy',
    show: true,
  },
  {
    icon: 'mdi-open-in-new',
    hint: t('TT_open_rule'),
    event: 'open',
    show: true,
  },
]);

onMounted(() => {
  reloadInternal();
});

watch(projectModel, async (value) => {
  if (value && value._key) {
    if (value.isGroup) {
      restartWithCorrectView(value);
      return;
    }
    rules.value = (await projectService.getPolicyRules(value._key)).data;
  }
});
</script>

<template>
  <TableLayout has-title has-tab>
    <template #buttons>
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <div ref="tablePolicyRules" class="fill-height">
        <v-data-table
          density="compact"
          fixed-header
          class="striped-table fill-height"
          :search="search"
          :headers="headers"
          :items="rules"
          :footer-props="{
            'items-per-page-options': [10, 50, 100, -1],
          }"
          :item-class="getCssClassForReadonlyRow">
          <template #[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="actionButtons"
              @copy="copyRuleToClipboard(item)"
              @open="openRule(item)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
</template>
