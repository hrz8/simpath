package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/hrz8/simpath/config"
	"github.com/hrz8/simpath/database"
	"github.com/hrz8/simpath/handler"
	"github.com/hrz8/simpath/internal/authcode"
	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/consent"
	"github.com/hrz8/simpath/internal/introspect"
	"github.com/hrz8/simpath/internal/scope"
	"github.com/hrz8/simpath/internal/token"
	"github.com/hrz8/simpath/internal/tokengrant"
	"github.com/hrz8/simpath/internal/user"
	"github.com/hrz8/simpath/session"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/cors"
)

func newServer(db *sql.DB) *chi.Mux {
	mux := chi.NewRouter()

	sessionSvc := session.NewService()
	userSvc := user.NewService(db)
	clientSvc := client.NewService(db)
	scopeSvc := scope.NewService(db)
	tokenSvc := token.NewService(db, userSvc, clientSvc, scopeSvc)
	authCodeSvc := authcode.NewService(db)
	tokenGrantSvc := tokengrant.NewService(db, scopeSvc, userSvc, tokenSvc, authCodeSvc)
	introspectSvc := introspect.NewService(db, userSvc, tokenSvc)
	consentSvc := consent.NewService(db)

	hdl := handler.NewHandler(
		db,
		sessionSvc,
		userSvc,
		clientSvc,
		scopeSvc,
		tokenSvc,
		authCodeSvc,
		tokenGrantSvc,
		introspectSvc,
		consentSvc,
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.AllowClient},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	mux.Use(c.Handler)
	addRoutes(mux, hdl)

	return mux
}

func main() {
	execCtx, execCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer execCancel()

	db, err := database.ConnectDB(config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if config.AutoMigrate {
		err = database.RunMigrations(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	srv := &http.Server{Addr: ":5001", Handler: newServer(db)}
	srvErr := make(chan error)
	go func() {
		fmt.Println("server started")
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case e := <-srvErr:
		if e != http.ErrServerClosed {
			log.Fatalf("http server listen error: %+v", err)
		}
	case <-execCtx.Done():
		fmt.Println("shutdown...")
	}

	cleanup(srv, db)
}
