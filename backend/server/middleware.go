// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"compress/flate"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/helper/middlewareDisco"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) setupMW() {
	s.r.Use(middleware.RealIP)
	s.r.Use(middlewareDisco.RequestSession)
	s.r.Use(middlewareDisco.Logging)
	s.r.Use(middlewareDisco.CustomHeaders)
	s.r.Use(middlewareDisco.Recoverer)

	compressor := middleware.NewCompressor(flate.DefaultCompression)
	s.r.Use(compressor.Handler)

	s.r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(conf.Config.Server.AllowedOrigins, ","),
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "api_key", "X-Client-Version"},
		ExposedHeaders:   []string{"Link", "X-Notification"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	s.r.Use(middleware.Timeout(3 * time.Minute))

	s.basicauthMW = middlewareDisco.InitInternalTokenMW(s.repos.basicauth)
}
