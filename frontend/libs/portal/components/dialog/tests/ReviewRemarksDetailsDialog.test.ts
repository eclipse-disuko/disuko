// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import GridReviewComments from '@disclosure-portal/components/grids/GridReviewComments.vue';
import {ReviewRemark, ReviewRemarkStatus} from '@disclosure-portal/model/Quality';
import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {flushPromises, mount} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import ReviewRemarksDetailsDialog from '../ReviewRemarksDetailsDialog.vue';

const mocks = vi.hoisted(() => ({
  commentOnRemark: vi.fn(),
  setRemarkStatus: vi.fn(),
  fetchRemarks: vi.fn(),
  snack: vi.fn(),
  error: vi.fn(),
  validate: vi.fn(),
}));

vi.mock('vuetify/components', () => ({VForm: {}}));
vi.mock('@disclosure-portal/stores/reviewRemarks.store', () => ({
  useReviewRemarksStore: () => ({...mocks, getRemark: () => undefined}),
}));
vi.mock('@disclosure-portal/stores/project.store', () => ({
  useProjectStore: () => ({currentProject: {isDeprecated: false}}),
}));
vi.mock('@shared/user/stores/user.store', () => ({
  useUserStore: () => ({getRights: {allowLicense: {read: true}}}),
}));
vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({info: mocks.snack}),
}));

const slotStub = {template: '<div><slot /></div>'};

const createWrapper = async (status = ReviewRemarkStatus.OPEN) => {
  const wrapper = mount(ReviewRemarksDetailsDialog, {
    props: {projectUuid: 'project', versionUuid: 'version'},
    global: {
      config: {errorHandler: mocks.error},
      components: {GridReviewComments},
      stubs: {
        ...vuetifyStubs,
        DialogLayout: {
          template: '<div><slot name="title-right" /><slot /></div>',
        },
        Stack: slotStub,
        Tooltip: slotStub,
        'v-dialog': {props: ['modelValue'], template: '<div v-if="modelValue"><slot /></div>'},
        'v-tabs': slotStub,
        'v-tab': slotStub,
        'v-tabs-window': slotStub,
        'v-tabs-window-item': slotStub,
        'v-data-table': true,
        'v-form': {
          template: '<form><slot /></form>',
          methods: {validate: mocks.validate, resetValidation: vi.fn()},
        },
        DCActionButton: {
          props: ['text', 'disabled'],
          template: '<button type="button" :disabled="disabled">{{ text }}</button>',
        },
      },
    },
  });
  wrapper.vm.open(Object.assign(new ReviewRemark(), {key: 'remark', title: 'Review', status, events: []}));
  await flushPromises();
  return wrapper;
};

type Wrapper = Awaited<ReturnType<typeof createWrapper>>;
const button = (wrapper: Wrapper, text: string) => wrapper.findAll('button').find((entry) => entry.text() === text)!;

