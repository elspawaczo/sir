// Application HTTP Handler

package sir

import (
	"log"
	"net/http"

	"github.com/zenazn/goji/web"
)

// Base HTTP Handler Type emending our Applicaton Conext Type
type ApplicationHandler struct {
	*ApplicationContext
	Handler func(*ApplicationContext, web.C, http.ResponseWriter, *http.Request) (int, error)
}

// Satisfy http.Handler, which Goji's web.Handler extends.
func (a ApplicationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.ServeHTTPC(web.C{}, w, r)
}

// Custom Serve HTTP handler which provides our application context to our
// route handler functions
func (h ApplicationHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	status, err := h.Handler(h.ApplicationContext, c, w, r)
	// Handler errors
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}
