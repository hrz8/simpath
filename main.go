package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/hrz8/simpath/database"
	"github.com/hrz8/simpath/handler"
	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/scope"
	"github.com/hrz8/simpath/internal/token"
	"github.com/hrz8/simpath/internal/user"
	"github.com/hrz8/simpath/session"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	db, err := database.ConnectDB("postgres://postgres:toor@localhost:5432/simpath?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = database.RunMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	sessionSvc := session.NewService()
	userSvc := user.NewService(db)
	clientSvc := client.NewService(db)
	scopeSvc := scope.NewService(db)
	tokenSvc := token.NewService(db)
	hdl := handler.NewHandler(
		db,
		sessionSvc,
		userSvc,
		clientSvc,
		scopeSvc,
		tokenSvc,
	)

	// frontend
	mux.Handle("GET /v1/login", hdl.ShouldHaveClientID(hdl.GuestOnly(http.HandlerFunc(hdl.LoginFormHandler))))
	mux.HandleFunc("GET /v1/register", hdl.RegisterFormHandler)
	mux.Handle("GET /v1/authorize", hdl.ShouldHaveClientID(hdl.LoggedInOnly(http.HandlerFunc(hdl.AuthorizeFormHandler))))
	mux.Handle("GET /v1/logout", hdl.ShouldHaveClientID(hdl.LoggedInOnly(http.HandlerFunc(hdl.LogoutPage))))

	// backend
	mux.HandleFunc("POST /v1/login", hdl.LoginHandler)
	mux.HandleFunc("POST /v1/authorize", hdl.AuthorizeHandler)

	execCtx, execCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer execCancel()

	srv := &http.Server{
		Addr:    ":5001",
		Handler: mux,
	}
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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("err shutdown http server: %+v", err)
	}
}
