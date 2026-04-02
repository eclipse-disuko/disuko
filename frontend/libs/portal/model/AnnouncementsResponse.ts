// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export enum AnnouncementTypes {
  licenseChange = 'license_change',
}

export interface IAnnouncementDto {
  _key: string;
  when: string;
  type: AnnouncementTypes;
  content: string;
}

export type LicenseAnnouncement = {
  licenseName: string;
  licenseId: string;
  changeType: 'license_forbidden' | 'license_chart' | 'license_family' | 'license_type' | 'custom_license_deleted';
  oldVal: string;
  newVal: string;
};

export type IAnnouncement = {
  _key: string;
  when: string;
} & {
  type: AnnouncementTypes.licenseChange;
  content: LicenseAnnouncement;
};

export class AnnouncementDto implements IAnnouncementDto {
  public _key!: string;
  public when!: string;
  public type!: AnnouncementTypes;
  public content!: string;
}

export class AnnouncementsResponse extends AnnouncementDto {
  constructor(dto: AnnouncementDto) {
    super();
    Object.assign(this, dto);
  }
}
