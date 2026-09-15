// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import PolicyRule from '@disclosure-portal/model/PolicyRule';
import {describe, expect, it} from 'vitest';

describe('PolicyRule', () => {
  it('preserves calculated policy rule configuration from imported data', () => {
    const importedData = Object.assign(new PolicyRule(), {
      name: 'calculated-rule',
      calculated: true,
      calculatedConfig: {
        bucketDefinition: {
          deniedClassifications: ['denied'],
          warnedClassifications: ['warned'],
          allowedClassifications: ['allowed'],
        },
        licenseScope: {
          isLicenseChart: [true],
          approvalState: ['approved'],
          family: ['permissive'],
          licenseType: ['open source'],
          source: ['spdx'],
        },
      },
    });
    const rule = new PolicyRule(importedData);

    expect(rule.calculated).toBe(true);
    expect(rule.calculatedConfig).toEqual({
      bucketDefinition: {
        deniedClassifications: ['denied'],
        warnedClassifications: ['warned'],
        allowedClassifications: ['allowed'],
      },
      licenseScope: {
        isLicenseChart: [true],
        approvalState: ['approved'],
        family: ['permissive'],
        licenseType: ['open source'],
        source: ['spdx'],
      },
    });
  });

  it('normalizes missing calculated policy rule configuration from imported data', () => {
    const importedData = Object.assign(new PolicyRule(), {
      calculated: true,
      calculatedConfig: {
        bucketDefinition: null,
        licenseScope: {
          source: ['custom'],
        },
      },
    });
    const rule = new PolicyRule(importedData);

    expect(rule.calculated).toBe(true);
    expect(rule.calculatedConfig).toEqual({
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
        source: ['custom'],
      },
    });
  });
});