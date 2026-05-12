// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import * as Help from '@disclosure-portal/assets/documents/help';
import {Rights} from '@shared/types/Rights';
import profileService from '@shared/services/profile';
import {useUserStore} from '@disclosure-portal/stores/user';
import Dashboard from '@disclosure-portal/views/Dashboard.vue';
import {createRouter, createWebHashHistory, RouteRecordRaw} from 'vue-router';

// Dynamic imports for all components
const Home = () => import('@disclosure-portal/views/Home.vue');
const Projects = () => import('@disclosure-portal/views/projects/Projects.vue');
const ProjectsDetail = () => import('@disclosure-portal/views/projects/ProjectsDetail.vue');
const ProjectsVersions = () => import('@disclosure-portal/views/projects/ProjectsVersions.vue');
const Licenses = () => import('@disclosure-portal/views/licenses/Licenses.vue');
const PolicyRules = () => import('@disclosure-portal/views/policies/PolicyRules.vue');
const PolicyRulesDetail = () => import('@disclosure-portal/views/policies/PolicyRulesDetail.vue');
const Announcements = () => import('@disclosure-portal/views/announcements/Announcements.vue');
const Tasks = () => import('@disclosure-portal/views/tasks/Tasks.vue');
const TaskApprovalDialog = () => import('@disclosure-portal/components/dialog/TaskApprovalDialog.vue');
const AnalyticMain = () => import('@disclosure-portal/views/analytics/AnalyticMain.vue');
const LicenseMain = () => import('@disclosure-portal/views/licenses/LicenseMain.vue');
const UserProfile = () => import('@disclosure-portal/views/user/UserProfile.vue');
const ReviewTemplates = () => import('@disclosure-portal/views/admin/ReviewTemplates.vue');
const AdminDashboard = () => import('@disclosure-portal/views/admin/AdminDashboard.vue');
const Tools = () => import('@disclosure-portal/views/admin/AdminTools.vue');
const Labels = () => import('@disclosure-portal/views/admin/Labels.vue');
const Jobs = () => import('@disclosure-portal/views/admin/tools/Jobs.vue');
const AdminClassifications = () => import('@disclosure-portal/views/admin/AdminClassifications.vue');
const Users = () => import('@disclosure-portal/views/admin/Users.vue');
const AdminProfile = () => import('@disclosure-portal/views/admin/AdminProfile.vue');
const Schemas = () => import('@disclosure-portal/views/admin/schema/Schemas.vue');
const CustomIds = () => import('@disclosure-portal/views/admin/CustomIds.vue');
const CheckList = () => import('@disclosure-portal/views/admin/checklist/Checklist.vue');
const CheckListMain = () => import('@disclosure-portal/views/admin/checklist/ChecklistMain.vue');
const SchemaMain = () => import('@disclosure-portal/views/admin/schema/SchemaMain.vue');
const AdminProjects = () => import('@disclosure-portal/views/admin/AdminProjects.vue');
const InternalToken = () => import('@disclosure-portal/views/admin/InternalToken.vue');
const Newsbox = () => import('@disclosure-portal/views/admin/Newsbox.vue');
const FeatureFlags = () => import('@disclosure-portal/views/admin/FeatureFlags.vue');
const UpcomingDeletions = () => import('@disclosure-portal/views/admin/UpcomingDeletions.vue');
const UserManagement = () => import('@disclosure-portal/views/admin/tools/UserManagement.vue');

const baseUrl = import.meta.env.BASE_URL;

