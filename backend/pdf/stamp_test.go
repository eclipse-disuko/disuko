// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package pdf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"mercedes-benz.ghe.com/foss/disuko/conf"
)

func Test_createWaterMarkFile(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"test create water mark file",
			args{
				text: "test",
			},
		},
		{
			"test create multiline water mark file",
			args{
				text: "test\ntest\ntest",
			},
		},
	}
	customFontFile := filepath.Join(conf.Config.Server.BasePath, "../", conf.DocumentResources, conf.Config.Fonts.Stamp)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test-*.png")
			if err != nil {
				assert.Fail(t, "Error creating temp file:", err)
				return
			}
			createWaterMarkFile(tmpFile.Name(), tt.args.text, &customFontFile)
			if _, err := os.Stat(tmpFile.Name()); os.IsNotExist(err) {
				t.Errorf("Temp file does not exist: %s", tmpFile.Name())
			}
			os.Remove(tmpFile.Name()) // Clean up
		})
	}
}
