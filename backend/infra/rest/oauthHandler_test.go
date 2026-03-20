// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"mercedes-benz.ghe.com/foss/disuko/domain/oauth"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/test"
)

func TestOAuthHandler_convertUserInfoWithGroupsDto(t *testing.T) {
	// users with multiple roles
	entitlementGroupArray := []byte("{\"sub\":\"SOME_SUB\",\"email_verified\":true,\"entitlement_group\":[\"FOSSDP.policy_admin\", \"FOSSDP.test\"],\"app_id\":\"FOSSDP\",\"email\":\"some_email@email.com\",\"personal_data_hint\":\"You have requested scopes allowing access to personal data. Be aware that you are not allowed to give personal data to other applications without permission by HR. Please request an access token without personal data scopes for accessing other applications. Personal data scopes are marked in the GAS-OIDC integration guide.\"}\n")

	var userInfoClaimsDto = oauth.UserInfoWithGroupsDto{}
	_ = json.Unmarshal(entitlementGroupArray, &userInfoClaimsDto)

	want := []string{"FOSSDP.policy_admin", "FOSSDP.test"}
	got := convertUserInfoWithGroupsDto(&userInfoClaimsDto).EntitlementGroup
	assert.Equal(t, want, got, "EntitlementGroups were not correctly parsed from an array")

	// users with single role as array
	entitlementGroupSingleElementArray := []byte("{\"sub\":\"SOME_SUB\",\"email_verified\":true,\"entitlement_group\":[\"FOSSDP.policy_admin\"],\"app_id\":\"FOSSDP\",\"email\":\"some_email@email.com\",\"personal_data_hint\":\"You have requested scopes allowing access to personal data. Be aware that you are not allowed to give personal data to other applications without permission by HR. Please request an access token without personal data scopes for accessing other applications. Personal data scopes are marked in the GAS-OIDC integration guide.\"}\n")

	userInfoClaimsDto = oauth.UserInfoWithGroupsDto{}
	_ = json.Unmarshal(entitlementGroupSingleElementArray, &userInfoClaimsDto)

	want = []string{"FOSSDP.policy_admin"}
	got = convertUserInfoWithGroupsDto(&userInfoClaimsDto).EntitlementGroup
	assert.Equal(t, want, got, "EntitlementGroups were not correctly parsed from a single element array")

	// users with 1 role
	entitlementGroupString := []byte("{\"sub\":\"SOME_SUB\",\"email_verified\":true,\"entitlement_group\":\"FOSSDP.policy_admin\",\"app_id\":\"FOSSDP\",\"email\":\"some_email@email.com\",\"personal_data_hint\":\"You have requested scopes allowing access to personal data. Be aware that you are not allowed to give personal data to other applications without permission by HR. Please request an access token without personal data scopes for accessing other applications. Personal data scopes are marked in the GAS-OIDC integration guide.\"}\n")

	userInfoClaimsDto = oauth.UserInfoWithGroupsDto{}
	_ = json.Unmarshal(entitlementGroupString, &userInfoClaimsDto)

	want = []string{"FOSSDP.policy_admin"}
	got = convertUserInfoWithGroupsDto(&userInfoClaimsDto).EntitlementGroup
	assert.Equal(t, want, got, "EntitlementGroups were not correctly parsed from a string")

	// users without a role
	entitlementGroupNil := []byte("{\"sub\":\"SOME_SUB\",\"email_verified\":true,\"app_id\":\"FOSSDP\",\"email\":\"some_email@email.com\",\"personal_data_hint\":\"You have requested scopes allowing access to personal data. Be aware that you are not allowed to give personal data to other applications without permission by HR. Please request an access token without personal data scopes for accessing other applications. Personal data scopes are marked in the GAS-OIDC integration guide.\"}\n")

	userInfoClaimsDto = oauth.UserInfoWithGroupsDto{}
	_ = json.Unmarshal(entitlementGroupNil, &userInfoClaimsDto)

	want = []string{}
	got = convertUserInfoWithGroupsDto(&userInfoClaimsDto).EntitlementGroup
	assert.Equal(t, want, got, "EntitlementGroups were not correctly parsed from a string")

	// malformed JSON
	entitlementGroupMalformed := []byte("{\"sub\":\"SOME_SUB\",\"email_verified\":true,\"entitlement_group\":false,\"app_id\":\"FOSSDP\",\"email\":\"some_email@email.com\",\"personal_data_hint\":\"You have requested scopes allowing access to personal data. Be aware that you are not allowed to give personal data to other applications without permission by HR. Please request an access token without personal data scopes for accessing other applications. Personal data scopes are marked in the GAS-OIDC integration guide.\"}\n")

	userInfoClaimsDto = oauth.UserInfoWithGroupsDto{}
	_ = json.Unmarshal(entitlementGroupMalformed, &userInfoClaimsDto)

	test.ExpectException(t, message.ErrorUnexpectedType, func() { convertUserInfoWithGroupsDto(&userInfoClaimsDto) })
}
