// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts" setup>
import ProjectPostRequest from '@disclosure-portal/model/ProjectPostRequest';
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import {LabelsTools} from '@disclosure-portal/utils/Labels';
import useRules from '@disclosure-portal/utils/Rules';
import DialogLayout from '@shared/layouts/DialogLayout.vue';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const props = defineProps({
  fromList: {
    type: Boolean,
    required: false,
    default: true,
  },
});
const {t} = useI18n();
const project = ref(new ProjectPostRequest());
const isEditor = ref(false);
const show = ref(false);
const {minMax} = useRules();
const activeRules = ref({
  description: minMax(t('NP_DIALOG_GROUP_DESCRIPTION'), 3, 1000, true),
  name: minMax(t('NP_DIALOG_TF_DEVELOPER'), 3, 80, true),
});
const groupAddForm = ref<VForm | null>(null);
const labelTools = new LabelsTools();
const projectStore = useProjectStore();

const emit = defineEmits(['modified']);

const edit = (modelItem: ProjectSlim) => {
  const request = new ProjectPostRequest();
  isEditor.value = true;
  request.fillWithProjectSlim(modelItem);
  project.value = request;
  show.value = true;
  nextTick(() => {
    if (groupAddForm.value) {
      groupAddForm.value.resetValidation();
    }
  });
};

const close = () => {
  show.value = false;
};

const showDialog = async () => {
  await labelTools.loadLabels();
  project.value = new ProjectPostRequest();
  project.value.freeLabels = [];
  project.value.owner = useUserStore().getProfile.user;
  project.value.policyLabels = [];
  project.value.children = [];
  project.value.schemaLabel = labelTools.schemaLabelsMapByName['common standard']?._key;
  project.value.isGroup = true;
  isEditor.value = false;
  show.value = true;
};

const save = async () => {
  await nextTick();
  const isValid = await groupAddForm.value?.validate();
  if (!isValid?.valid || project.value.name.length === 0) {
    return;
  }
  const isUpdate = isEditor.value;
  if (isUpdate) {
    await projectStore.updateProject(project.value);
    if (props.fromList) {
      emit('modified');
    }
  } else {
    await projectStore.createProject(project.value);
    emit('modified');
  }
  close();
};
const dialogConfig = computed(() => {
  return isEditor.value
    ? {
        title: t('Title_Edit_Group'),
        primaryButton: {text: t('NP_DIALOG_BTN_EDIT'), disabled: projectStore.loading, loading: projectStore.loading},
        secondaryButton: {text: t('BTN_CANCEL'), disabled: projectStore.loading},
      }
    : {
        title: t('Title_New_Group'),
        primaryButton: {text: t('NP_DIALOG_BTN_CREATE'), disabled: projectStore.loading, loading: projectStore.loading},
        secondaryButton: {text: t('BTN_CANCEL'), disabled: projectStore.loading},
      };
});
defineExpose({showDialog, edit});
</script>

<template>
  <slot name="default" :showDialog="showDialog" />
  <slot name="edit" :editSlim="edit" />
  <v-dialog v-model="show" content-class="medium" scrollable width="500">
    <DialogLayout :config="dialogConfig" @secondary-action="close" @primary-action="save" @close="close">
      <v-form ref="groupAddForm" @submit.prevent="save">
        <Stack>
          <v-text-field
            required
            autocomplete="off"
            variant="outlined"
            density="compact"
            :rules="activeRules.name"
            v-model="project.name"
            class="required"
            :label="t('NP_DIALOG_GROUP_NAME')"
            autofocus />
          <v-textarea
            no-resize
            required
            variant="outlined"
            density="compact"
            v-model="project.description"
            :label="t('NP_DIALOG_GROUP_DESCRIPTION')"
            hide-details="auto"
            counter="1000"
            :rules="activeRules.description" />
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
</template>
