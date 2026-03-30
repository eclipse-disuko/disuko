// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts"></script>

<template>
  <v-btn plain size="small" variant="text" icon="mdi-dots-horizontal" color="primary" @click.stop class="size-10">
    <v-icon>mdi-dots-horizontal</v-icon>
    <v-menu activator="parent" open-on-hover>
      <v-list>
        <div class="flex flex-row p-1">
          <slot />
        </div>
      </v-list>
    </v-menu>
    <v-tooltip location="bottom" text="Actions"></v-tooltip>
  </v-btn>
</template>
