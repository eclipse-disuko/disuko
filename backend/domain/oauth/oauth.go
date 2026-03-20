// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package oauth

type UserInfo struct {
	Sub         string `json:"sub"`
	AppId       string `json:"app_id"`
	AccessGroup string `json:"authorization_group"`
}

type UserInfoWithGroupsDto struct {
	UserInfo
	EntitlementGroup any    `json:"entitlement_group"` // string or []string
	GroupType        string `json:"type_group"`
}

type UserInfoWithGroups struct {
	UserInfo
	EntitlementGroup []string `json:"entitlement_group"`
	GroupType        string   `json:"type_group"`
}
