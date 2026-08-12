// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {mount} from '@vue/test-utils';
import {afterEach, beforeEach, describe, expect, it, vi} from 'vitest';
import {nextTick} from 'vue';
import DiscoFileUpload from '../DiscoFileUpload.vue';

const {uploadFormDataFileMock} = vi.hoisted(() => ({uploadFormDataFileMock: vi.fn()}));

vi.mock('@disclosure-portal/composables/useUpload', () => ({
  useUpload: () => ({uploadFormDataFile: uploadFormDataFileMock}),
}));

const buildFileList = (files: File[]): FileList => {
  const list = files as unknown as FileList & File[];
  Object.defineProperty(list, 'item', {value: (i: number) => files[i]});
  return list;
};

const setInputFiles = (input: HTMLInputElement, files: File[]) => {
  Object.defineProperty(input, 'files', {value: buildFileList(files), configurable: true});
};

describe('DiscoFileUpload', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    uploadFormDataFileMock.mockResolvedValue({data: {ok: true}});
  });

  const createWrapper = (props: Record<string, unknown> = {}) => {
    return mount(DiscoFileUpload, {
      props: {acceptTypes: '.zip', uploadTargetUrl: '/api/upload', ...props},
      attachTo: document.body,
    });
  };

  afterEach(() => {
    document.body.innerHTML = '';
  });

  it('forwards acceptTypes and multiple to the underlying input', () => {
    const wrapper = createWrapper({multiple: true});

    const input = wrapper.find('input').element as HTMLInputElement;
    expect(input.accept).toBe('.zip');
    expect(input.multiple).toBe(true);
  });

  it('emits filesChanged when a file is selected', async () => {
    const wrapper = createWrapper({directUpload: false});
    const input = wrapper.find('input').element as HTMLInputElement;
    const file = new File(['content'], 'test.zip');
    setInputFiles(input, [file]);

    await wrapper.find('input').trigger('change');

    expect(wrapper.emitted('filesChanged')?.[0][0]).toHaveLength(1);
  });

  it('auto-uploads a single selected file when directUpload is true and multiple is false', async () => {
    const wrapper = createWrapper();
    const input = wrapper.find('input').element as HTMLInputElement;
    const file = new File(['content'], 'test.zip');
    setInputFiles(input, [file]);

    await wrapper.find('input').trigger('change');
    await flushPromises();

    expect(uploadFormDataFileMock).toHaveBeenCalledTimes(1);
    expect(uploadFormDataFileMock.mock.calls[0][0]).toMatchObject({uploadUrl: '/api/upload', file});
  });

  it('does not auto-upload when multiple is true', async () => {
    const wrapper = createWrapper({multiple: true});
    const input = wrapper.find('input').element as HTMLInputElement;
    setInputFiles(input, [new File(['content'], 'test.zip')]);

    await wrapper.find('input').trigger('change');
    await flushPromises();

    expect(uploadFormDataFileMock).not.toHaveBeenCalled();
  });

  it('does not auto-upload when directUpload is false', async () => {
    const wrapper = createWrapper({directUpload: false});
    const input = wrapper.find('input').element as HTMLInputElement;
    setInputFiles(input, [new File(['content'], 'test.zip')]);

    await wrapper.find('input').trigger('change');
    await flushPromises();

    expect(uploadFormDataFileMock).not.toHaveBeenCalled();
  });

  it('emits reqProgress with a rounded percentage while uploading', async () => {
    uploadFormDataFileMock.mockImplementation(({onUploadProgress}) => {
      onUploadProgress({loaded: 50, total: 200});
      return Promise.resolve({data: {}});
    });
    const wrapper = createWrapper();
    const input = wrapper.find('input').element as HTMLInputElement;
    setInputFiles(input, [new File(['content'], 'test.zip')]);

    await wrapper.find('input').trigger('change');
    await flushPromises();

    expect(wrapper.emitted('reqProgress')?.[0][1]).toBe(25);
  });

  it('emits reqFinished with the file and response data on success, then clears the input', async () => {
    uploadFormDataFileMock.mockResolvedValue({data: {id: 'abc'}});
    const wrapper = createWrapper();
    const input = wrapper.find('input').element as HTMLInputElement;
    const file = new File(['content'], 'test.zip');
    setInputFiles(input, [file]);

    await wrapper.find('input').trigger('change');
    await flushPromises();

    expect(wrapper.emitted('reqFinished')?.[0]).toEqual([file, {id: 'abc'}]);
    expect(wrapper.emitted('filesChanged')).toHaveLength(2);
    expect(input.value).toBe('');
  });

  it('emits reqFailed when the upload rejects', async () => {
    uploadFormDataFileMock.mockRejectedValue(new Error('network error'));
    const wrapper = createWrapper();
    const input = wrapper.find('input').element as HTMLInputElement;
    setInputFiles(input, [new File(['content'], 'test.zip')]);

    await wrapper.find('input').trigger('change');
    await flushPromises();

    expect(wrapper.emitted('reqFailed')).toHaveLength(1);
  });

  it('passes the schema prop through to the upload call', async () => {
    const schema = {name: 's', version: '1', description: '', label: ''};
    const wrapper = createWrapper({schema});
    const input = wrapper.find('input').element as HTMLInputElement;
    setInputFiles(input, [new File(['content'], 'test.zip')]);

    await wrapper.find('input').trigger('change');
    await flushPromises();

    expect(uploadFormDataFileMock.mock.calls[0][0]).toMatchObject({schema});
  });

  it('clicks the hidden input when uploadClick is invoked', async () => {
    const wrapper = createWrapper();
    const input = wrapper.find('input').element as HTMLInputElement;
    const clickSpy = vi.spyOn(input, 'click');

    (wrapper.vm as unknown as {uploadClick: () => void}).uploadClick();
    await nextTick();

    expect(clickSpy).toHaveBeenCalled();
  });
});

const flushPromises = () => new Promise((resolve) => setTimeout(resolve));
