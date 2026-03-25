<script setup lang="ts">
import {ErrorDialogInterface} from '@disclosure-portal/components/dialog/DialogInterfaces';
import ErrorDialog from '@disclosure-portal/components/dialog/ErrorDialog.vue';
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import {PolicyDecisionRequest} from '@disclosure-portal/model/PolicyDecision';
import {ComponentInfoSlim, PolicyRuleStatus} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import useRules from '@disclosure-portal/utils/Rules';
import {getIconColorForPolicyType, getIconForPolicyType} from '@disclosure-portal/utils/View';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';
import {DialogPolicyDecisionConfig} from '@disclosure-portal/components/dialog/DialogConfigs';

const {t} = useI18n();
const {info} = useSnackbar();
const rules = useRules();
const appStore = useAppStore();
const userStore = useUserStore();
const emit = defineEmits(['reload', 'triggerBulk']);
const projectStore = useProjectStore();

const form = ref<VForm | null>(null);
const isVisible = ref(false);

const comment = ref<string | undefined>(undefined);
const verification = ref(false);
const errorDialog = ref<ErrorDialogInterface | null>(null);
const selectedPolicy = ref<PolicyRuleStatus | undefined>(undefined);

const policyDecisionRule = rules.required(t('POLICY_DECISION'));
const commentRulesMinMaxRule = rules.minMax(t('LICENSE_RULE_COMMENT'), 0, 80, true);
const commentRequiredRule = rules.required(t('LICENSE_RULE_COMMENT'));

const config = ref<DialogPolicyDecisionConfig>({
  component: new ComponentInfoSlim(),
  policies: [],
  type: 'warn',
});

const projectKey = computed(() => projectStore.currentProject!._key);
const currentVersionKey = computed(() => appStore.getCurrentVersion._key);
const currentSbomId = computed(() => appStore.getSelectedSpdx._key);
const currentSbomName = computed(() => appStore.getSelectedSpdx.MetaInfo.Name);
const currentSbomUploaded = computed(() => appStore.getSelectedSpdx.Uploaded);

const selectedComponent = computed(() => config.value.component);
const policies = computed(() => config.value.policies);

const selectedComponentStr = computed(() => {
  if (selectedComponent.value) {
    return `${selectedComponent.value.name} (${selectedComponent.value.version})`;
  }
  return '';
});

const licenseExpression = computed(() => selectedComponent.value?.licenseExpression || '');

const licenseMatched = computed(() => selectedPolicy.value?.licenseMatched || '');
const policyName = computed(() => selectedPolicy.value?.name || '');
const policyType = computed(() => selectedPolicy.value?.type || '');
const isWarned = computed(() => config.value.type === 'warn');
const verificationText = computed(() =>
  isWarned.value
    ? t('WARNED_POLICY_DECISION_VERIFICATION_NOTE_TEXT')
    : t('DENIED_POLICY_DECISION_VERIFICATION_NOTE_TEXT'),
);

const commentRequired = computed(() => config.value.type === 'deny');
const commentRules = computed(() => {
  return commentRequired.value ? [...commentRulesMinMaxRule, ...commentRequiredRule] : [...commentRulesMinMaxRule];
});

const open = async (
  newConfig: DialogPolicyDecisionConfig = {
    component: new ComponentInfoSlim(),
    policies: [],
    type: 'warn',
  },
) => {
  config.value = newConfig;
  if (policies.value.length === 1) {
    selectedPolicy.value = policies.value[0];
  }
  isVisible.value = true;
};

const doDialogAction = async (decision: 'allow' | 'deny') => {
  if (!(await form.value?.validate())?.valid) {
    return;
  }

  const policyDecisionRequest: PolicyDecisionRequest = {
    sbomId: currentSbomId.value,
    sbomName: currentSbomName.value,
    sbomUploaded: currentSbomUploaded.value,
    componentSpdxId: selectedComponent.value!.spdxId,
    componentName: selectedComponent.value!.name,
    componentVersion: selectedComponent.value!.version,
    licenseExpression: licenseExpression.value,
    licenseId: licenseMatched.value,
    policyId: selectedPolicy.value!.key,
    policyEvaluated: policyType.value,
    policyDecision: decision,
    comment: comment.value ?? '',
    creator: userStore.getProfile.user,
  };

  const response = (
    await projectService.createPolicyDecision(projectKey.value, currentVersionKey.value, policyDecisionRequest)
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

const closeAndTriggerBulk = () => {
  close();
  emit('triggerBulk');
};

defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" width="650" persistent>
    <DialogLayout :config="dialogConfig" @close="close">
      <template v-if="isWarned" #left>
        <DCActionButton
          size="small"
          is-dialog-button
          @click="closeAndTriggerBulk"
          :text="t('BTN_BULK_POLICY_DECISION')"
          icon="mdi-checkbox-marked-circle-plus-outline" />
      </template>
      <template #right>
        <DCActionButton
          is-dialog-button
          size="small"
          variant="flat"
          @click="doDialogAction('allow')"
          :disabled="!verification"
          :color="!verification ? 'gray' : 'success'"
          :text="t('ALLOW')" />

        <DCActionButton
          v-if="isWarned"
          is-dialog-button
          size="small"
          variant="flat"
          @click="doDialogAction('deny')"
          :disabled="!verification"
          :color="!verification ? 'gray' : 'error'"
          :text="t('DENY')" />
      </template>
      <v-form ref="form">
        <Stack>
          <v-text-field
            autocomplete="off"
            v-model="selectedComponentStr"
            disabled
            variant="outlined"
            density="compact"
            hide-details
            :label="t('RELATED_COMPONENT')" />
          <v-textarea
            auto-grow
            rows="1"
            variant="outlined"
            density="compact"
            disabled
            :label="t('LICENSE_EXPRESSION')"
            v-model="licenseExpression"
            hide-details />
          <v-textarea
            auto-grow
            rows="1"
            variant="outlined"
            density="compact"
            disabled
            :label="t('COL_POLICY_NAME')"
            v-model="policyName"
            hide-details />
          <v-select
            class="mt-5"
            v-model="selectedPolicy"
            :clearable="policies.length > 1"
            :label="t('LICENSE_WITH_POLICY_TYPE')"
            :disabled="policies.length === 1"
            :items="policies"
            return-object
            variant="outlined"
            density="compact"
            hide-details
            required
            :rules="policyDecisionRule">
            <template #item="{item, props}">
              <v-list-item v-bind="props" title="">
                <v-icon size="small" :color="getIconColorForPolicyType(item.raw.type)">
                  {{ getIconForPolicyType(item.raw.type) }}
                </v-icon>
                <span class="d-subtitle-2 ml-2">{{ item.raw.licenseMatched }}</span>
              </v-list-item>
            </template>
            <template #selection="{item}">
              <div class="d-inline">
                <v-icon size="small" :color="getIconColorForPolicyType(item.raw.type)">
                  {{ getIconForPolicyType(item.raw.type) }}
                </v-icon>
                <span class="d-subtitle-2 ml-2">{{ item.raw.licenseMatched }}</span>
              </div>
            </template>
          </v-select>
          <v-textarea
            v-model="comment"
            variant="outlined"
            density="compact"
            :label="t('LICENSE_RULE_COMMENT')"
            hide-details="auto"
            persistent-placeholder
            :class="commentRequired ? 'required' : ''"
            :rules="commentRules" />
          <v-checkbox v-model="verification" :label="verificationText" hide-details />
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
  <ErrorDialog ref="errorDialog" />
</template>
