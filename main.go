package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/hrz8/simpath/database"
	"github.com/hrz8/simpath/handler"
	"github.com/hrz8/simpath/internal/authcode"
	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/scope"
	"github.com/hrz8/simpath/internal/token"
	"github.com/hrz8/simpath/internal/tokengrant"
	"github.com/hrz8/simpath/internal/user"
	"github.com/hrz8/simpath/session"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func newServer(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	sessionSvc := session.NewService()
	userSvc := user.NewService(db)
	clientSvc := client.NewService(db)
	scopeSvc := scope.NewService(db)
	tokenSvc := token.NewService(db)
	authCodeSvc := authcode.NewService(db)
	tokenGrantSvc := tokengrant.NewService(db, userSvc, tokenSvc, authCodeSvc)

	hdl := handler.NewHandler(
		db,
		sessionSvc,
		userSvc,
		clientSvc,
		scopeSvc,
		tokenSvc,
		authCodeSvc,
		tokenGrantSvc,
	)

	addRoutes(mux, hdl)

	return mux
}

func main() {
	execCtx, execCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer execCancel()

	db, err := database.ConnectDB("postgres://postgres:toor@localhost:5432/simpath?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = database.RunMigrations(db)
	if err != nil {
		log.Fatal(err)
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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("err shutdown http server: %+v", err)
	}
}
