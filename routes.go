package main

import (
	"net/http"

	"github.com/hrz8/simpath/handler"
)

func addRoutes(mux *http.ServeMux, hdl *handler.Handler) {
	// web
	mux.Handle("GET /v1/login", hdl.UseSession(hdl.UseUserSession(hdl.UseForm(hdl.CheckClientID(hdl.GuestOnly(http.HandlerFunc(hdl.LoginFormHandler)))))))
	mux.Handle("GET /v1/register", hdl.UseSession(hdl.UseUserSession(hdl.UseForm(hdl.CheckClientID(hdl.GuestOnly(http.HandlerFunc(hdl.RegisterFormHandler)))))))
	mux.Handle("GET /v1/authorize", hdl.UseSession(hdl.UseUserSession(hdl.UseForm(hdl.CheckClientID(hdl.LoggedInOnly(http.HandlerFunc(hdl.AuthorizeFormHandler)))))))
	mux.Handle("GET /v1/logout", hdl.UseSession(hdl.UseUserSession(hdl.UseForm(hdl.CheckClientID(hdl.LoggedInOnly(http.HandlerFunc(hdl.LogoutPage)))))))

	// backend - form handler
	mux.Handle("POST /v1/login", hdl.UseSession(hdl.UseUserSession(hdl.UseForm(hdl.CheckCSRFToken(hdl.CheckClientID(hdl.GuestOnly(http.HandlerFunc(hdl.LoginHandler))))))))
	mux.Handle("POST /v1/authorize", hdl.UseSession(hdl.UseUserSession(hdl.UseForm(hdl.CheckCSRFToken(hdl.CheckClientID(hdl.LoggedInOnly(http.HandlerFunc(hdl.AuthorizeHandler))))))))

	// backend - json
	mux.HandleFunc("POST /v1/oauth/tokens", hdl.TokenHandler)
	mux.HandleFunc("POST /v1/oauth/introspect", hdl.IntrospectHandler)
}
