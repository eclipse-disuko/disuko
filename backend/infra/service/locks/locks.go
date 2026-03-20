// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package locks

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/valkey-io/valkey-go/valkeylock"

	"github.com/valkey-io/valkey-go"

	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type Options struct {
	Key      string
	Blocking bool
	Timeout  time.Duration
}

type Service struct {
	cacheClient valkey.Client
	locker      valkeylock.Locker
}

type Lock struct {
	cancel context.CancelFunc
}

func InitService(requestSession *logy.RequestSession) *Service {
	var res Service
	clientOption := valkey.ClientOption{
		InitAddress: []string{conf.Config.Cache.Host + ":" + strconv.Itoa(conf.Config.Cache.Port)},
		Password:    conf.Config.Cache.Password,
	}
	cacheClient, err := valkey.NewClient(clientOption)
	if err != nil {
		logy.Fatalf(requestSession, "valkey client creation failed: %s", err)
	}
	res.cacheClient = cacheClient
	locker, err := valkeylock.NewLocker(valkeylock.LockerOption{
		ClientOption:   clientOption,
		KeyMajority:    2,
		NoLoopTracking: true,
	})
	if err != nil {
		logy.Fatalf(requestSession, "valkey locker creation failed: %s", err)
	}
	res.locker = locker
	return &res
}

func (s *Service) Acquire(opt Options) (Lock, bool) {
	var (
		cancel context.CancelFunc
		err    error
	)
	if opt.Blocking {
		lockCtx, cancelLock := context.WithCancel(context.Background())
		lockErrCh := make(chan error)
		go func() {
			_, _, err := s.locker.WithContext(lockCtx, opt.Key)
			lockErrCh <- err
		}()

		to := time.NewTicker(opt.Timeout)
		select {
		case <-to.C:
			cancelLock()
			<-lockErrCh
			return Lock{}, false
		case err := <-lockErrCh:
			if err != nil {
				exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorAcquiringLock))
			}
			return Lock{
				cancel: cancelLock,
			}, true
		}

	}
	ctx := context.Background()
	_, cancel, err = s.locker.TryWithContext(ctx, opt.Key)

	if errors.Is(err, valkeylock.ErrNotLocked) || errors.Is(err, valkeylock.ErrLockerClosed) || errors.Is(err, context.DeadlineExceeded) {
		return Lock{}, false
	} else if err != nil {
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorAcquiringLock))
	}
	return Lock{
		cancel: cancel,
	}, true
}

func (s *Service) Release(lock Lock) {
	lock.cancel()
}

func (s *Service) Cleanup() {
	s.cacheClient.Close()
}
