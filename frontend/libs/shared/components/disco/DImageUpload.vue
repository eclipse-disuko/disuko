<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <div class="w-full">
    <!-- Hidden file input -->
    <input ref="fileInput" type="file" accept="image/*" style="display: none" @change="handleFileSelect" />

    <!-- Drag and Drop Area -->
    <div
      class="border-2 border-dashed border-gray-300 rounded-lg p-5 text-center cursor-pointer transition-all duration-300 min-h-[150px] flex items-center justify-center relative"
      :class="{'border-blue-500 border-solid': isDragOver}"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
      @click="openFileDialog">
      <div v-if="!modelValue" class="w-full">
        <v-icon size="48" color="grey-lighten-1">mdi-cloud-upload</v-icon>
        <p class="text-body-2 text-grey-darken-1 mt-2 mb-1">{{ uploadText || 'Drag and drop an image here' }}</p>
        <p class="text-caption text-grey">{{ clickText || 'or click to select' }}</p>
        <v-btn variant="outlined" color="primary" size="small" class="mt-2" @click.stop="openFileDialog">
          <v-icon left>mdi-upload</v-icon>
          {{ buttonText || 'Select Image' }}
        </v-btn>
      </div>

      <!-- Image Preview -->
      <div v-else class="w-full">
        <img :src="modelValue" alt="Preview" class="max-w-full max-h-[200px] rounded shadow-lg mb-2" />
        <div class="flex justify-center gap-2">
          <v-btn
            icon="mdi-pencil"
            size="small"
            variant="tonal"
            color="primary"
            @click.stop="openFileDialog"
            :disabled="disabled">
          </v-btn>
          <v-btn
            icon="mdi-delete"
            size="small"
            variant="tonal"
            color="error"
            @click.stop="removeImage"
            :disabled="disabled">
          </v-btn>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue';

interface Props {
  modelValue?: string | null;
  maxSizeBytes?: number;
  uploadText?: string;
  clickText?: string;
  buttonText?: string;
  disabled?: boolean;
}

interface Emits {
  (e: 'update:modelValue', value: string | null): void;
  (e: 'error', message: string): void;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: null,
  maxSizeBytes: 5 * 1024 * 1024, // 5MB default
  disabled: false,
});

const emit = defineEmits<Emits>();

const fileInput = ref<HTMLInputElement | null>(null);
const isDragOver = ref(false);

const convertImageToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const result = reader.result as string;
      resolve(result);
    };
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });
};

const handleImageUpload = async (file: File) => {
  if (!file.type.startsWith('image/')) {
    emit('error', 'Please select a valid image file');
    return;
  }

  if (file.size > props.maxSizeBytes) {
    const sizeMB = Math.round(props.maxSizeBytes / (1024 * 1024));
    emit('error', `Image size should be less than ${sizeMB}MB`);
    return;
  }

  try {
    const base64 = await convertImageToBase64(file);
    emit('update:modelValue', base64);
  } catch (error) {
    console.error('Error converting image to base64:', error);
    emit('error', 'Failed to process image');
  }
};

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (file) {
    handleImageUpload(file);
  }
};

const openFileDialog = () => {
  if (!props.disabled) {
    fileInput.value?.click();
  }
};

const handleDragOver = (event: DragEvent) => {
  if (props.disabled) return;
  event.preventDefault();
  isDragOver.value = true;
};

const handleDragLeave = (event: DragEvent) => {
  if (props.disabled) return;
  event.preventDefault();
  isDragOver.value = false;
};

const handleDrop = async (event: DragEvent) => {
  if (props.disabled) return;
  event.preventDefault();
  isDragOver.value = false;

  const files = event.dataTransfer?.files;
  if (files && files.length > 0) {
    const file = files[0];
    await handleImageUpload(file);
  }
};

const removeImage = () => {
  if (!props.disabled) {
    emit('update:modelValue', null);
    if (fileInput.value) {
      fileInput.value.value = '';
    }
  }
};
</script>
