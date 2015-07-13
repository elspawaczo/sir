// Application HTTP Handler

package sir

import (
	"encoding/json"
	"net/http"

	"github.com/zenazn/goji/web"
)

type ErrorRespone struct {
	StatusCode int    `json:"code"`
	StatusText string `json:"message"`
}

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
		ErrorHandler(status, err, w)
	}
}

// Error Handler
func ErrorHandler(status int, err error, w http.ResponseWriter) {
	body, _ := json.Marshal(&ErrorRespone{
		StatusCode: status,
		StatusText: http.StatusText(status),
	})
	w.WriteHeader(status)
	w.Write(body)
}
