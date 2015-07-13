// States Requests & Response Handling
// This is a GET request on the root /

package registry

import (
	"encoding/json"
	"net/http"

	"github.com/thisissoon/sir"
	"github.com/zenazn/goji/web"
)

// JSON Data structure for a stats reqyest response
type statsResponse struct {
	Total     int64 `json:"total"`
	Taken     int64 `json:"taken"`
	Remaining int64 `json:"remaining"`
}

// Return basic stats about the registry, such as remaining hostnames
// total taken and total available
func statsHandler(a *sir.ApplicationContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	// Number of available names in the pool
	remaining, _ := a.Redis.SCard(a.PoolKey).Result()
	// Number of taken names in the pool
	taken, _ := a.Redis.SCard(a.AllocatedKey).Result()
	// Remaining
	total := remaining + taken

	resp, _ := json.Marshal(&statsResponse{
		Total:     total,
		Taken:     taken,
		Remaining: remaining,
	})

	status := 200

	w.WriteHeader(status)
	w.Write(resp)

	return status, nil
}
