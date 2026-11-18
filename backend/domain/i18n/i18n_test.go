// SPDX-FileCopyrightText: 2026 Mercedes-Benz Group AG and Mercedes-Benz AG
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestToDTOMetaContainsTimestampsAndAuthors(t *testing.T) {
	locale := NewI18nLocale("en")
	entry := NewI18nEntry("GREETING", "Hello", "Greeting description")
	entry.Created = time.Date(2026, time.November, 1, 10, 0, 0, 0, time.UTC)
	entry.Updated = entry.Created.Add(time.Hour)
	entry.CreatedBy = "creator"
	entry.UpdatedBy = "editor"
	locale.SetEntry(entry)

	dto := locale.ToDTO(false)
	require.Equal(t, "Hello", dto.Entries["GREETING"])

	data, err := json.Marshal(dto.Meta)
	require.NoError(t, err)
	require.JSONEq(t, `{"GREETING":{"created":"2026-11-01T10:00:00Z","updated":"2026-11-01T11:00:00Z","createdBy":"creator","updatedBy":"editor"}}`, string(data))
}
