// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import mitt from 'mitt';
import {type Events} from '@shared/types/eventBus';

const eventBus = mitt<Events>();

export default eventBus;
