<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import {SelfBuildingSquareSpinner} from 'epic-spinners';
import {VProgressCircular} from 'vuetify/components';

const idle = useIdleStore();
</script>

<template>
  <Transition>
    <div v-if="idle.showIdle" class="idleFullContainer" data-testid="idle">
      <v-layout row>
        <div class="idle top-16">
          <Stack align="center">
            <SelfBuildingSquareSpinner v-if="idle.progress === -1" />
            <v-progress-circular v-else :model-value="idle.progress" />

            <div v-if="idle.idleMessage" class="msg pt-4 text-nowrap">
              {{ idle.idleMessage }}
              <span v-if="idle.progress !== -1">{{ idle.progress }}{{ idle.progressUnit }}</span>
            </div>
          </Stack>
        </div>
      </v-layout>
    </div>
  </Transition>
</template>

<style scoped lang="scss">
.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}

:global(.idleFullContainer) {
  position: absolute;
  left: 0;
  right: 0;
  z-index: 100000;
  top: 0;
  bottom: 0;
}
:global(.v-theme--light .idleFullContainer) {
  background-color: rgba(128, 128, 128, 0.8);
}

:global(.v-theme--dark .idleFullContainer) {
  background-color: rgba(45, 45, 45, 0.8);
}

.idle {
  top: 40%;
  margin-left: auto;
  margin-right: auto;
  display: flex;
  align-items: center;
  position: relative;
  margin-top: 14%;
  z-index: 100000;
  min-height: 240px;
}
</style>
