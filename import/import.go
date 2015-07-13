// Imports a txt file that will contain line seperated pool of
// server naames.

package importer

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/thisissoon/sir"
)

// Imports a txt file, reading each line and saving each line
// into a redis set
func Import(a *sir.ApplicationContext, p *string) {

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
		exists, _ := a.Redis.SIsMember(a.AllocatedKey, v).Result()
		if !exists {
			exists, _ := a.Redis.SIsMember(a.PoolKey, v).Result()
			if !exists {
				err := a.Redis.SAdd(a.PoolKey, v).Err()
				if err != nil {
					log.Println(err)
				} else {
					log.Println(fmt.Sprintf("Added: %s", v))
				}
			}
		}
	}

	log.Println("Done")

}
