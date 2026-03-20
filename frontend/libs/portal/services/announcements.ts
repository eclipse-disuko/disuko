// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {getApi} from '@disclosure-portal/api';
import {AnnouncementsResponse} from '@disclosure-portal/model/AnnouncementsResponse';

const {api} = getApi();

const modelName = 'announcements';

class AnnouncementService {
  public getAll = () => api.get<AnnouncementsResponse[]>(`/api/v1/${modelName}`);
}

const announcementService = new AnnouncementService();
export default announcementService;
