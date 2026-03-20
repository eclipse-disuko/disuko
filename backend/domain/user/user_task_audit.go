// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package user

type TaskAudit struct {
	Type TaskType `json:"type"`

	TargetGuid  string `json:"targetGuid"`
	ProjectGuid string `json:"projectGuid"`
	Creator     string `json:"creator"`

	Status TaskStatus `json:"status"`
}

func (task Task) ToAudit() TaskAudit {
	return TaskAudit{
		Creator:     task.Creator,
		ProjectGuid: task.ProjectGuid,
		Status:      task.Status,
		TargetGuid:  task.TargetGuid,
		Type:        task.Type,
	}
}
