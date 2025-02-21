package settings

import (
	"encoding/json"
	"log"
	"os"
)

func loadProfileCollection(path string) *ProfileCollection {
	ctn, err := os.ReadFile(path)
	if err != nil {
		log.Printf("error reading profile-collection file: %v\n", err)
		return NewProfileCollection()
	}

	var profiles *ProfileCollection

	if err := json.Unmarshal(ctn, &profiles); err != nil {
		log.Printf("error unmarshalling profile-collection file: %v\n", err)
		profiles = NewProfileCollection()
	}

	profiles.Normalize()
	return profiles
}
