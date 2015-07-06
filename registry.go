// Main REST API Service Functionality

package sir

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/zenazn/goji/web"

	"gopkg.in/redis.v3"
	"gopkg.in/zenazn/goji.v0"
)

// Redis Keys
var (
	POOL_KEY  = "sir:pool"
	TAKEN_KEY = "sir:key"
)

type AllocateResponse struct {
	Name string `json:"name"`
}

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
func Register(w http.ResponseWriter, r *http.Request) {
	name := allocate()
	resp, _ := json.Marshal(&AllocateResponse{
		Name: name,
	})

	w.Write(resp)
}

// Put the name back in the pool
func DeRegister(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func JsonContentTypeMW(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
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

	goji.Use(JsonContentTypeMW)

	// Register Routes
	goji.Post("/", Register)
	goji.Delete("/:name", DeRegister)

	// Serve the Application
	goji.Serve()
}
