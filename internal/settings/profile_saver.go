package settings

import (
	"encoding/json"
	"log"
	"os"
)

func saveProfileCollection(path string, profiles *ProfileCollection) {
	data, err := json.MarshalIndent(profiles, "", "    ")
	if err != nil {
		log.Printf("error marshalling profile-collection file: %v\n", err)
		return
	}
	err = os.WriteFile(path, data, 0777)
	if err != nil {
		log.Printf("error writing profile-collection file: %v\n", err)
	}
}
