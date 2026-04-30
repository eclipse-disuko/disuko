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
  public Status = '';
  public Name = '';
  public Description = '';
  public ComponentsAllow: string[] = [];
  public ComponentsDeny: string[] = [];
  public ComponentsWarn: string[] = [];
  public LabelSets: string[][] = [];
  public Auxiliary: boolean = false;
  public Deprecated: boolean = false;
  public DeprecatedDate = '';
  public Active: boolean = true;
  public ApplyToAll: boolean = false;
  public Calculated: boolean = false;
  public CalculatedConfig: CalculatedPolicyConfig = {
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
    if (!this.ComponentsAllow) {
      this.ComponentsAllow = [];
    }
    if (!this.ComponentsDeny) {
      this.ComponentsDeny = [];
    }
    if (!this.ComponentsWarn) {
      this.ComponentsWarn = [];
    }
    if (!this.LabelSets) {
      this.LabelSets = [];
    }
    if (this.LabelSets.length < 1) {
      this.LabelSets[0] = [];
    }

    const config = this.CalculatedConfig ?? ({} as CalculatedPolicyConfig);

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

    this.CalculatedConfig = {
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
