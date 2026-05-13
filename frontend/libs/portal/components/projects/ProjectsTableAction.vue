<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import {Rights} from '@shared/types/Rights';
import adminService from '@disclosure-portal/services/admin';
import projectService from '@disclosure-portal/services/projects';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import {canDeleteProject, getDeleteTooltip} from '@disclosure-portal/utils/project-deletion-error';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import {useClipboard} from '@shared/utils/clipboard';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

interface Props {
  item: ProjectSlim;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  reload: [];
}>();

const {copyToClipboard} = useClipboard();

const {t} = useI18n();
const appStore = useAppStore();
const rights = useUserStore().getRights as Rights;
const projectStore = useProjectStore();

// Component Refs
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const editproject = ref();
const confirmVisible = ref(false);
const dlgNewGroup = ref();
const cloneDialogVisible = ref(false);
const cloneDialogProjectKey = ref('');
const cloneDialogProjectName = ref('');
const cloneDialogInitialCount = ref(1);

const labelTools = computed(() => appStore.getLabelsTools);

const actionButtons = computed((): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-pencil',
      hint: t('TT_edit_project'),
      event: 'edit',
      show: rights && props.item.accessRights.allowProject.update,
    },
    {
      icon: 'mdi-content-copy',
      hint: t('TT_COPY_REFERENCE_INFO'),
      event: 'copy',
      show: rights.allowProject.read,
    },
    {
      icon: 'mdi-archive-outline',
      hint: t('TT_deprecate_project'),
      event: 'deprecate',
      show: props.item.accessRights.allowProject.delete,
    },
    {
      icon: 'mdi-plus-circle-multiple-outline',
      hint: t('BTN_CLONE'),
      event: 'clone',
      show: props.item.accessRights.allowProject.create && !props.item.isGroup,
    },
    {
      icon: 'mdi-delete',
      hint: t(getDeleteTooltip(props.item)),
      event: 'delete',
      disabled: !canDeleteProject(props.item),
    },
  ];
});

const onConfirm = async (config: IConfirmationDialogConfig) => {
  if (config.okButtonIsDisabled) return;
  if (config.type === ConfirmationType.DELETE) {
    await projectStore.deleteProject(config.key);
    await projectStore.fetchProjects();
  } else if (config.type === ConfirmationType.DEPRECATE) {
    await projectService.deprecate(config.key);
    await projectStore.fetchProjects();
  } else {
    console.error(`Unhandled confirmation type: ${config.type}`);
    return;
  }
  emit('reload');
};

const getReferenceInfoForClipboard = async (item: ProjectSlim): Promise<string> => {
  const schemaLabelName = labelTools.value.schemaLabelsMap[item.schemaLabel]
    ? labelTools.value.schemaLabelsMap[item.schemaLabel].name
    : 'UNKNOWN_LABEL';
  const policyLabelNames = item.policyLabels
    .map((l: string) =>
      labelTools.value.policyLabelsMap[l] ? labelTools.value.policyLabelsMap[l].name : 'UNKNOWN_LABEL',
    )
    .join(', ');
  const projectLink = `${window.location.origin}/#/dashboard/projects/${encodeURIComponent(item._key)}`;

  const refInfo = `Disclosure Portal Project Reference

Project Name: ${item.name}
Project Identifier: ${item._key}
Project Schema Label: ${schemaLabelName}
Project Policy Labels: ${policyLabelNames}
Project Link: ${projectLink}
Application Name: ${item.applicationMeta.name}
Application Secondary ID: ${item.applicationMeta.secondaryId}
Application Link: ${item.applicationMeta.externalLink}`;

  if (rights.allowUsers.read || rights.allowAllProjectUserManagement.read) {
    const userMail = await adminService.getUserMailById(item.responsible);
    return `${refInfo}
Project Responsible with Mail: ${item.responsible} (${userMail.email})`;
  } else {
    return `${refInfo}
Project Responsible: ${item.responsible}`;
  }
};

const showDeletionConfirmationDialog = async (project: ProjectSlim) => {
  if (project.isDummy) {
    confirmConfig.value = {
      type: ConfirmationType.DELETE,
      key: project._key,
      name: project.name,
      description: 'DLG_CONFIRMATION_DESCRIPTION_DUMMY',
      okButton: 'Btn_delete',
    } as IConfirmationDialogConfig;
    confirmVisible.value = true;
  } else {
    await projectService.getApprovalOrReviewUsage(project._key).then((r) => {
      const isInUse = r.data.success;
      if (isInUse) {
        confirmConfig.value = {
          type: ConfirmationType.DEPRECATE,
          title: 'DLG_WARNING_TITLE',
          key: project._key,
          okButton: 'BTN_DEPRECATE',
          description: 'PROJECT_IN_APPROVAL_DEPRECATION',
          emphasiseText: 'PROJECT_DEPRECATION_UNREVERTABLE',
          emphasiseConfirmationText: 'PROJECT_DEPRECATION_UNREVERTABLE_CONFIRM',
        } as IConfirmationDialogConfig;
      } else {
        confirmConfig.value = {
          type: ConfirmationType.DELETE,
          key: project._key,
          name: project.name,
          description: 'DLG_CONFIRMATION_DESCRIPTION',
          okButton: 'Btn_delete',
        } as IConfirmationDialogConfig;
      }
      confirmVisible.value = true;
    });
  }
};

const showDeprecationConfirmationDialog = async (project: ProjectSlim) => {
  confirmConfig.value = {
    type: ConfirmationType.DEPRECATE,
    key: project._key,
    name: project.name,
    description: 'DLG_DEPRECATION_CONFIRMATION_DESCRIPTION',
    emphasiseText: 'PROJECT_DEPRECATION_UNREVERTABLE',
    emphasiseConfirmationText: 'PROJECT_DEPRECATION_UNREVERTABLE_CONFIRM',
    okButton: 'BTN_DEPRECATE',
  };
  confirmVisible.value = true;
};

const copyReferenceInfoToClipboard = async (item: ProjectSlim) => {
  const content = await getReferenceInfoForClipboard(item as ProjectSlim);
  copyToClipboard(content);
};

const openEdit = (item: ProjectSlim) => {
  if (item.isGroup) {
    dlgNewGroup.value!.edit(item);
  } else {
    editproject.value!.open({
      projectSlim: item,
    });
  }
};

const showCopyConfirmationDialog = (item: ProjectSlim) => {
  cloneDialogProjectKey.value = item._key;
  cloneDialogProjectName.value = item.name;
  cloneDialogInitialCount.value = 1;
  cloneDialogVisible.value = true;
};
</script>

<template>
  <TableActionButtons
    variant="compact"
    :buttons="actionButtons"
    @edit="openEdit(props.item)"
    @copy="copyReferenceInfoToClipboard(props.item)"
    @deprecate="showDeprecationConfirmationDialog(props.item)"
    @clone="showCopyConfirmationDialog(props.item)"
    @delete="showDeletionConfirmationDialog(props.item)" />
  <template>
    <NewGroupDialog ref="dlgNewGroup" @modified="emit('reload')"> </NewGroupDialog>
    <EditProjectDialog ref="editproject" @edited="emit('reload')"></EditProjectDialog>
    <CloneProjectDialog
      v-model:showDialog="cloneDialogVisible"
      :project-key="cloneDialogProjectKey"
      :project-name="cloneDialogProjectName"
      :initial-count="cloneDialogInitialCount"
      @reload="emit('reload')"></CloneProjectDialog>
    <ConfirmationDialog
      v-model:showDialog="confirmVisible"
      :config="confirmConfig"
      @confirm="onConfirm"></ConfirmationDialog>
  </template>
</template>
