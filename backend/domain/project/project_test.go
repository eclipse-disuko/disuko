// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"testing"
	"time"

	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/test"
	"github.com/stretchr/testify/assert"
)

func compareUsers(a, b *UserManagementEntity) bool {
	if len(a.Users) != len(b.Users) {
		return false
	}
	for _, au := range a.Users {
		found := false
		for _, bu := range b.Users {
			if au.UserId == bu.UserId && au.UserType == bu.UserType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name  string
		input UserManagementEntity
		want  UserManagementEntity
		arg   string
		fail  bool
	}{
		{
			name: "Delete only owner",
			input: UserManagementEntity{
				Users: []*ProjectMemberEntity{
					{
						UserId:   "testOwner",
						UserType: OWNER,
					},
				},
			},
			want: UserManagementEntity{
				Users: []*ProjectMemberEntity{
					{
						UserId:   "testOwner",
						UserType: OWNER,
					},
				},
			},
			fail: true,
			arg:  "testOwner",
		},
		{
			name: "Delete one of two owners",
			input: UserManagementEntity{
				Users: []*ProjectMemberEntity{
					{
						UserId:   "testOwner",
						UserType: OWNER,
					},
					{
						UserId:   "testOwner2",
						UserType: OWNER,
					},
				},
			},
			want: UserManagementEntity{
				Users: []*ProjectMemberEntity{
					{
						UserId:   "testOwner",
						UserType: OWNER,
					},
				},
			},
			fail: false,
			arg:  "testOwner2",
		},
		{
			name: "Delete regular user",
			input: UserManagementEntity{
				Users: []*ProjectMemberEntity{
					{
						UserId:   "testOwner",
						UserType: OWNER,
					},
					{
						UserId:   "testUser",
						UserType: VIEWER,
					},
				},
			},
			want: UserManagementEntity{
				Users: []*ProjectMemberEntity{
					{
						UserId:   "testOwner",
						UserType: OWNER,
					},
				},
			},
			fail: false,
			arg:  "testUser",
		},
		{
			name: "Delete owner not regular user",
			input: UserManagementEntity{
				Users: []*ProjectMemberEntity{
					{
						UserId:   "testOwner",
						UserType: OWNER,
					},
					{
						UserId:   "testUser",
						UserType: VIEWER,
					},
				},
			},
			want: UserManagementEntity{
				Users: []*ProjectMemberEntity{
					{
						UserId:   "testOwner",
						UserType: OWNER,
					},
					{
						UserId:   "testUser",
						UserType: VIEWER,
					},
				},
			},
			fail: true,
			arg:  "testOwner",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				UserManagement: tt.input,
			}
			if tt.fail {
				test.ExpectException(t, message.ErrorProjectDeleteLastMember, func() {
					p.DeleteUser(tt.arg)
				})
			} else {
				p.DeleteUser(tt.arg)
			}
			if !compareUsers(&p.UserManagement, &tt.want) {
				t.Errorf("DeleteUser didnt work as expected. Wanted: %v Got: %v", p.UserManagement, tt.want)
			}
		})
	}
}

func TestProject_ExpireTokens(t *testing.T) {
	tests := []struct {
		name        string
		Token       []*Token
		want        bool
		tokenStatus TokenStatus
	}{{
		name: "no changes", Token: []*Token{{
			Expiry: time.Now().Add(1 * time.Minute).Format("2006-01-02T15:04:05Z07:00"),
			Status: ACTIVE,
		}}, want: false, tokenStatus: ACTIVE}, {
		name: "expired", Token: []*Token{{
			Expiry: time.Now().Add(-1 * time.Minute).Format("2006-01-02T15:04:05Z07:00"),
			Status: ACTIVE,
		}}, want: true, tokenStatus: EXPIRED}, {
		name: "not expire revoked", Token: []*Token{{
			Expiry: time.Now().Add(-1 * time.Minute).Format("2006-01-02T15:04:05Z07:00"),
			Status: REVOKED,
		}}, want: false, tokenStatus: REVOKED},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := &Project{
				Token: tt.Token,
			}
			assert.Equalf(t, tt.want, project.ExpireTokens(), "ExpireTokens()")
			assert.Equalf(t, tt.tokenStatus, project.Token[0].Status, "ExpireTokens()")
		})
	}
}
