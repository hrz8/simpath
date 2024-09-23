package main

import (
	"net/http"

	"github.com/hrz8/simpath/handler"
)

func addRoutes(mux *http.ServeMux, hdl *handler.Handler) {
	mux.Handle("GET /v1/login", hdl.ShouldHaveClientID(hdl.GuestOnly(http.HandlerFunc(hdl.LoginFormHandler))))
	mux.HandleFunc("GET /v1/register", hdl.RegisterFormHandler)
	mux.Handle("GET /v1/authorize", hdl.ShouldHaveClientID(hdl.LoggedInOnly(http.HandlerFunc(hdl.AuthorizeFormHandler))))
	mux.Handle("GET /v1/logout", hdl.ShouldHaveClientID(hdl.LoggedInOnly(http.HandlerFunc(hdl.LogoutPage))))

	// backend
	mux.Handle("POST /v1/login", hdl.ShouldHaveClientID(hdl.GuestOnly(http.HandlerFunc(hdl.LoginHandler))))
	mux.Handle("POST /v1/authorize", hdl.ShouldHaveClientID(hdl.LoggedInOnly(http.HandlerFunc(hdl.AuthorizeHandler))))
}
