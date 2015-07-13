// Main REST API Service Functionality

package registry

import (
	"net/http"

	"github.com/thisissoon/sir"
	"github.com/thisissoon/sir/registry/register"
	"github.com/thisissoon/sir/registry/unregister"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func JSONContentType(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Serves the HTTP Application
func Serve(a *sir.ApplicationContext) {
	// Create new Web Client
	r := web.New()

	// Register Routes
	r.Get("/", sir.ApplicationHandler{a, statsHandler})
	r.Post("/", sir.ApplicationHandler{a, register.RegisterHandler})
	r.Get("/:instance_id", sir.ApplicationHandler{a, register.DetailsHandler})
	r.Delete("/:instance_id", sir.ApplicationHandler{a, unregister.UnRegisterHandler})

	// Use Json Middleware
	r.Use(JSONContentType)

	// Serve the Application
	graceful.ListenAndServe(":8000", r)
}
