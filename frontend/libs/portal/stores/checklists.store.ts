// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Checklist, ChecklistItem} from '@disclosure-portal/model/Checklist';
import adminService from '@disclosure-portal/services/admin';
import {defineStore} from 'pinia';
import {computed, onMounted, reactive, toRefs} from 'vue';
import {useRoute} from 'vuetify/lib/composables/router';

export const useChecklistsStore = defineStore('checklists', () => {
  const route = useRoute();

  const state = reactive({
    checklists: [] as Checklist[],
    isLoading: false,
    loaded: false,
  });

  const fetchAll = async () => {
    state.isLoading = true;
    state.checklists = (await adminService.getChecklist()).data;
    state.isLoading = false;
    state.loaded = true;
  };

  const createChecklist = async (checklist: Checklist) => {
    await adminService.createChecklist(checklist);
    await fetchAll();
  };

  const editChecklist = async (checklist: Checklist) => {
    await adminService.editChecklist(checklist);
    await fetchAll();
  };

  const deleteChecklist = async (key: string) => {
    await adminService.deleteChecklistById(key);
    await fetchAll();
  };

  const checklistId = computed(() =>
    Array.isArray(route.value?.params?.id) ? route.value.params.id[0] : route.value?.params?.id,
  );

  const checklist = computed(() => {
    if (!checklistId.value) {
      return null;
    }

    const checklist = state.checklists.find((c) => c._key === checklistId.value);

    if (!checklist) {
      return null;
    }

    return checklist;
  });

  const addItem = async (item: ChecklistItem) => {
    const newChecklist = (await adminService.createChecklistItem(checklist.value!._key, item)).data;
    replaceChecklist(newChecklist);
  };

  const editItem = async (item: ChecklistItem) => {
    const newChecklist = (await adminService.editChecklistItem(checklist.value!._key, item)).data;
    replaceChecklist(newChecklist);
  };

  const deleteItem = async (key: string) => {
    const newChecklist = (await adminService.deleteChecklistItem(checklist.value!._key, key)).data;
    replaceChecklist(newChecklist);
  };

  const replaceChecklist = (repl: Checklist) => {
    const i = state.checklists.findIndex((cl) => cl._key === repl._key);
    if (i === -1) {
      return;
    }
    state.checklists[i] = repl;
  };

  onMounted(async () => {
    if (state.loaded) {
      return;
    }
    await fetchAll();
  });

  return {
    ...toRefs(state),
    fetchAll,
    createChecklist,
    editChecklist,
    deleteChecklist,
    checklist,
    addItem,
    editItem,
    deleteItem,
  };
});
