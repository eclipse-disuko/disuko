// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {
  DialogBulkPolicyDecisionEntry,
  DialogBulkPolicyDecisionsConfig,
} from '@disclosure-portal/components/dialog/DialogConfigs';
import {ErrorDialogInterface} from '@disclosure-portal/components/dialog/DialogInterfaces';
import ErrorDialog from '@disclosure-portal/components/dialog/ErrorDialog.vue';
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import {PolicyDecisionRequest} from '@disclosure-portal/model/PolicyDecision';
import projectService from '@disclosure-portal/services/projects';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import useRules from '@disclosure-portal/utils/Rules';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import {
  getIconColorForPolicyType,
  getIconForPolicyType,
  policyStateToTranslationKey,
} from '@disclosure-portal/utils/View';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader} from '@shared/types/table';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';

type TableItem = DialogBulkPolicyDecisionEntry & {
  key: string;
};

const {t} = useI18n();
const {info} = useSnackbar();
const rules = useRules();
const appStore = useAppStore();
const userStore = useUserStore();
const emit = defineEmits(['reload']);
const projectStore = useProjectStore();

const form = ref<VForm | null>(null);
const isVisible = ref(false);

const comment = ref<string | undefined>(undefined);
const verification = ref(false);
const errorDialog = ref<ErrorDialogInterface | null>(null);

const tableItems = ref<TableItem[]>([]);
const selected = ref<TableItem[]>([]);

const commentRules = rules.minMax(t('LICENSE_RULE_COMMENT'), 0, 80, true);

const config = ref<DialogBulkPolicyDecisionsConfig>({
  items: [],
});

const headers: DataTableHeader[] = [
  {
    title: t('COL_SPDX_STATUS'),
    sortable: true,
    align: 'center',
    value: 'policy.type',
    width: 110,
  },
  {
    title: t('COL_POLICY_NAME'),
    align: 'start',
    width: 200,
    value: 'policy.name',
    sortable: true,
  },
  {
    title: t('COL_LICENSE_MATCHED'),
    sortable: true,
    align: 'center',
    value: 'policy.licenseMatched',
    width: 200,
  },
  {
    title: t('COL_COMPONENT_NAME'),
    align: 'start',
    value: 'component.name',
    sortable: true,
    width: 250,
  },
  {
    title: t('COL_COMPONENT_VERSION'),
    align: 'start',
    value: 'component.version',
    width: 130,
    sortable: true,
  },
  {
    title: t('COL_LICENSE_EXPRESSION'),
    align: 'start',
    width: 200,
    value: 'component.licenseExpression',
    sortable: true,
  },
];

const projectKey = computed(() => projectStore.currentProject!._key);
const currentVersionKey = computed(() => appStore.getCurrentVersion._key);
const currentSbomId = computed(() => appStore.getSelectedSpdx._key);
const currentSbomName = computed(() => appStore.getSelectedSpdx.MetaInfo.Name);
const currentSbomUploaded = computed(() => appStore.getSelectedSpdx.Uploaded);
const buttonsDisabled = computed(() => !verification.value || selected.value.length === 0);

const open = async (
  newConfig: DialogBulkPolicyDecisionsConfig = {
    items: [],
  },
) => {
  config.value = newConfig;

  tableItems.value = (config.value.items ?? []).map((item) => ({...item, key: crypto.randomUUID()}));
  selected.value = tableItems.value.slice();

  isVisible.value = true;
};

const doDialogAction = async (decision: 'allow' | 'deny') => {
  if (!(await form.value?.validate())?.valid) {
    return;
  }

  const requestItems: PolicyDecisionRequest[] = [];
  for (const item of selected.value) {
    const policyDecisionRequest: PolicyDecisionRequest = {
      sbomId: currentSbomId.value,
      sbomName: currentSbomName.value,
      sbomUploaded: currentSbomUploaded.value,
      componentSpdxId: item.component.spdxId,
      componentName: item.component.name,
      componentVersion: item.component.version,
      licenseExpression: item.component.licenseExpression,
      licenseId: item.policy.licenseMatched,
      policyId: item.policy.key,
      policyEvaluated: item.policy.type,
      policyDecision: decision,
      comment: comment.value ?? '',
      creator: userStore.getProfile.user,
    };
    requestItems.push(policyDecisionRequest);
  }

  const response = (
    await projectService.createBulkPolicyDecision(projectKey.value, currentVersionKey.value, requestItems)
  ).data;
  if (!response.success) {
    const dialog = new ErrorDialogConfig();
    dialog.title = t('policy_decision_create_error_title');
    dialog.description = t(response.message);
    errorDialog.value?.open(dialog);
    return;
  }

  form.value?.reset();
  emit('reload');
  close();
  info(t('POLICY_DECISION_CREATED'));
};

const dialogConfig = computed(() => ({
  title: t('POLICY_DECISION_CREATE'),
}));

const close = () => {
  form.value?.reset();
  isVisible.value = false;
};

const formatText = (text: string): string => {
  text = escapeHtml(text);
  if (text.includes(' AND ') || text.includes(' OR ')) {
    return text
      .replace(/ AND /g, ' <strong class="db-highlight">AND</strong> ')
      .replace(/ OR /g, ' <strong class="db-highlight">OR</strong> ');
  }
  return text;
};

defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" width="1300" persistent>
    <DialogLayout :config="dialogConfig" @close="close">
      <template #right>
        <DCActionButton
          is-dialog-button
          size="small"
          :variant="buttonsDisabled ? 'flat' : 'tonal'"
          @click="doDialogAction('deny')"
          :disabled="buttonsDisabled"
          :color="buttonsDisabled ? 'gray' : 'error'"
          icon="mdi-minus-circle"
          :text="t('DENY')" />

        <DCActionButton
          is-dialog-button
          size="small"
          :variant="buttonsDisabled ? 'flat' : 'tonal'"
          @click="doDialogAction('allow')"
          :disabled="buttonsDisabled"
          :color="buttonsDisabled ? 'gray' : 'success'"
          icon="mdi-check-circle"
          :text="t('ALLOW')" />
      </template>
      <v-form ref="form">
        <Stack>
          <v-data-table
            :headers="headers"
            fixed-header
            density="compact"
            hide-default-footer
            :items-per-page="-1"
            :items="tableItems"
            class="striped-table custom-data-table"
            height="380"
            show-select
            v-model="selected"
            return-object
            item-value="key">
            <template v-slot:[`item.policy.type`]="{item}">
              <v-icon :color="getIconColorForPolicyType(item.policy.type)">
                {{ getIconForPolicyType(item.policy.type) }}
              </v-icon>
              <tooltip :text="policyStateToTranslationKey(item.policy.type)"></tooltip>
            </template>
            <template v-slot:[`item.component.licenseExpression`]="{item}">
              <span v-html="formatText(item.component.licenseExpression)"></span>
            </template>
          </v-data-table>
          <v-textarea
            v-model="comment"
            variant="outlined"
            density="compact"
            :label="t('LICENSE_RULE_COMMENT')"
            hide-details="auto"
            persistent-placeholder
            :rules="commentRules" />
          <v-checkbox v-model="verification" :label="t('WARNED_POLICY_DECISION_VERIFICATION_NOTE_TEXT')" hide-details />
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
  <ErrorDialog ref="errorDialog" />
</template>
