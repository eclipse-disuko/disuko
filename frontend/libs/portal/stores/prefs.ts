// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {IMap} from '@disclosure-portal/utils/View';
import {defineStore} from 'pinia';

const STORE_NAME = 'uiprefs';

interface WithTTL<Type> {
  expiry: Date | undefined;
  value: Type;
}
interface sortPref {
  col: string;
  desc: boolean;
}
interface PermPrefs {
  appLanguage: string;
  tableFilter: IMap<WithTTL<string[]>>;
  tableSort: IMap<WithTTL<sortPref>>;
  tableSearch: IMap<WithTTL<string>>;
}
interface SessionPrefs {
  tableFilter: IMap<WithTTL<string[]>>;
  tableSort: IMap<WithTTL<sortPref>>;
  tableSearch: IMap<WithTTL<string>>;
}
interface UiPrefs {
  perm: PermPrefs;
  sess: SessionPrefs;
}

const getDefaultPrefs: () => UiPrefs = () => ({
  perm: {
    appLanguage: 'en',
    tableSort: {} as IMap<WithTTL<sortPref>>,
    tableFilter: {} as IMap<WithTTL<string[]>>,
    tableSearch: {} as IMap<WithTTL<string>>,
  },
  sess: {
    tableSort: {} as IMap<WithTTL<sortPref>>,
    tableFilter: {} as IMap<WithTTL<string[]>>,
    tableSearch: {} as IMap<WithTTL<string>>,
  },
});

function cleanupPerm(perm: PermPrefs): boolean {
  let changed = false;
  for (const k in perm.tableFilter) {
    if (perm.tableFilter[k].expiry && perm.tableFilter[k].expiry < new Date()) {
      delete perm.tableFilter[k];
      changed = true;
    }
  }
  for (const k in perm.tableSort) {
    if (perm.tableSort[k].expiry && perm.tableSort[k].expiry < new Date()) {
      delete perm.tableSort[k];
      changed = true;
    }
  }
  for (const k in perm.tableSearch) {
    if (perm.tableSearch[k].expiry && perm.tableSearch[k].expiry < new Date()) {
      delete perm.tableSearch[k];
      changed = true;
    }
  }
  return changed;
}

const getPrefs = (): UiPrefs => {
  const dateReviver = (key: string, value: any) => {
    if (typeof value === 'string' && key === 'expiry') {
      return new Date(value);
    }
    return value;
  };
  const res = getDefaultPrefs();
  const permPrefs = localStorage.getItem(STORE_NAME);
  if (permPrefs) {
    res.perm = JSON.parse(permPrefs, dateReviver);
    if (cleanupPerm(res.perm)) {
      localStorage.setItem(STORE_NAME, JSON.stringify(res.perm));
    }
  }
  const sessPrefs = sessionStorage.getItem(STORE_NAME);
  if (sessPrefs) {
    res.sess = JSON.parse(sessPrefs, dateReviver);
  }
  return res;
};

