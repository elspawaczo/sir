// Main REST API Service Functionality

package sir

import (
	"io"
	"log"
	"net/http"

	"gopkg.in/redis.v3"
	"gopkg.in/zenazn/goji.v0"
)

// Redis Keys
var (
	POOL_KEY  = "sir:pool"
	TAKEN_KEY = "sir:key"
)

// Holds the current redis client
var RedisClient *redis.Client

// Allocates a name to the server, ensuring we always get a unique one
func allocate() string {
	for {
		// Get a random name from the pool
		member, err := RedisClient.SRandMember(POOL_KEY).Result()
		if err != nil {
			log.Println(err)
		}
		// Exists in the taken set?
		exists, err := RedisClient.SIsMember(TAKEN_KEY, member).Result()
		if err != nil {
			log.Println(err)
		}
		// If not taken add to remove from the pool and add to taken set
		if !exists {
			err = RedisClient.SRem(POOL_KEY, member).Err()
			if err != nil {
				log.Println(err)
			}
			err = RedisClient.SAdd(TAKEN_KEY, member).Err()
			if err != nil {
				log.Println(err)
			}
			return member
		}
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World\n")
}

// Get a random name from the pool and place it in the taken pool,
// if it is already in the taken pool add a number to the name by the numer
// of times the name has been used
func AllocateName(w http.ResponseWriter, r *http.Request) {
	log.Println(allocate())
}

// Serves the HTTP Application
func Serve(r *string) {
	// Connect to Redis
	RedisClient = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    *r,
	})

	// Ensure we close redis
	defer RedisClient.Close()

	// Register Routes
	goji.Post("/", AllocateName)

	// Serve the Application
	goji.Serve()
}
