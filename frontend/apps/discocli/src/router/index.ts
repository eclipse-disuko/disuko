import {createRouter, createWebHashHistory, RouteRecordRaw} from 'vue-router';
import {authService} from '@cli/services/authService';
import {useAppStore} from '@cli/stores/app';

const baseUrl = import.meta.env.BASE_URL || '/supplierportal/';

const routes: RouteRecordRaw[] = [
  {
    path: '/auth',
    name: 'Login',
    component: () => import('@cli/views/Login.vue'),
    meta: {requiresAuth: false, layout: 'auth'},
  },
  {
    path: '/',
    redirect: () => {
      const authProjectUuid = authService.getAuth();
      if (authProjectUuid) {
        return {name: 'ProjectDetails', params: {id: authProjectUuid}};
      }
      return {name: 'Login'};
    },
  },
  {
    path: '/projects',
    redirect: '/',
  },
  {
    path: '/projects/:id/:tab?',
    name: 'ProjectDetails',
    component: () => import('@cli/views/ProjectDetails.vue'),
    props: true,
    meta: {requiresAuth: true, layout: 'app'},
  },
  {
    path: '/projects/:id/versions/:version/:tab?/:spdx?',
    name: 'VersionDetails',
    component: () => import('@cli/views/ViewVersion.vue'),
    props: true,
    meta: {requiresAuth: true, layout: 'app'},
  },
];

const router = createRouter({
  history: createWebHashHistory(baseUrl),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const requiresAuth = Boolean(to.meta.requiresAuth);
  const appStore = useAppStore();
  let isLoggedIn = authService.isAuthenticated();
  let authProjectUuid = authService.getAuth();

  if (requiresAuth && !isLoggedIn) {
    const info = await authService.fetchAuthInfo();
    if (info) {
      appStore.setAuthProjectUuid(info);
      isLoggedIn = true;
      authProjectUuid = info;
    } else {
      appStore.setAuthProjectUuid(null);
    }
  }

  if (to.name === 'Login' && isLoggedIn) {
    if (authProjectUuid) {
      next({name: 'ProjectDetails', params: {id: authProjectUuid}});
      return;
    }
    next({name: 'Login'});
    return;
  }

  if (requiresAuth && !isLoggedIn) {
    next({name: 'Login'});
    return;
  }

  next();
});

export default router;
