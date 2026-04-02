<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {usePageTitle} from '@disclosure-portal/composables/usePageTitle';
import {ProjectSubscriptions} from '@disclosure-portal/model/Project';
import {ProjectSlim, ProjectSlimDto} from '@disclosure-portal/model/ProjectsResponse';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import _ from 'lodash';
import {storeToRefs} from 'pinia';
import {computed, nextTick, onUnmounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';

const route = useRoute();
const router = useRouter();
const {t, locale} = useI18n();
const {dashboardCrumbs, projectsCrumb, ...breadcrumbs} = useBreadcrumbsStore();
const {useReactiveTitle} = usePageTitle();
const appStore = useAppStore();
const projectStore = useProjectStore();
const idleStore = useIdleStore();

const {currentProject} = storeToRefs(projectStore);

const projectSubscriptionsDialogOpen = ref(false);
const projectSubscriptionsDialogVisible = ref(false);
const selectedTab = ref('overview');

const projectId = computed(() => (Array.isArray(route.params?.uuid) ? route.params.uuid[0] : route.params?.uuid || ''));
const tab = computed(() => (Array.isArray(route.params?.tab) ? route.params.tab[0] : route.params?.tab || ''));
const itemVersion = computed(() =>
  Array.isArray(route.params?.version) ? route.params.version[0] : route.params?.version || '',
);
const encodedCurrentProjectParent = computed(() => encodeURIComponent(currentProject.value?.parent ?? ''));
const encodedProjectId = computed(() => encodeURIComponent(projectId.value));

const slimFromProject = (): ProjectSlim => {
  const projectSlim = new ProjectSlim(new ProjectSlimDto());
  projectSlim.fromProjectModel(currentProject.value!);
  return projectSlim;
};
const tabUrl = computed(() => {
  const type = currentProject.value?.isGroup ? 'groups' : 'projects';
  return `/dashboard/${type}/${encodedProjectId.value}`;
});

// Function to get tab display name from route
const getTabDisplayName = () => {
  const path = route.path;
  if (path.includes('/overview')) return t('TAB_PROJECT_OVERVIEW');
  if (path.includes('/children')) return t('TAB_PROJECT_CHILDREN');
  if (path.includes('/sbomlist')) return t('SBOM_DELIVERIES');
  if (path.includes('/versionlist')) return t('TAB_PROJECT_VERSION');
  if (path.includes('/users')) return t('TAB_PROJECT_USERMANAGEMENT');
  if (path.includes('/childrenUsers')) return t('TAB_PROJECT_USERMANAGEMENT');
  if (path.includes('/tokens')) return t('TAB_PROJECT_TOKENMANAGEMENT');
  if (path.includes('/policyrules')) return t('TAB_POLICYRULES');
  if (path.includes('/licenserules')) return t('TAB_LICENSERULES');
  if (path.includes('/approvals')) return t('TAB_PROJECT_APPROVALS');
  if (path.includes('/auditLog')) return t('TAB_PROJECT_AUDIT');
  return t('TAB_PROJECT_OVERVIEW');
};

// Set up reactive page title
watch(
  () => [currentProject.value?.name, route.name, locale.value],
  ([projectName]) => {
    if (projectName) {
      const tabName = getTabDisplayName();
      useReactiveTitle(`${String(projectName)} | ${tabName}`);
    }
  },
  {immediate: true},
);

const openParent = async () => {
  await router.push(`/dashboard/groups/${encodedCurrentProjectParent.value}/children`);
};

const initBreadcrumbs = () => {
  const currentGroupCrumb = {
    title: currentProject.value?.parentName ?? '',
    href: `/dashboard/groups/${encodedCurrentProjectParent.value}/children`,
  };
  const currentProjectCrumb = {
    title: currentProject.value?.name ?? '',
    href: `/dashboard/projects/${encodedProjectId.value}/overview`,
  };
  const groupProjectCrumbs = currentProject.value?.parent
    ? [currentGroupCrumb, currentProjectCrumb]
    : [currentProjectCrumb];
  breadcrumbs.setCurrentBreadcrumbs([...dashboardCrumbs, projectsCrumb, ...groupProjectCrumbs]);
};

const reload = async () => {
  await projectStore.fetchProjectByKey(projectId.value);
};

const initPage = async () => {
  await nextTick();
  appStore.setDummyDesignMode(currentProject.value?.isDummy ?? false);
  initBreadcrumbs();
};

const saveProjectSubscriptions = async (item: ProjectSubscriptions) => {
  if (_.isEqual(item, currentProject.value?.subscriptions) || !currentProject.value) {
    projectSubscriptionsDialogOpen.value = false;
    return;
  }
  try {
    await projectStore.updateProjectSubscriptions(currentProject.value._key, item);
    projectSubscriptionsDialogOpen.value = false;
  } catch (e) {}
};

watch(
  projectId,
  async () => {
    await reload();
  },
  {immediate: true},
);

watch(currentProject, async (cp) => {
  if (cp) {
    await initPage();
    idleStore.hide();
  }
});

watch(tab, (newTab) => {
  if (newTab) {
    selectedTab.value = newTab;
  } else {
    selectedTab.value = 'overview';
  }
});

onUnmounted(() => {
  appStore.unsetDummyDesignMode();
});
</script>

<template>
  <v-container v-if="currentProject" fluid class="h-full px-6" data-testid="projects-details">
    <div v-if="!itemVersion" class="flex flex-row align-center pb-3 ga-2 flex-wrap">
      <v-chip v-if="currentProject.isDummy" class="dummy-tag" label>DUMMY</v-chip>
      <span class="text-h5 inline-block" :class="{statusDeprecated: currentProject.status === 'deprecated'}">
        {{ currentProject.isGroup ? t('GROUP') : t('PROJECT') }} <q>{{ currentProject.name }}</q>
      </span>
      <span :class="'pStatus' + (!currentProject.status ? 'new' : currentProject.status)">
        {{ !currentProject.status ? 'new' : t('STATUS_' + currentProject.status) }}
      </span>
      <DIconButton
        v-if="currentProject.parent"
        @clicked="openParent"
        :hint="currentProject.parentName"
        icon="mdi-table-multiple"></DIconButton>
      <v-spacer></v-spacer>
      <DCActionButton
        v-if="currentProject.showSubscriptionButton"
        :icon="currentProject.hasSubscriptions ? 'mdi-bell-ring' : 'mdi-bell'"
        :hint="t('TT_PROJECT_SUBSCRIPTIONS')"
        :text="t('BTN_SUBSCRIPTIONS')"
        large
        @click="projectSubscriptionsDialogOpen = true"></DCActionButton>
      <ProjectSettings v-if="currentProject.allowProjectEdit" v-slot="{showDialog}">
        <DCActionButton
          icon="mdi-cog"
          :hint="t('TT_settings_project')"
          :text="t('BTN_SETTINGS')"
          large
          data-testid="projects-edit-button"
          @click.stop="showDialog"></DCActionButton>
      </ProjectSettings>
      <ProjectMenu></ProjectMenu>
    </div>
    <v-row class="expand" v-if="!itemVersion">
      <v-col cols="12">
        <v-card>
          <v-tabs v-model="selectedTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
            <v-tab value="overview" :to="`${tabUrl}/overview`">
              {{ t('TAB_PROJECT_OVERVIEW') }}
            </v-tab>
            <v-tab value="children" v-if="currentProject.isGroup" :to="`${tabUrl}/children`">
              {{ t('TAB_PROJECT_CHILDREN') }}
            </v-tab>
            <v-tab value="sbomlist" v-if="!currentProject.isGroup" :to="`${tabUrl}/sbomlist`">
              {{ t('TAB_SBOM_DELIVERIES') }}
            </v-tab>
            <v-tab value="versionlist" v-if="!currentProject.isGroup" :to="`${tabUrl}/versionlist`">
              {{ t('TAB_PROJECT_VERSION') }}
            </v-tab>
            <v-tab
              value="users"
              v-if="!currentProject.isGroup && currentProject.isUserManagementAllowed"
              :to="`${tabUrl}/users`">
              {{ t('TAB_PROJECT_USERMANAGEMENT') }}
            </v-tab>
            <v-tab
              value="childrenUsers"
              v-if="currentProject.isGroup && currentProject.isUserManagementAllowed"
              :to="`${tabUrl}/childrenUsers`">
              {{ t('TAB_PROJECT_USERMANAGEMENT') }}
            </v-tab>
            <v-tab value="tokens" v-if="currentProject.isTokenManagementAllowed" :to="`${tabUrl}/tokens`">
              {{ t('TAB_PROJECT_TOKENMANAGEMENT') }}
            </v-tab>
            <v-tab value="policyrules" v-if="!currentProject.isGroup" :to="`${tabUrl}/policyrules`">
              {{ t('TAB_POLICYRULES') }}
            </v-tab>
            <v-tab value="decisions" v-if="!currentProject.isGroup" :to="`${tabUrl}/decisions`">
              {{ t('TAB_DECISIONS') }}
            </v-tab>
            <v-tab value="approvals" v-if="currentProject.accessRights" :to="`${tabUrl}/approvals`">
              {{ t('TAB_PROJECT_APPROVALS') }}
            </v-tab>
            <v-tab
              value="auditLog"
              v-if="currentProject.accessRights.allowProjectAudit.read"
              :to="`${tabUrl}/auditLog`">
              {{ t('TAB_PROJECT_AUDIT') }}
            </v-tab>
          </v-tabs>
          <v-tabs-window v-model="selectedTab">
            <v-tabs-window-item value="overview">
              <TabProjectOverview></TabProjectOverview>
            </v-tabs-window-item>
            <v-tabs-window-item value="children">
              <GridChildren></GridChildren>
            </v-tabs-window-item>
            <v-tabs-window-item value="sbomlist">
              <GridSBOM></GridSBOM>
            </v-tabs-window-item>
            <v-tabs-window-item value="versionlist">
              <GridVersions></GridVersions>
            </v-tabs-window-item>
            <v-tabs-window-item value="users">
              <TabProjectUserManagement v-if="currentProject._key"></TabProjectUserManagement>
            </v-tabs-window-item>
            <v-tabs-window-item value="childrenUsers">
              <TabProjectChildrenUserManagement v-if="currentProject._key"></TabProjectChildrenUserManagement>
            </v-tabs-window-item>
            <v-tabs-window-item value="tokens">
              <TabProjectTokenManagement></TabProjectTokenManagement>
            </v-tabs-window-item>
            <v-tabs-window-item value="policyrules">
              <TabPolicyrules @reload="reload"></TabPolicyrules>
            </v-tabs-window-item>
            <v-tabs-window-item value="decisions">
              <TabDecisions></TabDecisions>
            </v-tabs-window-item>
            <v-tabs-window-item value="approvals">
              <TabProjectApprovals></TabProjectApprovals>
            </v-tabs-window-item>
            <v-tabs-window-item value="auditLog">
              <TabAuditLog :project-uuid="currentProject._key"></TabAuditLog>
            </v-tabs-window-item>
          </v-tabs-window>
        </v-card>
      </v-col>
    </v-row>
    <DFormDialog v-model:dialog="projectSubscriptionsDialogVisible" v-model="projectSubscriptionsDialogOpen" persistent>
      <ProjectSubscriptionsDialog
        :title="t('TITLE_PROJECT_SUBSCRIPTIONS')"
        :confirm-text="t('NP_DIALOG_BTN_EDIT')"
        :item="currentProject.subscriptions"
        @confirm="saveProjectSubscriptions"
        @close="projectSubscriptionsDialogOpen = false">
      </ProjectSubscriptionsDialog>
    </DFormDialog>
  </v-container>
</template>

<style scoped lang="scss">
.dummy-tag {
  margin-left: 2px;
  border: 1px solid rgb(var(--v-theme-chartYellow));
}

.statusDeprecated {
  color: rgb(var(--v-theme-projectDeprecated));
}
</style>