describe('Review remark close with comment', () => {
  beforeEach(() => {
    vi.resetAllMocks();
    mocks.validate.mockResolvedValue({valid: true});
  });

  it('requires a non-blank draft, not an existing comment', async () => {
    const wrapper = await createWrapper();
    expect(button(wrapper, 'Close with comment').attributes('disabled')).toBeDefined();
    await wrapper.get('textarea').setValue('   ');
    expect(button(wrapper, 'Close with comment').attributes('disabled')).toBeDefined();
    await wrapper.get('textarea').setValue('Resolved in latest SBOM');
    expect(button(wrapper, 'Close with comment').attributes('disabled')).toBeUndefined();
  });

  it.each([ReviewRemarkStatus.OPEN, ReviewRemarkStatus.IN_PROGRESS])(
    'saves the final comment before closing a %s remark',
    async (status) => {
      const wrapper = await createWrapper(status);
      let finishComment!: () => void;
      mocks.commentOnRemark.mockImplementation(
        () =>
          new Promise<void>((resolve) => {
            finishComment = resolve;
          }),
      );
      await wrapper.get('textarea').setValue('Resolved in latest SBOM');
      await button(wrapper, 'Close with comment').trigger('click');
      await flushPromises();

      expect(mocks.commentOnRemark).toHaveBeenCalledWith('project', 'version', 'remark', {
        content: 'Resolved in latest SBOM',
      });
      expect(mocks.setRemarkStatus).not.toHaveBeenCalled();
      expect(button(wrapper, 'Close with comment').attributes('disabled')).toBeDefined();
      expect(button(wrapper, 'Comment').attributes('disabled')).toBeDefined();
      expect(button(wrapper, 'Close remark').attributes('disabled')).toBeDefined();

      finishComment();
      await flushPromises();
      expect(mocks.setRemarkStatus).toHaveBeenCalledWith('project', 'version', 'remark', {
        status: ReviewRemarkStatus.CLOSED,
      });
      expect(mocks.snack).toHaveBeenCalledWith('remark closed');
      expect(wrapper.find('textarea').exists()).toBe(false);
    },
  );

  it('does not close when comment validation fails', async () => {
    mocks.validate.mockResolvedValue({valid: false});
    const wrapper = await createWrapper();
    await wrapper.get('textarea').setValue('No');
    await button(wrapper, 'Close with comment').trigger('click');
    await flushPromises();
    expect(mocks.validate).toHaveBeenCalled();
    expect(mocks.commentOnRemark).not.toHaveBeenCalled();
    expect(mocks.setRemarkStatus).not.toHaveBeenCalled();
  });

  it('does not close when saving the comment fails', async () => {
    const error = new Error('Comment failed');
    mocks.commentOnRemark.mockRejectedValue(error);
    const wrapper = await createWrapper();
    await wrapper.get('textarea').setValue('Resolved in latest SBOM');
    await button(wrapper, 'Close with comment').trigger('click');
    await flushPromises();
    expect(mocks.error).toHaveBeenCalledWith(error, expect.anything(), expect.anything());
    expect(mocks.setRemarkStatus).not.toHaveBeenCalled();
    expect(wrapper.find('textarea').exists()).toBe(true);
    expect(button(wrapper, 'Comment').attributes('disabled')).toBeUndefined();
  });

  it('keeps the dialog open when closing fails after saving the comment', async () => {
    const error = new Error('Close failed');
    mocks.setRemarkStatus.mockRejectedValue(error);
    const wrapper = await createWrapper();
    await wrapper.get('textarea').setValue('Resolved in latest SBOM');
    await button(wrapper, 'Close with comment').trigger('click');
    await flushPromises();
    expect(mocks.commentOnRemark).toHaveBeenCalledTimes(1);
    expect(mocks.error).toHaveBeenCalledWith(error, expect.anything(), expect.anything());
    expect(mocks.snack).not.toHaveBeenCalled();
    expect(wrapper.find('textarea').exists()).toBe(true);
    expect(button(wrapper, 'Close remark').attributes('disabled')).toBeUndefined();
  });

  it('preserves comment-only and close-without-comment actions', async () => {
    const wrapper = await createWrapper();
    await wrapper.get('textarea').setValue('Still investigating');
    await button(wrapper, 'Comment').trigger('click');
    await flushPromises();
    expect(mocks.commentOnRemark).toHaveBeenCalledTimes(1);
    expect(mocks.setRemarkStatus).not.toHaveBeenCalled();
    await button(wrapper, 'Close remark').trigger('click');
    expect(wrapper.emitted('closeRemark')?.[0]?.[0]).toMatchObject({key: 'remark'});
  });

  it.each([ReviewRemarkStatus.CLOSED, ReviewRemarkStatus.CANCELLED])(
    'hides comment actions for %s remarks',
    async (status) => {
      const wrapper = await createWrapper(status);
      expect(wrapper.find('textarea').exists()).toBe(false);
      expect(wrapper.findAll('button').some((entry) => entry.text() === 'Close with comment')).toBe(false);
    },
  );
});
