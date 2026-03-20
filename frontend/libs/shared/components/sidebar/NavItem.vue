<template>
  <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip" v-if="item">
    <template #activator="{props}">
      <v-list-item
        v-if="item?.externalPath"
        :href="item?.externalPath"
        @mouseover="hover = true"
        @mouseleave="hover = false"
        :title="t(item.title) + '121'"
        v-bind="props"
      >
        <template v-slot:prepend>
          <v-icon v-if="hover || item.active">{{ item.iconName }}</v-icon>
          <v-icon v-else>{{ `${item.iconName}-outline` }}</v-icon>
        </template>
        <template v-slot:title>
          <span>{{ item.tooltip ? t(item.tooltip) : '' }}</span>
        </template>
      </v-list-item>
      <v-list-item
        v-else
        :to="item.path"
        :active="item.active"
        @mouseover="hover = true"
        @mouseleave="hover = false"
        :title="t(item.title)"
        v-bind="props"
      >
        <template v-slot:prepend>
          <v-icon v-if="hover || item.active">{{ item.iconName }}</v-icon>
          <v-icon v-else>{{ `${item.iconName}-outline` }}</v-icon>
        </template>
        <template v-slot:title>
          <span>{{ item.tooltip ? t(item.tooltip) : '' }}</span>
        </template>
      </v-list-item>
    </template>
    <span>{{ item.tooltip ? t(item.tooltip) : '' }}</span>
  </v-tooltip>
</template>

<script setup lang="ts">
import INavItem from '@disclosure-portal/model/INavItem';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

defineProps<{
  item: INavItem;
}>();
const hover = ref(false);
const {t} = useI18n();
</script>
