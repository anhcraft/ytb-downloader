package settings

import (
	"encoding/json"
	"log"
	"os"
)

func persistSettings(path string, settings *Settings) {
	data, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		log.Printf("error marshalling settings file: %v\n", err)
		return
	}
	err = os.WriteFile(path, data, 0777)
	if err != nil {
		log.Printf("error writing settings file: %v\n", err)
	}
}
