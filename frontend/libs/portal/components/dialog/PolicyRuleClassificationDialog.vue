<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import PolicyRule from '@disclosure-portal/model/PolicyRule';
import IObligation from '@disclosure-portal/model/IObligation';
import adminService from '@disclosure-portal/services/admin';
import {
  toBucketDefinition,
  toRuleStatusMap,
  RuleStatus,
  statusConfig,
} from '@disclosure-portal/utils/calculatedPolicyRule';
import {DiscoForm} from '@disclosure-portal/types/discobasics';

type StatusColumn = {key: RuleStatus | null; label: string; icon: string | null; color: string | null};
import useRules from '@disclosure-portal/utils/Rules';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, nextTick, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info: snack, error} = useSnackbar();
const {minMax} = useRules();

const emit = defineEmits(['reload']);

const isVisible = ref(false);
const saving = ref(false);
const dialog = ref<DiscoForm | null>(null);
const obligations = ref<IObligation[]>([]);

const name = ref('');
const rules = ref<Record<string, RuleStatus>>({});
const originalItem = ref<PolicyRule | null>(null);

const buildDto = (base: PolicyRule): Partial<PolicyRule> => ({
  _key: base._key,
  name: name.value,
  description: base.description,
  labelSets: base.labelSets,
  auxiliary: base.auxiliary,
  active: base.active,
  applyToAll: base.applyToAll,
  calculated: true,
  calculatedConfig: {
    bucketDefinition: toBucketDefinition(rules.value),
    licenseScope: base.calculatedConfig.licenseScope,
  },
});

const statusColumns: StatusColumn[] = [
  ...(['allowed', 'warned', 'denied', 'forbidden'] as const).map((key) => ({
    key,
    label: statusConfig[key].labelKey,
    icon: statusConfig[key].icon,
    color: statusConfig[key].color,
  })),
  {key: null, label: 'MATRIX_STATUS_NOT_SET', icon: null, color: null},
];

const validationRules = {
  name: minMax(t('PRC_DIALOG_TF_NAME'), 1, 200, false),
};

const dialogConfig = computed(() => ({
  title: t('PRC_DIALOG_TITLE_EDIT'),
  primaryButton: {
    text: t('NP_DIALOG_BTN_EDIT'),
    loading: saving.value,
  },
  secondaryButton: {
    text: t('BTN_CANCEL'),
  },
}));

const open = (existing: PolicyRule & {rules?: Record<string, RuleStatus>}) => {
  name.value = existing.name;
  rules.value = existing.rules ? {...existing.rules} : toRuleStatusMap(existing);
  originalItem.value = existing;
  dialog.value?.reset();
  isVisible.value = true;
};

const setRule = (obligationKey: string, status: RuleStatus | undefined) => {
  if (status) {
    rules.value[obligationKey] = status;
  } else {
    delete rules.value[obligationKey];
  }
};

const doDialogAction = async () => {
  await nextTick();
  const info = await dialog.value?.validate();
  if (!info?.valid) return;
  if (!originalItem.value) return;

  saving.value = true;
  try {
    await adminService.putPolicyRule(buildDto(originalItem.value) as PolicyRule);
    snack(t('DIALOG_prc_update_success'));
    emit('reload');
    isVisible.value = false;
  } catch {
    error(t('DIALOG_prc_save_error'));
  } finally {
    saving.value = false;
  }
};

const close = () => {
  isVisible.value = false;
};

onMounted(async () => {
  const res = await adminService.getAllObligations();
  obligations.value = res.data.items ?? res.data ?? [];
});

defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" scrollable width="800">
    <DialogLayout :config="dialogConfig" @primary-action="doDialogAction" @secondary-action="close" @close="close">
      <v-form ref="dialog" @submit.prevent="doDialogAction">
        <Stack>
          <v-text-field
            v-model="name"
            :rules="validationRules.name"
            :label="t('PRC_DIALOG_TF_NAME')"
            autofocus
            autocomplete="off"
            variant="outlined"
            density="compact" />

          <v-table density="compact" class="striped-table mt-4" fixed-header height="400">
            <thead>
              <tr>
                <th class="w-[250px] text-left">{{ t('COL_CLASSIFICATION') }}</th>
                <th v-for="col in statusColumns" :key="col.key ?? 'not-set'" class="w-[80px] text-center">
                  <div class="flex items-center justify-center gap-1">
                    <v-icon v-if="col.icon" :icon="col.icon" :color="col.color ?? undefined" size="small" />
                    <span class="text-caption">{{ t(col.label) }}</span>
                  </div>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="ob in obligations" :key="ob._key">
                <td class="w-[250px]">{{ ob.name }}</td>
                <td v-for="col in statusColumns" :key="col.key ?? 'not-set'" class="w-[80px] text-center">
                  <v-radio
                    :model-value="rules[ob._key] ?? null"
                    :value="col.key"
                    @click="setRule(ob._key, col.key ?? undefined)"
                    density="compact"
                    hide-details
                    class="d-flex justify-center" />
                </td>
              </tr>
            </tbody>
          </v-table>
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
</template>
