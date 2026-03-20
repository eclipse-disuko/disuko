// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/mail"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/middlewareDisco"
	"github.com/eclipse-disuko/disuko/infra/service/startup"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
	"github.com/go-chi/chi/v5"
)

var isShuttingDown atomic.Bool

type Server struct {
	r          chi.Router
	repos      dbRepos
	connectors connectors
	handlers   handlers
	services   services

	scheduler *scheduler.Scheduler

	// s3 rest.S3

	mailClient mail.Client

	basicauthMW *middlewareDisco.InternalTokenMW
}

func StartServer(ctx context.Context) {
	defer exception.CatchException(exception.NewExceptionHandlerFatal(nil))

	requestSession := &logy.RequestSession{ReqID: logy.InternalReqID}
	logy.Infof(requestSession, "Starting disclosure application")

	server := createServer(ctx, requestSession, nil, nil)
	server.serve(requestSession)
}

func createServer(ctx context.Context, rs *logy.RequestSession, routeExt []RouteExtender, migExt []startup.Step) Server {
	server := Server{
		r: chi.NewRouter(),
		mailClient: mail.NewClient(
			conf.Config.Smtp.Host,
			conf.Config.Smtp.Port,
			conf.Config.Smtp.Sender,
			conf.Config.Smtp.User,
			conf.Config.Smtp.Pass,
		),
	}
	server.setupDatabase(rs)
	server.setupS3(rs) // is maybe used during db migration
	server.setupConnector(rs)
	server.setupServices(rs)
	server.setUpOAuth(rs)
	server.setupHandlers()
	server.setupMW()
	setupI18N()
	server.setupRoutes(routeExt...)
	server.setupScheduling(ctx, rs)
	server.registerObserver()
	if conf.Config.Database.MigrateOnly || conf.Config.Server.Env == "local" {
		server.migrateDatabase(rs, migExt...)
		if conf.Config.Database.MigrateOnly {
			os.Exit(0)
		}
	}
	return server
}

func (s *Server) serve(requestSession *logy.RequestSession) {
	isShuttingDown.Store(false)
	server := &http.Server{Addr: ":" + conf.Config.Server.Port, Handler: s.r}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		isShuttingDown.Store(true)
		shutdownCtx, _ := context.WithTimeout(serverCtx, time.Duration(conf.Config.Server.TerminationGracePeriodSeconds)*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				logy.Fatalf(requestSession, "graceful shutdown timed out after %v seconds", conf.Config.Server.TerminationGracePeriodSeconds)
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logy.Fatalf(requestSession, "server shutdown failed : %s", err)
		} else {
			logy.Infof(requestSession, "graceful shutdown of server complete")
		}
		s.scheduler.Stop()
		logy.Infof(requestSession, "waiting for scheduler to finish")
		s.scheduler.Wait()
		logy.Infof(requestSession, "scheduler finished")
		serverStopCtx()
	}()

	logy.Infof(requestSession, "Start WebServer on port: %s tls:%v", conf.Config.Server.Port, conf.Config.Server.Tls)
	err := listenAndServe(server)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logy.Fatalf(requestSession, "listen: %s\n", err)
	}

	<-serverCtx.Done()
}

func listenAndServe(server *http.Server) error {
	if conf.Config.Server.Tls {
		return server.ListenAndServeTLS("server.crt", "server.key")
	} else {
		return server.ListenAndServe()
	}
}

func setupI18N() {
	message.InitI18N()
}
