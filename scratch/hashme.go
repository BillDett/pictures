package scratch

// Test for collisions on sample space of all paths to photos
// Verdict: FNV-a hashing is fine for what we need

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"hash/fnv"
	"os"
	"path/filepath"
)

var keys map[string]string

func hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum64())
}

func murmurhash(s string) string {
	h := murmur3.New64()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum64())
}

func main() {
	keys = make(map[string]string)
	err := filepath.Walk("/Volumes/Data", func(path string, info os.FileInfo, err error) error {
		//fmt.Println(path)
		if err == nil {
			thekey := hash(path)
			val, exists := keys[thekey]
			if exists {
				fmt.Println("Collision between " + path + " and " + val)
			} else {
				keys[thekey] = path
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("We had %d keys\n", len(keys))
}
