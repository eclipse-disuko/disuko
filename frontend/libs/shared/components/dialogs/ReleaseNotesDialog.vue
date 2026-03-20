<script lang="ts" setup>
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

interface Props {
  releaseNotes: string;
}

const props = defineProps<Props>();

const isDialogOpen = ref(false);
const releaseNotesDialogText = ref(props.releaseNotes);
const {t} = useI18n();

const showDialog = () => {
  isDialogOpen.value = true;
};

const closeDialog = () => {
  isDialogOpen.value = false;
};
</script>

<template>
  <slot :showDialog="showDialog">
    <v-btn text="Replace me" size="small" color="primary" @click.stop="showDialog"></v-btn>
  </slot>

  <v-dialog v-model="isDialogOpen" width="auto" max-width="700px" scrim="#010101" scrollable>
    <v-card class="pa-8" variant="elevated">
      <v-card-title>
        <v-row>
          <v-col cols="10" class="d-flex align-center">
            <span class="text-h5">{{ t('RELEASE_NOTES') }}</span>
          </v-col>
          <v-col cols="2" class="text-right px-0">
            <DCloseButton @click="closeDialog" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" xs="12">
            <Markdown :text="releaseNotesDialogText" :id="new Date().getTime().toString()" />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.v-card-text {
  background-color: rgb(var(--v-theme-markdownBackground));
}
</style>
