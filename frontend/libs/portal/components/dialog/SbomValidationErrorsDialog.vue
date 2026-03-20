<template>
  <v-form>
    <v-dialog v-model="show" content-class="large" scrollable width="700" persistent max-height="500">
      <v-card class="pa-8">
        <v-card-title>
          <div class="d-flex flex-row">
            <span class="text-h5">{{ t('VALIDATE_SCHEMA') }}</span>
            <v-spacer></v-spacer>
            <DCloseButton @click="closeDialog" />
          </div>
        </v-card-title>

        <v-card-text class="position-relative scrollable-errors">
          <v-list class="pa-0 noneBorder">
            <v-list-item v-for="(error, index) in formattedErrors" :key="index">
              {{ error }}
            </v-list-item>
          </v-list>
          <DCopyClipboardButton
            :tableButton="true"
            class="mr-4 position-absolute top-0 right-0"
            :hint="t('TT_COPY_TO_CLIPBOARD')"
            :content="rawErrors"
          />
        </v-card-text>
        <v-card-actions class="justify-end mr-7">
          <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip">
            <template v-slot:activator="{props}">
              <v-btn v-bind="props" outlined class="mr-2 card-border" @click="openHelp">
                <v-icon color="primary">mdi mdi-help</v-icon>
                {{ t('HELP_BTN') }}
              </v-btn>
            </template>
            <span>{{ t('HELP_SHOW') }}</span>
          </v-tooltip>
          <DCActionButton size="small" variant="flat" @click="closeDialog" :text="t('BTN_OK')" />
        </v-card-actions>
      </v-card>
    </v-dialog>
    <DSimpleDialog v-if="helpText" v-model="helpOpen" :title="t('HELP')">
      <Markdown :text="helpText" />
    </DSimpleDialog>
  </v-form>
</template>

<script setup lang="ts">
import Markdown from '@shared/components/Markdown.vue';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import DCopyClipboardButton from '@shared/components/disco/DCopyClipboardButton.vue';
import DSimpleDialog from '@shared/components/disco/DSimpleDialog.vue';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';

import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const show = ref(false);
const helpText = ref<string | null>(null);
const helpOpen = ref(false);
const rawErrors = ref('');

const formattedErrors = computed(() => (rawErrors.value ? rawErrors.value.split('\n').filter(Boolean) : []));

const open = (errors: string, help: string | null) => {
  show.value = true;
  rawErrors.value = errors;
  helpText.value = help;
};

const closeDialog = () => {
  show.value = false;
};

const openHelp = () => {
  helpOpen.value = true;
};
defineExpose({open});
</script>

<style scoped>
.scrollable-errors {
  max-height: 400px; /* Adjust the height as needed */
  overflow-y: auto;
}
</style>
