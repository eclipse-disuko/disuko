// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {BaseDto} from '@disclosure-portal/model/BaseClass';

export class SystemStats extends BaseDto {
  public projectCount = 0;
  public projectActiveCount = 0;
  public projectDeletedCount = 0;
  public licenseActiveCount = 0;
  public licenseCount = 0;
  public licenseChartCount = 0;
  public licenseDeletedCount = 0;
  public policyRuleCount = 0;
  public policyRuleActiveCount = 0;
  public policyRuleDeletedCount = 0;
  public labelCount = 0;
  public schemaCount = 0;
  public obligationCount = 0;
  public obligationActiveCount = 0;
  public obligationDeletedCount = 0;
  public uploadFileCnt = 0;
  public uploadFileCntPDF = 0;
  public uploadFileCntJSON = 0;
  public uploadFileCntSBOM = 0;
  public dbBackupFileCnt = 0;

  public userCount = 0;
  public userActiveCount = 0;
  public userDeactivateCount = 0;
  public userTermsNotAcceptedCount = 0;
  public userDeprovisionedCount = 0;

  public maxVersionsInOneProject = 0;
  public projectsOverOrAtVersionLimit = 0;
  public versionLimit = 0;

  public missingProjects = false;
  public missingLicenses = false;
  public missingPolicyRules = false;
  public missingObligations = false;
  public missingUploadFiles = false;
  public missingUsers = false;

  public constructor(dto: SystemStats | null | undefined = null) {
    super(dto);
  }
}

export default class SystemStatsResponse {
  dayStats: SystemStats[] = [];
  monthsStats: SystemStats[] = [];
}
