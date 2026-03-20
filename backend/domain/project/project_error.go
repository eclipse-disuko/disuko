// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

type ProjectError struct {
	ErrorType int
	Err       error
}

func (err *ProjectError) Error() string {
	return err.Err.Error()
}
