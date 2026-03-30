// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<template>
  <div ref="dataTableAsElement">
    <v-data-table
      fixed-header
      :headers="headers"
      :items="items"
      class="striped-table"
      :height="tableHeight"
      @click:row="onRowClick"
    >
      <template v-slot:item.actions="{item}">
        <v-icon :id="item.id">mdi-home</v-icon>
      </template>
      <template v-slot:item.status="{item}">
        <v-icon :id="item.id">mdi-plus</v-icon>
      </template>
    </v-data-table>
  </div>
</template>

<script lang="ts">
import useDimensions from '@disclosure-portal/composables/useDimensions';
import {PropType, defineComponent, nextTick, onMounted, onUnmounted, ref} from 'vue';
import {useRouter} from 'vue-router';

interface Item {
  id: number;
  title: string;
  key: string;
  value: string;
  price: number;
  quantity: number;
}

export default defineComponent({
  name: 'GridMock',
  props: {
    parentId: {
      type: String as PropType<string>,
      required: false,
    },
  },
  setup(props) {
    const tableHeight = ref(0);
    const dataTableAsElement = ref<HTMLElement | null>(null);
    const {calculateHeight} = useDimensions();
    const router = useRouter();
    const headers = ref([
      {title: 'Status', key: 'status'},
      {title: 'quantity', key: 'quantity'},
      {title: 'Price', key: 'price', value: (item: Item) => '' + parseFloat(item.price).toFixed(2)},
      {
        title: 'Full Name',
        key: 'fullName',
        value: (item: Item) => `${item.name.first} ${item.name.last}`,
      },
      {
        title: 'Actions',
        key: 'actions',
      },
    ]);

    const items = ref<Item[]>(
      Array.from({length: 100}, (v, k) => ({
        id: k + 1,
        name: {
          first: 'John',
          last: 'Doe ' + (k + 1),
        },
        quantity: 10,
        title: `Item ${k + 1}`,
        price: Math.random() * 100,
        key: `This is the content for item ${k + 1}.`,
        value: `Extra 1 for item ${k + 1}`,
      })),
    );

    const onRowClick = (event: Event, item: Item) => {
      router.push({
        path: '/dashboard/projects/1234/overview',
      });
    };

    const updateTableHeight = () => {
      nextTick(() => {
        tableHeight.value = calculateHeight(dataTableAsElement.value, true);
      });
    };

    onMounted(() => {
      updateTableHeight();
      window.addEventListener('resize', updateTableHeight);
    });
    onUnmounted(() => {
      window.removeEventListener('resize', updateTableHeight);
    });
    return {
      dataTableAsElement,
      tableHeight,
      headers,
      items,
      onRowClick,
    };
  },
});
</script>

<style scoped></style>
