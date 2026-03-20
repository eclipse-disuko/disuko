<script setup lang="ts">
const props = defineProps<{
  text: string;
  icon?: string;
  tooltip?: string;
  disabled?: boolean;
}>();

const emit = defineEmits<{
  click: [event: Event];
}>();

const handleClick = (event: Event) => {
  if (props.disabled) {
    event.preventDefault();
    event.stopPropagation();
    return;
  }
  emit('click', event);
};
</script>

<template>
  <v-list-item @click="handleClick">
    <Stack direction="row" :class="{disabledText: disabled}">
      <v-icon style="width: 30px" :color="disabled ? '' : 'primary'">{{ icon }}</v-icon>
      <span class="d-subtitle">{{ text }}</span>
    </Stack>
    <Tooltip v-if="tooltip || $slots.tooltip" location="left">
      <span v-if="tooltip">{{ tooltip }}</span>
      <slot v-if="$slots.tooltip" name="tooltip"></slot>
    </Tooltip>
  </v-list-item>
</template>
