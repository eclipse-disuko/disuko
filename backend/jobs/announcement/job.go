// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package announcement

import (
	"regexp"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain/announcement"
	"mercedes-benz.ghe.com/foss/disuko/domain/audit"
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/license"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	announcementRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/announcements"
	licenseRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
)

type Job struct {
	announcementsRepository announcementRepo.IAnnouncementsRepository
	licenseRepository       licenseRepo.ILicensesRepository
}

func Init(ar announcementRepo.IAnnouncementsRepository, lr licenseRepo.ILicensesRepository) *Job {
	return &Job{
		announcementsRepository: ar,
		licenseRepository:       lr,
	}
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log

	log.AddEntry(job.Info, "started")

	var (
		licenses = j.licenseRepository.FindAllWithDeleted(rs, false)
		failed   = 0
		ok       = 0
		result   = announcement.AnnouncementGenerationResultDto{
			ReqID:              rs.ReqID,
			AnnouncementsCount: 0,
			Errors:             make([]string, 0),
		}
	)
	for _, l := range licenses {
		m := make(map[string]*announcement.LicenseAnnouncement)
		runChecks(l, m, info.Updated)
		if len(m) == 0 {
			continue
		}

		v := make([]*announcement.Announcement, 0, len(m))
		for _, value := range m {
			a, err := value.ToAnnouncement()
			if err == nil {
				v = append(v, a)
				ok++
			} else {
				log.AddEntry(job.Error, "could not marshall announcement content %s", err.Error())
				result.Errors = append(result.Errors, "could not marshall announcement content"+err.Error())
				failed++
			}
		}
		log.AddEntry(job.Info, "saved map %s", l.Name)
		for k, v := range m {
			log.AddEntry(job.Info, "key / value = %s / %v", k, v)
		}
		j.announcementsRepository.SaveList(rs, v, false)
	}
	result.AnnouncementsCount = ok
	log.AddEntry(job.Info, "finished: ok %d / failed %d", ok, failed)
	return scheduler.ExecutionResult{
		Success:   !(failed > 0),
		Log:       log,
		CustomRes: result,
	}
}

// runs all license announcement checks given a license and fills given map m
// only runs for changes (audit logs) later than given time lastRun
func runChecks(l *license.License, m map[string]*announcement.LicenseAnnouncement, lastRun time.Time) {
	for _, check := range licenseAnnouncementChecks {
		unfilteredAnnouncements := make([]*announcement.LicenseAnnouncement, 0)
		for _, auditTrail := range l.AuditTrail {
			if auditTrail.Created.Before(lastRun) {
				// ignore older auditTrails
				continue
			}
			if r, a := check(l, auditTrail); r {
				unfilteredAnnouncements = append(unfilteredAnnouncements, a)
			}
		}

		for _, a := range unfilteredAnnouncements {
			// by using the date + ChangeType as keys in our map we make sure only 1 announcement per ChangeType and day
			// is created
			mapKey := a.When.Format(time.DateOnly) + ":" + a.ChangeType

			oldA, exists := m[mapKey]
			if exists && oldA.When.Format(time.DateOnly) == a.When.Format(time.DateOnly) {
				if oldA.OldVal == a.NewVal && oldA.NewVal == a.OldVal {
					// if there are "reverting" changes on the same day, we can drop the announcement altogether
					delete(m, mapKey)
				} else {
					// otherwise we can "combine" the changes to one announcement
					oldA.NewVal = a.NewVal
				}
			} else {
				// no announcement on the same day with the same change type exists, thus we just add it to the map
				m[mapKey] = a
			}
		}
	}
}

func createNewLicenseAnnouncement(l *license.License, when time.Time, changeType string, oldVal string, newVal string) *announcement.LicenseAnnouncement {
	a := announcement.NewLicenseAnnouncement()
	a.When = when
	a.LicenseName = l.Name
	a.Key = l.Key
	a.LicenseId = l.LicenseId
	a.ChangeType = changeType
	a.OldVal = oldVal
	a.NewVal = newVal
	return a
}

var licenseAnnouncementChecks = []func(l *license.License, auditLog *audit.Audit) (bool, *announcement.LicenseAnnouncement){
	func(l *license.License, auditLog *audit.Audit) (bool, *announcement.LicenseAnnouncement) {
		r, _ := regexp.Compile("\\{\\*license.LicenseAudit}\\.Meta\\.IsLicenseChart:\\n\\t-: (false|true)\\n\\t\\+: (false|true)\\n")

		matches := r.FindStringSubmatch(auditLog.MetaJSON)
		if len(matches) == 3 {
			return true, createNewLicenseAnnouncement(l, auditLog.Created, "license_chart", matches[1], matches[2])
		} else {
			return false, nil
		}
	},
	func(l *license.License, auditLog *audit.Audit) (bool, *announcement.LicenseAnnouncement) {
		r, _ := regexp.Compile("\\{\\*license.LicenseAudit}\\.Meta\\.Family:\\n\\t([-+]): ([a-z A-Z]*)(\\n\\t[-+]: ([a-z A-Z]*))*\\n")

		matches := r.FindStringSubmatch(auditLog.MetaJSON)
		if !l.Meta.IsLicenseChart || len(matches) != 5 {
			return false, nil
		}
		oldVal, newVal := "", ""
		if matches[3] == "" {
			if matches[1] == "+" {
				newVal = matches[2]
			} else {
				oldVal = matches[2]
			}
		} else {
			oldVal = matches[2]
			newVal = matches[4]
		}
		return true, createNewLicenseAnnouncement(l, auditLog.Created, "license_family", oldVal, newVal)
	},
	func(l *license.License, auditLog *audit.Audit) (bool, *announcement.LicenseAnnouncement) {
		r, _ := regexp.Compile("\\{\\*license.LicenseAudit}\\.Meta\\.LicenseType:\\n\\t([-+]): ([a-z A-Z]*)(\\n\\t[-+]: ([a-z A-Z]*))*\\n")

		matches := r.FindStringSubmatch(auditLog.MetaJSON)
		if !l.Meta.IsLicenseChart || len(matches) != 5 {
			return false, nil
		}
		oldVal, newVal := "", ""
		if matches[3] == "" {
			if matches[1] == "+" {
				newVal = matches[2]
			} else {
				oldVal = matches[2]
			}
		} else {
			oldVal = matches[2]
			newVal = matches[4]
		}
		return true, createNewLicenseAnnouncement(l, auditLog.Created, "license_type", oldVal, newVal)
	},
	func(l *license.License, auditLog *audit.Audit) (bool, *announcement.LicenseAnnouncement) {
		r, _ := regexp.Compile("\\{\\*license.LicenseAudit}\\.Meta\\.ApprovalState:(\\n\\t-: ([a-z A-Z]*))*\\n\\t\\+: forbidden\\n")

		matches := r.FindStringSubmatch(auditLog.MetaJSON)
		if len(matches) != 3 {
			return false, nil
		}
		return true, createNewLicenseAnnouncement(l, auditLog.Created, "license_forbidden", matches[2], "forbidden")
	},
	func(l *license.License, auditLog *audit.Audit) (bool, *announcement.LicenseAnnouncement) {
		if auditLog.Message != message.LicenseDeleted || l.Source != license.CUSTOM {
			return false, nil
		}
		return true, createNewLicenseAnnouncement(l, auditLog.Created, "custom_license_deleted", "", "")
	},
}
