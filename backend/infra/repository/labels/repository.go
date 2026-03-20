// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package labels

import (
	"sync"

	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type labelsRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*label.Label]
	mem   map[string]*label.Label
	mutex sync.RWMutex
}

func NewLabelsRepository(requestSession *logy.RequestSession) ILabelRepository {
	innerRepo := &labelsRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*label.Label](
			requestSession,
			LabelCollectionName,
			func() *label.Label {
				return &label.Label{}
			},
			nil,
			"",
			nil,
			nil),
	}
	innerRepo.LoadFromDb(requestSession)
	return innerRepo
}

func (lr *labelsRepositoryStruct) FindByNameAndType(requestSession *logy.RequestSession, name string, labelType label.LabelType) *label.Label {
	lr.mutex.RLock()
	defer lr.mutex.RUnlock()
	for _, label := range lr.mem {
		if label.Name == name && label.Type == labelType {
			return label
		}
	}
	return nil
}

func (lr *labelsRepositoryStruct) FindAllByType(requestSession *logy.RequestSession, labelType label.LabelType) []*label.Label {
	lr.mutex.RLock()
	defer lr.mutex.RUnlock()
	var res []*label.Label
	for _, label := range lr.mem {
		if label.Type == labelType {
			res = append(res, label)
		}
	}
	return res
}

func (lr *labelsRepositoryStruct) LoadFromDb(requestSession *logy.RequestSession) int {
	lr.mutex.Lock()
	defer lr.mutex.Unlock()
	lr.mem = make(map[string]*label.Label)
	all := lr.BaseRepositoryWithHardDelete.FindAll(requestSession, false)
	for _, label := range all {
		lr.mem[label.Key] = label
	}
	return len(all)
}

func (lr *labelsRepositoryStruct) FindAll(requestSession *logy.RequestSession, optimized bool) []*label.Label {
	lr.mutex.RLock()
	defer lr.mutex.RUnlock()
	var res []*label.Label
	for _, l := range lr.mem {
		res = append(res, l)
	}
	return res
}

func (lr *labelsRepositoryStruct) FindByKey(requestSession *logy.RequestSession, key string, optimized bool) *label.Label {
	lr.mutex.RLock()
	defer lr.mutex.RUnlock()
	if lr.mem == nil {
		return nil
	}
	if res, ok := lr.mem[key]; ok {
		return res
	}
	return nil
}
