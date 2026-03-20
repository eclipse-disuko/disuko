// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import Label, {type LabelType} from '@disclosure-portal/model/Label';
import AdminService from '@disclosure-portal/services/admin';
import {defineStore} from 'pinia';
import {computed, reactive, toRefs} from 'vue';

/**
 * New Label store that should replace labels in app store
 * TODO: Migrate existing usages from app store to this store and remove labels from app store
 */
export const useLabelStore = defineStore('label', () => {
  const state = reactive({
    labels: [] as Label[],
  });

  const fetchAllLabels = async () => {
    try {
      state.labels = (await AdminService.getLabels()).data;
    } catch (e) {
      console.error('Error fetching labels:', e);
    }
  };

  const getLabelByKey = (key: string): Partial<Label> => {
    const notFoundLabel: Partial<Label> = {
      name: 'UNKNOWN_LABEL',
      description: '',
    };
    const foundLabel = state.labels.find((label) => label._key === key);
    return foundLabel ? foundLabel : notFoundLabel;
  };

  const getLabelByNameAndType = (name: string, type: LabelType): Label | undefined => {
    return state.labels.find((label) => label.name === name && label.type === type);
  };

  const policyLabels = computed(() =>
    state.labels.filter((label) => label.type === 'POLICY').sort((a, b) => a.name.localeCompare(b.name)),
  );

  return {
    ...toRefs(state),
    policyLabels,
    fetchAllLabels,
    getLabelByKey,
    getLabelByNameAndType,
  };
});
