<script setup lang="ts">
import GridReviewRemarksPublic from '@cli/components/grids/GridReviewRemarksPublic.vue';
import type {Project} from '@cli/models/Project';
import {useAppStore} from '@cli/stores/app';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import {useTabsWindows} from '@shared/composables/useTabsWindows';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {projectService} from '../services/projectService';

const {t} = useI18n();
const router = useRouter();
const breadcrumbs = useBreadcrumbsStore();

const props = defineProps<{
  id?: string;
  tab?: string;
}>();

const currentProject = ref<Project | null>(null);
const loading = ref(false);
const error = ref<string | null>(null);

const appStore = useAppStore();

const projectTabs = ['overview', 'sbomlist', 'reviewRemarks', 'policyRules'];

const baseUrl = computed(() => `/projects/${props.id}`);
const {selectedTab} = useTabsWindows(baseUrl, projectTabs, props.tab || 'overview');

const updateBreadcrumbs = () => {
  if (currentProject.value) {
    const crumbs = [];

    // Add parent group to breadcrumbs if it exists
    if (appStore.parentProject && appStore.parentProject.uuid !== currentProject.value.uuid) {
      crumbs.push({
        title: appStore.parentProject.name,
        href: `/projects/${appStore.parentProject.uuid}/overview`,
        disabled: false,
      });
    }

    // Add current project
    crumbs.push({
      title: currentProject.value.name || t('PROJECT'),
      href: `/projects/${currentProject.value.uuid}/overview`,
      disabled: true,
    });

    // Only add tab breadcrumb if not overview and project is not a group (group has only overview)
    if (!currentProject.value.isGroup && selectedTab.value !== 'overview') {
      const tabKey = selectedTab.value;
      const tabTitleKeyMap: Record<string, string> = {
        sbomlist: 'TAB_PROJECT_SBOM_LIST',
        reviewRemarks: 'TAB_PROJECT_REVIEW_REMARKS',
        policyRules: 'TAB_PROJECT_POLICY_RULES',
      };
      const titleKey = tabTitleKeyMap[tabKey];

      if (titleKey) {
        crumbs.push({
          title: t(titleKey),
          href: `/projects/${currentProject.value.uuid}/${encodeURIComponent(tabKey)}`,
          disabled: true,
        });
      }
    }

    breadcrumbs.setCurrentBreadcrumbs(crumbs);
  }
};

// Fetch project data
const getProject = async () => {
  if (!props.id) {
    currentProject.value = null;
    appStore.setCurrentProject(null);
    return;
  }

  loading.value = true;
  error.value = null;

  try {
    const currentProjectData = await projectService.getProject(props.id, {
      parentProjectUuid: appStore.parentProject?.uuid ?? null,
    });
    if (currentProjectData) {
      currentProject.value = currentProjectData;
      appStore.setCurrentProject(currentProjectData);
      updateBreadcrumbs();
    } else {
      throw new Error(t('NO_PROJECT_DATA'));
    }
  } catch (err: unknown) {
    console.error('Error fetching project data:', err);
    // Error message will be passed from projectService
    error.value = err instanceof Error ? err.message : t('ERROR_LOADING_PROJECT');
    currentProject.value = null;
    appStore.setCurrentProject(null);
    currentProject.value = null;
    appStore.setCurrentProject(null);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  if (props.id) {
    getProject();
  }
});

watch(
  () => props.id,
  async (newId, oldId) => {
    if (newId && newId !== oldId) {
      await getProject();
    } else if (!newId) {
      currentProject.value = null;
    }
  },
);

const projectStatusClasses = computed(() => {
  const status = currentProject.value?.status || 'new';
  const baseClass = 'pStatus';
  return baseClass + (status === 'new' ? 'new' : status);
});

const projectStatusText = computed(() => {
  const status = currentProject.value?.status;
  return status ? t(`STATUS_${status}`) : 'new';
});
</script>

<template>
  <v-container fluid class="pa-4" data-testid="projects-details">
    <v-spacer></v-spacer>
    <div v-if="loading" class="d-flex align-center justify-center py-8">
      <v-progress-circular indeterminate color="primary" size="48" />
      <span class="text-body-1 ml-4">{{ t('Loading...') }}</span>
    </div>
    <div v-else>
      <div v-if="error" class="d-flex flex-column align-center justify-center py-4">
        <v-alert type="error" variant="tonal" width="100%" max-width="600" border="start" closable>
          <template v-slot:prepend>
            <v-icon color="error">mdi-alert-circle</v-icon>
          </template>
          {{ error }}
        </v-alert>
        <div class="d-flex align-center mt-8 w-full justify-center">
          <DCActionButton
            large
            variant="flat"
            icon="mdi-arrow-left"
            :text="t('BTN_BACK_TO_PROJECTS')"
            @click="router.push('/')" />
        </div>
      </div>
      <div v-if="!props.id" class="d-flex flex-column align-center justify-center py-8">
        <v-alert type="info" variant="tonal" width="100%" max-width="600" border="start">
          <template v-slot:prepend>
            <v-icon color="info">mdi-information</v-icon>
          </template>
          {{ t('NO_PROJECT_DATA') }}
        </v-alert>
        <div class="d-flex align-center mt-8 w-full justify-center">
          <DCActionButton
            large
            variant="flat"
            icon="mdi-arrow-left"
            :text="t('BTN_BACK_TO_PROJECTS')"
            @click="router.push('/')" />
        </div>
      </div>
      <template v-else-if="currentProject">
        <div class="d-flex align-center gap-2 py-3">
          <span class="text-h5">
            {{ currentProject.isGroup ? t('GROUP') : t('PROJECT') }} <q>{{ currentProject.name }}</q>
          </span>
          <span :class="projectStatusClasses">{{ projectStatusText }}</span>
        </div>
        <Stack class="expand">
          <v-card>
            <v-tabs v-model="selectedTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
              <v-tab value="overview">
                {{ t('TAB_PROJECT_OVERVIEW') }}
              </v-tab>
              <v-tab v-if="!currentProject.isGroup" value="sbomlist">
                {{ t('TAB_PROJECT_SBOM_LIST') }}
              </v-tab>
              <v-tab v-if="!currentProject.isGroup" value="reviewRemarks">
                {{ t('TAB_PROJECT_REVIEW_REMARKS') }}
              </v-tab>
              <v-tab v-if="!currentProject.isGroup" value="policyRules">
                {{ t('TAB_PROJECT_POLICY_RULES') }}
              </v-tab>
            </v-tabs>
            <v-window v-model="selectedTab">
              <v-window-item value="overview">
                <CliTabProjectOverview></CliTabProjectOverview>
              </v-window-item>
              <v-window-item v-if="!currentProject.isGroup" value="sbomlist">
                <GridSBOMPublic></GridSBOMPublic>
              </v-window-item>
              <v-window-item v-if="!currentProject.isGroup" value="reviewRemarks">
                <GridReviewRemarksPublic></GridReviewRemarksPublic>
              </v-window-item>
              <v-window-item v-if="!currentProject.isGroup" value="policyRules">
                <GridPolicyRulesPublic></GridPolicyRulesPublic>
              </v-window-item>
            </v-window>
          </v-card>
        </Stack>
      </template>
    </div>
  </v-container>
</template>
