// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package hash

import (
	"testing"

	"mercedes-benz.ghe.com/foss/disuko/logy"
)

var requestSessionTest = &logy.RequestSession{ReqID: "TEST"}

func TestHash(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want string
	}{
		{name: "hash empty string", args: "", want: ""},
		{name: "hash string", args: "lala", want: "37b1aed0d40aadbaff85c590e0e18f04d8f7bf87a74eec118e739755908b9b1f"},
		{name: "hash struct", args: struct{ name string }{name: "lala"}, want: "cd6bd19d0363b3d75151744719275d59e79ffbeeb332cc056c87e66ae0fd1868"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(requestSessionTest, tt.args); got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
