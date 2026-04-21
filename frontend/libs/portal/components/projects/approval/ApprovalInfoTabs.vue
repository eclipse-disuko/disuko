<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import GridSPDXList from '@disclosure-portal/components/grids/GridSPDXList.vue';
import {Approval} from '@disclosure-portal/model/Approval';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import DApprovalComponents from '@shared/components/disco/DApprovalComponents.vue';
import DApprovalDocuments from '@shared/components/disco/DApprovalDocuments.vue';
import DApprovalState from '@shared/components/disco/DApprovalState.vue';
import DExternalApprovalReview from '@shared/components/disco/DExternalApprovalReview.vue';
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
    <v-tabs v-model="currentTab" slider-color="mbti" show-arrows bg-color="tabsHeader">
      <v-tab v-for="(tab, tabIndex) in tabsList" :key="tabIndex">
        <span>{{ t(`TAB_TITLE_${tab.toUpperCase().replace('REVIEW', '').replace('EXTERNAL', '')}`) }}</span>
      </v-tab>
    </v-tabs>
    <v-tabs-window v-model="currentTab" grow class="pa-2" style="min-height: 350px">
      <v-tabs-window-item v-for="(tab, tabIndex) in tabsList" :key="tabIndex" class="py-4">
        <template v-if="tab == 'task'">
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
          v-if="tab == 'general'"
          :stats="item.info.stats"
          :showRedWarnDeniedDecisionsMessage="showRedWarnDeniedDecisionsMessage" />
        <template v-if="tab == 'generalReview'">
          <v-row>
            <v-col cols="6" class="px-0 pb-0">
              <v-text-field
                autocomplete="off"
                :label="t('TAD_USER_ID')"
                v-model="creator"
                readonly
                variant="outlined"
                hide-details></v-text-field>
            </v-col>
            <v-col cols="6" class="px-0 pb-0">
              <v-text-field
                autocomplete="off"
                :label="t('APPROVER_LABEL')"
                v-model="approver"
                readonly
                variant="outlined"
                hide-details></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field
                autocomplete="off"
                :label="t('Lbl_created')"
                v-model="requestCreated"
                readonly
                variant="outlined"
                hide-details></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field
                autocomplete="off"
                :label="t('Lbl_updated')"
                v-model="reviewUpdated"
                readonly
                variant="outlined"
                hide-details></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-textarea
                rows="3"
                auto-grow
                variant="outlined"
                readonly
                :label="t('TAD_USER_ID') + ' ' + t('TAD_COMMENT')"
                v-model="item.comment"
                hide-details></v-textarea>
            </v-col>
            <v-col cols="6" class="mb-4">
              <v-textarea
                rows="3"
                auto-grow
                variant="outlined"
                readonly
                :label="t('APPROVER_LABEL') + ' ' + t('TAD_COMMENT')"
                v-model="item.plausibility.comment"
                hide-details></v-textarea>
            </v-col>
          </v-row>
          <DApprovalComponents
            :stats="item.info.stats"
            :showRedWarnDeniedDecisionsMessage="showRedWarnDeniedDecisionsMessage" />
        </template>
        <template v-if="tab == 'generalExternal'">
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
