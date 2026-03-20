// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/infra/service/analytics"
	"github.com/eclipse-disuko/disuko/observermngmt"
)

type Analytics struct {
	service *analytics.Analytics
}

func Init(service *analytics.Analytics) *Analytics {
	return &Analytics{
		service: service,
	}
}

func (a *Analytics) RegisterHandlers() {
	observermngmt.RegisterHandler(observermngmt.SpdxAdded, a.OnSpdxAdded)
	observermngmt.RegisterHandler(observermngmt.SpdxUpdatedNewest, a.OnSpdxAdded)
	observermngmt.RegisterHandler(observermngmt.SpdxDeleted, a.OnSpdxDeleted)
	observermngmt.RegisterHandler(observermngmt.ProjectDeleted, a.OnProjectDeleted)
	observermngmt.RegisterHandler(observermngmt.ProjectVersionDeleted, a.OnVersionDeleted)
	observermngmt.RegisterHandler(observermngmt.LicenseAdded, a.OnLicenseAdded)
	observermngmt.RegisterHandler(observermngmt.LicenseDeleted, a.OnLicenseDeleted)
	observermngmt.RegisterHandler(observermngmt.LicenseAliasAdded, a.OnAliasAdded)
	observermngmt.RegisterHandler(observermngmt.LicenseAliasDeleted, a.OnAliasDeleted)
	observermngmt.RegisterHandler(observermngmt.ProjectUpdated, a.OnProjectUpdated)
}

func (a *Analytics) OnSpdxAdded(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.SpdxData)
	if !ok {
		return
	}
	go exception.TryCatchAndLog(data.RequestSession, func() {
		a.service.ExportSPDX(data.RequestSession, data.Project, data.Version, data.SpdxFile)
	})
}

func (a *Analytics) OnSpdxDeleted(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.SpdxData)
	if !ok {
		return
	}
	go exception.TryCatchAndLog(data.RequestSession, func() {
		a.service.Handler.HandleSpdxDeleted(data.RequestSession, data.SpdxFile.Key)
	})
}

func (a *Analytics) OnVersionDeleted(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.VersionData)
	if !ok {
		return
	}
	a.service.DeleteVersion(data.RequestSession, data.Project, data.Version)
}

func (a *Analytics) OnProjectDeleted(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.ProjectDeletedData)
	if !ok {
		return
	}
	a.service.DeleteProject(data.RequestSession, data.Project)
}

func (a *Analytics) OnLicenseAdded(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.LicenseData)
	if !ok {
		return
	}
	go exception.TryCatchAndLog(data.RequestSession, func() {
		a.service.Handler.HandleLicenseIdAdded(data.RequestSession, data.Id, data.Id)
	})
}

func (a *Analytics) OnLicenseDeleted(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.LicenseData)
	if !ok {
		return
	}
	go exception.TryCatchAndLog(data.RequestSession, func() {
		a.service.Handler.HandleLicenseIdDeleted(data.RequestSession, data.Id)
	})
}

func (a *Analytics) OnAliasAdded(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.AliasData)
	if !ok {
		return
	}
	go exception.TryCatchAndLog(data.RequestSession, func() {
		a.service.Handler.HandleLicenseIdAdded(data.RequestSession, data.Alias, data.Id)
	})
}

func (a *Analytics) OnAliasDeleted(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.AliasData)
	if !ok {
		return
	}
	go exception.TryCatchAndLog(data.RequestSession, func() {
		a.service.Handler.HandleLicenseIdDeleted(data.RequestSession, data.Id)
	})
}

func (a *Analytics) OnProjectUpdated(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.ProjectUpdatedData)
	if !ok {
		return
	}
	if data.New.Parent != data.Old.Parent {
		if data.New.Parent == "" {
			go exception.TryCatchAndLog(data.RequestSession, func() {
				a.service.Handler.HandleCompanyChanged(data.RequestSession, data.New.Key, data.New.CustomerMeta.DeptId)
			})
		} else {
			go exception.TryCatchAndLog(data.RequestSession, func() {
				a.service.Handler.HandleCompanyChanged(data.RequestSession, data.New.Key, data.NewParent.CustomerMeta.DeptId)
			})
		}
	} else if data.NewParent != nil {
		go exception.TryCatchAndLog(data.RequestSession, func() {
			a.service.Handler.HandleCompanyChanged(data.RequestSession, data.New.Key, data.NewParent.CustomerMeta.DeptId)
		})
	} else if data.Old.CustomerMeta.DeptId != data.New.CustomerMeta.DeptId {
		go exception.TryCatchAndLog(data.RequestSession, func() {
			a.service.Handler.HandleCompanyChanged(data.RequestSession, data.New.Key, data.New.CustomerMeta.DeptId)
		})
	}
	var (
		newRes = data.New.ProjectResponsible()
		oldRes = data.Old.ProjectResponsible()
	)
	if newRes == nil && oldRes == nil {
		return
	}
	if newRes != nil && oldRes == nil {
		go exception.TryCatchAndLog(data.RequestSession, func() {
			a.service.Handler.HandleResponsibleChanged(data.RequestSession, data.New.Key, newRes.UserId)
		})
	} else if newRes == nil && oldRes != nil {
		go exception.TryCatchAndLog(data.RequestSession, func() {
			a.service.Handler.HandleResponsibleChanged(data.RequestSession, data.New.Key, "")
		})
	} else {
		go exception.TryCatchAndLog(data.RequestSession, func() {

			a.service.Handler.HandleResponsibleChanged(data.RequestSession, data.New.Key, newRes.UserId)
		})
	}
}
