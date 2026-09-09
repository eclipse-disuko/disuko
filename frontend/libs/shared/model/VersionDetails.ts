// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export type Nullable<T> = T | null;

export class SpdxFileSlim {
  public _key = '';
  public projectVersionId = '';
  public uploaded: Date = new Date();
  public updated = '';
  public name = '';
}

export class OverallReview {
  public created = '';
  public updated = '';
  public state = '';
  public comment = '';
  public sbomId = '';
  public sbomName = '';
  public sbomUploaded = '';
  public creator = '';
  public creatorFullName = '';
}

export class VersionSlimDto {
  public _key = '';
  public parentKey = '';
  public name = '1.0';
  public description = '';
  public created = '';
  public updated = '';
  public status = '';
  public currentSpdxFile: SpdxFileSlim = new SpdxFileSlim();
  public spdxFileHistory: SpdxFileSlim[] = [];
  public isDeleted = false;
  public overallReviews: OverallReview[] = [];
}

export class VersionSlim extends VersionSlimDto {
  constructor(dto: VersionSlimDto) {
    super();
    Object.assign(this, dto);
  }
}
