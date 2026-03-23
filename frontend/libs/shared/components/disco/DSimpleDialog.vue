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
    <DialogLayout :config="dialogConfig" @close="$emit('update:modelValue', false)">
      <slot />
    </DialogLayout>
  </v-dialog>
</template>

<script lang="ts">
import DialogLayout, {type DialogLayoutConfig} from '@shared/layouts/DialogLayout.vue';
import {computed, defineComponent, ref, watch} from 'vue';

export default defineComponent({
  components: {
    DialogLayout,
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

    const dialogConfig = computed<DialogLayoutConfig>(() => ({
      title: props.title,
    }));

    return {
      localValue,
      dialogConfig,
    };
  },
});
</script>
