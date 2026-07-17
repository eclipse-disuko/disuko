<script setup lang="ts">
import {useAppStore} from '@cli/stores/app';
import {useTabsWindows} from '@shared/composables/useTabsWindows';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';
import SbomSelector from '../components/SbomSelector.vue';
import GridReviewRemarksPublic from '../components/grids/GridReviewRemarksPublic.vue';
import GridSBOMComponentsPublic from '../components/grids/GridSBOMComponentsPublic.vue';

const {t} = useI18n();
const route = useRoute();
const breadcrumbs = useBreadcrumbsStore();

const props = defineProps<{
  id?: string;
  version?: string;
  tab?: string;
  spdx?: string;
}>();

const appStore = useAppStore();
const currentProject = computed(() => appStore.getCurrentProject());

const versionTabs = ['components', 'sbom', 'reviewRemarks', 'sourcecode', 'thirdPartyNotices'];

const baseUrl = computed(() => `/projects/${props.id}/versions/${props.version}`);
const spdxSuffix = computed(() => props.spdx || '');
const {selectedTab} = useTabsWindows(baseUrl, versionTabs, props.tab || 'components', spdxSuffix);

const isLoading = ref(true);

// Initialize project data when component mounts
onMounted(async () => {
  try {
    if (!currentProject.value && props.id) {
      await appStore.fetchCurrentProject(props.id);
    }
  } catch (error) {
    console.error('Error fetching project:', error);
  } finally {
    isLoading.value = false;
  }
});

const updateBreadcrumbs = () => {
  if (currentProject.value) {
    const crumbs = [];

    if (appStore.parentProject) {
      crumbs.push({
        title: appStore.parentProject.name,
        href: `/projects/${appStore.parentProject.uuid}/overview`,
        disabled: false,
      });
    }

    crumbs.push({
      title: currentProject.value.name || t('PROJECT'),
      href: `/projects/${currentProject.value.uuid}/overview`,
      disabled: false,
    });

    if (props.spdx) {
      crumbs.push({
        title: t('TAB_PROJECT_SBOM_LIST'),
        href: `/projects/${currentProject.value.uuid}/sbomlist`,
        disabled: false,
      });
    }

    if (props.version) {
      crumbs.push({
        title: props.version,
        href: `/projects/${currentProject.value.uuid}/versions/${props.version}/components`,
        disabled: true,
      });
    }

    breadcrumbs.setCurrentBreadcrumbs(crumbs);
  }
};

watch(
  [() => route.params, currentProject],
  () => {
    updateBreadcrumbs();
  },
  {immediate: true},
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
    <div v-if="isLoading" class="d-flex align-center justify-center py-8">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </div>
    <div v-else-if="!currentProject" class="d-flex align-center justify-center py-8">
      <v-alert type="info" variant="tonal" width="100%" max-width="600" border="start">
        <template v-slot:prepend>
          <v-icon color="info">mdi-information</v-icon>
        </template>
        {{ t('SELECT_PROJECT_TOKEN') }}
      </v-alert>
    </div>
    <template v-else>
      <div class="d-flex align-center gap-2">
        <div class="d-flex align-center">
          <span class="text-h5">
            {{ currentProject.isGroup ? t('GROUP') : t('PROJECT') }} <q>{{ currentProject.name }}</q>
          </span>
          <span :class="projectStatusClasses">{{ projectStatusText }}</span>
        </div>
        <v-spacer></v-spacer>
        <SbomSelector
          :sboms="currentProject.sboms || []"
          :spdx="props.spdx"
          :project-id="props.id"
          :version="props.version" />
        <v-spacer></v-spacer>
      </div>
      <Stack class="expand">
        <v-card>
          <v-tabs v-model="selectedTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
            <!-- <v-tab value="overview">
              {{ t('TAB_VERSION_OVERVIEW') }}
            </v-tab> -->
            <v-tab value="components">
              {{ t('COMPONENTS') }}
            </v-tab>
            <v-tab value="sbom">
              {{ t('TAB_PROJECT_SBOM_LIST') }}
            </v-tab>
            <v-tab value="reviewRemarks">
              {{ t('TAB_PROJECT_REVIEW_REMARKS') }}
            </v-tab>
            <v-tab value="sourcecode">
              {{ t('SOURCE_CODE') }}
            </v-tab>
            <v-tab value="thirdPartyNotices">
              {{ t('THIRD_PARTY_NOTICES') }}
            </v-tab>
          </v-tabs>
          <v-window v-model="selectedTab">
            <v-window-item value="components">
              <GridSBOMComponentsPublic v-if="currentProject"></GridSBOMComponentsPublic>
            </v-window-item>
            <v-window-item value="sbom">
              <GridSBOMPublic v-if="currentProject"></GridSBOMPublic>
            </v-window-item>
            <v-window-item value="reviewRemarks">
              <GridReviewRemarksPublic v-if="currentProject"></GridReviewRemarksPublic>
            </v-window-item>
            <v-window-item value="sourcecode">
              <GridSourceCodePublic v-if="currentProject"></GridSourceCodePublic>
            </v-window-item>
            <v-window-item value="thirdPartyNotices">
              <CliTabNoticeFile></CliTabNoticeFile>
            </v-window-item>
          </v-window>
        </v-card>
      </Stack>
    </template>
  </v-container>
</template>
