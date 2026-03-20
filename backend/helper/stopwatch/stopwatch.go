// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package stopwatch

import (
	"strconv"
	"time"
)

type StopWatch struct {
	StartTime time.Time
	EndTime   time.Time
	DiffTime  time.Duration
}

func (sw *StopWatch) Start() {
	sw.StartTime = time.Now()
}

func (sw *StopWatch) Stop() {
	sw.EndTime = time.Now()
	sw.DiffTime = sw.EndTime.Sub(sw.StartTime)
}

func (sw *StopWatch) GetTransferRate(mb float64) string {
	return strconv.FormatFloat(mb/sw.DiffTime.Seconds(), 'f', 2, 64) + " MB/sec"
}

func (sw *StopWatch) FormatSeconds() string {
	return strconv.FormatFloat(sw.DiffTime.Seconds(), 'f', 2, 64) + "  sec"
}