export const useUiPrefsStore = defineStore(STORE_NAME, {
  state: () => ({
    prefs: getPrefs(),
  }),
  actions: {
    changeLanguage(newLang: string) {
      this.prefs.perm.appLanguage = newLang;
      localStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.perm));
    },
    setFilter(tableKey: string, filterCol: string, values: string[], sessionOnly: boolean, ttlSec = 0) {
      const key = tableKey + '.' + filterCol;
      if (!sessionOnly) {
        if (!this.prefs.perm.tableFilter) {
          this.prefs.perm.tableFilter = {} as IMap<WithTTL<string[]>>;
        }
        if (Object.prototype.hasOwnProperty.call(this.prefs.sess.tableFilter, key)) {
          delete this.prefs.sess.tableFilter[key];
          sessionStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.sess));
        }
        if (ttlSec === 0) {
          this.prefs.perm.tableFilter[key] = {
            value: values,
            expiry: undefined,
          };
        } else {
          const expiry = new Date();
          expiry.setSeconds(expiry.getSeconds() + ttlSec);

          this.prefs.perm.tableFilter[key] = {
            value: values,
            expiry: expiry,
          };
        }
        localStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.perm));
      } else {
        if (!this.prefs.sess.tableFilter) {
          this.prefs.sess.tableFilter = {} as IMap<WithTTL<string[]>>;
        }
        if (Object.prototype.hasOwnProperty.call(this.prefs.perm.tableFilter, key)) {
          delete this.prefs.perm.tableFilter[key];
          localStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.perm));
        }
        this.prefs.sess.tableFilter[key] = {
          value: values,
          expiry: undefined,
        };
        sessionStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.sess));
      }
    },
    setSort(tableKey: string, sortBy: string, sortDesc: boolean, sessionOnly: boolean, ttlSec = 0) {
      const key = tableKey;
      if (!sessionOnly) {
        if (!this.prefs.perm.tableSort) {
          this.prefs.perm.tableSort = {} as IMap<WithTTL<sortPref>>;
        }
        if (Object.prototype.hasOwnProperty.call(this.prefs.sess.tableSort, key)) {
          delete this.prefs.sess.tableSort[key];
          sessionStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.sess));
        }
        if (ttlSec === 0) {
          this.prefs.perm.tableSort[key] = {
            value: {
              col: sortBy,
              desc: sortDesc,
            },
            expiry: undefined,
          };
        } else {
          const expiry = new Date();
          expiry.setSeconds(expiry.getSeconds() + ttlSec);

          this.prefs.perm.tableSort[key] = {
            value: {
              col: sortBy,
              desc: sortDesc,
            },
            expiry: expiry,
          };
        }
        localStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.perm));
      } else {
        if (!this.prefs.sess.tableSort) {
          this.prefs.sess.tableSort = {} as IMap<WithTTL<sortPref>>;
        }
        if (Object.prototype.hasOwnProperty.call(this.prefs.perm.tableSort, key)) {
          delete this.prefs.perm.tableSort[key];
          localStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.perm));
        }
        this.prefs.sess.tableSort[key] = {
          value: {
            col: sortBy,
            desc: sortDesc,
          },
          expiry: undefined,
        };
        sessionStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.sess));
      }
    },
    setSearch(tableKey: string, search: string, sessionOnly: boolean, ttlSec = 0) {
      const key = tableKey;
      if (!sessionOnly) {
        if (!this.prefs.perm.tableSearch) {
          this.prefs.perm.tableSearch = {} as IMap<WithTTL<string>>;
        }
        if (Object.prototype.hasOwnProperty.call(this.prefs.sess.tableSearch, key)) {
          delete this.prefs.sess.tableSearch[key];
          sessionStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.sess));
        }
        if (ttlSec === 0) {
          this.prefs.perm.tableSearch[key] = {
            value: search,
            expiry: undefined,
          };
        } else {
          const expiry = new Date();
          expiry.setSeconds(expiry.getSeconds() + ttlSec);

          this.prefs.perm.tableSearch[key] = {
            value: search,
            expiry: expiry,
          };
        }
        localStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.perm));
      } else {
        if (!this.prefs.sess.tableSearch) {
          this.prefs.sess.tableSearch = {} as IMap<WithTTL<string>>;
        }
        if (Object.prototype.hasOwnProperty.call(this.prefs.perm.tableSearch, key)) {
          delete this.prefs.perm.tableSearch[key];
          localStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.perm));
        }
        this.prefs.sess.tableSearch[key] = {
          value: search,
          expiry: undefined,
        };
        sessionStorage.setItem(STORE_NAME, JSON.stringify(this.prefs.sess));
      }
    },
  },
  getters: {
    getFilter: (state) => {
      return (tableKey: string, filterCol: string) => {
        const now = new Date();
        const combined = {...state.prefs.perm.tableFilter, ...state.prefs.sess.tableFilter};
        const res = combined[tableKey + '.' + filterCol];
        if (res) {
          if (res.expiry && res.expiry < now) {
            return undefined;
          }
          return res.value;
        }
        return undefined;
      };
    },
    getSearch: (state) => {
      return (tableKey: string) => {
        const now = new Date();
        const combined = {...state.prefs.perm.tableSearch, ...state.prefs.sess.tableSearch};
        const res = combined[tableKey];
        if (res) {
          if (res.expiry && res.expiry < now) {
            return undefined;
          }
          return res.value;
        }
        return undefined;
      };
    },
    getSort: (state) => {
      return (tableKey: string) => {
        const now = new Date();
        const combined = {...state.prefs.perm.tableSort, ...state.prefs.sess.tableSort};
        const res = combined[tableKey];
        if (res) {
          if (res.expiry && res.expiry < now) {
            return undefined;
          }
          return res.value;
        }
        return undefined;
      };
    },
  },
});
