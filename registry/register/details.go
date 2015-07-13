// Handles retrieving Instance details from the Register

package register

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/redis.v3"

	"github.com/thisissoon/sir"
	"github.com/zenazn/goji/web"
)

// Detail HTTP Handler
func DetailsHandler(
	a *sir.ApplicationContext,
	c web.C,
	w http.ResponseWriter,
	r *http.Request) (int, error) {

	var err error
	instanceID := c.URLParams["instance_id"]
	instanceKey := fmt.Sprintf(a.InstanceKey, instanceID)
	status := 200

	// Get instance data
	data, err := a.Redis.Get(instanceKey).Result()
	if err == redis.Nil {
		return 404, errors.New(fmt.Sprintf("Instance %s not found", instanceID))
	} else if err != nil {
		return 500, err
	}

	// Write Response
	w.WriteHeader(status)
	w.Write([]byte(data))

	return status, err
}
