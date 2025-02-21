package settings

import (
	"encoding/json"
	"log"
	"os"
)

func retrieveSettings(path string) *Settings {
	ctn, err := os.ReadFile(path)
	if err != nil {
		log.Printf("error reading settings file: %v\n", err)
		return NewSettings()
	}

	var settings *Settings

	if err := json.Unmarshal(ctn, &settings); err != nil {
		log.Printf("error unmarshalling settings file: %v\n", err)
		settings = NewSettings()
	}

	settings.Normalize()
	return settings
}
