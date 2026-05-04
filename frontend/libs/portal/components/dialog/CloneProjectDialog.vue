<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import projectService from '@disclosure-portal/services/projects';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {VForm} from 'vuetify/components';
import {useUrls} from '@shared/composables/useUrls';

interface CloneProjectConfig {
  projectKey: string;
  projectName: string;
  count: number;
}

interface CloneProjectDialogProps {
  projectKey?: string;
  projectName?: string;
  initialCount?: number;
}

const props = withDefaults(defineProps<CloneProjectDialogProps>(), {
  projectKey: '',
  projectName: '',
  initialCount: 1,
});

const emit = defineEmits<{
  reload: [];
}>();

const {t} = useI18n();
const {info: infoSnackbar} = useSnackbar();
const showDialog = defineModel<boolean>('showDialog', {required: false});
const router = useRouter();
const {openProjectUrlByKey} = useUrls();

const projectName = ref('');
const projectKey = ref('');
const cloneCount = ref(1);
const cloneForm = ref<VForm | null>(null);

const countRules = [
  (v: number) => !!v || t('DLG_CLONE_PROJECT_COUNT_REQUIRED'),
  (v: number) => (v >= 1 && v <= 10) || t('DLG_CLONE_PROJECT_COUNT_RANGE'),
  (v: number) => Number.isInteger(v) || t('DLG_CLONE_PROJECT_COUNT_INTEGER'),
];

const isValidCount = computed(() => {
  return cloneCount.value >= 1 && cloneCount.value <= 10 && Number.isInteger(cloneCount.value);
});

const dialogConfig = computed(() => ({
  title: t('CLONE_PROJECT'),
  secondaryButton: {text: t('BTN_CANCEL')},
  primaryButton: {text: t('BTN_CLONE'), disabled: !isValidCount.value},
}));

const validateCount = () => {
  if (cloneCount.value < 1) {
    cloneCount.value = 1;
  } else if (cloneCount.value > 10) {
    cloneCount.value = 10;
  }
};

watch(
  [() => props.projectKey, () => props.projectName, () => props.initialCount],
  ([newProjectKey, newProjectName, newInitialCount]) => {
    if (newProjectKey) {
      open({projectKey: newProjectKey, projectName: newProjectName || '', count: newInitialCount});
    }
  },
);

watch(showDialog, (isOpen) => {
  if (isOpen && props.projectKey) {
    projectKey.value = props.projectKey;
    projectName.value = props.projectName || '';
    cloneCount.value = props.initialCount || 1;
  }
});

const open = (config: CloneProjectConfig) => {
  projectKey.value = config.projectKey;
  projectName.value = config.projectName;
  cloneCount.value = 1;
  showDialog.value = true;
};

const close = () => {
  showDialog.value = false;
  projectKey.value = '';
  projectName.value = '';
  cloneCount.value = 1;
};

const confirm = async () => {
  await nextTick();
  const isValid = await cloneForm.value?.validate();
  if (isValid?.valid) {
    doCopy({
      projectKey: projectKey.value,
      projectName: projectName.value,
      count: cloneCount.value,
    });
    close();
  }
};

const doCopy = async ({projectKey, count}: CloneProjectConfig) => {
  try {
    const clonePromises = [];
    for (let i = 0; i < count; i++) {
      clonePromises.push(projectService.cloneProject(projectKey));
    }

    const responses = await Promise.all(clonePromises);
    emit('reload');
    if (count === 1) {
      infoSnackbar(t('DIALOG_project_copy_success'));
      if (responses[0]?.id) {
        openProjectUrlByKey(responses[0].id, router);
      }
    } else {
      infoSnackbar(t('DIALOG_project_copy_multiple_success', {count: count}));
    }
  } catch (error) {
    infoSnackbar(t('DIALOG_project_copy_error'));
  }
};
</script>

<template>
  <v-dialog v-model="showDialog" content-class="medium" scrollable width="500">
    <v-form ref="cloneForm" @submit.prevent="confirm">
      <DialogLayout :config="dialogConfig" @close="close" @secondary-action="close" @primary-action="confirm">
        <Stack>
          <p>
            {{ t('DLG_CLONE_PROJECT_DESCRIPTION_1') }}
            <span>
              <q>{{ projectName }}</q>
            </span>
          </p>

          <p class="pt-3">
            {{ t('DLG_CLONE_PROJECT_DESCRIPTION_2') }}
          </p>
          <div class="errorBorder">
            <v-text-field
              v-model.number="cloneCount"
              :label="t('DLG_CLONE_PROJECT_COUNT_LABEL')"
              type="number"
              min="1"
              max="10"
              :rules="countRules"
              persistent-hint
              variant="outlined"
              density="compact"
              class="required"
              autofocus
              hide-details
              @input="validateCount"></v-text-field>
          </div>
        </Stack>
      </DialogLayout>
    </v-form>
  </v-dialog>
</template>
