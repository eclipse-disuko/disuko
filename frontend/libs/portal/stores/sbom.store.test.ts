import {createPinia, setActivePinia} from 'pinia';
import {
  type GeneralStats,
  type SbomStats,
  type SpdxFile,
  type VersionSlim,
} from '@disclosure-portal/model/VersionDetails';
import {type VersionSbomsFlat} from '@disclosure-portal/model/ProjectsResponse';
import {beforeEach, describe, expect, it, vi} from 'vitest';

const {projectStoreMock, versionServiceMock, projectServiceMock} = vi.hoisted(() => ({
  projectStoreMock: {
    currentProject: {
      _key: 'project-1',
      isGroup: false,
      versions: {},
    },
  },
  versionServiceMock: {
    getSBOMStats: vi.fn(),
    getGeneralVersionStats: vi.fn(),
  },
  projectServiceMock: {
    getAllSbomsFlat: vi.fn(),
  },
}));

vi.mock('@disclosure-portal/stores/project.store', () => ({
  useProjectStore: () => projectStoreMock,
}));

vi.mock('@disclosure-portal/services/version', () => ({
  default: versionServiceMock,
}));

vi.mock('@disclosure-portal/services/projects', () => ({
  default: projectServiceMock,
}));

import {useSbomStore} from './sbom.store';

const version = (key: string, name: string): VersionSlim => ({_key: key, name}) as VersionSlim;
const sbomStats = (allowed: number): SbomStats => ({policyState: {allowed: allowed}}) as SbomStats;
const generalStats = (acceptable: number): GeneralStats => ({reviewRemark: {acceptable: acceptable}}) as GeneralStats;
const flatItem = (key: string, versionKey: string, versionName: string): VersionSbomsFlat =>
  ({_key: key, versionKey, versionName}) as VersionSbomsFlat;

function deferred<T>() {
  let resolve!: (value: T) => void;
  let reject!: (reason?: unknown) => void;
  const promise = new Promise<T>((res, rej) => {
    resolve = res;
    reject = rej;
  });
  return {promise, resolve, reject};
}

