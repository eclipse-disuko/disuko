// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"mercedes-benz.ghe.com/foss/disuko/server"
)

type ProjectTest struct {
	Key  string
	Name string
}

func main() {
	// db := mongo.Database{}
	// db.Init(nil, "projects", nil)

	ctx := context.Background()
	server.StartServer(ctx)
}
