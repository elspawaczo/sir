// Application Context

package sir

import "gopkg.in/redis.v3"

// Application context, holds application configuration
type ApplicationContext struct {
	// Constants
	PoolKey      string
	AllocatedKey string
	InstanceKey  string
	// Connections
	Redis *redis.Client
}

// Constructs a new Application Context Instance
func NewApplicationContext(r *string) *ApplicationContext {
	// Connext to Redis
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    *r,
	})

	// Create Application Context
	return &ApplicationContext{
		PoolKey:      "sir:pool",
		AllocatedKey: "sir:allocated",
		InstanceKey:  "sir:instance:%s",
		Redis:        c,
	}
}
