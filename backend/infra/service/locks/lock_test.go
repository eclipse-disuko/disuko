// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package locks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

// TODO: move these tests into integration tests after they are moved from old repo

var requestSessionTest = &logy.RequestSession{ReqID: "INTEGRATION-TEST"}

func _TestAcquireRelease(t *testing.T) {
	conf.Config.Cache.Host = "localhost"
	conf.Config.Cache.Password = "cache-stack"

	service := InitService(requestSessionTest)
	defer service.Cleanup()

	// 1. Acquire must be successful
	l1, ok := service.Acquire(Options{
		Key:      "foo",
		Blocking: false,
	})
	require.True(t, ok, "Acquire should succeed")

	// 2. Acquire must return false if lock is already acquired
	_, ok2 := service.Acquire(Options{
		Key:      "foo",
		Blocking: false,
	})
	require.False(t, ok2, "Second Acquire should fail as the lock is already acquired")

	// 3. Acquire must be successful again after releasing the lock
	service.Release(l1)
	l3, ok3 := service.Acquire(Options{
		Key:      "foo",
		Blocking: false,
	})
	require.True(t, ok3, "Acquire should succeed again after releasing the lock")
	service.Release(l3)
}

func _TestBlockingAcquireTimeout(t *testing.T) {
	conf.Config.Cache.Host = "localhost"
	conf.Config.Cache.Password = "cache-stack"

	service := InitService(requestSessionTest)
	defer service.Cleanup()

	// Hold the lock
	l, ok := service.Acquire(Options{
		Key:      "bar",
		Timeout:  time.Second * 10,
		Blocking: false,
	})
	require.True(t, ok, "Initial Acquire should succeed")
	defer service.Release(l)

	// Blocking Acquire with short-term timeout
	start := time.Now()
	_, ok2 := service.Acquire(Options{
		Key:      "bar",
		Blocking: true,
		Timeout:  time.Millisecond * 1000,
	})
	dur := time.Since(start)
	require.False(t, ok2, "Blocking Acquire should fail after timeout")
	t.Logf("Acquire took %v", dur)
	require.GreaterOrEqual(t, dur.Milliseconds(), int64(900), "Blocking Acquire was aborted too early: %v", dur)
	t.Logf("Blocking Acquire was aborted as expected after timeout")
}

func _TestTTLExpiry(t *testing.T) {
	conf.Config.Cache.Host = "localhost"
	conf.Config.Cache.Password = "cache-stack"

	service := InitService(requestSessionTest)
	defer service.Cleanup()

	// Lock with short TTL
	_, ok := service.Acquire(Options{
		Key:      "baz",
		Blocking: false,
	})
	require.True(t, ok, "Initial Acquire should succeed")

	// Wait for TTL without to release
	time.Sleep(1100 * time.Millisecond)

	// Release must be called, otherwise the lock will not be released, as it has auto-renewal
	// service.Release(l1)

	// New Acquire should work
	l2, ok2 := service.Acquire(Options{
		Key:      "baz",
		Blocking: false,
	})
	require.True(t, ok2, "Acquire after TTL-Expire should be successful")
	service.Release(l2)
}

func _Test_DoLockEntity(t *testing.T) {
	conf.Config.Cache.Host = "localhost"
	conf.Config.Cache.Password = "cache-stack"

	servce := InitService(requestSessionTest)

	l, acquired := servce.Acquire(Options{
		Key:      "importData",
		Blocking: true,
		Timeout:  time.Minute,
	})
	if !acquired {
		assert.Fail(t, "lock failed")
	}
	servce.Release(l)
	_, acquired = servce.Acquire(Options{
		Key:      "importData",
		Blocking: true,
		Timeout:  time.Minute,
	})
	if !acquired {
		assert.Fail(t, "lock failed 2")
	}
	_, acquired = servce.Acquire(Options{
		Key:      "importData",
		Blocking: false,
	})
	if acquired {
		assert.Fail(t, "lock failed without blocking")
	}
	_, acquired = servce.Acquire(Options{
		Key:      "importData",
		Blocking: true,
		Timeout:  time.Millisecond * 400,
	})
	if acquired {
		assert.Fail(t, "lock failed with blocking. This should not happen")
	}
	servce.Release(l)
}
