// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"net/http"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/helper/client_utils"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/infra/rest"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func (s *Server) setUpOAuth(rs *logy.RequestSession) {
	ctx := context.Background()
	if conf.Config.OAuth2.InsecureProvider != "" {
		ctx = oidc.InsecureIssuerURLContext(ctx, conf.Config.OAuth2.InsecureProvider)
	}
	provider, err := oidc.NewProvider(ctx, conf.Config.OAuth2.Provider)
	if err != nil {
		logy.Fatalf(rs, "oauth setup failed: %s", err)
	}

	oidcConfig := oidc.Config{
		ClientID: conf.Config.OAuth2.ClientId,
	}

	oauth2Config := oauth2.Config{
		ClientID:     conf.Config.OAuth2.ClientId,
		ClientSecret: conf.Config.OAuth2.Secret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  conf.Config.OAuth2.RedirectURL,
		Scopes:       []string{"openid", "offline_access", "authorization_group", "entitlement_group", "email", "group_type", "personal_data", "organizational_data"},
	}

	s.handlers.auth = rest.OAuthHandler{
		HttpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: client_utils.GetTransport(false),
		},
		Ctx:         context.Background(),
		Verifier:    provider.Verifier(&oidcConfig),
		Provider:    provider,
		ProviderURL: conf.Config.OAuth2.Provider,
		Config:      oauth2Config,
	}
}
