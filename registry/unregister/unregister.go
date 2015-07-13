// Removes an instance from the registry

package unregister

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/thisissoon/sir"
	"github.com/zenazn/goji/web"
	"gopkg.in/redis.v3"
)

func UnRegisterHandler(
	a *sir.ApplicationContext,
	c web.C,
	w http.ResponseWriter,
	r *http.Request) (int, error) {

	var err error
	instanceID := c.URLParams["instance_id"]
	instanceKey := fmt.Sprintf(a.InstanceKey, instanceID)
	status := 204

	// Get instance data
	data, err := a.Redis.Get(instanceKey).Result()
	if err == redis.Nil {
		return 404, errors.New(fmt.Sprintf("Instance %s not found", instanceID))
	} else if err != nil {
		return 500, err
	}

	// Decode the instance data
	i := &sir.Instance{}
	err = json.Unmarshal([]byte(data), i)
	if err != nil {
		return 500, err
	}

	// Remove it from the allocated pool
	err = a.Redis.SRem(a.AllocatedKey, i.Name).Err()
	if err != nil {
		return 500, err
	}
	// Add it to the pool
	err = a.Redis.SAdd(a.PoolKey, i.Name).Err()
	if err != nil {
		return 500, err
	}
	// Remove instance key
	err = a.Redis.Del(instanceKey).Err()
	if err != nil {
		return 500, err
	}

	// Write an empty response
	w.WriteHeader(status)
	w.Write([]byte{})

	return 204, nil
}
