<script setup lang="ts">
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

const emit = defineEmits(['close', 'openSettings']);
const show = ref(false);
const {t} = useI18n();

const open = () => {
  show.value = true;
};

const close = () => {
  show.value = false;
};
defineExpose({open, close});

const openSettingsDialog = () => {
  emit('openSettings');
  close();
};
</script>
<template>
  <v-dialog v-model="show" width="500" scrollable max-height="500">
    <v-card class="pa-8">
      <v-card-title>
        <v-row no-gutters>
          <v-col class="d-flex justify-center" align-self="center" cols="2">
            <v-icon color="warning">mdi-alert</v-icon>
          </v-col>
          <v-col cols="10">
            <span class="text-h6 text-wrap">
              {{ t('SETTINGS_REQUIRED_MESSAGE') }}
            </span>
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <span class="paddingError d-subtitle-1 mr-2">
          {{ t('FILL_SETTINGS_PART_1') }}
          <span>
            <button
              style="color: rgb(var(--v-theme-primary)); text-decoration: underline"
              @click="openSettingsDialog()"
            >
              {{ t('SETTINGS_DIALOG_FROM_ERROR') }}
            </button>
          </span>
          {{ t('FILL_SETTINGS_PART_2') }}
        </span>
      </v-card-text>
      <v-card-actions class="justify-center">
        <DCActionButton isDialogButton size="small" variant="flat" @click="close" :text="t('BTN_OK')" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
