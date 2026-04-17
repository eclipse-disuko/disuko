<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <div class="pa-4">
    <v-row class="header d-flex align-center ga-4 flex-row" v-if="mode === 'edit'">
      <h2 class="d-headline">Alias</h2>
      <DCActionButton large :text="t('BTN_ADD')" icon="mdi-plus" :hint="t('TT_add_alias')" @click="createEmptyAlias" />
    </v-row>
    <v-row v-if="props.license?.aliases">
      <v-col cols="12" xs="12">
        <v-data-table
          :items="props.license?.aliases"
          :headers="headers"
          density="compact"
          class="striped-table"
          fixed-header
          item-key="_key"
          :hide-default-footer="true">
          <template v-slot:item.licenseId="{item}">
            <v-text-field
              autocomplete="off"
              class="my-1 pt-5"
              v-if="item._key === ''"
              density="compact"
              solo
              variant="outlined"
              :rules="licenseIdRules"
              v-model="item.licenseId"></v-text-field>
            <span v-else>{{ item.licenseId }}</span>
          </template>
          <template v-slot:item.description="{item}">
            <v-text-field
              autocomplete="off"
              class="my-1 pt-5"
              v-if="item._key === ''"
              density="compact"
              variant="outlined"
              v-model="item.description"></v-text-field>
            <span v-else>{{ item.description }}</span>
          </template>
          <template v-slot:item.actions="{_, index}" v-if="mode === 'edit'">
            <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <DeleteConfirmationDialog v-slot="{showDialog}" @confirmed="() => deleteAliasByIndex(index)">
                  <DCActionButton icon="mdi-close" @click="showDialog" v-bind="props" variant="plain"></DCActionButton>
                </DeleteConfirmationDialog>
              </template>
              <span>{{ t('TT_delete_alias') }}</span>
            </v-tooltip>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts" setup>
import License, {AliasDTO} from '@disclosure-portal/model/License';
import LicenseService from '@disclosure-portal/services/license';
import {isSpdxAliasIdentifier} from '@disclosure-portal/utils/Validation';
import DeleteConfirmationDialog from '@shared/components/dialogs/DeleteConfirmationDialog.vue';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import {DataTableHeader} from '@shared/types/table';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();

const emits = defineEmits(['update:license']);
const props = defineProps<{license: License; mode: string}>();
const computedLicense = computed({
  get: () => props.license,
  set: (value: License) => {
    emits('update:license', value);
  },
});
const isAliasAdded = ref(false);

const headers = computed<DataTableHeader[]>(() => {
  const res: DataTableHeader[] = [
    {
      title: t('COL_NAME'),
      align: 'start',
      width: '180',
      value: 'licenseId',
    },
    {
      title: t('COL_DESCRIPTION'),
      align: 'start',
      value: 'description',
    },
  ];
  if (props.mode === 'edit') {
    res.push({
      title: t('COL_ACTIONS'),
      align: 'end',
      width: 140,
      class: 'tableHeaderCell',
      value: 'actions',
      sortable: false,
    });
  }
  return res;
});

const deleteAliasByIndex = (index: number) => {
  computedLicense.value.aliases = computedLicense.value.aliases.filter((_, i) => i !== index);
};

const createEmptyAlias = async () => {
  const aliasDto = new AliasDTO();
  aliasDto._key = '';
  aliasDto.licenseId = '';
  aliasDto.description = '';

  if (!props.license.aliases) {
    props.license.aliases = [];
  }
  computedLicense.value.aliases.unshift(aliasDto);
  isAliasAdded.value = true;
};

const existsLicenseId = (licenseId: string): boolean => {
  return computedLicense.value?.aliases.filter((alias) => alias.licenseId === licenseId).length > 1;
};

const licenseIdRules = [
  (v: string) => !!v || 'Field is required',
  (v: string) => isSpdxAliasIdentifier(v) || t('ERROR_VALIDATION_LICENSE_ID_IS_NOT_VALID'),
  (v: string) => v !== props.license.licenseId || 'This alias matches the license ID',
  (v: string) => !existsLicenseId(v) || 'This alias already exists for this license',
  async (v: string) => {
    const res = await LicenseService.headAlias(v);
    return !res.data.found || 'This alias already exists for another license';
  },
];
</script>
