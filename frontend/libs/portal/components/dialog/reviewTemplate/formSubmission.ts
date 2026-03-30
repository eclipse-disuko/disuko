// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ReviewTemplate} from '@disclosure-portal/model/ReviewTemplate';
import adminService from '@disclosure-portal/services/admin';
import useSnackbar from '@shared/composables/useSnackbar';
import {AxiosError} from 'axios';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

export function useFormSubmission() {
  const isLoading = ref(false);
  const errorMessage = ref<string | null>(null);
  const {t} = useI18n();

  const {info} = useSnackbar();
  const submitForm = async (formData: ReviewTemplate, mode: 'create' | 'edit') => {
    isLoading.value = true;
    errorMessage.value = null;

    try {
      let response;
      if (mode === 'create') {
        response = await adminService.createReviewTemplate(formData);
        info(t(`DIALOG_REVIEW_TEMPLATE_${mode.toUpperCase()}_SUCCESS`));
      } else if (mode === 'edit') {
        response = await adminService.editReviewTemplate(formData);
        info(t(`DIALOG_REVIEW_TEMPLATE_${mode.toUpperCase()}_SUCCESS`));
      }

      return response?.data;
    } catch (error) {
      const axiosError = error as AxiosError;

      errorMessage.value = axiosError.message;
      throw error;
    } finally {
      isLoading.value = false;
    }
  };

  return {
    isLoading,
    errorMessage,
    submitForm,
  };
}
