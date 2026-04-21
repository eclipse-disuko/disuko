<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <Stack class="mt-4 gap-4">
    <v-form ref="form" v-if="!readonly">
      <v-textarea
        rows="3"
        no-resize
        density="compact"
        :label="t('TAD_COMMENT')"
        v-model="currentComment"
        :counter="500"
        variant="outlined"
        :rules="activeRules.comment"
        hide-details="auto"></v-textarea>
      <DCActionButton
        :hint="t('BTN_COMMENT')"
        icon="mdi-comment"
        :text="t('BTN_COMMENT')"
        @click="comment"
        :disabled="isDeprecated || disableClick" />
    </v-form>

    <v-data-table
      v-if="events"
      :headers="headers"
      fixed-header
      hide-default-footer
      sort-asc
      :sort-by="[{key: 'created', order: 'desc'}]"
      class="striped-table"
      item-key="key"
      :items="events"
      :item-class="getCssClassForTableRow">
      <template v-slot:item.authorFullName="{item}"> {{ item.authorFullName }} ({{ item.author }})</template>

      <template v-slot:item.created="{item}">
        <DDateCellWithTooltip :show-time="true" :value="item.created" />
      </template>
      <template v-slot:item.content="{item}">
        <span v-if="item.type === EventType.COMMENT">
          {{ item.content as Comment }}
        </span>
        <span v-else-if="item.type === EventType.CHANGED_LEVEL && item.content">
          <i
            >{{ t('RR_EVENT_' + item.type) }} {{ t('DIALOG_remark_from') }}
            {{ t('REMARK_LEVEL_' + (item.content as LevelChange)?.before) }} {{ t('DIALOG_remark_to') }}
            {{ t('REMARK_LEVEL_' + (item.content as LevelChange)?.after) }}
          </i>
        </span>
        <span v-else-if="item.type === EventType.CHANGED_TITLE && item.content">
          <i
            >{{ t('RR_EVENT_' + item.type) }} {{ t('DIALOG_remark_from') }} {{ (item.content as TitleChange)?.before }}
            {{ t('DIALOG_remark_to') }} {{ (item.content as TitleChange)?.after }}
          </i>
        </span>
        <span v-else-if="item.type === EventType.CHANGED_DESCRIPTION && item.content">
          <i
            >{{ t('RR_EVENT_' + item.type) }} {{ t('DIALOG_remark_from') }}
            {{ (item.content as DescriptionChange)?.before }} {{ t('DIALOG_remark_to') }}
            {{ (item.content as DescriptionChange)?.after }}
          </i>
        </span>
        <span v-else>
          <i>{{ t('RR_EVENT_' + item.type) }}</i>
        </span>
      </template>
    </v-data-table>
  </Stack>
</template>

<script lang="ts" setup>
import {Comment, DescriptionChange, Event, EventType, LevelChange, TitleChange} from '@disclosure-portal/model/Quality';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import useRules from '@disclosure-portal/utils/Rules';
import {getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import {DataTableHeader} from '@shared/types/table';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const {t} = useI18n();
const rules = useRules();
defineProps<{
  reviewId: string;
  events: Event[];
  readonly: boolean;
  disableClick?: boolean;
}>();

const activeRules = {
  comment: rules.minMax(t('ATTR_COMMENT'), 3, 500, false),
};

const form = ref<VForm | null>(null);
const isDeprecated = computed(() => useProjectStore().currentProject!.isDeprecated);
const headers: DataTableHeader[] = [
  {
    title: t('COL_CREATED'),
    width: 150,
    align: 'start',
    class: 'tableHeaderCell',
    value: 'created',
    sortable: true,
  },
  {
    title: t('COL_USER'),
    width: 230,
    align: 'start',
    class: 'tableHeaderCell',
    value: 'authorFullName',
  },
  {
    title: t('COL_UPDATE'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'content',
  },
];
const currentComment = ref('');

const emit = defineEmits<{
  (e: 'comment', content: string): void;
}>();

const comment = async () => {
  if (!(await form.value?.validate())?.valid) {
    return;
  }
  emit('comment', currentComment.value);
  currentComment.value = '';
  form.value?.resetValidation();
};
</script>