const currentRouteKey = 'dpCurrentRoute';
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: {name: 'Home'},
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    redirect: {name: 'Home'},
    children: [
      {
        path: 'home',
        name: 'Home',
        component: Home,
        meta: {
          title: {
            en: 'Dashboard',
            de: 'Dashboard',
          },
          helpText: {
            en: Help.HomeEn,
            de: Help.HomeDe,
          },
        },
      },
      {
        path: 'user',
        name: 'UserProfile',
        component: UserProfile,
        meta: {
          title: {
            en: 'User Profile',
            de: 'Benutzerprofil',
          },
        },
      },
      {
        path: 'policyrules/:id/:tab?',
        name: 'Policy rule',
        component: PolicyRulesDetail,
        meta: {
          helpText: {
            en: Help.PolicyRuleDetailsEn,
            de: Help.PolicyRuleDetailsDe,
          },
        },
      },
      {
        path: 'policyrules',
        name: 'Policy rules',
        component: PolicyRules,
        meta: {
          title: {
            en: 'Policy Rules',
            de: 'Policy Regeln',
          },
          helpText: {
            en: Help.PolicyRulesEn,
            de: Help.PolicyRulesDe,
          },
        },
      },
      {
        path: 'announcements',
        name: 'Announcements',
        component: Announcements,
        meta: {
          title: {
            en: 'Announcements',
            de: 'Ankündigungen',
          },
          helpText: {
            en: Help.AnnouncementsEn,
            de: Help.AnnouncementsDe,
          },
        },
      },
      {
        path: 'tasks',
        name: 'Tasks',
        component: Tasks,
        meta: {
          title: {
            en: 'Tasks',
            de: 'Aufgaben',
          },
        },
        children: [
          {
            path: ':id',
            name: 'TasksApprovalDialog',
            component: TaskApprovalDialog,
          },
        ],
      },
      {
        path: 'analytics/:tab?',
        name: 'Analytic',
        component: AnalyticMain,
        meta: {
          title: {
            en: 'Analytics',
            de: 'Analysen',
          },
          helpText: {
            en: Help.AnalyticsEn,
            de: Help.AnalyticsDe,
          },
        },
      },
      {
        path: 'projects',
        name: 'Projects',
        component: Projects,
        meta: {
          title: {
            en: 'Projects',
            de: 'Projekte',
          },
          helpText: {
            en: Help.ProjectsListEn,
            de: Help.ProjectsListDe,
          },
        },
      },
      {
        path: 'licenses',
        name: 'Licenses',
        component: Licenses,
        meta: {
          title: {
            en: 'Licenses',
            de: 'Lizenzen',
          },
          helpText: {
            en: Help.LicensesListEn,
            de: Help.LicensesListDe,
          },
        },
      },
      {
        path: 'licenses/compare',
        name: 'CompareLicenses',
        component: Licenses,
      },
      {
        path: 'licenses/:id/:tab?',
        name: 'License',
        component: LicenseMain,
        meta: {
          helpText: {
            en: Help.LicenseEn,
            de: Help.LicenseDe,
          },
        },
      },
      {
        path: 'licenses/:id/:tab?',
        name: 'LicenseClassifications',
        component: LicenseMain,
        meta: {
          helpText: {
            en: Help.LicenseEn,
            de: Help.LicenseDe,
          },
        },
      },
      {
        path: 'licenses/filtersets/:id',
        name: 'FilterSetsForLicenses',
        component: Licenses,
      },
      {
        path: 'admin',
        name: 'AdminDashboard',
        component: AdminDashboard,
        meta: {
          title: {
            en: 'Admin Dashboard',
            de: 'Admin Dashboard',
          },
        },
      },
      {
        path: 'admin/projects',
        name: 'adminProjects',
        component: AdminProjects,
        meta: {
          title: {
            en: 'All Projects',
            de: 'Alle Projekte',
          },
        },
      },
      {
        path: 'admin/tools/:tab?',
        name: 'DefaultTools',
        component: Tools,
        meta: {
          title: {
            en: 'Tools',
            de: 'Tools',
          },
        },
        children: [
          {
            path: 'analytics',
            name: 'AdminAnalytics',
            component: Tools,
          },
          {
            path: 'accessRights',
            name: 'AdminAccessRights',
            component: Tools,
          },
          {
            path: 'export_import',
            name: 'AdminExportImport',
            component: Tools,
          },
          {
            path: 'storageConsistency',
            name: 'AdminStorageConsistency',
            component: Tools,
            meta: {
              helpText: {
                en: Help.AdminToolStorageConsistencyEn,
                de: Help.AdminToolStorageConsistencyDe,
              },
            },
          },
          {
            path: 'sampleData',
            name: 'AdminSampleData',
            component: Tools,
          },
          {
            path: 'termsOfUseManagement',
            name: 'AdminTermsOfUse',
            component: Tools,
            meta: {
              helpText: {
                en: Help.AdminToolTermsOfUseManagementEn,
                de: Help.AdminToolTermsOfUseManagementDe,
              },
            },
          },
          {
            path: 'notificationBar',
            name: 'AdminNotificationBar',
            component: Tools,
          },
          {
            path: 'mail',
            name: 'AdminMail',
            component: Tools,
          },
        ],
      },
      {
        path: 'admin/schemas',
        name: 'Schemas',
        component: Schemas,
        meta: {
          title: {
            en: 'SBOM-Schemas',
            de: 'SBOM-Schemas',
          },
        },
      },
      {
        path: 'schemas/:id',
        name: 'SchemaReadonly',
        component: SchemaMain,
      },
      {
        path: 'admin/schemas/:id',
        name: 'Schema',
        component: SchemaMain,
        meta: {
          title: {
            en: 'Schema Details',
            de: 'Schema Details',
          },
        },
      },
      {
        path: 'admin/schemas/new',
        name: 'New schema',
        component: SchemaMain,
        meta: {
          title: {
            en: 'New Schema',
            de: 'Neues Schema',
          },
        },
      },
      {
        path: 'admin/jobs',
        name: 'AdminJobList',
        component: Jobs,
        meta: {
          title: {
            en: 'Jobs',
            de: 'Jobs',
          },
        },
      },
      {
        path: 'admin/userManagement',
        name: 'AdminUserManagement',
        component: UserManagement,
        meta: {
          title: {
            en: 'User Management',
            de: 'Benutzerverwaltung',
          },
        },
      },
      {
        path: 'admin/templates/:tab?',
        name: 'ReviewTemplates',
        component: ReviewTemplates,
        meta: {
          title: {
            en: 'Review Templates',
            de: 'Review Vorlagen',
          },
        },
        children: [
          {
            path: 'review',
            name: 'ReviewTemplatesReview',
            component: ReviewTemplates,
          },
        ],
      },
      {
        path: 'admin/labels/:tab?',
        name: 'Labels',
        component: Labels,
        meta: {
          title: {
            en: 'Labels',
            de: 'Labels',
          },
        },
        children: [
          {
            path: 'schema',
            name: 'LabelsSchema',
            component: Labels,
          },
          {
            path: 'policy',
            name: 'LabelsPolicy',
            component: Labels,
          },
        ],
      },
      {
        path: 'admin/classifications',
        name: 'Classification',
        component: AdminClassifications,
        meta: {
          title: {
            en: 'Classifications',
            de: 'Klassifizierungen',
          },
        },
      },
      {
        path: 'admin/users',
        name: 'AdminUsers',
        component: Users,
        meta: {
          title: {
            en: 'Users',
            de: 'Benutzer',
          },
        },
      },
      {
        path: 'admin/users/:uuid',
        name: 'AdminProfile',
        component: AdminProfile,
        meta: {
          title: {
            en: 'User Profile',
            de: 'Benutzerprofil',
          },
        },
      },
      {
        path: 'admin/customids',
        name: 'CustomIds',
        component: CustomIds,
        meta: {
          title: {
            en: 'Custom IDs',
            de: 'Custom IDs',
          },
        },
      },
      {
        path: 'admin/checklist',
        name: 'Checklist',
        component: CheckList,
        meta: {
          title: {
            en: 'Checklist',
            de: 'Checkliste',
          },
        },
      },
      {
        path: 'admin/checklist/:id',
        name: 'ChecklistMain',
        component: CheckListMain,
        meta: {
          title: {
            en: 'Checklist Details',
            de: 'Checkliste Details',
          },
        },
      },
      {
        path: 'admin/internaltoken',
        name: 'InternalToken',
        component: InternalToken,
        meta: {
          title: {
            en: 'Internal Token',
            de: 'Interner Token',
          },
        },
      },
      {
        path: 'admin/newsbox',
        name: 'Newsbox',
        component: Newsbox,
        meta: {
          title: {
            en: 'Newsbox',
            de: 'Newsbox',
          },
        },
      },
      {
        path: 'admin/featureflags',
        name: 'FeatureFlags',
        component: FeatureFlags,
        meta: {
          title: {
            en: 'Feature Flags',
            de: 'Feature Flags',
          },
        },
      },
      {
        path: 'admin/deletions',
        name: 'UpcomingDeletions',
        component: UpcomingDeletions,
        meta: {
          title: {
            en: 'Upcoming Deletions',
            de: 'Bevorstehende Löschungen',
          },
        },
      },
      {
        path: 'groups/:uuid/:tab?',
        name: 'Group',
        component: ProjectsDetail,
        meta: {
          helpText: {
            en: Help.GroupsOverviewEn,
            de: Help.GroupsOverviewDe,
          },
        },
        children: [
          {
            path: 'overview',
            name: 'GroupsOverview',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.GroupsOverviewEn,
                de: Help.GroupsOverviewDe,
              },
            },
          },
          {
            path: 'children',
            name: 'GroupsChildren',
            component: ProjectsDetail,
          },
          {
            path: 'users',
            name: 'GroupsUsers',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.GroupsUsersEn,
                de: Help.GroupsUsersDe,
              },
            },
          },
          {
            path: 'tokens',
            name: 'GroupsTokens',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.GroupsTokenManagementEn,
                de: Help.GroupsTokenManagementDe,
              },
            },
          },
          {
            path: 'approvals',
            name: 'GroupsApprovals',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.GroupsApprovalsEn,
                de: Help.GroupsApprovalsDe,
              },
            },
          },
          {
            path: 'auditLog',
            name: 'GroupsAuditLog',
            component: ProjectsDetail,
          },
        ],
      },
      {
        path: 'projects/:uuid/versions/:version',
        name: 'VersionSubTap',
        component: ProjectsVersions,
        meta: {
          helpText: {
            en: Help.OverviewEn,
            de: Help.OverviewEn,
          },
        },
        children: [
          {
            path: 'overview/:currentSbom?',
            name: 'VersionOverview',
            component: ProjectsVersions,
            meta: {
              helpText: {
                en: Help.OverviewEn,
                de: Help.OverviewDe,
              },
            },
          },
          {
            path: 'overallReviews/:currentSbom?',
            name: 'overallReviews',
            component: ProjectsVersions,
            meta: {
              helpText: {
                en: Help.OverviewEn,
                de: Help.OverviewDe,
              },
            },
          },
          {
            path: 'component/NOT_SET/:currentSbom?',
            redirect: {name: 'Component'},
          },
          {
            path: 'component/:currentSbom?/:componentId?',
            name: 'Component',
            component: ProjectsVersions,
            meta: {
              helpText: {
                en: Help.ComponentsEn,
                de: Help.ComponentsDe,
              },
            },
          },
          {
            path: 'history/:currentSbom?',
            name: 'Sbomhistory/:currentSbom?',
            component: ProjectsVersions,
            meta: {
              helpText: {
                en: Help.SBOMDeliveriesEn,
                de: Help.SBOMDeliveriesDe,
              },
            },
          },
          {
            path: 'sbomCompare/:currentSbom?',
            name: 'SbomCompare',
            component: ProjectsVersions,
            meta: {
              helpText: {
                en: Help.SBOMCompareEn,
                de: Help.SBOMCompareDe,
              },
            },
          },
          {
            path: 'source/:currentSbom?',
            name: 'Source',
            component: ProjectsVersions,
            meta: {
              helpText: {
                en: Help.SourceCodeEn,
                de: Help.SourceCodeDe,
              },
            },
          },
          {
            path: 'sbomQuality/:currentSbom?',
            name: 'SbomQuality',
            meta: {
              helpText: {
                en: Help.QualityScanRemarksEn,
                de: Help.QualityScanRemarksDe,
              },
            },
            component: ProjectsVersions,
            children: [
              {
                path: 'scanRemarks',
                name: 'scanRemarks',
                component: ProjectsVersions,
                meta: {
                  helpText: {
                    en: Help.QualityScanRemarksEn,
                    de: Help.QualityScanRemarksDe,
                  },
                },
              },
              {
                path: 'licenseRemarks',
                name: 'licenseRemarks',
                component: ProjectsVersions,
                meta: {
                  helpText: {
                    en: Help.QualityLicenseRemarksEn,
                    de: Help.QualityLicenseRemarksDe,
                  },
                },
              },
              {
                path: 'generalRemarks',
                name: 'generalRemarks',
                component: ProjectsVersions,
                meta: {
                  helpText: {
                    en: Help.QualityGeneralRemarksEn,
                    de: Help.QualityGeneralRemarksDe,
                  },
                },
              },
              {
                path: 'reviewRemarks',
                name: 'reviewRemarks',
                component: ProjectsVersions,
                meta: {
                  helpText: {
                    en: Help.QualityReviewRemarksEn,
                    de: Help.QualityReviewRemarksDe,
                  },
                },
              },
            ],
          },
          {
            path: 'notice/:currentSbom?',
            name: 'Notice',
            component: ProjectsVersions,
            meta: {
              helpText: {
                en: Help.ThirdPartyNoticeEn,
                de: Help.ThirdPartyNoticeDe,
              },
            },
          },
          {
            path: 'auditLog/:currentSbom?',
            name: 'VersionAuditLog',
            component: ProjectsVersions,
          },
        ],
      },
      {
        path: 'projects/:uuid/:tab?',
        name: 'Project',
        component: ProjectsDetail,
        meta: {
          helpText: {
            en: Help.ProjectOverviewEn,
            de: Help.ProjectOverviewDe,
          },
        },
        children: [
          {
            path: 'overview',
            name: 'Overview',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.ProjectOverviewEn,
                de: Help.ProjectOverviewDe,
              },
            },
          },
          {
            path: 'versionlist',
            name: 'VersionList',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.VersionListEn,
                de: Help.VersionListDe,
              },
            },
          },
          {
            path: 'users',
            name: 'ProjectUsers',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.ProjectUsersEn,
                de: Help.ProjectUsersDe,
              },
            },
          },
          {
            path: 'tokens',
            name: 'Tokens',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.TokenManagementEn,
                de: Help.TokenManagementDe,
              },
            },
          },
          {
            path: 'policyrules',
            name: 'PolicyRules',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.ProjectPolicyRulesEn,
                de: Help.ProjectPolicyRulesDe,
              },
            },
          },
          {
            path: 'decisions',
            name: 'Decisions',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.LicenseDecisionsEn,
                de: Help.LicenseDecisionsDe,
              },
            },
          },
          {
            path: 'licenserules',
            redirect: 'decisions',
          },
          {
            path: 'approvals',
            name: 'Approvals',
            component: ProjectsDetail,
            meta: {
              helpText: {
                en: Help.ApprovalsEn,
                de: Help.ApprovalsDe,
              },
            },
          },
          {
            path: 'auditLog',
            name: 'ProjectAuditLog',
            component: ProjectsDetail,
          },
        ],
      },
    ],
  },
  {
    path: '/oauth/callback',
    name: 'OAuthCallback',
    component: Home,
    beforeEnter: (to, from, next) => {
      profileService.getProfileData().then((profileData) => {
        useUserStore().setSimpleProfileData(profileData);
        let item = localStorage.getItem(currentRouteKey);
        if (!item) {
          item = '/dashboard/home';
        }
        return next({path: item});
      });
    },
  },
];

