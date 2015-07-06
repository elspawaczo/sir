// Imports a txt file that will contain line seperated pool of
// server naames.

package sir

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"gopkg.in/redis.v3"
)

// Imports a txt file, reading each line and saving each line
// into a redis set
func Import(p *string, r *string) {
	// Connect to Redis
	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    *r,
	})

	defer c.Close()

	// Open the file
	f, err := os.Open(*p)
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	// Scan the file
	scanner := bufio.NewScanner(f)
	// Split by lines
	scanner.Split(bufio.ScanLines)

	// Loop over the lines
	for scanner.Scan() {
		v := scanner.Text()
		err := c.SAdd("sir:pool", v).Err()
		if err != nil {
			log.Println(err)
		} else {
			log.Println(fmt.Sprintf("Added: %s", v))
		}
	}

}
