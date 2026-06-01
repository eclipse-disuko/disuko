<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {
  PolicyRuleClassificationDto,
  PolicyRuleClassificationRequestDto,
  RuleStatus,
} from '@disclosure-portal/model/PolicyRuleClassification';
import IObligation from '@disclosure-portal/model/IObligation';
import adminService from '@disclosure-portal/services/admin';
import policyRuleClassificationService from '@disclosure-portal/services/policyruleclassification';
import {statusConfig} from '@disclosure-portal/utils/classificationStatus';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, nextTick, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info: snack, error} = useSnackbar();
const {minMax} = useRules();

const emit = defineEmits(['reload']);

const isVisible = ref(false);
const isEdit = ref(false);
const saving = ref(false);
const dialog = ref<DiscoForm | null>(null);
const obligations = ref<IObligation[]>([]);

const name = ref('');
const itemKey = ref('');
const rules = ref<Record<string, RuleStatus>>({});

const statusColumns = (['allowed', 'warned', 'denied', 'forbidden'] as const).map((key) => ({
  key,
  label: statusConfig[key].labelKey,
  icon: statusConfig[key].icon,
  color: statusConfig[key].color,
}));

const validationRules = {
  name: minMax(t('PRC_DIALOG_TF_NAME'), 1, 200, false),
};

const dialogConfig = computed(() => ({
  title: t(isEdit.value ? 'PRC_DIALOG_TITLE_EDIT' : 'PRC_DIALOG_TITLE_ADD'),
  primaryButton: {
    text: t(isEdit.value ? 'NP_DIALOG_BTN_EDIT' : 'NP_DIALOG_BTN_CREATE'),
    loading: saving.value,
  },
  secondaryButton: {
    text: t('BTN_CANCEL'),
  },
}));

const open = (existing?: PolicyRuleClassificationDto) => {
  if (existing) {
    isEdit.value = true;
    itemKey.value = existing._key;
    name.value = existing.name;
    rules.value = {...existing.rules};
  } else {
    isEdit.value = false;
    itemKey.value = '';
    name.value = '';
    rules.value = {};
  }
  dialog.value?.reset();
  isVisible.value = true;
};

const setRule = (obligationKey: string, status: RuleStatus | 'forbidden') => {
  if (status && status !== 'forbidden') {
    rules.value[obligationKey] = status;
  } else {
    delete rules.value[obligationKey];
  }
};

const doDialogAction = async () => {
  await nextTick();
  const info = await dialog.value?.validate();
  if (!info?.valid) return;

  saving.value = true;
  try {
    const dto: PolicyRuleClassificationRequestDto = {name: name.value, rules: {...rules.value}};
    if (isEdit.value) {
      await policyRuleClassificationService.update(itemKey.value, dto);
      snack(t('DIALOG_prc_update_success'));
    } else {
      await policyRuleClassificationService.create(dto);
      snack(t('DIALOG_prc_create_success'));
    }
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
            class="required"
            v-model="name"
            :rules="validationRules.name"
            :label="t('PRC_DIALOG_TF_NAME')"
            autofocus
            variant="outlined"
            density="compact" />

          <v-table density="compact" class="striped-table mt-4" fixed-header height="400">
            <thead>
              <tr>
                <th class="text-left" style="width: 250px">{{ t('COL_CLASSIFICATION') }}</th>
                <th v-for="col in statusColumns" :key="col.key" class="text-center" style="width: 80px">
                  <div class="d-flex flex-column align-center">
                    <v-icon :color="col.color" size="small">{{ col.icon }}</v-icon>
                    <span class="text-caption">{{ t(col.label) }}</span>
                  </div>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="ob in obligations" :key="ob._key">
                <td style="width: 250px">{{ ob.name }}</td>
                <td v-for="col in statusColumns" :key="col.key" class="text-center" style="width: 80px">
                  <v-radio
                    :model-value="rules[ob._key]"
                    :value="col.key"
                    @click="setRule(ob._key, col.key)"
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
