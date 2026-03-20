// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"

	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const suffix = "go@d1sco$f0r@ever!"

func Hash(requestSession *logy.RequestSession, o interface{}) string {
	if s, ok := o.(string); ok && s == "" {
		logy.Infof(requestSession, "tried to hash an empty string")
		return ""
	}
	sum256 := sha256.Sum256([]byte(fmt.Sprintf("%v%s", o, suffix)))
	return hex.EncodeToString(sum256[:])
}

func GetSha256Hash(source []byte) string {
	sha256hash := sha256.New()
	sha256hash.Write(source)
	return EncodeSha256ToString(sha256hash)
}

func EncodeSha256ToString(hash hash.Hash) string {
	return hex.EncodeToString(hash.Sum(nil))
}
