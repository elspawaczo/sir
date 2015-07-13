// Registers an Instance with the Registry

package register

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/thisissoon/sir"
	"github.com/zenazn/goji/web"
	"gopkg.in/redis.v3"
	"gopkg.in/validator.v2"
)

// Request JSON Structure
type requestJSON struct {
	InstanceID string `json:"instance_id" validate:"nonzero"`
	PrivateIP  string `json:"private_ip" validate:"nonzero"`
}

// Response JSON Structure
type responseJSON struct {
	Name string `json:"name"`
}

// Instance data stored at sir:instance:i-124abc
type instance struct {
	InstanceID string `json:"instance_id"`
	PrivateIP  string `json"private_ip"`
	Name       string `json:"name"`
}

// Allocated a name from the Pool, returning a name for the new instance
func allocate(a *sir.ApplicationContext, d *requestJSON) (*instance, error) {
	var err error

	instanceData := &instance{}
	instanceKey := fmt.Sprintf(a.InstanceKey, d.InstanceID)

	// Does the instance already exist, if so lets just return the hostname it has
	// already been assigned
	data, err := a.Redis.Get(instanceKey).Result()
	if err == redis.Nil {
		err = nil // So we don't trigger an error later
		// Does not exist - allocate a new random name
		name, err := a.Redis.SRandMember(a.PoolKey).Result()
		if err != nil {
			log.Println("Failed to get Random Name")
			return nil, err
		}
		// Remove the name from the pool
		err = a.Redis.SRem(a.PoolKey, name).Err()
		if err != nil {
			log.Println("Failed to remove %s from pool", name)
			return nil, err
		}
		// Add it to the allocated pool
		err = a.Redis.SAdd(a.AllocatedKey, name).Err()
		if err != nil {
			log.Println("Failed to add %s to Allocated Pool", name)
			return nil, err
		}
		// Add instance data to the instance key
		instanceData.Name = name
		instanceData.PrivateIP = d.PrivateIP
		instanceData.InstanceID = d.InstanceID
		value, err := json.Marshal(instanceData)
		if err != nil {
			log.Println("Failed to Marshal Instance Data")
			return nil, err
		}
		err = a.Redis.Set(instanceKey, value, 0).Err()
		if err != nil {
			log.Println("Failed to Set Instance Key with Instance Data")
			return nil, err
		}
	} else if err != nil {
		// Error getting data from redis
		log.Panicln("Failed to get Get isntance data")
		return nil, err
	} else {
		// Decode the result from redis for the isntance
		err := json.Unmarshal([]byte(data), instanceData)
		if err != nil {
			log.Println("Failed to Unmarshal instance data")
			return nil, err
		}
	}

	return instanceData, err
}

// Register HTTP Handler Function
func RegisterHandler(
	a *sir.ApplicationContext,
	c web.C,
	w http.ResponseWriter,
	r *http.Request) (int, error) {

	// Decode the request JSON
	decoder := json.NewDecoder(r.Body)
	data := &requestJSON{}
	err := decoder.Decode(data)
	if err != nil {
		return 400, err
	}

	// Validate our data
	if err := validator.Validate(data); err != nil {
		return 422, errors.New("Invalid POST request")
	}

	i, err := allocate(a, data)
	if err != nil {
		return 500, err
	}

	resp, err := json.Marshal(&responseJSON{Name: i.Name})
	if err != nil {
		return 500, err
	}

	// Write the response
	w.Write(resp)

	return 204, nil
}
