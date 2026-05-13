<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {Approval} from '@disclosure-portal/model/Approval';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {computed, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

export type ApprovalTabs =
  | 'history'
  | 'general'
  | 'generalReview'
  | 'generalExternal'
  | 'details'
  | 'documents'
  | 'attributes'
  | 'task';

const {t} = useI18n();
const currentTab = ref(0);

const props = withDefaults(
  defineProps<{
    item: Approval;
    taskDescription: string;
    tabsList: ApprovalTabs[];
    showRedWarnDeniedDecisionsMessage?: boolean;
  }>(),
  {
    showRedWarnDeniedDecisionsMessage: false,
  },
);

const emit = defineEmits<{
  'reloads-approvals': [];
}>();

const reload = async () => {
  emit('reloads-approvals');
};

watch(
  () => props.item,
  () => {
    currentTab.value = 0;
  },
);

const creator = computed(() => `${props.item.creatorFullName} (${props.item.creator})`);
const approver = computed(() => `${props.item.plausibility.approverFullName} (${props.item.plausibility.approver})`);
const requestCreated = computed(() => formatDateAndTime(props.item.created));
const reviewUpdated = computed(() => formatDateAndTime(props.item.plausibility.state.updated));
</script>

<template>
  <div class="pt-0">
    <v-tabs v-model="currentTab" slider-color="brand" show-arrows bg-color="tabsHeader">
      <v-tab v-for="(tab, tabIndex) in tabsList" :key="tabIndex">
        <span>{{ t(`TAB_TITLE_${tab.toUpperCase().replace('REVIEW', '').replace('EXTERNAL', '')}`) }}</span>
      </v-tab>
    </v-tabs>
    <v-tabs-window v-model="currentTab" grow class="pa-2 min-h-[350px]">
      <v-tabs-window-item v-for="(tab, tabIndex) in tabsList" :key="tabIndex" class="py-4">
        <template v-if="tab === 'task'">
          <v-col cols="12" xs="12" class="pa-0">
            <blockquote class="taskMessage" v-html="taskDescription"></blockquote>
          </v-col>
          <v-col cols="12" xs="12" class="px-0 pt-8">
            <DApprovalComponents
              :stats="item.info.stats"
              :showRedWarnDeniedDecisionsMessage="showRedWarnDeniedDecisionsMessage" />
          </v-col>
        </template>
        <DApprovalComponents
          v-if="tab === 'general'"
          :stats="item.info.stats"
          :showRedWarnDeniedDecisionsMessage="showRedWarnDeniedDecisionsMessage" />
        <template v-if="tab === 'generalReview'">
          <Stack>
            <Stack direction="row" justify="between" align="center" class="gap-6">
              <v-text-field
                autocomplete="off"
                :label="t('TAD_USER_ID')"
                v-model="creator"
                readonly
                variant="outlined"
                hide-details />
              <v-text-field
                autocomplete="off"
                :label="t('APPROVER_LABEL')"
                v-model="approver"
                readonly
                variant="outlined"
                hide-details />
            </Stack>
            <Stack direction="row" justify="between" align="center" class="gap-6">
              <v-text-field
                autocomplete="off"
                :label="t('Lbl_created')"
                v-model="requestCreated"
                readonly
                variant="outlined"
                hide-details />
              <v-text-field
                autocomplete="off"
                :label="t('Lbl_updated')"
                v-model="reviewUpdated"
                readonly
                variant="outlined"
                hide-details />
            </Stack>
            <Stack direction="row" justify="between" align="center" class="gap-6">
              <v-textarea
                rows="3"
                auto-grow
                variant="outlined"
                readonly
                :label="t('TAD_USER_ID') + ' ' + t('TAD_COMMENT')"
                v-model="item.comment"
                hide-details />
              <v-textarea
                rows="3"
                auto-grow
                variant="outlined"
                readonly
                :label="t('APPROVER_LABEL') + ' ' + t('TAD_COMMENT')"
                v-model="item.plausibility.comment"
                hide-details />
            </Stack>
            <DApprovalComponents
              :stats="item.info.stats"
              :showRedWarnDeniedDecisionsMessage="showRedWarnDeniedDecisionsMessage" />
          </Stack>
        </template>
        <template v-if="tab === 'generalExternal'">
          <DExternalApprovalReview :external-approval="item" @reloading="reload"></DExternalApprovalReview>
          <DApprovalComponents
            :stats="item.info.stats"
            :showRedWarnDeniedDecisionsMessage="showRedWarnDeniedDecisionsMessage" />
        </template>
        <template v-if="tab == 'details'">
          <GridSPDXList :projects="item.info.projects" :approval-items="[item]"></GridSPDXList>
        </template>

        <DApprovalState v-if="tab == 'history'" :internal="item.internal"></DApprovalState>
        <DApprovalDocuments v-if="tab == 'documents'" :documents="item.documents" />
      </v-tabs-window-item>
    </v-tabs-window>
  </div>
</template>
