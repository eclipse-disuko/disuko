<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <v-snackbar v-model="visible" :timeout="timeout" bottom color="brand">
    <div class="text-center">{{ message }}</div>
  </v-snackbar>
</template>

<script setup lang="ts">
import eventBus from '@shared/utils/eventbus';
import {onMounted, onUnmounted, ref} from 'vue';

const visible = ref(false);
const message = ref('');
const timeout = ref(3000);
const level = ref('info');

const showSnackbar = ({message: msg, timeout: time = 3000, level: lvl = 'info'}) => {
  message.value = msg;
  timeout.value = time;
  level.value = lvl;
  visible.value = true;
};

onMounted(() => {
  eventBus.on('show-snackbar', showSnackbar);
});

onUnmounted(() => {
  eventBus.off('show-snackbar', showSnackbar);
});
</script>

<style scoped>
.v-snackbar__wrapper.borderSnackbar {
  border: 1px solid red !important; /* Beispiel: Füge eine graue Umrandung hinzu */
}
</style>
