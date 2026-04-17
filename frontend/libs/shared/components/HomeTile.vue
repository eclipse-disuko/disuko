<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

interface Props {
  type: 'action' | 'navigation';
  title: string;
  icon: string;
  description?: string;
  url?: string;
  showBadge?: boolean;
  badgeContent?: number | string;
  showExpand?: boolean;
}

withDefaults(defineProps<Props>(), {
  showExpand: true,
});

const emit = defineEmits<{
  click: [];
}>();

const {t} = useI18n();
const expand = ref(false);
</script>

<template>
  <v-card
    v-if="type === 'action'"
    min-width="368"
    max-width="368"
    color="primary"
    variant="outlined"
    link
    @click="emit('click')">
    <v-card-item class="text-center">
      <span class="text-h4 font-weight-light text-font">{{ title }}</span>
    </v-card-item>

    <v-card-text class="py-0">
      <Stack justify="center" align="center">
        <v-icon class="font-weight-light" :icon="icon" size="88"></v-icon>
      </Stack>
    </v-card-text>

    <v-divider></v-divider>

    <v-card-actions v-if="showExpand && description">
      <v-btn class="text-button-more" @click.stop="expand = !expand">
        <v-icon color="primary">
          {{ expand ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
        </v-icon>
        {{ expand ? t('LESS') : t('MORE') }}
      </v-btn>
    </v-card-actions>
    <v-expand-transition>
      <div v-if="expand && description" class="pa-4 text-font">
        {{ description }}
      </div>
    </v-expand-transition>
  </v-card>

  <v-btn
    v-else-if="type === 'navigation'"
    :href="url ? '#' + url : undefined"
    size="x-large"
    variant="text"
    @click="emit('click')"
    class="align-center flex justify-start">
    <v-icon color="primary" size="x-large">{{ icon }}</v-icon>
    <span class="text-h4 font-weight-light">{{ title }}</span>
    <v-badge
      v-if="showBadge && badgeContent"
      :content="badgeContent"
      color="primary"
      overlap
      floating
      class="custom-badge"
      size="x-large"></v-badge>
  </v-btn>
</template>

<style scoped lang="scss">
.custom-badge :deep(.v-badge__badge) {
  font-size: 1.5rem;
  height: 2rem;
  min-width: 2rem;
  line-height: 2rem;
  margin-left: 10px;
}
</style>
