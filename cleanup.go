package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func cleanup(srv *http.Server, db *sql.DB) {
	var wg sync.WaitGroup
	fmt.Println("cleaning up...")

	// cleanup process maximum 30s to complete
	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cleanupCancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// http shutdown process maximum 10s to complete
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("err shutdown http server: %+v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := db.Close(); err != nil {
			fmt.Printf("err close database connection: %+v", err)
		}
	}()

	cleanupDone := make(chan struct{})
	go func() {
		defer close(cleanupDone)
		wg.Wait()
	}()

	select {
	case <-cleanupCtx.Done():
		fmt.Println("clean up done partially, because it takes longer than it should")
	case <-cleanupDone:
		fmt.Println("clean up done")
	}
}
