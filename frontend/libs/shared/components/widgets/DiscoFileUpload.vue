// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts">
import {useUpload} from '@disclosure-portal/composables/useUpload';
import SchemaPostRequest from '@disclosure-portal/model/SchemaPostRequest';
import {defineComponent, nextTick} from 'vue';

const {uploadFormDataFile} = useUpload();

export default defineComponent({
  name: 'DiscoFileUpload',
  props: {
    acceptTypes: {
      type: String,
      required: true,
    },
    uploadTargetUrl: {
      type: String,
      required: true,
    },
    multiple: {
      type: Boolean,
      default: false,
      required: false,
    },
    directUpload: {
      type: Boolean,
      default: true,
      required: false,
    },
    schema: {
      type: SchemaPostRequest,
      required: false,
    },
  },
  emits: ['filesChanged', 'reqFinished', 'reqFailed', 'reqProgress'],
  setup(props, {emit}) {
    const onInputChange = (e: Event) => {
      const input = e.target as HTMLInputElement;
      if (input.files?.length && !/safari/i.test(navigator.userAgent)) {
        input.type = 'file';
      }
      emit('filesChanged', input.files);
      if (input.files && input.files.length > 0 && !props.multiple && props.directUpload) {
        doUpload();
      }
    };

    const uploadClick = () => {
      nextTick(() => {
        const el = document.getElementById('inId') as HTMLElement;
        if (el.onclick) {
          el.onclick(new MouseEvent('click', {}));
        } else if (el.click) {
          el.click();
        }
      });
    };

    const clearFiles = () => {
      const input = document.getElementById('inId') as HTMLInputElement;
      input.value = null;
      emit('filesChanged', input.files);
    };

    const doUpload = () => {
      const input = document.getElementById('inId') as HTMLInputElement;
      if (!input.files) {
        return;
      }
      for (let i = 0; i < input.files.length; i++) {
        uploadFile(input.files[i]);
      }
    };

    const uploadFile = (file: File) => {
      uploadFormDataFile({
        uploadUrl: props.uploadTargetUrl,
        file,
        onUploadProgress: (progressEvent) => {
          if (progressEvent.total) {
            const progressValue = Math.round((progressEvent.loaded / progressEvent.total) * 100);
            emit('reqProgress', file, progressValue);
          }
        },
        schema: props.schema,
      })
        .then((response) => {
          emit('reqFinished', file, response.data);
        })
        .catch(() => {
          emit('reqFailed');
        })
        .finally(() => {
          clearFiles();
        });
    };

    return {
      props,
      doUpload,
      onInputChange,
      uploadClick,
      clearFiles,
    };
  },
});
</script>
<template>
  <div>
    <input
      type="file"
      id="inId"
      :multiple="multiple"
      :accept="acceptTypes"
      @change="onInputChange"
      style="opacity: 0; position: absolute; z-index: -12121212"
    />
  </div>
</template>
