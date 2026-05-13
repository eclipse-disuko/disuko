<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import icons from '@disclosure-portal/constants/icons';
import PolicyRule from '@disclosure-portal/model/PolicyRule';
import AdminService from '@disclosure-portal/services/admin';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useSnackbar from '@shared/composables/useSnackbar';
import {nextTick, Ref, ref} from 'vue';
import {useI18n} from 'vue-i18n';

interface Props {
  policyRule?: PolicyRule;
  policyLabels?: any[];
}

const props = withDefaults(defineProps<Props>(), {
  policyLabels: () => [],
});

const emit = defineEmits(['reload']);

const {t} = useI18n();
const {info: snack} = useSnackbar();

const isVisible = ref(false);
const item = ref(new PolicyRule());
const title = ref('');
const isLoading = ref(false);
const policyRulesDialog: Ref<DiscoForm | null> = ref(null);

const activeRules = {
  required: [(v: any) => !!v || t('REQUIRED_FIELD', {field: 'name'})],
  description: [(v: any) => v.length <= 1500 || t('MAX_LENGTH_EXCEEDED', {field: 'Description', length: 1500})],
};

const doDialogAction = async () => {
  const info = await policyRulesDialog.value?.validate();
  if (!info?.valid) {
    return;
  }
  if (item.value.name) {
    item.value.name = item.value.name.trim();
  }
  isLoading.value = true;

  try {
    if (props.policyRule) {
      await AdminService.putPolicyRule(item.value);
    } else {
      await AdminService.postPolicyRule(item.value);
    }
    snack(t(props.policyRule ? 'DIALOG_policy_rule_edit_success' : 'DIALOG_policy_rule_create_success'));
    isVisible.value = false;
    emit('reload');
  } catch (error: any) {
    snack(t('ERROR_POLICYRULES_API'));
    console.error(error);
  } finally {
    isLoading.value = false;
  }
};
const addLabelSet = () => {
  if (item.value.labelSets[item.value.labelSets.length - 1].length > 0) {
    item.value.labelSets.push([]);
  }
};
const showDialog = () => {
  title.value = props.policyRule ? 'AL_DIALOG_TITLE_EDIT' : 'AL_DIALOG_TITLE';
  if (props.policyRule) {
    const cloned = new PolicyRule(Object.assign({}, props.policyRule));
    cloned.labelSets = (props.policyRule.labelSets || []).map((set) => [...set]);
    item.value = cloned;
  } else {
    item.value = new PolicyRule();
  }
  isVisible.value = true;
};
const closeDialog = () => {
  if (isLoading.value) {
    return;
  }
  isVisible.value = false;
};
const removeLabelSetIfEmpty = async (index: number) => {
  if (item.value.labelSets[index]?.length === 0) {
    await nextTick();
    item.value.labelSets = item.value.labelSets.toSpliced(index, 1);
  }
};

defineExpose({
  showDialog,
});
</script>

<template>
  <slot :showDialog="showDialog" />
  <v-dialog v-model="isVisible" persistent width="600" scrollable>
    <v-form ref="policyRulesDialog">
      <v-card class="pa-8 dDialog">
        <v-card-title class="d-flex align-center justify-space-between">
          <span class="text-h5">{{ t(title) }}</span>
          <DCloseButton @click="isVisible = false" />
        </v-card-title>
        <v-card-text class="pt-2">
          <v-row>
            <v-col cols="12" xs="12" class="errorBorder">
              <v-text-field
                autocomplete="off"
                variant="outlined"
                class="required"
                hide-details="auto"
                v-model="item.name"
                :rules="activeRules.required"
                :label="t('AL_DIALOG_TF_NAME')"
                autofocus />
            </v-col>
          </v-row>
          <v-row v-if="!item.applyToAll">
            <v-col cols="12" xs="12" class="pb-2">
              <v-select
                variant="outlined"
                :class="{
                  'pb-2': item.labelSets.length > 1 && index !== item.labelSets.length - 1,
                }"
                hide-details="auto"
                v-model="item.labelSets[index]"
                item-title="name"
                item-value="_key"
                clearable
                multiple
                :items="policyLabels"
                :label="t('AL_DIALOG_SB_LABELS')"
                @update:modelValue="() => removeLabelSetIfEmpty(index)"
                v-bind:menu-props="{location: 'bottom'}"
                v-for="(_, index) in item.labelSets"
                :key="index">
                <template v-slot:chip="{item, props}">
                  <DLabel closable :parentProps="props" :labelName="item.title" :iconName="icons.TAG" />
                </template>
              </v-select>
            </v-col>
            <v-col cols="12" xs="12" class="pt-0" v-if="!(item.labelSets[item.labelSets.length - 1]?.length <= 0)">
              <div class="d-flex align-center border-md border-opacity-25 border-dashed px-3 py-3" @click="addLabelSet">
                <v-icon color="primary">mdi-plus</v-icon>
                <span class="font-weight-light pl-1">{{ t('NP_DIALOG_MORE_LABEL_SET') }}</span>
              </div>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12">
              <v-textarea
                variant="outlined"
                class="required"
                hide-details="auto"
                no-resize
                v-model="item.description"
                :counter="1500"
                :label="t('AL_DIALOG_TF_DESCRIPTION')"
                :rules="activeRules.description"></v-textarea>
              <v-checkbox
                v-model="item.active"
                hide-details
                color="primary"
                :label="t('ACTIVE_FLAG')"
                class="mt-0 shrink pt-0" />
              <v-checkbox
                v-model="item.applyToAll"
                hide-details
                color="primary"
                :label="t('APPLY_TO_ALL_FLAG')"
                class="mt-0 shrink pt-0" />
              <v-checkbox
                v-model="item.auxiliary"
                hide-details
                color="primary"
                :label="t('AUXILIARY_FLAG')"
                class="mt-0 shrink pt-0" />

              <v-checkbox
                v-model="item.calculated"
                hide-details
                color="primary"
                :label="t('CALCULATED_POLICY_RULE_ENABLED')"
                class="mt-0 shrink pt-0" />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <DCActionButton
            isDialogButton
            class="mr-4"
            variant="text"
            @click="closeDialog"
            color="primary"
            :disabled="isLoading"
            :text="t('BTN_CANCEL')" />
          <DCActionButton
            isDialogButton
            size="small"
            variant="flat"
            @click="doDialogAction"
            :loading="isLoading"
            :text="policyRule ? t('NP_DIALOG_BTN_SAVE') : t('NP_DIALOG_BTN_CREATE')"
            color="primary" />
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
