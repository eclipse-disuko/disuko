// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package audit

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func ConvertDateTime(t time.Time) string {
	if t.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local)) {
		return ""
	}
	dateTimeString := t.Format("02.01.2006 15:04")
	return dateTimeString
}

func CreateAndAddAuditEntry(container *audit.Container, user string, message string, diffFunc audit.DiffFunc, after, before interface{}, ignoreFieldsOption ...cmp.Option) {
	auditEntry := CreateAuditEntry(user, message, diffFunc, after, before, ignoreFieldsOption...)
	container.AddAuditEntry(auditEntry)
}

func CreateAuditEntry(user string, message string, diffFunc audit.DiffFunc, after, before interface{}, ignoreFieldsOption ...cmp.Option) *audit.Audit {
	compareStr := createCompareString(diffFunc, after, before, ignoreFieldsOption...)
	return audit.NewAudit(user, message, compareStr)
}

func createCompareString(diffFunc audit.DiffFunc, after, before interface{}, ignoreFieldsOption ...cmp.Option) string {
	options := []cmp.Option{cmpopts.IgnoreTypes(audit.Container{})}
	if len(ignoreFieldsOption) > 0 {
		options = append(options, ignoreFieldsOption...)
	}
	return diffFunc(before, after, options...)
}
