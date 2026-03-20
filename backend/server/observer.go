// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/eclipse-disuko/disuko/observer/analytics"
	"github.com/eclipse-disuko/disuko/observer/approvalmail"
	"github.com/eclipse-disuko/disuko/observer/reviewmail"
	"github.com/eclipse-disuko/disuko/observer/spdxsubscribe"
	"github.com/eclipse-disuko/disuko/observer/userstats"
)

type observer interface {
	RegisterHandlers()
}

func (s *Server) registerObserver() {
	approvalMail := approvalmail.Init(s.mailClient, s.repos.user, s.repos.project)
	userStatsCon := userstats.Init(s.scheduler)
	analyticsCon := analytics.Init(&s.services.analytics)
	spdxMail := spdxsubscribe.Init(s.mailClient, s.repos.user)
	overallReview := reviewmail.Init(s.mailClient, s.repos.user)

	observers := []observer{
		approvalMail,
		analyticsCon,
		spdxMail,
		overallReview,
		userStatsCon,
	}
	for _, o := range observers {
		o.RegisterHandlers()
	}
}
