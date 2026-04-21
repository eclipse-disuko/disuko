<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {usePageTitle} from '@disclosure-portal/composables/usePageTitle';
import icons from '@disclosure-portal/constants/icons';
import Label from '@disclosure-portal/model/Label';
import PolicyRule from '@disclosure-portal/model/PolicyRule';
import AdminService from '@disclosure-portal/services/admin';
import policyRuleService from '@disclosure-portal/services/policyrules';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {formatDateTimeShort, IMap} from '@disclosure-portal/utils/View';
import {useTabsWindows} from '@shared/composables/useTabsWindows';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const {t} = useI18n();
const router = useRouter();
const breadcrumbs = useBreadcrumbsStore();
const {useReactiveTitle} = usePageTitle();

const rule = ref(new PolicyRule());
const isPolicyManager = ref(false);
const hasChanges = ref(false);
const selectedFilterClassificationsSelected = ref<string[]>([]);
const labelsMap = ref<IMap<Label>>({});
const policyLabels = ref<Label[]>([]);
const ruleId = ref('');

onMounted(() => {
  ruleId.value = router.currentRoute.value.params.id as string;
});
// This has to stay after declaration of ruleId
const baseUrl = computed(() => `/dashboard/policyrules/${encodeURIComponent(ruleId.value)}`);
// This has to stay after declaration of baseUrl
const {tabUrl, selectedTab} = useTabsWindows(baseUrl, [
  'tabPolicyRuleTable',
  'tabPolicyRuleOverview',
  'tabChangeLog',
  'tabAuditLog',
]);
const retrieveRule = async () => {
  rule.value = new PolicyRule((await policyRuleService.getPolicyRule(ruleId.value)).data);
};

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    {title: t('BC_Dashboard'), disabled: false, href: '/dashboard/home'},
    {title: t('POLICY_RULES'), disabled: false, href: '/dashboard/policyrules'},
    {
      title: '' + rule.value.Name,
      disabled: false,
      href: '/dashboard/policyrules/' + encodeURIComponent(ruleId.value),
    },
  ]);
};

const reload = async () => {
  await retrieveRule();
  initBreadcrumbs();
  hasChanges.value = false;
};

const getChangeLog = async () => {
  return await AdminService.getChangeLog(ruleId.value);
};

const getAuditTrail = async () => {
  return await AdminService.getAuditTrail(ruleId.value);
};

const reloadLabels = async () => {
  policyLabels.value = (await AdminService.getPolicyLabels()).data;
  createLabelsMap();
};

const createLabelsMap = () => {
  labelsMap.value = {};
  for (const lbl of policyLabels.value) {
    labelsMap.value[lbl._key] = lbl;
  }
};

watch(selectedFilterClassificationsSelected, reload);

onMounted(async () => {
  isPolicyManager.value = RightsUtils.isPolicyManager();
  await retrieveRule();
  initBreadcrumbs();

  await reloadLabels();
});

// Set up reactive title based on policy rule name
watch(
  () => rule.value?.Name,
  (rule) => {
    if (rule) {
      useReactiveTitle(rule + ' | Policy Rule');
    }
  },
  {immediate: true},
);
</script>

