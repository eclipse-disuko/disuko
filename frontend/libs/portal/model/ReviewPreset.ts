// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ProjectModel} from './Project';

export default class ReviewPreset {
  public EntryTitle = '';
  public Tooltip = '';
  public Reviewer = '';
  public Comment = '';
  public MustContainTags: string[] = [];
  public CheckCbs: ((p: ProjectModel) => Promise<boolean>)[] = [];
}
