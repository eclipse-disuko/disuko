// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Capabilities} from '@disclosure-portal/model/Capabilities';
import {defineStore} from 'pinia';

const STORE_NAME = 'capabilities';

export const useCapabilitiesStore = defineStore(STORE_NAME, {
  state: (): Capabilities => ({applicationConnector: true}) as Capabilities,
});
