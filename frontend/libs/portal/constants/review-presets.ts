// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ProjectModel} from '@disclosure-portal/model/Project';
import projectService from '@disclosure-portal/services/projects';
import ReviewPreset from '../model/ReviewPreset';

export const ReviewPresets: ReviewPreset[] = [
  {
    EntryTitle: 'PRESET_RD_TITLE',
    Tooltip: 'PRESET_RD_TT',
    Comment: 'PRESET_RD_COMM',
    MustContainTags: ['vehicle platform'],
    Reviewer: 'CUSTOMER1',
    CheckCbs: [],
  },
  {
    EntryTitle: 'PRESET_RD_TITLE',
    Tooltip: 'PRESET_RD_TT',
    Comment: 'PRESET_RD_COMM',
    Reviewer: 'CUSTOMER1',
    MustContainTags: [],
    CheckCbs: [
      (p: ProjectModel): Promise<boolean> => {
        if (!p.isGroup) {
          return Promise.resolve(false);
        }
        return (async () => {
          return (await projectService.getVehiclePlatformOnly(p._key)).data.found;
        })();
      },
    ],
  },
];
