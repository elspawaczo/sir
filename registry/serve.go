// Main REST API Service Functionality

package registry

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/thisissoon/sir"
	"github.com/thisissoon/sir/registry/register"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

type StatsResponse struct {
	Available int64  `json:"available"`
	Taken     int64  `json:"taken"`
	Remaining string `json:"remaining"`
}

type AllocateResponse struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// JSON request data structure for registering a new instance
type AllocateRequest struct {
	InstanceID string `json:instance_id`
	PrivateIP  string `josn:private_ip`
}

// Put the name back in the pool
func DeRegister(a *sir.ApplicationContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	var err error

	// Does it exist in the taken pool
	exists, err := a.Redis.SIsMember(a.AllocatedKey, c.URLParams["name"]).Result()
	if err != nil || !exists {
		return 404, errors.New(fmt.Sprintf("%s not allocated", c.URLParams["name"]))
	}

	// Remove it from the taken pool
	a.Redis.SRem(a.AllocatedKey, c.URLParams["name"])
	// Add it to the pool
	a.Redis.SAdd(a.PoolKey, c.URLParams["name"])

	// Write Response
	w.WriteHeader(204)
	w.Write([]byte{})

	return 204, nil
}

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
	r.Delete("/:name", sir.ApplicationHandler{a, DeRegister})

	// Use Json Middleware
	r.Use(JSONContentType)

	// Serve the Application
	graceful.ListenAndServe(":8000", r)
}
