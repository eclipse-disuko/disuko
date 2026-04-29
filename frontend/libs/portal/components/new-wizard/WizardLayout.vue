<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import Stepper from '@disclosure-portal/components/new-wizard/Stepper.vue';
import {useIdleStore} from '@shared/stores/idle.store';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {openProjectUrlByKey} from '@disclosure-portal/utils/url';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const {t} = useI18n();
const wizardStore = useWizardStore();
const snack = useSnackbar();
const idle = useIdleStore();
const router = useRouter();

const lastButtonText = computed(() => {
  if (wizardStore.mode === 'edit') {
    return t('NP_DIALOG_BTN_SAVE');
  } else {
    return t('NP_DIALOG_BTN_CREATE');
  }
});

const title = computed(() => {
  if (wizardStore.mode === 'edit') {
    if (wizardStore.project?.isDummy) {
      return t('WIZARD_title_edit_dummy');
    } else {
      return t('WIZARD_title_edit');
    }
  }
  if (wizardStore.project.isGroup) {
    return t('Title_New_Group');
  }
  if (wizardStore.project?.isDummy) {
    return t('WIZARD_title_dummy');
  } else {
    return t('WIZARD_title');
  }
});

const onCreateProject = async () => {
  idle.show(t('PROJECT_IS_CREATING'));

  try {
    const project = await wizardStore.createProject();
    wizardStore.close();
    if (project?._key) {
      await openProjectUrlByKey(project._key, router);
    }
  } finally {
    idle.hide();
  }
};

const onEditProject = async () => {
  try {
    await wizardStore.updateProject();
    wizardStore.close();
  } finally {
  }
};
</script>

<template>
  <v-card
    class="bg-[rgb(var(--v-theme-backgroundColor))] p-12"
    :class="{'barrier-tape-background': wizardStore.project?.isDummy}">
    <Stack>
      <Stack direction="row" align="center">
        <div class="w-full">
          <h4 class="text-h5 text-center">
            {{ title }}
          </h4>
        </div>
        <v-spacer></v-spacer>
        <DCloseButton :disabled="wizardStore.previewLoading" class="-mr-4" @click="wizardStore.close()" />
      </Stack>
      <Stepper>
        <Step
          v-for="step in wizardStore.steps"
          :key="step.id"
          :step="step"
          :current-index="wizardStore.currentStep.index"
          @click="wizardStore.selectStep(step.index)" />
      </Stepper>
    </Stack>

    <v-card-text class="p-0">
      <section v-if="wizardStore.mode === 'edit' && wizardStore.project.hasDeniedDecisions">
        <v-alert color="warning" type="warning" class="my-2">
          {{ t('PROJECT_HAS_DENIED_DECISIONS_TO_BE_CANCELLED') }}
        </v-alert>
      </section>
      <slot></slot>
    </v-card-text>

    <Stack direction="row" class="pt-2">
      <div
        v-if="wizardStore.currentStep.index === 0 && wizardStore.mode !== 'edit' && !wizardStore.project.isGroup"
        class="flex h-7">
        <v-checkbox-btn
          v-model="wizardStore.project.isDummy"
          color="primary"
          :label="t('DUMMY_MAKE')"
          density="compact"></v-checkbox-btn>
        <v-btn size="x-small" flat icon variant="text" class="cursor-default">
          <v-icon color="primary">mdi mdi-help</v-icon>
          <Tooltip :text="t('HOME_EXPLAIN_NEW_DUMMY_PROJECT')"></Tooltip>
        </v-btn>
      </div>
      <DCActionButton
        v-if="wizardStore.currentStep.index > 0"
        is-dialog-button
        size="small"
        variant="text"
        @click="wizardStore.previousStep()"
        :disabled="wizardStore.previewLoading"
        :text="t('BTN_BACK')" />
      <v-spacer></v-spacer>
      <DCActionButton
        v-if="!wizardStore.isFinalStep"
        is-dialog-button
        size="small"
        @click="wizardStore.nextStep()"
        variant="flat"
        :color="!wizardStore.isCurrentStepCompleted ? 'gray' : ''"
        :disabled="!wizardStore.isCurrentStepCompleted"
        :text="t('BTN_NEXT')" />
      <DCActionButton
        v-else
        is-dialog-button
        size="small"
        variant="flat"
        @click="wizardStore.mode === 'create' ? onCreateProject() : onEditProject()"
        :color="!wizardStore.allStepsCompleted || wizardStore.previewLoading ? 'gray' : ''"
        :disabled="!wizardStore.allStepsCompleted || wizardStore.previewLoading"
        :loading="wizardStore.previewLoading"
        :text="lastButtonText" />
    </Stack>
  </v-card>
</template>

<style scoped lang="scss">
.barrier-tape-background {
  &.v-theme--dark {
    background-image:
      linear-gradient(
        to bottom,
        transparent 0%,
        rgb(var(--v-theme-backgroundColor)) 1.5%,
        rgb(var(--v-theme-backgroundColor)) 98.5%,
        transparent 100%
      ),
      repeating-linear-gradient(
        -45deg,
        transparent,
        transparent 7px,
        rgb(var(--v-theme-chartYellow)) 8px,
        rgb(var(--v-theme-chartYellow)) 12px
      );
  }
  &.v-theme--light {
    background-image:
      linear-gradient(
        to bottom,
        transparent 0%,
        rgb(var(--v-theme-surface)) 2.5%,
        rgb(var(--v-theme-surface)) 97.5%,
        transparent 100%
      ),
      repeating-linear-gradient(
        -45deg,
        transparent,
        transparent 7px,
        rgb(var(--v-theme-chartYellow)) 8px,
        rgb(var(--v-theme-chartYellow)) 12px
      );
  }
}
</style>
