// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package client_utils

import (
	"net/http"
	"net/url"

	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func GetTransport(proxied bool) *http.Transport {
	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	defaultTransport.MaxIdleConns = 70
	defaultTransport.MaxConnsPerHost = 70
	defaultTransport.MaxIdleConnsPerHost = 70
	if !proxied {
		return defaultTransport
	}
	proxiedTransport := defaultTransport.Clone()
	// We don't want to use this proxy locally as it demands authentication using non-CaaS IP ranges.
	if conf.Config.Server.Env != "local" {
		if len(conf.Config.Proxy.HttpProxy) > 0 {
			pUrl, _ := url.Parse(conf.Config.Proxy.HttpProxy)
			proxiedTransport.Proxy = http.ProxyURL(pUrl)
		} else {
			logy.Warnf(nil, "Environment is not local and no proxy set!")
		}
	}
	return proxiedTransport
}
