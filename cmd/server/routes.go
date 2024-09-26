package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/hrz8/simpath/handler"
)

func addRoutes(r *chi.Mux, hdl *handler.Handler) {
	// /v1 router
	r.Route("/v1", func(r chi.Router) {
		// non json handlers
		r.Route("/", func(r chi.Router) {
			r.Use(hdl.UseSession, hdl.UseUserSession, hdl.UseForm, hdl.CheckClientID)

			// web - guest only
			r.Group(func(r chi.Router) {
				r.Use(hdl.GuestOnly)
				r.Get("/login", hdl.LoginFormHandler)
				r.Get("/register", hdl.RegisterFormHandler)
			})

			// web - logged in only
			r.Group(func(r chi.Router) {
				r.Use(hdl.LoggedInOnly)
				r.Get("/authorize", hdl.AuthorizeFormHandler)
				r.Get("/logout", hdl.LogoutPage)
			})

			// backend - form handler - csrf protection
			r.Group(func(r chi.Router) {
				r.Use(hdl.CheckCSRFToken)

				// guest only
				r.With(hdl.GuestOnly).Post("/login", hdl.LoginHandler)
				r.With(hdl.GuestOnly).Post("/register", hdl.RegisterHandler)

				// logged in only
				r.With(hdl.LoggedInOnly).Post("/authorize", hdl.AuthorizeHandler)
			})
		})

		// backend - form-urlencoded
		r.Route("/oauth2", func(r chi.Router) {
			r.Use(hdl.UseForm)
			r.Post("/token", hdl.TokenHandlerForm)
		})

		// backend - json
		r.Post("/oauth/tokens", hdl.TokenHandlerJSON)
		r.Post("/oauth/introspect", hdl.IntrospectHandler)
	})
}
