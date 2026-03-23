import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useTableActionSliderStore = defineStore('tableActionSlider', () => {
  const slideInTimeout = ref<ReturnType<typeof setTimeout> | null>(null);

  return {
    slideInTimeout,
  };
});
