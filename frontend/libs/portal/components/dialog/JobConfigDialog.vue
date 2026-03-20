<script lang="ts" setup>
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';
import JsonViewer3 from 'vue-json-viewer';

const {t} = useI18n();
const isVisible = ref(false);
const config = ref<object>({});

const open = (newConf: string) => {
  try {
    config.value = JSON.parse(newConf);
  } catch {
    config.value = {string: newConf};
  }
  isVisible.value = true;
  isVisible.value = true;
};

const close = () => {
  isVisible.value = false;
};

defineExpose({open, close});
</script>

<template>
  <v-dialog v-model="isVisible" content-class="small" width="800" max-width="500">
    <DialogLayout
      :config="{
        title: t('JOB_CONFIG'),
      }"
      @close="close">
      <Stack>
        <JsonViewer3 :value="config" theme="jv-dark" sort />
      </Stack>
    </DialogLayout>
  </v-dialog>
</template>
