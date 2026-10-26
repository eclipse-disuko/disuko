// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {createPinia, setActivePinia} from 'pinia';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import {useRouter} from 'vue-router';
import {useEventKeysStore} from '@shared/stores/eventKeys.store';
import {openUrlInNewTab} from '@shared/utils/url';
import {useUrls} from './useUrls';

vi.mock('@shared/utils/url', () => ({openUrlInNewTab: vi.fn()}));

describe('openProjectUrlByKey', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it.each([
    {isGroup: false, key: 'project-id', path: '/dashboard/projects/project-id/overview'},
    {isGroup: true, key: 'group-id', path: '/dashboard/groups/group-id/overview'},
    {isGroup: false, key: 'id/with space', path: '/dashboard/projects/id%2Fwith%20space/overview'},
    {isGroup: true, key: 'id/with space', path: '/dashboard/groups/id%2Fwith%20space/overview'},
  ])('navigates to $path', ({isGroup, key, path}) => {
    const router = useRouter();
    const {openProjectUrlByKey} = useUrls();

    openProjectUrlByKey(key, router, isGroup);

    expect(router.push).toHaveBeenCalledWith({path});
    expect(openUrlInNewTab).not.toHaveBeenCalled();
  });

  it.each([
    {isGroup: false, path: '/dashboard/projects/item-id/overview'},
    {isGroup: true, path: '/dashboard/groups/item-id/overview'},
  ])('opens $path in a new tab when Ctrl is pressed', ({isGroup, path}) => {
    const router = useRouter();
    useEventKeysStore().controlIsPressed = true;
    const {openProjectUrlByKey} = useUrls();

    openProjectUrlByKey('item-id', router, isGroup);

    expect(openUrlInNewTab).toHaveBeenCalledWith(path);
    expect(router.push).not.toHaveBeenCalled();
  });

  it('defaults to project navigation when isGroup is omitted', () => {
    const router = useRouter();
    const {openProjectUrlByKey} = useUrls();

    openProjectUrlByKey('project-id', router);

    expect(router.push).toHaveBeenCalledWith({path: '/dashboard/projects/project-id/overview'});
  });
});
