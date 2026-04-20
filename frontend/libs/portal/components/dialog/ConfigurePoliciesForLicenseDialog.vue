<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {PolicyRulesAssignmentsDto, PolicyRulesForLicenseDto} from '@disclosure-portal/model/PolicyRule';
import licenseService from '@disclosure-portal/services/license';
import {useAppStore} from '@disclosure-portal/stores/app';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader} from '@shared/types/table';
import _lodash from 'lodash';
import {nextTick, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const emit = defineEmits<{
  (e: 'policyRulesSaved'): void;
}>();

// State variables
const {t} = useI18n();
const show = ref(false);
const title = 'LM_DIALOG_TITLE_CONFIGURE_POLICIES_FOR_LICENSE';
const licenseId = ref('');
const licenseName = ref('');
const formDialog = ref<DiscoForm | null>(null);
const model = ref<PolicyRulesForLicenseDto>({} as PolicyRulesForLicenseDto);
const itemsOriginal = ref<PolicyRulesAssignmentsDto[]>([]);
const items = ref<PolicyRulesAssignmentsDto[]>([]);
const headers = ref<DataTableHeader[]>([]);
const search = ref('');
const dataLoaded = ref(false);
const typeWidth = ref(420);

headers.value = [
  {
    title: 'Name',
    align: 'start',
    filterable: true,
    class: 'tableHeaderCell',
    value: 'name',
  },
  {
    title: 'Description',
    align: 'start',
    filterable: false,
    class: 'tableHeaderCell',
    value: 'description',
    sortable: false,
  },
];

const load = async () => {
  let value = await licenseService.getPolicyRuleAssignmentsForThisLicence(licenseId.value);
  model.value = value;
  items.value = value.policyRulesAssignments;
  itemsOriginal.value = _lodash.cloneDeep(items.value);
  dataLoaded.value = true;
};

const doDialogAction = async () => {
  const changedItems = [] as PolicyRulesAssignmentsDto[];
  _lodash.each(items.value, (item) => {
    const itemOriginal = _lodash.find(itemsOriginal.value, {key: item.key});
    if (itemOriginal !== undefined && item.type !== itemOriginal.type) {
      changedItems.push(item);
    }
  });

  if (changedItems.length > 0) {
    model.value.policyRulesAssignments = changedItems;
    await licenseService.updatePolicyRulesAssignmentsForLicense(model.value);
    emit('policyRulesSaved');
  }

  useSnackbar().info(t('DIALOG_policy_rule_edit_success'));
  close();
};

const close = () => {
  if (formDialog.value) {
    formDialog.value.reset();
    formDialog.value.resetValidation();
  }
  show.value = false;
};

watch(show, (value) => {
  if (!value) {
    if (formDialog.value) {
      formDialog.value.reset();
    }
  }
});

onMounted(() => {
  const appStore = useAppStore();
  typeWidth.value = appStore.getAppLanguage === 'en' ? 420 : 470;
  headers.value.push({
    title: 'Type',
    align: 'center',
    filterable: false,
    class: 'tableHeaderCell',
    value: 'type',
    width: typeWidth.value,
    sortable: false,
  });
});

const open = (id: string, name: string) => {
  show.value = true;
  licenseId.value = id;
  licenseName.value = name;
  nextTick(() => {
    if (formDialog.value) {
      formDialog.value.reset();
    }
  });
  load();
};

defineExpose({open});
</script>

<template>
  <v-dialog v-model="show" scrollable persistent width="80%">
    <v-form ref="formDialog">
      <v-card class="pa-8 dDialog biggerDialog" min-height="800px">
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <span class="text-h5">
                {{ t(title) }}
                <template v-if="licenseName">
                  <q>{{ licenseName }}</q>
                </template>
              </span>
            </v-col>
            <v-col cols="2" class="d-flex justify-end">
              <DCloseButton @click="close" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-row class="justify-end">
            <v-col cols="12" xs="12" sm="6" md="5">
              <DSearchField v-model="search" />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12" v-if="items && dataLoaded">
              <PolicyRulesTable v-model="items" :loading="!dataLoaded" :isDialog="true" edit></PolicyRulesTable>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="pt-4">
          <v-spacer></v-spacer>
          <DCActionButton
            isDialogButton
            size="small"
            variant="text"
            @click="close"
            class="mr-5"
            :text="t('BTN_CANCEL')"></DCActionButton>
          <DCActionButton
            isDialogButton
            size="small"
            variant="flat"
            @click="doDialogAction"
            :text="t('NP_DIALOG_BTN_EDIT')"></DCActionButton>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