<template>
  <v-container fluid data-testid="policy-rules-details">
    <v-row>
      <v-col xs="12" md="auto" class="d-flex align-center">
        <span class="text-h5 pr-2">{{ t('POLICY_RULE') }}</span>
        <span class="text-h5 px-2">{{ rule.Name }}</span>
        <NewPolicyRuleDialog v-slot="{showDialog}" :policy-labels="policyLabels" :policy-rule="rule" @reload="reload">
          <DCActionButton
            v-if="isPolicyManager && !rule.Deprecated"
            large
            class="align-content-center mx-2"
            icon="mdi-pencil"
            :text="t('BTN_EDIT')"
            :hint="t('TT_edit_rule')"
            data-testid="policy-rules-edit-button"
            @click="showDialog"></DCActionButton>
        </NewPolicyRuleDialog>
      </v-col>
    </v-row>
    <v-row class="">
      <v-col class="pt-0">
        <v-card class="card-border">
          <v-tabs v-model="selectedTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
            <v-tab :to="tabUrl.tabPolicyRuleTable" value="tabPolicyRuleTable">
              {{ t('TAB_PROJECT_SELECTED') }}
            </v-tab>
            <v-tab :to="tabUrl.tabPolicyRuleOverview" value="tabPolicyRuleOverview">
              {{ t('TAB_PROJECT_DETAIL') }}
            </v-tab>
            <v-tab :to="tabUrl.tabChangeLog" value="tabChangeLog">
              {{ t('TAB_CHANGE_LOG') }}
            </v-tab>
            <v-tab v-if="isPolicyManager" value="tabAuditLog" :to="tabUrl.tabAuditLog">
              {{ t('TAB_PROJECT_AUDIT') }}
            </v-tab>
          </v-tabs>
          <v-tabs-window v-model="selectedTab">
            <v-tabs-window-item value="tabPolicyRuleTable">
              <GridTabPolicyRules></GridTabPolicyRules>
            </v-tabs-window-item>
            <v-tabs-window-item value="tabPolicyRuleOverview" class="pa-3">
              <v-row class="pa-4">
                <v-col md="6" xs="12" sm="12">
                  <v-sheet class="d-flex flex-row">
                    <v-card-text>
                      <span class="text-caption text-grey-lighten-1">{{ t('CREATED') }}</span
                      ><br />
                      <span class="text-body-2">{{ formatDateTimeShort(rule.created) }}</span
                      ><br />
                    </v-card-text>
                    <v-card-text>
                      <span class="text-caption text-grey-lighten-1">{{ t('UPDATED') }}</span
                      ><br />
                      <span class="text-body-2">{{ formatDateTimeShort(rule.updated) }}</span>
                    </v-card-text>
                    <v-card-text>
                      <span class="text-caption text-grey-lighten-1">{{ t('DEPRECATED') }}</span
                      ><br />
                      <span class="text-body-2" v-if="rule.Deprecated">{{
                        formatDateTimeShort(rule.DeprecatedDate)
                      }}</span>
                    </v-card-text>
                  </v-sheet>
                  <v-sheet>
                    <v-card-text>
                      <span class="text-caption text-grey-lighten-1">{{ t('DESCRIPTION') }}</span
                      ><br />
                      <span>{{ String(rule.Description) }}</span>
                    </v-card-text>
                  </v-sheet>
                </v-col>
                <v-col md="6" xs="12" sm="12" v-if="!rule.ApplyToAll">
                  <v-sheet class="d-flex flex-row">
                    <v-card-text>
                      <v-icon color="grey-lighten-1" size="x-small" class="mr-2">{{ icons.POLICY }}</v-icon>
                      <span class="text-caption text-grey-lighten-1">{{ t('POLICY_LABELS_SETS') }}</span>
                      <div class="d-flex flex-row flex-wrap">
                        <div
                          v-for="(labelSets, lsi) in rule.LabelSets"
                          :key="lsi"
                          class="ma-2 d-flex justify-space-between">
                          <div
                            v-if="rule.LabelSets.length >= 1 && rule.LabelSets[0].length >= 1"
                            class="policyLabelBorder pa-2 rounded">
                            <v-tooltip
                              :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
                              bottom
                              v-for="(l, i) in labelSets"
                              :key="i">
                              <template v-slot:activator="{props}">
                                <v-sheet width="180" class="justify-center">
                                  <v-card-text class="pa-1 d-text d-secondary-text pl-5">
                                    {{ labelsMap?.[l] ? labelsMap?.[l].name : 'UNKNOWN_LABEL' }}
                                  </v-card-text>
                                </v-sheet>
                              </template>
                              <span v-if="labelsMap?.[l]?.description.trim()">
                                {{ labelsMap[l].description }}
                              </span>
                              <span v-else>{{ t('NO_DESCRIPTION') }}</span>
                            </v-tooltip>
                          </div>
                        </div>
                      </div>
                    </v-card-text>
                  </v-sheet>
                </v-col>
                <v-col md="6" xs="12" sm="12" v-else>
                  <v-sheet class="d-flex flex-row">
                    <v-card-text>
                      <span class="text-caption text-grey-lighten-1">{{ t('APPLY_TO_ALL_FLAG') }}</span>
                      <br />
                      <span>{{ t('BTN_YES') }}</span>
                    </v-card-text>
                  </v-sheet>
                </v-col>
              </v-row>
            </v-tabs-window-item>
            <v-tabs-window-item value="tabChangeLog">
              <GridChangeLog ref="changeLog" :fetch-method="() => getChangeLog(ruleId)" />
            </v-tabs-window-item>
            <v-tabs-window-item v-if="isPolicyManager" value="tabAuditLog">
              <GridAuditLog ref="auditLog" :fetch-method="() => getAuditTrail(ruleId)" />
            </v-tabs-window-item>
          </v-tabs-window>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
