// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/eclipse-disuko/disuko/observer/analytics"
	"github.com/eclipse-disuko/disuko/observer/approvalmail"
	"github.com/eclipse-disuko/disuko/observer/confirmtransfer"
	"github.com/eclipse-disuko/disuko/observer/reviewmail"
	"github.com/eclipse-disuko/disuko/observer/spdxsubscribe"
	"github.com/eclipse-disuko/disuko/observer/userstats"
)

type observer interface {
	RegisterHandlers()
}

func (s *Server) registerObserver() {
	approvalMail := approvalmail.Init(s.services.mail, s.repos.user, s.repos.project)
	userStatsCon := userstats.Init(s.scheduler)
	analyticsCon := analytics.Init(&s.services.analytics)
	spdxMail := spdxsubscribe.Init(s.services.mail, s.repos.user)
	overallReview := reviewmail.Init(s.services.mail, s.repos.user)
	confirmTransfer := confirmtransfer.Init(
		s.repos.sbomList,
		s.repos.policyRules,
		s.repos.policyDecisions,
		&s.services.projectLabelService,
		s.services.spdx,
		&s.services.overallReview,
	)

	observers := []observer{
		approvalMail,
		analyticsCon,
		spdxMail,
		overallReview,
		userStatsCon,
		confirmTransfer,
	}
	for _, o := range observers {
		o.RegisterHandlers()
	}
}