const router = createRouter({
  history: createWebHashHistory(baseUrl),
  routes,
});

router.beforeEach((to, from, next) => {
  const rights = useUserStore().getRights as Rights;
  if (rights && Object.keys(rights).length > 0) {
    if (to.path.includes('dashboard/licenses')) {
      if (rights.allowLicense.read) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/classifications')) {
      if (rights.hasClassificationsAccess()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/projects')) {
      if (rights.hasAllProjectsReadonly()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/basicauth')) {
      if (rights.isDomainAdmin()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/labels')) {
      if (rights.hasLabelAccess()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/schemas')) {
      if (rights.hasSchemaAccess()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/userManagement')) {
      if (rights.isDomainAdmin()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/tools')) {
      if (rights.hasToolsAccess() || rights.hasSampleDataAccess() || rights.hasUsersAccess()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/newsbox')) {
      if (rights.isApplicationAdmin() || rights.isDomainAdmin()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/featureflags')) {
      if (rights.isApplicationAdmin()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/deletions')) {
      if (rights.isDomainAdmin()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/templates')) {
      if (!rights.hasReviewTemplatesAcces()) {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/templates/review')) {
      if (!rights.hasReviewTemplatesAcces()) {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin/users')) {
      if (!rights.hasUsersAccess()) {
        return next({path: '/dashboard/home'});
      }
    }
    if (to.path.includes('admin')) {
      if (rights.isAnyOfAdmin()) {
        return next();
      } else {
        return next({path: '/dashboard/home'});
      } //
    }
  }

  return next();
});

router.beforeEach((to, from, next) => {
  if (!to.fullPath.includes('oauth/callback')) {
    localStorage.setItem(currentRouteKey, to.fullPath);
  }

  next();
});
export default router;
