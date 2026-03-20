<script setup lang="ts">
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import DIconButton from '@shared/components/disco/DIconButton.vue';
import {DataTableHeader} from '@shared/types/table';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps<{
  allHeaders: DataTableHeader[];
}>();
const selectedHeaders = defineModel<number[]>({required: true});

const {t} = useI18n();
const menuSettings = ref(false);

const selectableHeaders = computed(() => {
  return [...Array(props.allHeaders.length).keys()];
});
</script>

<template>
  <v-menu offset-y :close-on-content-click="false" v-model="menuSettings">
    <template v-slot:activator="{props}">
      <DIconButton
        :parentProps="props"
        icon="mdi-cog"
        :hint="t('TT_TABLE_SETTINGS')"
        color="primary"
        class="pl-0.5 z-10" />
    </template>
    <div class="bg-background" style="width: 280px">
      <div class="flex justify-end m-1">
        <DCloseButton @click="menuSettings = false" />
      </div>
      <v-select
        v-model="selectedHeaders"
        :items="selectableHeaders"
        class="mx-2 pa-2 dp-select"
        :label="t('LABEL_SELECT_COL')"
        multiple
        return-object
        v-bind:menu-props="{location: 'bottom'}"
        variant="outlined"
        density="compact">
        <template v-slot:item="{props, item}">
          <v-list-item v-bind="props" class="py-0 px-2">
            <template v-slot:prepend="{isSelected}">
              <v-checkbox hide-details :model-value="isSelected" />
            </template>
            <template v-slot:title>
              <span class="pFilterEntry"> {{ allHeaders[item.value].title }}</span>
            </template>
          </v-list-item>
        </template>
        <template v-slot:selection="{item, index}">
          <div v-if="index === 0" class="d-flex align-center">
            <span class="pFilterEntry">{{ allHeaders[item.value].title }}</span>
          </div>
          <span v-if="index === 1" class="pAdditionalFilter"> +{{ selectedHeaders.length - 1 }} others </span>
        </template>
      </v-select>
    </div>
  </v-menu>
</template>
