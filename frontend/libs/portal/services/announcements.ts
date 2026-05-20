// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import {AnnouncementsResponse} from '@disclosure-portal/model/AnnouncementsResponse';

const {api} = useApi();

const modelName = 'announcements';

class AnnouncementService {
  public getAll = () => api.get<AnnouncementsResponse[]>(`/api/v1/${modelName}`);
}

const announcementService = new AnnouncementService();
export default announcementService;
