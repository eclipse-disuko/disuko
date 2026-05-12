// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Group} from '@shared/types/Rights';
import {useUserStore} from '@disclosure-portal/stores/user';
import {useRouter} from 'vue-router';

type AccessChecker = () => boolean;

export class RightsUtils {
  static rights = function () {
    // TODO: Risky to use store here, use a composable instead
    const userStore = useUserStore();
    return userStore.getRights;
  };

  static redirectNonAdmins = function () {
    if (!RightsUtils.isAnyOfAdmin()) {
      // TODO: Risky to use router here, use a composable instead
      const router = useRouter();
      router.replace({path: '/dashboard/home'}).finally();
    }
  };

  static redirectRestrictedAccess = function (accessChecker: AccessChecker) {
    if (!accessChecker()) {
      // TODO: Risky to use router here, use a composable instead
      const router = useRouter();
      router.replace({path: '/dashboard/home'}).finally();
    }
  };

  static redirectDisabledUser = function () {
    // TODO: Risky to use router here, use a composable instead
    const router = useRouter();
    router.replace({path: '/dashboard'}).finally();
  };

  static hasLicenseAccess = function (): boolean {
    const rights = RightsUtils.rights();
    return (
      rights.allowLicense && (rights.allowLicense.create || rights.allowLicense.update || rights.allowLicense.delete)
    );
  };

  static hasClassificationsAccess = function (): boolean {
    const rights = RightsUtils.rights();
    return (
      rights.allowObligation &&
      (rights.allowObligation.create || rights.allowObligation.update || rights.allowObligation.delete)
    );
  };

  static hasPolicyAccess = function (): boolean {
    const rights = RightsUtils.rights();
    return rights.allowPolicy && (rights.allowPolicy.create || rights.allowPolicy.update || rights.allowPolicy.delete);
  };

  static hasAllProjectsReadonly = function (): boolean {
    const rights = RightsUtils.rights();
    return rights.allowProject && rights.allowProject.read;
  };

  static hasLabelAccess = function (): boolean {
    const rights = RightsUtils.rights();
    return rights.allowLabel && (rights.allowLabel.create || rights.allowLabel.update || rights.allowLabel.delete);
  };

  static hasSchemaAccess = function (): boolean {
    const rights = RightsUtils.rights();
    return rights.allowSchema && (rights.allowSchema.create || rights.allowSchema.update || rights.allowSchema.delete);
  };

  static hasToolsAccess = function (): boolean {
    const rights = RightsUtils.rights();
    return rights.allowTools && (rights.allowTools.create || rights.allowTools.update || rights.allowTools.delete);
  };

  static hasSampleDataAccess = function (): boolean {
    const rights = RightsUtils.rights();
    return (
      rights.allowSampleData &&
      (rights.allowSampleData.create || rights.allowSampleData.update || rights.allowSampleData.delete)
    );
  };

  static hasStyleguideAccess = function (): boolean {
    const rights = RightsUtils.rights();
    return (
      rights.allowStyleguide &&
      (rights.allowStyleguide.create || rights.allowStyleguide.update || rights.allowStyleguide.delete)
    );
  };

  static hasUsersAccess = function (): boolean {
    const rights = RightsUtils.rights();
    return (
      rights.allowUsers &&
      rights.allowUsers.create &&
      rights.allowUsers.read &&
      rights.allowUsers.update &&
      rights.allowUsers.delete
    );
  };

  static hasReviewTemplatesAccess = function (): boolean {
    const rights = RightsUtils.rights();
    return rights.hasReviewTemplatesAcces();
  };

  static isAnyOfAdmin = function (): boolean {
    return (
      RightsUtils.hasClassificationsAccess() ||
      RightsUtils.hasPolicyAccess() ||
      RightsUtils.hasAllProjectsReadonly() ||
      RightsUtils.hasLabelAccess() ||
      RightsUtils.hasSchemaAccess() ||
      RightsUtils.hasToolsAccess() ||
      RightsUtils.hasSampleDataAccess() ||
      RightsUtils.hasStyleguideAccess() ||
      RightsUtils.hasUsersAccess()
    );
  };

  static hasRole = function (group: Group): boolean {
    const rights = RightsUtils.rights();
    return rights.groups.includes(group);
  };

  static isLicenseManager = (): boolean => {
    return this.hasRole(Group.UserLicenseManager);
  };

  static isPolicyManager = (): boolean => {
    return this.hasRole(Group.UserPolicyManager);
  };

  static isProjectAnalyst = (): boolean => {
    return this.hasRole(Group.UserProjectAnalyst);
  };

  static isDomainAdmin = (): boolean => {
    return this.hasRole(Group.UserDomainAdmin);
  };

  static isApplicationAdmin = (): boolean => {
    return this.hasRole(Group.UserApplicationAdmin);
  };

  static isInternalUser = (): boolean => {
    return this.hasRole(Group.UserInternal);
  };

  static isNonInternalUser = (): boolean => {
    return this.hasRole(Group.UserNonInternal);
  };

  static isFOSSOffice = (): boolean => {
    return this.hasRole(Group.UserFOSSOffice);
  };
}
