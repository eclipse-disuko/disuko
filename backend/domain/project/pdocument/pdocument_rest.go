// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package pdocument

import "time"

type PDocumentDto struct {
	Key         string    `json:"_key"`
	Description string    `json:"description"`
	ApprovalId  string    `json:"approvalId"`
	Type        string    `json:"type"`
	Lang        string    `json:"lang"`
	Created     time.Time `json:"created,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
	Hash        string    `json:"hash"`
	FileName    string    `json:"fileName"`
	Version     int       `json:"version"`
}

func (entity *PDocument) ToDto() PDocumentDto {
	fileName := GetFileName(entity.Type, entity.ApprovalId, LangStrToTag(entity.Lang))

	docVersion := int(NONE_VERSION)
	if entity.VersionIndex != nil {
		docVersion = int(*entity.VersionIndex)
	}

	return PDocumentDto{
		Key:         entity.Key,
		Description: entity.Description,
		ApprovalId:  entity.ApprovalId,
		Type:        string(entity.Type),
		Lang:        string(entity.Lang),
		Created:     entity.Created,
		Updated:     entity.Updated,
		Hash:        entity.Hash,
		FileName:    fileName,
		Version:     docVersion,
	}
}
