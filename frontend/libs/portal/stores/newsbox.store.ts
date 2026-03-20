// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {NewsboxItemCreateDto, NewsboxItem, NewsboxItems, default as Newsbox} from '@disclosure-portal/model/Newsbox';
import {UserLastSeenDto} from '@disclosure-portal/model/Users';
import {default as newsboxService} from '@disclosure-portal/services/newsbox.service';
import useSnackbar from '@shared/composables/useSnackbar';
import {defineStore} from 'pinia';
import {reactive, toRefs} from 'vue';
import {useI18n} from 'vue-i18n';

export const useNewsboxStore = defineStore('newsbox', () => {
  const {t} = useI18n();
  const {info} = useSnackbar();

  const state = reactive({
    newsItems: null as Newsbox | null,
    hasNewNewsboxItem: false,
    loading: false,
    showNewsbox: false,
    adminNewsItems: null as NewsboxItems | null,
    adminLoading: false,
  });

  const fetchItems = async () => {
    try {
      state.loading = true;
      const newsItems = await newsboxService.getNewsboxItems();
      state.newsItems = newsItems.data;
    } catch (error) {
      console.error('Error fetching newsbox items:', error);
    } finally {
      state.loading = false;
    }
  };

  const fetchItemsAdmin = async () => {
    try {
      state.adminLoading = true;
      const newsItems = await newsboxService.getAllNewsboxItems();
      state.adminNewsItems = { items: newsItems.data };
    } catch (error) {
      console.error('Error fetching newsbox items:', error);
    } finally {
      state.adminLoading = false;
    }
  };

  const createItemsAdmin = async (item: NewsboxItemCreateDto) => {
    try {
      state.adminLoading = true;
      const res = await newsboxService.createNewsboxItem(item);
      info(t('NEWSBOX_ITEM_CREATE_SUCCESS'));
      return res.data;
    } catch (error) {
      console.error('Error creating newsbox item:', error);
      throw error;
    } finally {
      state.adminLoading = false;
    }
  };

  const updateItemsAdmin = async (id: string, item: NewsboxItem) => {
    try {
      state.adminLoading = true;
      const res = await newsboxService.updateNewsboxItem(id, item);
      info(t('NEWSBOX_ITEM_UPDATE_SUCCESS'));
      return res.data;
    } catch (error) {
      console.error('Error updating newsbox item:', error);
      throw error;
    } finally {
      state.adminLoading = false;
    }
  };

  const updateLastSeen = async (id: string, itemId: UserLastSeenDto) => {
    try {
      state.loading = true;
      const res = await newsboxService.updateLastSeen(id, itemId);
      return res.data;
    } catch (error) {
      console.error('Error updating newsbox item:', error);
    } finally {
      state.loading = false;
    }
  };

  const getNewsItemByKey = (key: string): NewsboxItem | undefined => {
    return state.adminNewsItems?.items.find(item => item._key === key);
  };

  const deleteItemsAdmin = async (id: string) => {
    try {
      state.adminLoading = true;
      const res = await newsboxService.deleteItemsAdmin(id);
      info(t('NEWSBOX_ITEM_DELETE_SUCCESS'));
      return res.data;
    } catch (error) {
      console.error('Error deleting newsbox item:', error);
      throw error;
    } finally {
      state.adminLoading = false;
    }
  };

  return {
    ...toRefs(state),
    // actions
    fetchItems,
    fetchItemsAdmin,
    createItemsAdmin,
    updateItemsAdmin,
    deleteItemsAdmin,
    updateLastSeen,
    // getters
    getNewsItemByKey,
  };
});
