<template>
  <v-snackbar v-model="visible" :timeout="timeout" bottom color="mbti">
    <div class="text-center">{{ message }}</div>
  </v-snackbar>
</template>

<script setup lang="ts">
import eventBus from '@disclosure-portal/utils/eventbus';
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
