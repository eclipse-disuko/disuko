// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"mercedes-benz.ghe.com/foss/disuko/observer/analytics"
	"mercedes-benz.ghe.com/foss/disuko/observer/approvalmail"
	"mercedes-benz.ghe.com/foss/disuko/observer/reviewmail"
	"mercedes-benz.ghe.com/foss/disuko/observer/spdxsubscribe"
	"mercedes-benz.ghe.com/foss/disuko/observer/userstats"
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
