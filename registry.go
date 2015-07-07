// Main REST API Service Functionality

package sir

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zenazn/goji/web"

	"gopkg.in/redis.v3"
	"gopkg.in/zenazn/goji.v0"
)

// Redis Keys
var (
	POOL_KEY  = "sir:pool"
	TAKEN_KEY = "sir:allocated"
)

type StatsResponse struct {
	Available int64  `json:"available"`
	Taken     int64  `json:"taken"`
	Remaining string `json:"remaining"`
}

type AllocateResponse struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
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

// Return basic stats
func Stats(w http.ResponseWriter, r *http.Request) {
	// Number of available names in the pool
	avail, _ := RedisClient.SCard(POOL_KEY).Result()
	// Number of taken names in the pool
	taken, _ := RedisClient.SCard(TAKEN_KEY).Result()
	// Remaining
	remaining := float64(taken) / float64(avail+taken) * float64(100)

	resp, _ := json.Marshal(&StatsResponse{
		Available: avail,
		Taken:     taken,
		Remaining: fmt.Sprintf("%.2f%%", remaining),
	})

	w.WriteHeader(200)
	w.Write(resp)
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
	var err error
	// Does it exist in the taken pool
	exists, err := RedisClient.SIsMember(TAKEN_KEY, c.URLParams["name"]).Result()
	if err != nil || !exists {
		HTTPError(w, r, 404)
	}

	// Remove it from the taken pool
	RedisClient.SRem(TAKEN_KEY, c.URLParams["name"])
	// Add it to the pool
	RedisClient.SAdd(POOL_KEY, c.URLParams["name"])

	// Write Response
	w.WriteHeader(204)
	w.Write([]byte{})
}

func HTTPError(w http.ResponseWriter, r *http.Request, status int) {
	resp, _ := json.Marshal(&ErrorResponse{
		Error: http.StatusText(status),
	})

	w.WriteHeader(status)
	w.Write(resp)
}

func JSONContentType(c *web.C, h http.Handler) http.Handler {
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

	goji.Use(JSONContentType)

	// Register Routes
	goji.Get("/", Stats)
	goji.Post("/", Register)
	goji.Delete("/:name", DeRegister)

	// Serve the Application
	goji.Serve()
}
