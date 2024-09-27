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
				r.Get("/logout", hdl.LogoutPage)
			})

			// backend - form handler - csrf protection
			r.Group(func(r chi.Router) {
				r.Use(hdl.CheckCSRFToken)

				// guest only
				r.With(hdl.GuestOnly).Post("/login", hdl.LoginHandler)
				r.With(hdl.GuestOnly).Post("/register", hdl.RegisterHandler)
			})
		})

		// oauth2 endpoints - mostly used for oauth2 library
		r.Route("/oauth2", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(hdl.UseSession, hdl.UseUserSession, hdl.UseForm, hdl.CheckClientID)
				r.Use(hdl.LoggedInOnly)

				// authorize page that serve allow - deny button
				// this used by oauth2.Endpoint.AuthURL
				r.Get("/authorize", hdl.AuthorizeFormHandler)
				r.With(hdl.CheckCSRFToken).Post("/authorize", hdl.AuthorizeHandler)
			})

			// token handler endpoint used by oauth2.Endpoint.TokenURL
			r.With(hdl.UseForm).Post("/token", hdl.TokenHandler)
			// token handler but for external purpose
			r.Post("/external/token", hdl.TokenHandlerJSON)
			r.Get("/userinfo", hdl.UserInfoHandler)
			r.Post("/introspect", hdl.IntrospectHandler)
		})

		r.Route("/.well-known", func(r chi.Router) {
			r.Get("/jwks.json", hdl.JWKSHandler)
			r.Get("/openid-configuration", hdl.OIDCConfig)
		})
	})
}
