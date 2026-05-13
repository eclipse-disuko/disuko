// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {BaseDto} from '@disclosure-portal/model/BaseClass';
import type {
  BucketDefinition,
  CalculatedPolicyConfig,
  CalculatedPolicyScope,
} from '@disclosure-portal/model/CalculatedPolicyRules';

export default class PolicyRule extends BaseDto {
  public status = '';
  public name = '';
  public description = '';
  public componentsAllow: string[] = [];
  public componentsDeny: string[] = [];
  public componentsWarn: string[] = [];
  public labelSets: string[][] = [];
  public auxiliary: boolean = false;
  public deprecated: boolean = false;
  public deprecatedDate = '';
  public active: boolean = true;
  public applyToAll: boolean = false;
  public calculated: boolean = false;
  public calculatedConfig: CalculatedPolicyConfig = {
    bucketDefinition: {
      deniedClassifications: [],
      warnedClassifications: [],
      allowedClassifications: [],
    },
    licenseScope: {
      isLicenseChart: [],
      approvalState: [],
      family: [],
      licenseType: [],
      source: [],
    },
  };

  public constructor(dto: PolicyRule | null | undefined = null) {
    super(dto);
    if (dto) {
      Object.assign(this, dto);
    }
    if (!this.componentsAllow) {
      this.componentsAllow = [];
    }
    if (!this.componentsDeny) {
      this.componentsDeny = [];
    }
    if (!this.componentsWarn) {
      this.componentsWarn = [];
    }
    if (!this.labelSets) {
      this.labelSets = [];
    }
    if (this.labelSets.length < 1) {
      this.labelSets[0] = [];
    }

    const config = this.calculatedConfig ?? ({} as CalculatedPolicyConfig);

    const getLicenseScope = (): CalculatedPolicyScope => {
      return {
        isLicenseChart: config.licenseScope?.isLicenseChart ?? [],
        approvalState: config.licenseScope?.approvalState ?? [],
        family: config.licenseScope?.family ?? [],
        licenseType: config.licenseScope?.licenseType ?? [],
        source: config.licenseScope?.source ?? [],
      };
    };

    const getBucketDefinition = (): BucketDefinition => {
      return {
        deniedClassifications: config.bucketDefinition?.deniedClassifications ?? [],
        warnedClassifications: config.bucketDefinition?.warnedClassifications ?? [],
        allowedClassifications: config.bucketDefinition?.allowedClassifications ?? [],
      };
    };

    this.calculatedConfig = {
      bucketDefinition: getBucketDefinition(),
      licenseScope: getLicenseScope(),
    };
  }
}

export class PolicyRuleDto {
  public name = '';
  public description = '';
  public key = '';
  public created: number;
  public updated: number;

  public constructor() {
    this.created = new Date().getTime();
    this.updated = new Date().getTime();
  }
}

export enum PolicyState {
  ALLOW = 'allow',
  DENY = 'deny',
  WARN = 'warn',
  NOT_SET = 'NOT_SET',
  NOASSERTION = 'noassertion',
  QUESTIONED = 'questioned',
}

export const PolicyRules: PolicyState[] = [PolicyState.ALLOW, PolicyState.WARN, PolicyState.DENY];
export const PolicyStates: PolicyState[] = [
  PolicyState.NOT_SET,
  PolicyState.DENY,
  PolicyState.NOASSERTION,
  PolicyState.WARN,
  PolicyState.QUESTIONED,
  PolicyState.ALLOW,
];

export class PolicyRulesAssignmentsDto {
  public status = '';
  public key = '';
  public name = '';
  public description = '';
  public type = '';
}

export class PolicyRulesForLicenseDto {
  public id = '';
  public policyRulesAssignments = [] as PolicyRulesAssignmentsDto[];
}
