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
