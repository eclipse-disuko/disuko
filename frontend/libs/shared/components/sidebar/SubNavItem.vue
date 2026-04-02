<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip">
    <template #activator="{props}">
      <v-list-item
        class="pl-8"
        :to="item.path"
        @mouseover="hover = true"
        @mouseleave="hover = false"
        :title="t(item.title)"
        v-bind="props">
        <template v-slot:prepend>
          <v-icon>{{ item.iconName }}</v-icon>
        </template>
        <template v-slot:title>
          <span>{{ t(item.title) }}</span>
        </template>
      </v-list-item>
    </template>
    <span>{{ t(item.title) }}</span>
  </v-tooltip>
</template>

<script lang="ts">
import INavItem from '@disclosure-portal/model/INavItem';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {defineComponent, ref} from 'vue';
import {useI18n} from 'vue-i18n';

export default defineComponent({
  props: {
    item: {
      type: Object as () => INavItem,
      required: true,
    },
  },
  setup() {
    const hover = ref(false);
    const {t} = useI18n();
    return {
      t,
      hover,
      TOOLTIP_OPEN_DELAY_IN_MS,
    };
  },
});
</script>
