import {createPinia, setActivePinia} from 'pinia';
import {
  type GeneralStats,
  type SbomStats,
  type SpdxFile,
  type VersionSlim,
} from '@disclosure-portal/model/VersionDetails';
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
    getSbomHistory: vi.fn(),
  },
  projectServiceMock: {
    getAllSbomsFlat: vi.fn(),
    getAllSboms: vi.fn(),
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
const spdx = (key: string): SpdxFile => ({_key: key}) as SpdxFile;
const sbomStats = (allowed: number): SbomStats => ({PolicyState: {Allowed: allowed}}) as SbomStats;
const generalStats = (acceptable: number): GeneralStats => ({ReviewRemark: {Acceptable: acceptable}}) as GeneralStats;

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
    versionServiceMock.getSbomHistory.mockReset();
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
    store.setCurrentVersion(version('versionA', 'Version A'));
    store.setSelectedSpdx(spdx('spdxA'));

    versionServiceMock.getSBOMStats.mockResolvedValueOnce({data: sbomStats(3)});

    await store.fetchSBOMStats('spdxA');

    expect(versionServiceMock.getSBOMStats).toHaveBeenCalledTimes(1);
    expect(store.getSbomStats).toEqual({PolicyState: {Allowed: 3}});

    await store.fetchSBOMStats('spdxA');

    expect(versionServiceMock.getSBOMStats).toHaveBeenCalledTimes(1);
    expect(store.getSbomStats).toEqual({PolicyState: {Allowed: 3}});
  });

  it('does not deduplicate concurrent sbom requests before stats are loaded', async () => {
    const store = useSbomStore();
    store.setCurrentVersion(version('versionA', 'Version A'));
    store.setSelectedSpdx(spdx('spdxA'));

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

    expect(store.getSbomStats).toEqual({PolicyState: {Allowed: 2}});
  });

  it('clears only sbom stats when the selected SPDX changes', () => {
    const store = useSbomStore();
    store.setCurrentVersion(version('versionA', 'Version A'));
    store.sbomStats = sbomStats(1);
    store.generalStats = generalStats(2);

    store.setSelectedSpdx(spdx('spdxB'));

    expect(store.getSbomStats).toEqual({});
    expect(store.getGeneralStats).toEqual({ReviewRemark: {Acceptable: 2}});
  });

  it('clears both stat payloads when the version changes', () => {
    const store = useSbomStore();
    store.setCurrentVersion(version('versionA', 'Version A'));
    store.sbomStats = sbomStats(1);
    store.generalStats = generalStats(2);

    store.setCurrentVersion(version('versionB', 'Version B'));

    expect(store.getSbomStats).toEqual({});
    expect(store.getGeneralStats).toEqual({});
  });

  it('ignores stale sbom responses after the selected SPDX changes', async () => {
    const store = useSbomStore();
    store.setCurrentVersion(version('versionA', 'Version A'));
    store.setSelectedSpdx(spdx('spdxA'));

    const oldRequest = deferred<{data: SbomStats}>();
    versionServiceMock.getSBOMStats.mockReturnValueOnce(oldRequest.promise);
    const oldPromise = store.fetchSBOMStats('spdxA');

    store.setSelectedSpdx(spdx('spdxB'));

    const newRequest = deferred<{data: SbomStats}>();
    versionServiceMock.getSBOMStats.mockReturnValueOnce(newRequest.promise);
    const newPromise = store.fetchSBOMStats('spdxB');

    oldRequest.resolve({data: sbomStats(1)});
    await oldPromise;
    expect(store.getSbomStats).toEqual({});

    newRequest.resolve({data: sbomStats(9)});
    await newPromise;
    expect(store.getSbomStats).toEqual({PolicyState: {Allowed: 9}});
  });

  it('ignores stale general stats responses after the version changes', async () => {
    const store = useSbomStore();
    store.setCurrentVersion(version('versionA', 'Version A'));

    const oldRequest = deferred<{data: GeneralStats}>();
    versionServiceMock.getGeneralVersionStats.mockReturnValueOnce(oldRequest.promise);
    const oldPromise = store.fetchGeneralVersionStats();

    store.setCurrentVersion(version('versionB', 'Version B'));

    const newRequest = deferred<{data: GeneralStats}>();
    versionServiceMock.getGeneralVersionStats.mockReturnValueOnce(newRequest.promise);
    const newPromise = store.fetchGeneralVersionStats();

    oldRequest.resolve({data: generalStats(1)});
    await oldPromise;
    expect(store.getGeneralStats).toEqual({});

    newRequest.resolve({data: generalStats(4)});
    await newPromise;
    expect(store.getGeneralStats).toEqual({ReviewRemark: {Acceptable: 4}});
  });

  it('reuses current general stats after they are loaded', async () => {
    const store = useSbomStore();
    store.setCurrentVersion(version('versionA', 'Version A'));

    versionServiceMock.getGeneralVersionStats.mockResolvedValueOnce({data: generalStats(7)});

    await store.fetchGeneralVersionStats();

    expect(versionServiceMock.getGeneralVersionStats).toHaveBeenCalledTimes(1);
    expect(store.getGeneralStats).toEqual({ReviewRemark: {Acceptable: 7}});

    await store.fetchGeneralVersionStats();

    expect(versionServiceMock.getGeneralVersionStats).toHaveBeenCalledTimes(1);
    expect(store.getGeneralStats).toEqual({ReviewRemark: {Acceptable: 7}});
  });
});
