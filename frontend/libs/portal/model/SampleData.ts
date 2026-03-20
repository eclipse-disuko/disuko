// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export default class SampleDataCreationState {
  public _key = '';
  public _rev = '';
  public isRunning = false;
  public reqID = '';
  public hasErrors = false;
  public lastError = '';
  public withFileUpload = false;
  public targetCount = 0;
  public createdCount = 0;
  public startTime = '';
  public endTime = '';

  public constructor(dto: SampleDataCreationState | null | undefined = null) {
    this.set(dto);
  }

  public set(dto: SampleDataCreationState | null | undefined) {
    if (dto) {
      Object.assign(this, dto);
    }
  }
}
