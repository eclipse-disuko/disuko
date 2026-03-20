// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package exception

import (
	"log"
	"testing"

	"github.com/eclipse-disuko/disuko/logy"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	data string
}

func Test_CatchNilPanic(t *testing.T) {
	defer CatchPanic(func(a any) {
		log.Printf("PANIC: %s", a)
		log.Printf("PANIC: %#v", a)
	})

	var nilStruct *testStruct

	nilStruct.data = "Test" //null pointer

	assert.Fail(t, "panic before should!")
}

func Test_CatchPanic(t *testing.T) {
	defer CatchPanic(func(a any) {
		log.Printf("PANIC: %s", a)
	})

	panic("I HAVE PANIC")
	assert.Fail(t, "panic before should!")
}

func Test_CatchException(t *testing.T) {
	defer CatchException(ExceptionHandler{
		RequestSession: &logy.RequestSession{ReqID: "TEST"},
	})

	ThrowException2("CODE", "MESSAGE", "RAW")
	assert.Fail(t, "panic before should!")
}

func Test_CatchExceptionWithDefer(t *testing.T) {
	count := 0
	defer func() {
		assert.Equal(t, 1, count)
		count++
	}()

	defer CatchExceptionWithCustom(ExceptionHandler{
		RequestSession: &logy.RequestSession{ReqID: "TEST"},
	}, func(exception Exception) {
		assert.Equal(t, 1, count)
	})

	defer func() {
		count++
	}()

	ThrowException2("CODE", "MESSAGE", "RAW")
	assert.Fail(t, "panic before should!")
}
