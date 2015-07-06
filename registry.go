// Main REST API Service Functionality

package sir

import (
	"io"
	"net/http"

	"gopkg.in/zenazn/goji.v0"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World\n")
}

func Serve() {
	goji.Get("/", hello)
	goji.Serve()
}
