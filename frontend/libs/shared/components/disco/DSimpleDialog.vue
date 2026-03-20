<template>
  <v-dialog
    v-model="localValue"
    @input="$emit('update:modelValue', localValue)"
    @click:outside="$emit('update:modelValue', false)"
    width="auto"
    max-width="1000px"
    content-class="d-simple-dialog"
    scrollable
  >
    <v-card class="pa-8">
      <v-card-title v-if="title">
        <v-row>
          <v-col cols="10">
            <h1>{{ title }}</h1>
          </v-col>
          <v-col cols="2" class="d-flex justify-end">
            <DCloseButton @click="$emit('update:modelValue', false)"></DCloseButton>
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <slot />
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import {defineComponent, ref, watch} from 'vue';

export default defineComponent({
  components: {
    DCloseButton,
  },
  props: {
    title: {
      type: String,
      required: true,
    },
    modelValue: {
      type: Boolean,
      required: true,
    },
  },
  setup(props) {
    const localValue = ref(false);

    watch(
      () => props.modelValue,
      (newValue) => {
        localValue.value = newValue;
      },
    );

    return {
      localValue,
    };
  },
});
</script>