describe('useSbomStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    versionServiceMock.getSBOMStats.mockReset();
    versionServiceMock.getGeneralVersionStats.mockReset();
    projectStoreMock.currentProject = {
      _key: 'project-1',
      isGroup: false,
      versions: {
        versionA: {_key: 'versionA', name: 'Version A'},
        versionB: {_key: 'versionB', name: 'Version B'},
      },
    };
  });

  it('reuses current sbom stats after they are loaded', async () => {
    const store = useSbomStore();
    store.setCurrentVersion('versionA');
    store.setSelectedSBOMKey('spdxA');

    versionServiceMock.getSBOMStats.mockResolvedValueOnce({data: sbomStats(3)});

    await store.fetchSBOMStats('spdxA');

    expect(versionServiceMock.getSBOMStats).toHaveBeenCalledTimes(1);
    expect(store.getSbomStats).toEqual({policyState: {allowed: 3}});

    await store.fetchSBOMStats('spdxA');

    expect(versionServiceMock.getSBOMStats).toHaveBeenCalledTimes(1);
    expect(store.getSbomStats).toEqual({policyState: {allowed: 3}});
  });

  it('does not deduplicate concurrent sbom requests before stats are loaded', async () => {
    const store = useSbomStore();
    store.setCurrentVersion('versionA');
    store.setSelectedSBOMKey('spdxA');

    const firstPending = deferred<{data: SbomStats}>();
    const secondPending = deferred<{data: SbomStats}>();
    versionServiceMock.getSBOMStats.mockReturnValueOnce(firstPending.promise);
    versionServiceMock.getSBOMStats.mockReturnValueOnce(secondPending.promise);

    const first = store.fetchSBOMStats('spdxA');
    const second = store.fetchSBOMStats('spdxA');

    expect(versionServiceMock.getSBOMStats).toHaveBeenCalledTimes(2);

    firstPending.resolve({data: sbomStats(1)});
    secondPending.resolve({data: sbomStats(2)});
    await Promise.all([first, second]);

    expect(store.getSbomStats).toEqual({policyState: {allowed: 2}});
  });

  it('clears only sbom stats when the selected SPDX changes', () => {
    const store = useSbomStore();
    store.setCurrentVersion('versionA');
    store.sbomStats = sbomStats(1);
    store.generalStats = generalStats(2);

    store.setSelectedSBOMKey('spdxB');

    expect(store.getSbomStats).toEqual({});
    expect(store.getGeneralStats).toEqual({reviewRemark: {acceptable: 2}});
  });

  it('clears both stat payloads when the version changes', () => {
    const store = useSbomStore();
    store.setCurrentVersion('versionA');
    store.sbomStats = sbomStats(1);
    store.generalStats = generalStats(2);

    store.setCurrentVersion('versionB');

    expect(store.getSbomStats).toEqual({});
    expect(store.getGeneralStats).toEqual({});
  });

  it('ignores stale sbom responses after the selected SPDX changes', async () => {
    const store = useSbomStore();
    store.setCurrentVersion('versionA');
    store.setSelectedSBOMKey('spdxA');

    const oldRequest = deferred<{data: SbomStats}>();
    versionServiceMock.getSBOMStats.mockReturnValueOnce(oldRequest.promise);
    const oldPromise = store.fetchSBOMStats('spdxA');

    store.setSelectedSBOMKey('spdxB');

    const newRequest = deferred<{data: SbomStats}>();
    versionServiceMock.getSBOMStats.mockReturnValueOnce(newRequest.promise);
    const newPromise = store.fetchSBOMStats('spdxB');

    oldRequest.resolve({data: sbomStats(1)});
    await oldPromise;
    expect(store.getSbomStats).toEqual({});

    newRequest.resolve({data: sbomStats(9)});
    await newPromise;
    expect(store.getSbomStats).toEqual({policyState: {allowed: 9}});
  });

  it('ignores stale general stats responses after the version changes', async () => {
    const store = useSbomStore();
    store.setCurrentVersion('versionA');

    const oldRequest = deferred<{data: GeneralStats}>();
    versionServiceMock.getGeneralVersionStats.mockReturnValueOnce(oldRequest.promise);
    const oldPromise = store.fetchGeneralVersionStats();

    store.setCurrentVersion('versionB');

    const newRequest = deferred<{data: GeneralStats}>();
    versionServiceMock.getGeneralVersionStats.mockReturnValueOnce(newRequest.promise);
    const newPromise = store.fetchGeneralVersionStats();

    oldRequest.resolve({data: generalStats(1)});
    await oldPromise;
    expect(store.getGeneralStats).toEqual({});

    newRequest.resolve({data: generalStats(4)});
    await newPromise;
    expect(store.getGeneralStats).toEqual({reviewRemark: {acceptable: 4}});
  });

  it('reuses current general stats after they are loaded', async () => {
    const store = useSbomStore();
    store.setCurrentVersion('versionA');

    versionServiceMock.getGeneralVersionStats.mockResolvedValueOnce({data: generalStats(7)});

    await store.fetchGeneralVersionStats();

    expect(versionServiceMock.getGeneralVersionStats).toHaveBeenCalledTimes(1);
    expect(store.getGeneralStats).toEqual({reviewRemark: {acceptable: 7}});

    await store.fetchGeneralVersionStats();

    expect(versionServiceMock.getGeneralVersionStats).toHaveBeenCalledTimes(1);
    expect(store.getGeneralStats).toEqual({reviewRemark: {acceptable: 7}});
  });

  describe('channelSpdxs (derived)', () => {
    it('returns only items for the current version', () => {
      const store = useSbomStore();
      store.allSBOMSFlat = [
        flatItem('spdx-1', 'versionA', 'Version A'),
        flatItem('spdx-2', 'versionA', 'Version A'),
        flatItem('spdx-3', 'versionB', 'Version B'),
      ];
      store.setCurrentVersion('versionA');

      expect(store.channelSpdxs.map((s) => s._key)).toEqual(['spdx-1', 'spdx-2']);
    });

    it('sets isRecent on the first item only', () => {
      const store = useSbomStore();
      store.allSBOMSFlat = [flatItem('spdx-1', 'versionA', 'Version A'), flatItem('spdx-2', 'versionA', 'Version A')];
      store.setCurrentVersion('versionA');

      expect(store.channelSpdxs[0].isRecent).toBe(true);
      expect(store.channelSpdxs[1].isRecent).toBe(false);
    });

    it('returns an empty array when no items match the current version', () => {
      const store = useSbomStore();
      store.allSBOMSFlat = [flatItem('spdx-3', 'versionB', 'Version B')];
      store.setCurrentVersion('versionA');

      expect(store.channelSpdxs).toEqual([]);
    });
  });

  describe('getSelectedSBOM (derived)', () => {
    it('returns the selected SBOM with isRecent true when it is the most recent', () => {
      const store = useSbomStore();
      store.allSBOMSFlat = [flatItem('spdx-1', 'versionA', 'Version A'), flatItem('spdx-2', 'versionA', 'Version A')];
      store.setCurrentVersion('versionA');
      store.setSelectedSBOMKey('spdx-1');

      expect(store.getSelectedSBOM?._key).toBe('spdx-1');
      expect(store.getSelectedSBOM?.isRecent).toBe(true);
    });

    it('returns the selected SBOM with isRecent false when it is not the most recent', () => {
      const store = useSbomStore();
      store.allSBOMSFlat = [flatItem('spdx-1', 'versionA', 'Version A'), flatItem('spdx-2', 'versionA', 'Version A')];
      store.setCurrentVersion('versionA');
      store.setSelectedSBOMKey('spdx-2');

      expect(store.getSelectedSBOM?._key).toBe('spdx-2');
      expect(store.getSelectedSBOM?.isRecent).toBe(false);
    });
  });

  describe('getAllSBOMs (derived)', () => {
    it('groups flat items into VersionSboms by versionKey', () => {
      const store = useSbomStore();
      store.allSBOMSFlat = [
        flatItem('spdx-1', 'versionA', 'Version A'),
        flatItem('spdx-2', 'versionA', 'Version A'),
        flatItem('spdx-3', 'versionB', 'Version B'),
      ];

      const result = store.getAllSBOMs;

      expect(result).toHaveLength(2);
      expect(result[0].versionKey).toBe('versionA');
      expect(result[0].versionName).toBe('Version A');
      expect(result[0].spdxFileHistory.map((s) => s._key)).toEqual(['spdx-1', 'spdx-2']);
      expect(result[1].versionKey).toBe('versionB');
      expect(result[1].spdxFileHistory.map((s) => s._key)).toEqual(['spdx-3']);
    });

    it('preserves insertion order of versions', () => {
      const store = useSbomStore();
      store.allSBOMSFlat = [flatItem('spdx-3', 'versionB', 'Version B'), flatItem('spdx-1', 'versionA', 'Version A')];

      const result = store.getAllSBOMs;

      expect(result[0].versionKey).toBe('versionB');
      expect(result[1].versionKey).toBe('versionA');
    });

    it('returns an empty array when allSBOMSFlat is empty', () => {
      const store = useSbomStore();
      expect(store.getAllSBOMs).toEqual([]);
    });
  });
});
